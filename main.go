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

// server running on PORT 8080
func main() {

	DB = database.Connect()

	defer DB.Close()

	database.AutoMigrateAll()
	cache.ConnectRedisCache()

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", routes.LoadCors(router)))
}
