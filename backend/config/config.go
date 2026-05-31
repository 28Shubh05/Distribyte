package config

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10 MB
)

func LoadEnv() {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println(".env loaded successfully")
}
