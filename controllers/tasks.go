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
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	json.NewEncoder(w).Encode(&task)
}
