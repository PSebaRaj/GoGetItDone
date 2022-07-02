package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Loading enviroment variables
func LoadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Could not load env variables")
		panic(err)
	}
}
