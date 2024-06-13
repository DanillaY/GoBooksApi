package server

import (
	"fmt"
	"math"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/DanillaY/GoScrapper/cmd/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (d *Repository) InitAPIRoutes() {

	booksApi := gin.New()
	booksApi.Use(gin.Logger())
	booksApi.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	booksApi.SetTrustedProxies([]string{"localhost"})

	booksApi.GET("/getBooks", d.GetBooks)
	booksApi.GET("/getBooksById", d.GetBookByID)
	booksApi.GET("/getProperties", d.GetProperties)
	booksApi.GET("/getMinMaxPrice", d.GetMinMaxPrice)
	booksApi.GET("/getBooksByEmail", d.GetBooksByEmail)
	booksApi.GET("/getBookByStructuredSearch", d.GetBookByStructuredSearch)
	booksApi.DELETE("/deleteBookSubscriber", d.DeleteBookSubscriber)
	booksApi.POST("/addNewBookSubscriber", d.AddNewBookSubscriber)
	booksApi.Run(":" + d.Config.API_PORT)
}

func (d *Repository) GetBooks(context *gin.Context) {

	category := context.DefaultQuery("category", "")
	author := context.DefaultQuery("author", "")
	bookInfoType := context.DefaultQuery("bookInfoType", "partial")

	vendor := context.DefaultQuery("vendor", "")
	yearPublished, _ := strconv.Atoi(context.DefaultQuery("year", "0"))

	pageNumber := context.DefaultQuery("pageNum", "1")
	limit := context.DefaultQuery("limit", "30")
	search := strings.TrimSpace(strings.ToLower(context.DefaultQuery("search", "")))
	stockText := context.DefaultQuery("stockText", "")
	sort := context.DefaultQuery("sortOrder", "")
	sortField := context.DefaultQuery("sortField", "")
	if sortField != "" {
		sortField += " " + sort + ","
	}

	minPrice, _ := strconv.Atoi(context.DefaultQuery("minPrice", "50"))
	maxPrice, _ := strconv.Atoi(context.DefaultQuery("maxPrice", "100000"))

	limitInt, errLim := strconv.Atoi(limit)
	pageNumberInt, errPageNum := strconv.Atoi(pageNumber)
	if errLim != nil || errPageNum != nil {
		context.JSON(http.StatusInternalServerError, errLim.Error()+" "+errPageNum.Error())
	}

	books := &[]models.Book{}

	db := d.Db.Scopes(FilterBooks(maxPrice,
		minPrice, category,
		search, author,
		vendor, yearPublished, stockText, bookInfoType)).Order(sortField + "id").
		Find(&books)

	total := len(*books)
	lastpage := math.Ceil(float64(len(*books)) / float64(limitInt))

	if db.Offset((pageNumberInt - 1) * limitInt).Limit(limitInt).Find(&books); db.Error != nil {
		context.JSON(http.StatusBadRequest, db.Error)
	} else {
		pagination := &Pagination{
			Total:       total,
			PerPage:     len(*books),
			CurrentPage: pageNumberInt,
			LastPage:    int(lastpage),
		}

		context.JSON(http.StatusOK, gin.H{"Pagination": pagination, "Data": books})
	}
}

func (d *Repository) GetBookByStructuredSearch(context *gin.Context) {
	search := context.DefaultQuery("search", "")
	category := context.DefaultQuery("category", "")
	author := context.DefaultQuery("author", "")

	var filter []string
	filter = AppendToSearchIfNotEmpty(search, filter)
	filter = AppendToSearchIfNotEmpty(category, filter)
	filter = AppendToSearchIfNotEmpty(author, filter)

	books := []models.Book{}
	ts := "ts_rank(search, websearch_to_tsquery('simple', '" + search + " " + category + " " + author + "' )) + ts_rank(search, websearch_to_tsquery('russian', '" + search + " " + category + " " + author + "' )) as rank"
	dbErr := d.Db.Table("books").Select("*", ts).
		Where(strings.Join(filter, " and ")).
		Order("rank DESC").Find(&books).Error

	if dbErr != nil {
		context.JSON(http.StatusBadRequest, dbErr.Error())
	} else {
		context.JSON(http.StatusOK, gin.H{"books": &books})
	}
}

func (d *Repository) GetBookByID(context *gin.Context) {
	bookId := context.DefaultQuery("id", "1")
	book := &models.Book{}
	err := d.Db.First(&book, bookId).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	} else {
		context.JSON(http.StatusOK, &book)
	}
}

func (d *Repository) GetMinMaxPrice(context *gin.Context) {
	min := 0
	max := 0
	errMin := d.Db.Table("books").Select("MIN(current_price)").Where("current_price <> 0").Find(&min).Error
	errMax := d.Db.Table("books").Select("MAX(current_price)").Find(&max).Error

	mapResult := make(map[string]int)
	mapResult["minPrice"] = min
	mapResult["maxPrice"] = max

	if errMin != nil || errMax != nil {
		context.JSON(http.StatusBadRequest, errMin)
	} else {
		context.JSON(http.StatusOK, mapResult)
	}
}

func (d *Repository) GetProperties(context *gin.Context) {
	property := strings.ToLower(context.DefaultQuery("property", "author"))
	var result []string
	var err error

	if property == "year_publish" || property == "current_price" || property == "old_price" {
		err = d.Db.Model(&models.Book{}).Distinct(property).Order(property+" DESC").Where(property+" <> 0").Pluck(property, &result).Error
	} else {
		err = d.Db.Model(&models.Book{}).Distinct(property).Order(property+" ASC").Where(property+" <> ''").Pluck(property, &result).Error
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	} else {
		context.JSON(http.StatusOK, result)
	}

}

func (d *Repository) GetBooksByEmail(context *gin.Context) {
	email, errMail := mail.ParseAddress(context.DefaultQuery("userEmail", "@"))
	user := models.User{}
	books := []models.Book{}
	errUser := d.Db.Model(models.User{}).Find(&user, "email = ?", email.Address).Error
	d.Db.Model(&user).Association("Book").Find(&books)

	if errUser != nil {
		context.JSON(http.StatusInternalServerError, "No such user")
	} else if errMail != nil {
		context.JSON(http.StatusInternalServerError, "Email is not valid")
	} else {
		context.JSON(http.StatusOK, gin.H{"books": &books})
	}
}

func (d *Repository) DeleteBookSubscriber(context *gin.Context) {
	bookId := context.DefaultQuery("bookId", "1")
	email, errMail := mail.ParseAddress(context.DefaultQuery("userEmail", "@"))

	book := models.Book{}
	user := models.User{}
	errBook := d.Db.Find(&book, "ID = ?", bookId).Error
	errUser := d.Db.Find(&user, "email = ?", email).Error

	fmt.Println(d.Db.Model(&book).Association("User").Count())
	fmt.Println(book.ID)
	fmt.Println(user.Email)

	if errBook != nil || errUser != nil || errMail != nil {
		context.JSON(http.StatusBadRequest, "Server error")

	} else if (book.ID == 0 || email.Address == "") || d.Db.Model(&book).Association("User").Count() == 0 {
		context.JSON(http.StatusBadRequest, "No such user or book")

	} else if book.Vendor == "Book24" || book.Vendor == "Читай город" || book.Vendor == "Лабиринт" {
		d.Db.Model(&book).Association("User").Delete(&user)
		context.JSON(http.StatusOK, "User was successfully unsubscribed")

	} else {
		context.JSON(http.StatusMethodNotAllowed, "Subscribing/Unsubscribing function for this vendor is prohibited")
	}
}

func (d *Repository) AddNewBookSubscriber(context *gin.Context) {
	bookId := context.DefaultQuery("bookId", "0")
	email, errMail := mail.ParseAddress(context.DefaultQuery("userEmail", "@"))

	var errBook error
	var errUser error
	book := models.Book{}
	user := models.User{}

	if errMail == nil {
		user = models.User{Email: email.Address}
		errBook = d.Db.Find(&book, "ID = ?", bookId).Error
		errUser = d.Db.Where(models.User{Email: email.Address}).FirstOrCreate(&user).Error
	}

	if errBook != nil || bookId == "0" || errUser != nil || errMail != nil {
		context.JSON(http.StatusBadRequest, "Server error")
	} else if book.Vendor == "Book24" || book.Vendor == "Читай город" || book.Vendor == "Лабиринт" {
		d.Db.Model(&book).Association("User").Append(&user)
		context.JSON(http.StatusOK, "User subscription was successfully complete")
	} else {
		context.JSON(http.StatusMethodNotAllowed, "Subscribing/Unsubscribing function for this vendor is prohibited")
	}

}
