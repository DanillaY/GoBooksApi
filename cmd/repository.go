package server

import (
	"fmt"
	"strings"

	"github.com/DanillaY/GoScrapper/cmd/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	HOST     string
	DB_PORT  string
	API_PORT string
	PASSWORD string
	USER     string
	DB       string
	SSLMODE  string
}

type Repository struct {
	Db     *gorm.DB
	Config *Config
}

type Pagination struct {
	Total       int
	PerPage     int
	CurrentPage int
	LastPage    int
}
type ApiAnswer struct {
	Pagination *Pagination    `json:"pagination"`
	Data       *[]models.Book `json:"data"`
}

func NewPostgresConnection(c *Config) (db *gorm.DB, e error) {
	db, err := gorm.Open(postgres.Open(
		"host="+c.HOST+
			" port="+c.DB_PORT+
			" password="+c.PASSWORD+
			" user="+c.USER+
			" dbname="+c.DB+
			" sslmode="+c.SSLMODE), &gorm.Config{})
	if err != nil {
		fmt.Println("Error while opening the connection to database")
		return db, err
	}
	return db, nil
}

func FilterBooks(title string, author string, maxPrice string, minPrice string, category string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(category) SIMILAR TO ?", category).
			Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%").
			Where("LOWER(author) LIKE ?", "%"+strings.ToLower(author)+"%").
			Where("current_price >= ?", minPrice).
			Where("current_price <= ?", maxPrice)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "HEAD, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
