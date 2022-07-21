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

var DB *gorm.DB // global connection to DB
var err error

// connect to DB using GORM
func Connect() *gorm.DB {
	// Loading environment variables
	utils.LoadEnvVars()

	// Setting environment variables
	dialect := os.Getenv("DB_DIALECT")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	dbpassword := os.Getenv("DB_PASSWORD")

	//postgresURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)
	postgresURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, dbPort, user, dbname, dbpassword)

	// Openning connection to database
	DB, err = gorm.Open(dialect, postgresURI)

	/*** PREPARE STMT
	DB, err = gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		PrepareStmt: true,
	})
	***/


	if err != nil {
		fmt.Println("Could not connect to the database")
		panic(err)
	}

	fmt.Println("Connected to database successfully")

	return DB
}

// Make migrations to the database if they haven't been made already
func AutoMigrateAll() {

	err = DB.AutoMigrate(&models.Person{}, &models.Project{}, &models.Task{}, &models.ExpiringTask{}, &models.PriorityTask{}).Error
	if err != nil {
		fmt.Printf("Unable to AutoMigrate model(s) %s, %s, %s, %s, %s to Postgres DB", "Person", "Project", "Task", "ExpiringTask", "PriorityTask")
		panic(err)
	}

}

// update both the time left before sending to client
// and update on server for accuracy
// server update run on separate thread via go routines
func UpdateExpiringTask(expiringTasks []models.ExpiringTask) {
	go updateExpiringTasksInDB(expiringTasks)
	models.UpdateExpiringTaskTimeLeft(expiringTasks)
}

// note, looping over twice and not a func to update each task
// b/c spawning another process using a go routine for each entry
// doesn't make sense, too many processes/threads required
func updateExpiringTasksInDB(expiringTasks []models.ExpiringTask) {
	for i := 0; i < len(expiringTasks); i++ {
		psqlStatement := `
		UPDATE expiring_tasks
		SET time_left = $2
		WHERE id = $1;`

		err = DB.Exec(psqlStatement, expiringTasks[i].ID, time.Duration(expiringTasks[i].ExpiringAt.Sub(time.Now()).Minutes())).Error
		if err != nil {
			fmt.Printf("Error updating Expiring Task: %s in DB", expiringTasks[i].Title)
			panic(err) // panic here to stop this Go routine
			// lazy error handling, I know
		}
	}
}
