package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Task struct {
	gorm.Model

	Title       string
	Description string
	Complete    bool
	PersonID    int
}
