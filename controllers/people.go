package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person
	var tasks []models.Task

	database.DB.First(&person, "email = ?", params["email"])
	database.DB.Model(&person).Related(&tasks)

	person.Tasks = tasks

	json.NewEncoder(w).Encode(&person)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person

	database.DB.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := database.DB.Create(&person)
	err := createdPerson.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdPerson)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person models.Person

	database.DB.First(&person, "email = ?", params["email"])
	database.DB.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}
