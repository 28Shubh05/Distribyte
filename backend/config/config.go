package config

import (
	"log"
	"os"

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

	JWTSecret = os.Getenv("JWT_SECRET")

	log.Println(".env loaded successfully")
}

var JWTSecret string
