package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var defaultValues = map[string]string{
	"TODO_PORT":   "7540",
	"TODO_DBFILE": "./storage/scheduler.db",
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}
}

func Get(key string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValues[key]
	}

	return value
}
