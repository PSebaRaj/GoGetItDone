package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

// swagger:route GET /project/{id} Project getProject
//
//
// Produces:
// - application/json
//
// responses:
//   200: Project
//   404: nil
//
// controller: get project with all of its tasks
// res: one project as JSON
func GetProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project models.Project
	var tasks []models.Task
	var expiringTasks []models.ExpiringTask
	var priorityTasks []models.PriorityTask

	database.DB.First(&project, "id = ?", params["id"])
	if project.Title == "" { // i.e. project not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&project) // still want to send as json
		return
	}

	database.DB.Model(&project).Related(&tasks)
	database.DB.Model(&project).Related(&expiringTasks)
	database.DB.Model(&project).Related(&priorityTasks)

	database.UpdateExpiringTask(expiringTasks)

	project.Tasks = tasks
	project.ExpiringTasks = expiringTasks
	project.PriorityTasks = priorityTasks

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&project)
}

// swagger:route GET /projects/{email} Project getProjects
//
//
// Produces:
// - application/json
//
// responses:
//   200: []Project
//   204: nil
//   404: nil
//
// controller: get project with all of its tasks
// res: one project as JSON
func GetProjectsByPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person
	var projects []models.Project

	database.DB.First(&person, "email = ?", params["email"])
	if person.Email == "" { // i.e. person not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&person) // still want to send as json
		return
	}

	database.DB.Model(&person).Related(&projects)

	if len(projects) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&projects) // still want to send as json
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&projects)
}

// swagger:route POST /create/project Project createProject
//
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// responses:
//   201: Project
//   507: nil
//
// controller: create singular project to help user soft tasks
// res: created (empty) project as JSON
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	json.NewDecoder(r.Body).Decode(&project)

	createdProject := database.DB.Create(&project)
	err := createdProject.Error
	if err != nil {
		fmt.Printf("Error creating project %s, error: %s", project.Title, err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdProject)
}

// swagger:route DELETE /delete/project/{id} Project deleteProject
//
//
// Produces:
// - application/json
//
// responses:
//   200: Project
//   404: nil
//   500: nil
//
// controller: delete project (but not the project's tasks)
// res: deleted project as JSON
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project models.Project

	database.DB.First(&project, "id = ?", params["id"])
	if project.Title == "" {
		fmt.Printf("Error finding project %s before deletion", project.Title)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deleted := database.DB.Delete(&project)

	err := deleted.Error
	if err != nil {
		fmt.Printf("Error deleting project %s, error: %s", project.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&project)
}
