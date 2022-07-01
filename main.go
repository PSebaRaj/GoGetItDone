package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
	"github.com/psebaraj/gogetitdone/routes"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB
var err error

func main() {

	DB = database.Connect()

	// Close the databse connection when the main function closes
	defer DB.Close()

	// Make migrations to the database if they haven't been made already
	DB.AutoMigrate(&models.Person{})
	DB.AutoMigrate(&models.Task{})

	/*----------- API routes ------------*/
	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", routes.LoadCors(router)))
}
