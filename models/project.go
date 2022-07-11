package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model

	Title    string
	PersonID int

	Tasks         []Task
	ExpiringTasks []ExpiringTask
	PriorityTasks []PriorityTask
}
