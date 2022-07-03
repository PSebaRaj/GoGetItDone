package utils

import (
	"time"

	"github.com/psebaraj/gogetitdone/models"
)

func UpdateExpiringTaskTimeLeft(expiringTasks []models.ExpiringTask) {

	for i := 0; i < len(expiringTasks); i++ {
		expiringTasks[i].TimeLeft = time.Duration(expiringTasks[i].ExpiringAt.Sub(time.Now()).Minutes())
	}
	// add error handling
}
