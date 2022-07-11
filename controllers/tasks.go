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

// swagger:route GET /task/{id} Task getTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: Task
//   404: nil
//
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

	if task.Title == "" { // i.e. task not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&task) // still want to send as json
		return
	}

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

// swagger:route GET /tasks Task getTasks
//
//
// Produces:
// - application/json
//
// responses:
//   200: []Task
//   204: nil
//
// controller: get all (regular) tasks
// res: list of tasks as JSON
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	database.DB.Find(&tasks)

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&tasks) // still want to send as json
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&tasks)
}

// swagger:route POST /create/task Task createTask
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   201: Task
//   507: nil
//
// controller: create singular (regular) task
// res: created (regular) task as JSON
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	createdTask := database.DB.Create(&task)
	err := createdTask.Error
	if err != nil {
		fmt.Printf("Error creating task %s, error: %s", task.Title, err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdTask)
}

// swagger:route DELETE /delete/task/{id} Task deleteTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: Task
//   404: nil
//   500: nil
//
// controller: delete singular (regular) task
// res: deleted task as JSON
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var task models.Task

	database.DB.First(&task, params["id"])
	if task.Title == "" {
		fmt.Printf("Error finding task %s before deletion", task.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	deleted := database.DB.Delete(&task)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting task %s, error: %s", task.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

// swagger:route PATCH /complete/task/{id} Task toggleTaskComplete
//
//
// Produces:
// - application/json
//
// responses:
//   200: Task
//   400: nil
//   404: nil
//
// controller: toggle complete boolean for a (regular) task
// res: updated (regular) task with toggled completion status as JSON
func ToggleCompleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	if task.Title == "" {
		fmt.Printf("Error finding task %s before toggle", task.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ToggleTaskComplete(database.TYPE_TASK, task.ID, task.Complete) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.Complete = !task.Complete // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

// swagger:route PATCH /changetitle/task/{id} Task changeTaskTitle
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: Task
//   400: nil
//   404: nil
//
// controller: changes title of a (regular) task
// res: updated (regular) task with new title as JSON
func ChangeTitleTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	if task.Title == "" {
		fmt.Printf("Error finding task %s before title change", task.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ChangeTaskTitle(database.TYPE_TASK, task.ID, newTitle.Title) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.Title = newTitle.Title // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

// swagger:route PATCH /changedescrition/task/{id} Task changeTaskDescription
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: Task
//   400: nil
//   404: nil
//
// controller: changes description of a (regular) task
// res: updated (regular) task with new description as JSON
func ChangeDescriptionTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	if task.Title == "" {
		fmt.Printf("Error finding task %s before description change", task.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ChangeTaskDescription(database.TYPE_TASK, task.ID, newDescription.Description) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.Description = newDescription.Description // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

// swagger:route PATCH /changeproject/task/{id} Task changeTaskProject
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: Task
//   400: nil
//   404: nil
//
// controller: changes project of a (regular) task
// res: updated (regular) task with new project assignment as JSON
func ChangeProjectTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var task models.Task

	// Ignore cache, go straight to DB
	database.DB.First(&task, params["id"])
	if task.Title == "" {
		fmt.Printf("Error finding task %s before project assignment change", task.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ChangeTaskProject(database.TYPE_TASK, task.ID, newDescription.ProjectID) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.ProjectID = newDescription.ProjectID // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}
