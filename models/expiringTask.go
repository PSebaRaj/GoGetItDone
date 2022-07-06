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

// used to update the expiry time parameter before the response is encoded
// into JSON and sent
func UpdateExpiringTaskTimeLeft(expiringTasks []ExpiringTask) {

	for i := 0; i < len(expiringTasks); i++ {
		expiringTasks[i].TimeLeft = time.Duration(expiringTasks[i].ExpiringAt.Sub(time.Now()).Minutes())
	}
	// add error handling
}
