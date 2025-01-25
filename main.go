package main

import (
	"log"

	"github.com/xneverov/todo-list/internal/config"
	"github.com/xneverov/todo-list/internal/db"
	"github.com/xneverov/todo-list/internal/handlers"
)

func main() {
	config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer db.Get().Close()

	if err := handlers.StartServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
