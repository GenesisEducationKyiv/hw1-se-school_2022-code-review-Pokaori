package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type settings struct {
	EmailName         string
	EmailPass         string
	EmailHost         string
	EmailPort         int
	EmailsStoragePath string
}

var Settings settings

func LoadEnv() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Settings.EmailName = os.Getenv("EMAIL_NAME")
	Settings.EmailPass = os.Getenv("EMAIL_PASSWORD")
	Settings.EmailHost = os.Getenv("EMAIL_HOST")
	Settings.EmailPort, err = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatal("Incoorect value for EmailPort. Should be int.")
	}
	Settings.EmailsStoragePath = os.Getenv("EMAIL_STORAGE_PATH")
}
