package server

import (
	"fmt"
	"strings"

	"github.com/DanillaY/GoScrapper/cmd/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Db     *gorm.DB
	Config *repository.Config
}

type Pagination struct {
	Total       int
	PerPage     int
	CurrentPage int
	LastPage    int
}

func NewPostgresConnection(c *repository.Config) (db *gorm.DB, e error) {
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

func (r *Repository) PrepareDatabase() (e error) {
	updateExistingBooks := "update books set search = setweight(to_tsvector('simple',title), 'A') || ' ' || setweight(to_tsvector('simple',author), 'B') || ' ' || setweight(to_tsvector('simple',category), 'C'):: tsvector;"

	createFunc := "CREATE OR REPLACE FUNCTION books_trigger() RETURNS trigger AS $$ begin new.search := setweight(to_tsvector('simple',coalesce(new.title,'')), 'A') || ' ' || setweight(to_tsvector('simple',coalesce(new.author,'')), 'B') || ' ' || setweight(to_tsvector('simple',coalesce(new.category,'')), 'C'):: tsvector; return new; end $$ LANGUAGE plpgsql;"
	createTrigger := "CREATE OR REPLACE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON books FOR EACH ROW EXECUTE FUNCTION books_trigger(); "
	createBookTitleIndex := "CREATE INDEX IF NOT EXISTS books_title ON books USING GIN(to_tsvector('simple', title)); "
	createBookAuthorIndex := "CREATE INDEX IF NOT EXISTS books_author ON books USING GIN(to_tsvector('simple', author)); "
	createBookCategoryIndex := "CREATE INDEX IF NOT EXISTS books_category ON books USING GIN(to_tsvector('simple', category)); "
	createBookVendorIndex := "CREATE INDEX IF NOT EXISTS books_vendor ON books USING GIN(to_tsvector('simple', vendor)); "
	createBookStockIndex := "CREATE INDEX IF NOT EXISTS books_stock ON books USING GIN(to_tsvector('simple', in_stock_text)); "

	err := r.Db.Exec(createBookTitleIndex + createBookAuthorIndex +
		createBookCategoryIndex + createBookVendorIndex +
		createBookStockIndex + updateExistingBooks +
		createFunc + createTrigger).Error
	return err
}

func FilterBooks(
	maxPrice int,
	minPrice int,
	category string, search string,
	author string, vendor string,
	yearPublished int, stockText string,
	bookInfoType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		db = db.Where("current_price >= ?", minPrice).Where("current_price <= ?", maxPrice)
		selectTables := []string{"*"}

		if bookInfoType == "partial" {
			selectTables = []string{"id", "title", "current_price", "old_price", "author", "category", "img_path", "page_book_path", "vendor", "age_restriction"}
			db.Select(selectTables)
		}

		if search != "" {
			ts := "ts_rank(search, websearch_to_tsquery('simple', '" + search + "' )) + ts_rank(search, websearch_to_tsquery('russian', '" + search + "' )) as rank"
			selectTables = append(selectTables, ts)
			db = db.Table("books").Select(selectTables).
				Where("search @@ websearch_to_tsquery('simple', ?) or search @@ websearch_to_tsquery('simple', ?) or search @@ websearch_to_tsquery('simple', ?)", search, category, author).
				Order("rank DESC")
		}

		db = applyFilter("category", category, db)
		db = applyFilter("vendor", vendor, db)
		db = applyFilter("author", author, db)
		db = applyFilter("in_stock_text", stockText, db)

		if yearPublished != 0 {
			db.Where("year_publish = ?", yearPublished)
		}

		return db
	}
}

func applyFilter(field string, value string, db *gorm.DB) *gorm.DB {
	if value != "" {
		db = db.Where(field+" IN (?)", strings.Split(value, ","))
	}
	return db
}

func AppendToSearchIfNotEmpty(field string, filter []string) []string {
	if field != "" {
		filter = append(filter, "search @@ websearch_to_tsquery('simple', '"+field+"')")
	}
	return filter
}
