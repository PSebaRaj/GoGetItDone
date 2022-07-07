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

// swagger:route GET /prioritytask/{id} PriorityTask getPriorityTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: PriorityTask
//   404: nil
//
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

	if priorityTask.Title == "" { // i.e. task not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&priorityTask) // still want to send as json
		return
	}
	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)
}

// swagger:route GET /prioritytasks PriorityTask getPriorityTasks
//
//
// Produces:
// - application/json
//
// responses:
//   200: []PriorityTask
//   204: nil
//
// controller: get all priorityTasks
// res: list of priorityTasks as JSON
func GetPriorityTasks(w http.ResponseWriter, r *http.Request) {
	var priorityTasks []models.PriorityTask

	database.DB.Find(&priorityTasks)
	if len(priorityTasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&priorityTasks) // still want to send as json
		return
	}

	json.NewEncoder(w).Encode(&priorityTasks)
}

// swagger:route POST /create/prioritytask PriorityTask createPriorityTask
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   201: PriorityTask
//   400: nil
//   507: nil
//
// controller: create singular priorityTask
// res: created priorityTask as JSON
func CreatePriorityTask(w http.ResponseWriter, r *http.Request) {
	var priorityTask models.PriorityTask
	json.NewDecoder(r.Body).Decode(&priorityTask)

	if priorityTask.PriorityLevel == 0 {
		fmt.Println("Priority level of priorityTask is undefined")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !models.IsValidPriorityLevel(priorityTask.PriorityLevel) {
		fmt.Println("Priority level of priorityTask does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTask := database.DB.Create(&priorityTask)
	err := createdTask.Error
	if err != nil {
		fmt.Printf("Error creating priorityTask %s, error: %s", priorityTask.Title, err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdTask)
}

// swagger:route DELETE /delete/prioritytask/{id} PriorityTask deletePriorityTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: PriorityTask
//   404: nil
//   500: nil
//
// controller: delete singular priorityTask
// res: deleted priorityTask as JSON
func DeletePriorityTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var priorityTask models.PriorityTask

	database.DB.First(&priorityTask, params["id"])
	if priorityTask.Title == "" {
		fmt.Printf("Error finding priorityTask %s before deletion", priorityTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deleted := database.DB.Delete(&priorityTask)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting priorityTask %s, error: %s", priorityTask.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)
}

// swagger:route PATCH /complete/prioritytask/{id} PriorityTask togglePriorityTaskComplete
//
//
// Produces:
// - application/json
//
// responses:
//   200: PriorityTask
//   404: nil
//
// controller: toggle complete boolean for a priorityTask
// res: updated priorityTask with toggled completion status as JSON
func ToggleCompletePriorityTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	if priorityTask.Title == "" {
		fmt.Printf("Error finding priorityTask %s before toggle", priorityTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	database.ToggleTaskComplete(database.TYPE_PRIORITYTASK, priorityTask.ID, priorityTask.Complete) // updates DB
	priorityTask.Complete = !priorityTask.Complete                                                  // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)
}

// swagger:route PATCH /prioritytask/plevel/{id}/{new_level} PriorityTask changePriorityLevel
//
//
// Produces:
// - application/json
//
// responses:
//   200: PriorityTask
//   400: nil
//   404: nil
//
// controller: changes priority level for a priorityTask
// res: updated priorityTask with new priority level as JSON
func ChangePriorityLevel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	plevel, err := strconv.Atoi(params["plevel"])
	if err != nil {
		fmt.Printf("Invalid priority level")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	if priorityTask.Title == "" {
		fmt.Printf("Error finding priorityTask %s before priority level change", priorityTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	database.ChangePriority(priorityTask, models.PriorityLevelType(plevel))
	priorityTask.PriorityLevel = models.PriorityLevelType(plevel)

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)

}

// swagger:route PATCH /changetitle/prioritytask/{id} PriorityTask changePriorityTaskTitle
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: PriorityTask
//   404: nil
//
// controller: changes title of a priority task
// res: updated priority task with new title as JSON
func ChangeTitlePriorityTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	if priorityTask.Title == "" {
		fmt.Printf("Error finding priorityTask %s before title change", priorityTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	database.ChangeTaskTitle(database.TYPE_PRIORITYTASK, priorityTask.ID, newTitle.Title) // updates DB
	priorityTask.Title = newTitle.Title                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)
}

// swagger:route PATCH /changedescrition/prioritytask/{id} PriorityTask changePriorityTaskDescription
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: PriorityTask
//   404: nil
//
// controller: changes description of a priority task
// res: updated priority task with new description as JSON
func ChangeDescriptionPriorityTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var priorityTask models.PriorityTask

	// Ignore cache, go straight to DB
	database.DB.First(&priorityTask, params["id"])
	if priorityTask.Title == "" {
		fmt.Printf("Error finding priorityTask %s before description change", priorityTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	database.ChangeTaskDescription(database.TYPE_PRIORITYTASK, priorityTask.ID, newDescription.Description) // updates DB
	priorityTask.Description = newDescription.Description                                                   // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], priorityTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&priorityTask)
}
