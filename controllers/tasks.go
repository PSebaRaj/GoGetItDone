package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/cache"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// controller: get singular (regular) task
// res: one task as JSON
func GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task models.Task

	// If the element is found in the redis cache, directly return it
	res := cache.GetFromCache(cache.REDIS, params["id"])
	if res != nil {
		fmt.Println("Cache hit")
		io.WriteString(w, res.(string))
		return
	}
	fmt.Println("Cache miss")
	database.DB.First(&task, params["id"])

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], task)
	json.NewEncoder(w).Encode(&task)
}

// controller: get all (regular) tasks
// res: list of tasks as JSON
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	database.DB.Find(&tasks)

	json.NewEncoder(w).Encode(&tasks)
}

// controller: create singular (regular) task
// res: created (regular) task as JSON
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	createdTask := database.DB.Create(&task)
	err := createdTask.Error
	if err != nil {
		fmt.Printf("Error creating task %s, error: %s", task.Title, err)
		return
	}

	json.NewEncoder(w).Encode(&createdTask)
}

// controller: delete singular (regular) task
// res: deleted task as JSON
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var task models.Task

	database.DB.First(&task, params["id"])
	deleted := database.DB.Delete(&task)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting task %s, error: %s", task.Title, err)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	json.NewEncoder(w).Encode(&task)
}

// controller: toggle complete boolean for a (regular) task
// res: updated (regular) task with toggled completion status as JSON
func ToggleCompleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	database.ToggleTaskComplete(database.TYPE_TASK, task.ID, task.Complete) // updates DB
	task.Complete = !task.Complete                                          // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	json.NewEncoder(w).Encode(&task)
}

// controller: changes title of a (regular) task
// res: updated (regular) task with new title as JSON
func ChangeTitleTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	database.ChangeTaskTitle(database.TYPE_TASK, task.ID, newTitle.Title) // updates DB
	task.Title = newTitle.Title                                           // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	json.NewEncoder(w).Encode(&task)
}

// controller: changes description of a (regular) task
// res: updated (regular) task with new description as JSON
func ChangeDescriptionTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	database.ChangeTaskDescription(database.TYPE_TASK, task.ID, newDescription.Description) // updates DB
	task.Description = newDescription.Description                                           // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	json.NewEncoder(w).Encode(&task)
}
