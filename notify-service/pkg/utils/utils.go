package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//GetEnv func
func GetEnv(key string) string {
	// load .env file
	switch godotenv.Load() {
	case godotenv.Load("../.env"):
		log.Println("Error loading .env file")
	}
	return os.Getenv(key)
}
