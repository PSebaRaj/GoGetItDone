package database

import (
	"fmt"
	"os"
	"time"

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

	err = DB.AutoMigrate(&models.Person{}, &models.Task{}, &models.ExpiringTask{}).Error
	if err != nil {
		fmt.Printf("Unable to AutoMigrate model(s) %s, %s. %s to Postgres DB", "Person", "Task", "ExpiringTask")
		panic(err)
	}

}

// update both the time left before sending to client
// and update on server for accuracy
// server update run on seperate thread via go routines
func UpdateExpiringTask(expiringTasks []models.ExpiringTask) {
	go updateExpiringTasksInDB(expiringTasks)
	utils.UpdateExpiringTaskTimeLeft(expiringTasks)
}

// note, looping over twice and not a func to update each task
// b/c spawning another process using a go routine for each entry
// doesn't make sense, too many processes/threads required
func updateExpiringTasksInDB(expiringTasks []models.ExpiringTask) {
	for i := 0; i < len(expiringTasks); i++ {
		sqlStatement := `
		UPDATE expiring_tasks
		SET time_left = $2
		WHERE id = $1;`

		err = DB.Exec(sqlStatement, expiringTasks[i].ID, time.Duration(expiringTasks[i].ExpiringAt.Sub(time.Now()).Minutes())).Error
		if err != nil {
			fmt.Printf("Error updating Expiring Task: %s", expiringTasks[i].Title)
			panic(err)
		}
	}
}
