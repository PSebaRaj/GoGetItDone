package routes

import (
	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/controllers"
)

// used to init a new gmux
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	peopleRoutes(router)
	projectRoutes(router)
	taskRoutes(router)
	expiringTaskRoutes(router)
	priorityTaskRoutes(router)

	return router
}

// person/user routes
func peopleRoutes(router *mux.Router) {
	router.HandleFunc("/people", controllers.GetPeople).Methods("GET")
	router.HandleFunc("/person/{email}", controllers.GetPerson).Methods("GET")
	router.HandleFunc("/create/person", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/delete/person/{email}", controllers.DeletePerson).Methods("DELETE")
}

// project routes
func projectRoutes(router *mux.Router) {
	router.HandleFunc("/project/{id}", controllers.GetProject).Methods("GET")
	router.HandleFunc("/projects/{email}", controllers.GetProjectsByPerson).Methods("GET")
	router.HandleFunc("/create/project", controllers.CreateProject).Methods("POST")
	router.HandleFunc("/delete/project/{id}", controllers.DeleteProject).Methods("DELETE")
}

// (regular) task routes
func taskRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", controllers.GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", controllers.GetTask).Methods("GET")
	router.HandleFunc("/create/task", controllers.CreateTask).Methods("POST")
	router.HandleFunc("/delete/task/{id}", controllers.DeleteTask).Methods("DELETE")
	router.HandleFunc("/complete/task/{id}", controllers.ToggleCompleteTask).Methods("PATCH")

	router.HandleFunc("/changetitle/task/{id}", controllers.ChangeTitleTask).Methods("PATCH")
	router.HandleFunc("/changedescription/task/{id}", controllers.ChangeDescriptionTask).Methods("PATCH")
	router.HandleFunc("/changeproject/task/{id}", controllers.ChangeProjectTask).Methods("PATCH")
}

// expiring task routes
func expiringTaskRoutes(router *mux.Router) {
	router.HandleFunc("/expiringtasks", controllers.GetExpiringTasks).Methods("GET")
	router.HandleFunc("/expiringtask/{id}", controllers.GetExpiringTask).Methods("GET")
	router.HandleFunc("/create/expiringtask", controllers.CreateExpiringTask).Methods("POST")
	router.HandleFunc("/delete/expiringtask/{id}", controllers.DeleteExpiringTask).Methods("DELETE")
	router.HandleFunc("/complete/expiringtask/{id}", controllers.ToggleCompleteExpiringTask).Methods("PATCH")

	router.HandleFunc("/changetitle/expiringtask/{id}", controllers.ChangeTitleExpiringTask).Methods("PATCH")
	router.HandleFunc("/changedescription/expiringtask/{id}", controllers.ChangeDescriptionExpiringTask).Methods("PATCH")
	router.HandleFunc("/changeproject/expiringtask/{id}", controllers.ChangeProjectExpiringTask).Methods("PATCH")
}

// priority task routes
func priorityTaskRoutes(router *mux.Router) {
	router.HandleFunc("/prioritytasks", controllers.GetPriorityTasks).Methods("GET")
	router.HandleFunc("/prioritytask/{id}", controllers.GetPriorityTask).Methods("GET")
	router.HandleFunc("/create/prioritytask", controllers.CreatePriorityTask).Methods("POST")
	router.HandleFunc("/delete/prioritytask/{id}", controllers.DeletePriorityTask).Methods("DELETE")
	router.HandleFunc("/complete/prioritytask/{id}", controllers.ToggleCompletePriorityTask).Methods("PATCH")
	router.HandleFunc("/prioritytask/plevel/{id}/{plevel}", controllers.ChangePriorityLevel).Methods("PATCH")

	router.HandleFunc("/changetitle/prioritytask/{id}", controllers.ChangeTitlePriorityTask).Methods("PATCH")
	router.HandleFunc("/changedescription/prioritytask/{id}", controllers.ChangeDescriptionPriorityTask).Methods("PATCH")
	router.HandleFunc("/changeproject/prioritytask/{id}", controllers.ChangeProjectPriorityTask).Methods("PATCH")
}
