package main

import (
	"log"

	"github.com/xneverov/todo-list/internal/config"
	"github.com/xneverov/todo-list/internal/handlers"
)

func main() {
	config.Load()

	err := handlers.StartServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
