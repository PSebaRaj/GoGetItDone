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

// swagger:route GET /expiringtask/{id} ExpiringTask getExpiringTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: ExpiringTask
//   404: nil
//
// controller: get singular expiring task
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

	if expiringTask.Title == "" { // i.e. task not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&expiringTask) // still want to send as json
		return
	}

	// utils/clientLoading
	expiringTask.TimeLeft = time.Duration(expiringTask.ExpiringAt.Sub(time.Now()).Minutes())

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTask)
}

// swagger:route GET /expiringtasks ExpiringTask getExpiringTasks
//
//
// Produces:
// - application/json
//
// responses:
//   200: []ExpiringTask
//   204: nil
//
// controller: get all expiring tasks
// res: list of tasks as JSON
func GetExpiringTasks(w http.ResponseWriter, r *http.Request) {
	var expiringTasks []models.ExpiringTask

	database.DB.Find(&expiringTasks)

	if len(expiringTasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&expiringTasks) // still want to send as json
		return
	}

	database.UpdateExpiringTask(expiringTasks)
	//models.UpdateExpiringTaskTimeLeft(expiringTasks)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTasks)
}

// swagger:route POST /create/expiringtask ExpiringTask createExpiringTask
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   201: ExpiringTask
//   507: nil
//
// controller: create singular expiring task
// res: created expiring task as JSON
func CreateExpiringTask(w http.ResponseWriter, r *http.Request) {
	var expiringTask models.ExpiringTask
	json.NewDecoder(r.Body).Decode(&expiringTask)

	createdExpiringTask := database.DB.Create(&expiringTask)
	err := createdExpiringTask.Error
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdExpiringTask)
}

// swagger:route DELETE /delete/expiringtask/{id} ExpiringTask deleteExpiringTask
//
//
// Produces:
// - application/json
//
// responses:
//   200: ExpiringTask
//   404: nil
//   500: nil
//
// controller: delete singular expiring task
// res: deleted task as JSON
func DeleteExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var expiringTask models.ExpiringTask

	database.DB.First(&expiringTask, params["id"])
	if expiringTask.Title == "" {
		fmt.Printf("Error finding expiringTask %s before deletion", expiringTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deleted := database.DB.Delete(&expiringTask)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting expiringTask %s, error: %s", expiringTask.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// also delete from cache
	cache.DeleteFromCache(cache.REDIS, params["id"])

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTask)
}

// swagger:route PATCH /complete/expiringtask/{id} ExpiringTask toggleExpiringTaskComplete
//
//
// Produces:
// - application/json
//
// responses:
//   200: ExpiringTask
//   400: nil
//   404: nil
//
// controller: toggle complete boolean for an expiringTask
// res: updated expiringTask with toggled completion status as JSON
func ToggleCompleteExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	if expiringTask.Title == "" {
		fmt.Printf("Error finding expiringTask %s before toggle", expiringTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ToggleTaskComplete(database.TYPE_EXPIRINGTASK, expiringTask.ID, expiringTask.Complete) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expiringTask.Complete = !expiringTask.Complete // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTask)
}

// swagger:route PATCH /changetitle/expiringtask/{id} ExpiringTask changeExpiringTaskTitle
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: ExpiringTask
//   400: nil
//   404: nil
//
// controller: changes title of a expiring task
// res: updated expiring task with new title as JSON
func ChangeTitleExpiringTask(w http.ResponseWriter, r *http.Request) {
	var newTitle models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newTitle)

	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	if expiringTask.Title == "" {
		fmt.Printf("Error finding expiringTask %s before title change", expiringTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ChangeTaskTitle(database.TYPE_EXPIRINGTASK, expiringTask.ID, newTitle.Title) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expiringTask.Title = newTitle.Title // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTask)
}

// swagger:route PATCH /changedescrition/expiringtask/{id} ExpiringTask changeExpiringTaskDescription
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   200: ExpiringTask
//   400: nil
//   404: nil
//
// controller: changes description of a expiring task
// res: updated expiring task with new description as JSON
func ChangeDescriptionExpiringTask(w http.ResponseWriter, r *http.Request) {
	var newDescription models.TaskModifier
	json.NewDecoder(r.Body).Decode(&newDescription)

	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// Ignore cache, go straight to DB
	database.DB.First(&expiringTask, params["id"])
	if expiringTask.Title == "" {
		fmt.Printf("Error finding expiringTask %s before description change", expiringTask.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := database.ChangeTaskDescription(database.TYPE_EXPIRINGTASK, expiringTask.ID, newDescription.Description) // updates DB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expiringTask.Description = newDescription.Description // updates response

	// Update element in the redis cache before returning the result
	cache.SetInCache(cache.REDIS, params["id"], expiringTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&expiringTask)
}
