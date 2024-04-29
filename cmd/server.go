package server

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/DanillaY/GoScrapper/cmd/models"
	"github.com/gin-gonic/gin"
)

func (d *Repository) InitAPIRoutes() {

	booksApi := gin.New()
	booksApi.Use(gin.Logger())
	booksApi.Use(CORSMiddleware())

	booksApi.SetTrustedProxies([]string{"localhost"})

	booksApi.GET("/getBooks", d.GetBooks)
	booksApi.GET("/getBooksById", d.GetBookByID)
	booksApi.Run(":" + d.Config.API_PORT)
}

func (d *Repository) GetBooks(context *gin.Context) {

	category := context.DefaultQuery("category", "%")
	if category != "%" {
		category = strings.ToLower(category)
	}

	pageNumber := context.DefaultQuery("pageNum", "1")
	limit := context.DefaultQuery("limit", "30")
	title := context.DefaultQuery("title", "%")
	author := context.DefaultQuery("author", "%")

	minPrice := context.DefaultQuery("minPrice", "50")
	maxPrice := context.DefaultQuery("maxPrice", "100000")

	limitInt, errLim := strconv.Atoi(limit)
	pageNumberInt, errPageNum := strconv.Atoi(pageNumber)
	if errLim != nil || errPageNum != nil {
		context.JSON(http.StatusBadRequest, errLim.Error()+" "+errPageNum.Error())
	}

	books := &[]models.Book{}

	err := d.Db.Scopes(FilterBooks(title, author, maxPrice, minPrice, category)).
		Find(&books).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}

	total := len(*books)
	lastpage := math.Ceil(float64(len(*books)) / float64(limitInt))

	err = d.Db.Scopes(FilterBooks(title, author, maxPrice, minPrice, category)).
		Order("id").Offset((pageNumberInt - 1) * limitInt).Limit(limitInt).Find(&books).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
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
