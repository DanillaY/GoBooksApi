package server

import (
	"fmt"
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
	booksApi.Run(":" + d.Config.API_PORT)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
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
	if author != "%" {
		author = strings.ToLower(author)
		fmt.Println(author)
	}
	minPrice := context.DefaultQuery("minPrice", "50")
	maxPrice := context.DefaultQuery("maxPrice", "100000")

	limitInt, errLim := strconv.Atoi(limit)
	pageNumberInt, errPageNum := strconv.Atoi(pageNumber)
	if errLim != nil || errPageNum != nil {
		context.JSON(http.StatusBadRequest, errLim.Error()+" "+errPageNum.Error())
	}

	books := &[]models.Book{}

	err := d.Db.
		Where("LOWER(category) LIKE ?", "%"+category+"%").
		Where("LOWER(title) LIKE ?", "%"+title+"%").
		Where("LOWER(author) LIKE ?", "%"+author+"%").
		Where("current_price >= ?", minPrice).
		Where("current_price <= ?", maxPrice).
		Find(&books).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	total := len(*books)
	lastpage := math.Ceil(float64(len(*books)) / float64(limitInt))

	err = d.Db.
		Where("ID >= ?", pageNumberInt*limitInt-limitInt).
		Where("ID <= ?", pageNumberInt*limitInt).Find(&books).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	} else {

		pagination := &Pagination{
			Total:       total,
			PerPage:     limitInt,
			CurrentPage: pageNumberInt,
			LastPage:    int(lastpage),
		}
		result := ApiAnswer{Data: books, Pagination: pagination}

		context.JSON(http.StatusOK, &result)
	}
}
