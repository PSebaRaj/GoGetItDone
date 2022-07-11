package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name          string
	Email         string `gorm:"typevarchar(100);unique_index"`
	Projects      []Project
	Tasks         []Task
	ExpiringTasks []ExpiringTask
	PriorityTasks []PriorityTask
}
