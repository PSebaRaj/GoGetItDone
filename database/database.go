package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/psebaraj/gogetitdone/models"
	"github.com/psebaraj/gogetitdone/utils"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB
var err error

func Connect() *gorm.DB {
	// Loading enviroment variables
	utils.LoadEnvVars()

	// Setting environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	//postgresURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)
	postgresURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, dbPort, user, dbname, dbpassword)

	// Openning connection to database
	DB, err = gorm.Open(dialect, postgresURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}

	return DB
}

func AutoMigrateAll() {

	err = DB.AutoMigrate(&models.Person{}, &models.Task{}).Error
	if err != nil {
		fmt.Printf("Unable to AutoMigrate model(s) %s, %s to Postgres DB", "Person", "Task")
		panic(err)
	}

}
