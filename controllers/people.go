package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// controller: get singular person with all of their tasks
// res: one user as JSON
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person
	var tasks []models.Task
	var expiringTasks []models.ExpiringTask

	database.DB.First(&person, "email = ?", params["email"])
	database.DB.Model(&person).Related(&tasks)
	database.DB.Model(&person).Related(&expiringTasks)

	// update expiry time in both JSON response and Postgres DB
	database.UpdateExpiringTask(expiringTasks)
	//utils.UpdateExpiringTaskTimeLeft(expiringTasks)

	person.Tasks = tasks
	person.ExpiringTasks = expiringTasks

	json.NewEncoder(w).Encode(&person)
}

// controller: get all people, without loading tasks
// to avoid massive json responses
// res: all users as JSON
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person

	database.DB.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

// controller: create singular person
// res: created user as JSON
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := database.DB.Create(&person)
	err := createdPerson.Error
	if err != nil {
		fmt.Printf("Error creating person %s, error: %s", person.Name, err)
	}

	json.NewEncoder(w).Encode(&createdPerson)
}

// controller: delete singular person
// res: deleted user as JSON
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person models.Person

	database.DB.First(&person, "email = ?", params["email"])
	deleted := database.DB.Delete(&person)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting person %s, error: %s", person.Name, err)
	}

	json.NewEncoder(w).Encode(&person)
}
