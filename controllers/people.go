package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// swagger:response peopleResponse
// A list of all people
type peopleResponse struct {
	// All people in db
	// in: body
	Body []models.Person
}

// swagger:response personResponse
// A singular person
type personResponse struct {
	// One person in db
	// in: body
	Body models.Person
}

// swagger:route GET /person/{email} Person getPerson
//
//
// Produces:
// - application/json
//
// responses:
//   200: personResponse
//   404: nil
//
// controller: get singular person with all of their tasks
// res: one user as JSON
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person
	var tasks []models.Task
	var expiringTasks []models.ExpiringTask
	var priorityTasks []models.PriorityTask

	database.DB.First(&person, "email = ?", params["email"])
	if person.Email == "" { // i.e. person not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&person) // still want to send as json
		return
	}

	database.DB.Model(&person).Related(&tasks)
	database.DB.Model(&person).Related(&expiringTasks)
	database.DB.Model(&person).Related(&priorityTasks)

	// update expiry time in both JSON response and Postgres DB
	database.UpdateExpiringTask(expiringTasks)
	//models.UpdateExpiringTaskTimeLeft(expiringTasks)

	person.Tasks = tasks
	person.ExpiringTasks = expiringTasks
	person.PriorityTasks = priorityTasks

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&person)
}

// swagger:route GET /people Person getPeople
//
//
// Produces:
// - application/json
//
// responses:
//   200: peopleResponse
//   204: nil
//
// controller: get all people, without loading tasks
// to avoid massive json responses
// res: all users as JSON
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person

	database.DB.Find(&people)

	if len(people) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&people) // still want to send as json
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&people)
}

// swagger:route POST /create/person Person createPerson
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   201: personResponse
//   507: nil
//
// controller: create singular person
// res: created user as JSON
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := database.DB.Create(&person)
	err := createdPerson.Error
	if err != nil {
		fmt.Printf("Error creating person %s, error: %s", person.Name, err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdPerson)
}

// swagger:route DELETE /delete/person/{email} Person deletePerson
//
//
// Produces:
// - application/json
//
// responses:
//   200: personResponse
//   404: nil
//   500: nil
//
// controller: delete singular person
// res: deleted user as JSON
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person models.Person

	database.DB.First(&person, "email = ?", params["email"])
	if person.Email == "" {
		fmt.Printf("Error finding person %s before deletion", person.Email)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deleted := database.DB.Delete(&person)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting person %s, error: %s", person.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&person)
}
