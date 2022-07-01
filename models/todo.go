package models

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

var DB *gorm.DB
var REDIS *redis.Client

type Todo struct {
	gorm.Model
	TaskID          string `json:"taskId"`
	TaskDescription string `json:"taskDescription"`
}

func getDatabaseURI() string {

	var dbUser = os.Getenv("POSTGRES_USER")
	var dbPassword = os.Getenv("POSTGRES_USER")
	var db = os.Getenv("POSTGRES_USER")
	var dbHost = os.Getenv("POSTGRES_USER")
	var dbPort = os.Getenv("POSTGRES_USER")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, db, dbPassword)
}
