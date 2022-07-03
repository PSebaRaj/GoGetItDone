package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ExpiringTask struct {
	gorm.Model

	Title       string
	Description string
	Complete    bool
	PersonID    int

	ExpiringAt time.Time
	TimeLeft   time.Duration
}
