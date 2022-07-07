package utils

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Loading environment variables
func LoadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Could not load env variables")
		panic(err)
	}
}

func LoadSwagger(r *mux.Router) {
	// swagger
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	// end swagger
}
