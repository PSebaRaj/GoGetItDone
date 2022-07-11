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

// Loading swagger into /docs
func LoadSwagger(r *mux.Router) {
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	fmt.Println("Swagger can be found /docs")

	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	// end swagger
}
