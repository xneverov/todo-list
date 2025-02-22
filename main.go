package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/xneverov/todo-list/internal/api"
	"github.com/xneverov/todo-list/internal/config"
	"github.com/xneverov/todo-list/internal/db"
)

func main() {
	config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db.Get())

	port := config.Get("TODO_PORT")
	if port[0] != ':' {
		port = ":" + port
	}

	router := api.SetupRouter()

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
