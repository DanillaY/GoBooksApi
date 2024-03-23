package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DanillaY/GoScrapper/cmd/models"
	"github.com/gin-gonic/gin"
)

func (d *Repository) InitAPIRoutes() {

	booksApi := gin.New()
	booksApi.Use(gin.Logger())

	booksApi.SetTrustedProxies([]string{"localhost"})

	booksApi.GET("/getBooks", d.GetBooks)
	booksApi.Run(":" + d.Config.API_PORT)
}

func (d *Repository) GetBooks(context *gin.Context) {

	category := context.DefaultQuery("category", "%")
	if category != "%" {
		category = strings.ToLower(category)
	}

	limit := context.DefaultQuery("limit", "30") // you could send -1 if you want to get all records
	title := context.DefaultQuery("title", "%")
	author := context.DefaultQuery("author", "%")
	if author != "%" {
		author = strings.ToLower(author)
		fmt.Println(author)
	}
	minPrice := context.DefaultQuery("minPrice", "50")
	maxPrice := context.DefaultQuery("maxPrice", "100000")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}

	books := &[]models.Book{}

	err = d.Db.
		Where("LOWER(category) LIKE ?", "%"+category+"%").
		Where("LOWER(title) LIKE ?", "%"+title+"%").
		Where("LOWER(author) LIKE ?", "%"+author+"%").
		Where("current_price >= ?", minPrice).
		Where("current_price <= ?", maxPrice).
		Limit(limitInt).
		Find(&books).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	} else {
		context.JSON(http.StatusOK, &books)
	}
}
