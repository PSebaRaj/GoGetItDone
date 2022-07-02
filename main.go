package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/psebaraj/gogetitdone/cache"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/routes"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func main() {

	DB = database.Connect()

	// Close the database connection when the main function closes
	defer DB.Close()

	// Make migrations to the database if they haven't been made already
	database.AutoMigrateAll()
	cache.ConnectRedisCache()

	/*----------- API routes ------------*/
	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", routes.LoadCors(router)))
}
