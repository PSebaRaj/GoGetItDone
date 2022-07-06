package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/cache"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// controller: get singular expiring task
// res: one task as JSON
func GetExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// If the element is found in the redis cache, directly return it
	res := cache.GetFromCache(cache.REDIS, params["id"])
	if res != nil {
		fmt.Println("Cache hit")
		io.WriteString(w, res.(string))
		return
	}
	fmt.Println("Cache miss")
	database.DB.First(&expiringTask, params["id"])

	// utils/clientLoading
	expiringTask.TimeLeft = time.Duration(expiringTask.ExpiringAt.Sub(time.Now()).Minutes())

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)
	json.NewEncoder(w).Encode(&expiringTask)
}

// controller: get all expiring tasks
// res: list of tasks as JSON
func GetExpiringTasks(w http.ResponseWriter, r *http.Request) {
	var expiringTasks []models.ExpiringTask

	database.DB.Find(&expiringTasks)

	database.UpdateExpiringTask(expiringTasks)
	//models.UpdateExpiringTaskTimeLeft(expiringTasks)

	json.NewEncoder(w).Encode(&expiringTasks)
}

// controller: create singular expiring task
// res: created expiring task as JSON
func CreateExpiringTask(w http.ResponseWriter, r *http.Request) {
	var expiringTask models.ExpiringTask
	json.NewDecoder(r.Body).Decode(&expiringTask)

	createdExpiringTask := database.DB.Create(&expiringTask)
	err := createdExpiringTask.Error
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(&createdExpiringTask)
}

// controller: delete singular expiring task
// res: deleted task as JSON
func DeleteExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var expiringTask models.ExpiringTask

	database.DB.First(&expiringTask, params["id"])
	deleted := database.DB.Delete(&expiringTask)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting expiringTask %s, error: %s", expiringTask.Title, err)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	json.NewEncoder(w).Encode(&expiringTask)
}

// controller: toggle complete boolean for an expiringTask
// res: updated expiringTask with toggled completion status as JSON
func ToggleCompleteExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	database.ToggleTaskComplete(database.TYPE_EXPIRINGTASK, expiringTask.ID, expiringTask.Complete) // updates DB
	expiringTask.Complete = !expiringTask.Complete                                                  // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	json.NewEncoder(w).Encode(&expiringTask)
}

// controller: changes title of a expiring task
// res: updated expiring task with new title as JSON
func ChangeTitleExpiringTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	database.ChangeTaskTitle(database.TYPE_EXPIRINGTASK, expiringTask.ID, newTitle.Title) // updates DB
	expiringTask.Title = newTitle.Title                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	json.NewEncoder(w).Encode(&expiringTask)
}

// controller: changes description of a expiring task
// res: updated expiring task with new description as JSON
func ChangeDescriptionExpiringTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	database.ChangeTaskDescription(database.TYPE_EXPIRINGTASK, expiringTask.ID, newDescription.Description) // updates DB
	expiringTask.Description = newDescription.Description                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	json.NewEncoder(w).Encode(&expiringTask)
}
