package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var defaultValues = map[string]string{
	"TODO_PORT":     "7540",
	"TODO_DBFILE":   "./storage/scheduler.db",
	"TODO_PASSWORD": "",
}

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются стандартные значения")
	}
}

func Get(key string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValues[key]
	}

	return value
}
