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
	router.HandleFunc("/person/{email}", controllers.GetPerson).Methods("GET")

	router.HandleFunc("/create/person", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/create/task", controllers.CreateTask).Methods("POST")

	router.HandleFunc("/delete/person/{email}", controllers.DeletePerson).Methods("DELETE")
	router.HandleFunc("/delete/task/{id}", controllers.DeleteTask).Methods("DELETE")

	// expiring tasks
	router.HandleFunc("/expiringtasks", controllers.GetExpiringTasks).Methods("GET")
	router.HandleFunc("/expiringtask/{id}", controllers.GetExpiringTask).Methods("GET")
	router.HandleFunc("/create/expiringtask", controllers.CreateExpiringTask).Methods("POST")
	router.HandleFunc("/delete/expiringtask/{id}", controllers.DeleteExpiringTask).Methods("DELETE")

	return router
}
