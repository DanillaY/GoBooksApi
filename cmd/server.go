package server

import (
	"math"
	"net/http"
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
	booksApi.Run(":" + d.Config.API_PORT)
}

func (d *Repository) GetBooks(context *gin.Context) {

	category := strings.ToLower(context.DefaultQuery("category", "%"))
	categoryRawSql := AddRegexToQuery(category, ",", strings.Contains(category, ","))

	author := strings.ToLower(context.DefaultQuery("author", "%"))
	authorRawSql := AddRegexToQuery(author, ",", strings.Contains(author, ","))

	pageNumber := context.DefaultQuery("pageNum", "1")
	limit := context.DefaultQuery("limit", "30")
	search := context.DefaultQuery("search", "%")
	sort := context.DefaultQuery("priceSort", "")

	search = AddRegexToQuery(search, " ", search != "%")

	minPrice := context.DefaultQuery("minPrice", "50")
	maxPrice := context.DefaultQuery("maxPrice", "100000")

	limitInt, errLim := strconv.Atoi(limit)
	pageNumberInt, errPageNum := strconv.Atoi(pageNumber)
	if errLim != nil || errPageNum != nil {
		context.JSON(http.StatusBadRequest, errLim.Error()+" "+errPageNum.Error())
	}

	books := &[]models.Book{}

	if db := d.Db.Scopes(FilterBooks(maxPrice, minPrice, categoryRawSql, search, authorRawSql)).Order("id").
		Find(&books); db.Error != nil {
		context.JSON(http.StatusBadRequest, db.Error)
	}

	total := len(*books)
	lastpage := math.Ceil(float64(len(*books)) / float64(limitInt))

	if db := d.Db.Scopes(FilterBooks(maxPrice, minPrice, categoryRawSql, search, authorRawSql)).
		Order("id").
		Offset((pageNumberInt - 1) * limitInt).
		Limit(limitInt).
		Find(&books); db.Error != nil {
		context.JSON(http.StatusBadRequest, db.Error)
	} else {

		SortByCurrentPrice(sort, books)
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
