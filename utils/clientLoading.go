package utils

import (
	"time"

	"github.com/psebaraj/gogetitdone/models"
)

// used to update the expiry time parameter before the response is encoded
// into JSON and sent
func UpdateExpiringTaskTimeLeft(expiringTasks []models.ExpiringTask) {

	for i := 0; i < len(expiringTasks); i++ {
		expiringTasks[i].TimeLeft = time.Duration(expiringTasks[i].ExpiringAt.Sub(time.Now()).Minutes())
	}
	// add error handling
}
