package main

import (
	"github.com/gorilla/mux"
)

func InitRoutes() { // replace later with seperate routes folder / file to add CORS

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/todos", controllers.GetTodos).Methods("GET")
}

func main() {

	InitRoutes()
}

// test test test
