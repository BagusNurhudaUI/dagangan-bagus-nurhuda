package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// functions to get key from environment
func GetEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
