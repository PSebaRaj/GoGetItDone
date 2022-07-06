package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/cache"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// controller: get singular priorityTask
// res: one priorityTask as JSON
func GetPriorityTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// If the element is found in the redis cache, directly return it
	res := cache.GetFromCache(cache.REDIS, params["id"])
	if res != nil {
		fmt.Println("Cache hit")
		io.WriteString(w, res.(string))
		return
	}
	fmt.Println("Cache miss")
	database.DB.First(&priorityTask, params["id"])

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)
	json.NewEncoder(w).Encode(&priorityTask)
}

// controller: get all priorityTasks
// res: list of priorityTasks as JSON
func GetPriorityTasks(w http.ResponseWriter, r *http.Request) {
	var priorityTasks []models.PriorityTask

	database.DB.Find(&priorityTasks)

	json.NewEncoder(w).Encode(&priorityTasks)
}

// controller: create singular priorityTask
// res: created priorityTask as JSON
func CreatePriorityTask(w http.ResponseWriter, r *http.Request) {
	var priorityTask models.PriorityTask
	json.NewDecoder(r.Body).Decode(&priorityTask)

	if priorityTask.PriorityLevel == 0 {
		fmt.Println("Priority level of priorityTask is undefined")
		return
	}

	if !models.IsValidPriorityLevel(priorityTask.PriorityLevel) {
		fmt.Println("Priority level of priorityTask does not exist")
		return
	}

	createdTask := database.DB.Create(&priorityTask)
	err := createdTask.Error
	if err != nil {
		fmt.Printf("Error creating priorityTask %s, error: %s", priorityTask.Title, err)
		return
	}

	json.NewEncoder(w).Encode(&createdTask)
}

// controller: delete singular priorityTask
// res: deleted priorityTask as JSON
func DeletePriorityTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var priorityTask models.PriorityTask

	database.DB.First(&priorityTask, params["id"])
	deleted := database.DB.Delete(&priorityTask)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting priorityTask %s, error: %s", priorityTask.Title, err)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	json.NewEncoder(w).Encode(&priorityTask)
}

// controller: toggle complete boolean for a priorityTask
// res: updated priorityTask with toggled completion status as JSON
func ToggleCompletePriorityTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	database.ToggleTaskComplete(database.TYPE_PRIORITYTASK, priorityTask.ID, priorityTask.Complete) // updates DB
	priorityTask.Complete = !priorityTask.Complete                                                  // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	json.NewEncoder(w).Encode(&priorityTask)
}

// controller: changes priority level for a priorityTask
// res: updated priorityTask with new priority level as JSON
func ChangePriorityLevel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	plevel, err := strconv.Atoi(params["plevel"])
	if err != nil {
		fmt.Printf("Invalid priority level")
		return
	}

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	database.ChangePriority(priorityTask, models.PriorityLevelType(plevel))
	priorityTask.PriorityLevel = models.PriorityLevelType(plevel)

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	json.NewEncoder(w).Encode(&priorityTask)

}

// controller: changes title of a priority task
// res: updated priority task with new title as JSON
func ChangeTitlePriorityTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	database.ChangeTaskTitle(database.TYPE_PRIORITYTASK, priorityTask.ID, newTitle.Title) // updates DB
	priorityTask.Title = newTitle.Title                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	json.NewEncoder(w).Encode(&priorityTask)
}

// controller: changes description of a priority task
// res: updated priority task with new description as JSON
func ChangeDescriptionPriorityTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	database.ChangeTaskDescription(database.TYPE_PRIORITYTASK, priorityTask.ID, newDescription.Description) // updates DB
	priorityTask.Description = newDescription.Description                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	json.NewEncoder(w).Encode(&priorityTask)
}
