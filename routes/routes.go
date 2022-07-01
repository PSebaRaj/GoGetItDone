package routes

import (
	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/controllers"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/tasks", controllers.GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", controllers.GetTask).Methods("GET")
	router.HandleFunc("/people", controllers.GetPeople).Methods("GET")
	router.HandleFunc("/person/{id}", controllers.GetPerson).Methods("GET")

	router.HandleFunc("/create/person", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/create/task", controllers.CreateTask).Methods("POST")

	router.HandleFunc("/delete/person/{id}", controllers.DeletePerson).Methods("DELETE")
	router.HandleFunc("/delete/task/{id}", controllers.DeleteTask).Methods("DELETE")

	return router
}