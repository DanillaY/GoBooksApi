package server

import (
	"fmt"
	"strings"

	"github.com/DanillaY/GoScrapper/cmd/models"
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

func FilterBooks(
	maxPrice string,
	minPrice string,
	category string,
	search string,
	author string,
	vendor string,
	yearPublished string,
	stockText string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("current_price >= ?", minPrice).
			Where("current_price <= ?", maxPrice).
			Where("LOWER(category) SIMILAR TO ?", category).
			Where("LOWER(vendor) SIMILAR TO ?", "%"+vendor+"%").
			Where("LOWER(author) SIMILAR TO ?", author).
			Where("LOWER(title) SIMILAR TO ? OR LOWER(author) SIMILAR TO ? OR LOWER(category) SIMILAR TO ?", search, search, search).
			Where("LOWER(in_stock_text) SIMILAR TO ?", "%"+stockText+"%").
			Where("year_publish >= ?", yearPublished)
	}
}

func AddRegexToQuery(query string, separator string, entryCheck bool) string {
	result := ""
	if entryCheck {
		queryWithReg := strings.Split(query, separator)
		for i := 0; i < len(queryWithReg); i++ {
			if i == len(queryWithReg)-1 {
				result += "%" + queryWithReg[i] + "%"
			} else {
				result += "%" + queryWithReg[i] + "%|"
			}
		}

	} else {
		result = "%" + query + "%"
	}
	return result
}
