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
		AllowMethods:     []string{"GET"},
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
	booksApi.DELETE("/deleteBookSubscriber", d.DeleteBookSubscriber)
	booksApi.POST("/addNewBookSubscriber", d.AddNewBookSubscriber)
	booksApi.Run(":" + d.Config.API_PORT)
}

func (d *Repository) GetBooks(context *gin.Context) {

	category := strings.ToLower(context.DefaultQuery("category", "%"))
	categoryRawSql := AddRegexToQuery(category, ",", strings.Contains(category, ","))

	author := strings.ToLower(context.DefaultQuery("author", "%"))
	authorRawSql := strings.ToLower(AddRegexToQuery(author, ",", strings.Contains(author, ",")))

	vendor := strings.ToLower(context.DefaultQuery("vendor", "%"))
	yearPublished := context.DefaultQuery("year", "0")

	pageNumber := context.DefaultQuery("pageNum", "1")
	limit := context.DefaultQuery("limit", "30")
	search := strings.ToLower(context.DefaultQuery("search", "%"))
	stockText := strings.ToLower(context.DefaultQuery("stockText", "%"))
	sort := context.DefaultQuery("sortOrder", "")
	sortField := context.DefaultQuery("sortField", "")
	if sortField != "" {
		sortField += " " + sort + ","
	}
	fmt.Println(sortField)

	search = AddRegexToQuery(search, " ", search != "%")

	minPrice := context.DefaultQuery("minPrice", "50")
	maxPrice := context.DefaultQuery("maxPrice", "100000")

	limitInt, errLim := strconv.Atoi(limit)
	pageNumberInt, errPageNum := strconv.Atoi(pageNumber)
	if errLim != nil || errPageNum != nil {
		context.JSON(http.StatusInternalServerError, errLim.Error()+" "+errPageNum.Error())
	}

	books := &[]models.Book{}

	if db := d.Db.Scopes(FilterBooks(maxPrice,
		minPrice,
		categoryRawSql,
		search,
		authorRawSql,
		vendor,
		yearPublished, stockText)).Order("id").
		Find(&books); db.Error != nil {
		context.JSON(http.StatusBadRequest, db.Error)
	}

	total := len(*books)
	lastpage := math.Ceil(float64(len(*books)) / float64(limitInt))

	if db := d.Db.Scopes(FilterBooks(maxPrice, minPrice, categoryRawSql, search, authorRawSql, vendor, yearPublished, stockText)).
		Order(sortField + "id").
		Offset((pageNumberInt - 1) * limitInt).
		Limit(limitInt).
		Find(&books); db.Error != nil {
		context.JSON(http.StatusBadRequest, db.Error)
	} else {
		pagination := &Pagination{
			Total:       total,
			PerPage:     len(*books),
			CurrentPage: pageNumberInt,
			LastPage:    int(lastpage),
		}

		result := ApiAnswer{Data: books, Pagination: pagination}
		context.JSON(http.StatusOK, &result)
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
	errMin := d.Db.Table("books").Select("MIN(current_price)").Find(&min).Error
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
	err := d.Db.Model(&models.Book{}).Distinct(property).Pluck(property, &result).Error

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

	if errBook != nil || errUser != nil || errMail != nil {
		context.JSON(http.StatusBadRequest, "Server error")

	} else if (book.ID == 0 || user.Email == "") || d.Db.Model(&book).Association("User").Count() == 0 {
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

	book := models.Book{}
	user := models.User{Email: email.Address}

	errBook := d.Db.Find(&book, "ID = ?", bookId).Error
	errUser := d.Db.Where(models.User{Email: email.Address}).FirstOrCreate(&user).Error

	if errBook != nil || bookId == "0" || errMail != nil {
		context.JSON(http.StatusBadRequest, "Error while getting values")
	} else if errUser != nil {
		context.JSON(http.StatusInternalServerError, "Could not create user")
	} else if book.Vendor == "Book24" || book.Vendor == "Читай город" || book.Vendor == "Лабиринт" {
		d.Db.Model(&book).Association("User").Append(&user)
		context.JSON(http.StatusOK, "User subscription was successfully complete")
	} else {
		context.JSON(http.StatusMethodNotAllowed, "Subscribing/Unsubscribing function for this vendor is prohibited")
	}

}
