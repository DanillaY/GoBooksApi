package server

import (
	"fmt"

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
