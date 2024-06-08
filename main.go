package main

import (
	"fmt"

	server "github.com/DanillaY/BookApi/cmd"
	"github.com/DanillaY/GoScrapper/cmd/repository"
)

func main() {

	config, err := repository.GetConfigVariables()
	if err != nil {
		fmt.Println("Error while getting env data")
	}

	db, err := server.NewPostgresConnection(config)
	if err != nil {
		fmt.Println("Error while connecting to database")
	}

	repo := server.Repository{Db: db, Config: config}

	err = repo.PrepareDatabase()
	if err != nil {
		fmt.Println("Error while preparing the database: " + err.Error())
	}
	repo.InitAPIRoutes()
}
