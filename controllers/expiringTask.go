package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/psebaraj/gogetitdone/database"
	"github.com/psebaraj/gogetitdone/models"
)

func GetExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var expiringTask models.ExpiringTask

	// If the element is found in the redis cache, directly return it
	//res := cache.GetFromCache(cache.REDIS, params["id"])
	//if res != nil {
	//fmt.Println("Cache hit")
	//io.WriteString(w, res.(string))
	//return
	//}
	//fmt.Println("Cache miss")
	database.DB.First(&expiringTask, params["id"])

	// utils/clientLoading
	expiringTask.TimeLeft = time.Duration(expiringTask.ExpiringAt.Sub(time.Now()).Minutes())

	// Set element in the redis cache before returning the result
	// "id" is what I query with
	//cache.SetInCache(cache.REDIS, params["id"], expiringTask)
	json.NewEncoder(w).Encode(&expiringTask)
}

func GetExpiringTasks(w http.ResponseWriter, r *http.Request) {
	var expiringTasks []models.ExpiringTask

	database.DB.Find(&expiringTasks)

	database.UpdateExpiringTask(expiringTasks)
	//utils.UpdateExpiringTaskTimeLeft(expiringTasks)

	json.NewEncoder(w).Encode(&expiringTasks)
}

func CreateExpiringTask(w http.ResponseWriter, r *http.Request) {
	var expiringTask models.ExpiringTask
	json.NewDecoder(r.Body).Decode(&expiringTask)

	createdExpiringTask := database.DB.Create(&expiringTask)
	err := createdExpiringTask.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdExpiringTask)
}

func DeleteExpiringTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var expiringTask models.ExpiringTask

	database.DB.First(&expiringTask, params["id"])
	database.DB.Delete(&expiringTask)

	// also delete from cache
	//cache.DeleteFromCache(cache.REDIS, params["id"])

	json.NewEncoder(w).Encode(&expiringTask)
}
