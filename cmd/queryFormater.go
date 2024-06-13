package server

import (
	"strings"

	"gorm.io/gorm"
)

// method used in FilterBooks, meant to use with WHERE statement
func ApplyFilter(field string, value string, db *gorm.DB) *gorm.DB {
	if value != "" {
		db = db.Where(field+" IN (?)", strings.Split(value, ","))
	}
	return db
}

// method used in GetBookByStructuredSearch, meant to add another search condition to the existing one
func AppendToSearchIfNotEmpty(field string, filter []string) []string {
	if field != "" {
		filter = append(filter, "search @@ websearch_to_tsquery('simple', '"+field+"')")
	}
	return filter
}
