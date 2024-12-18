package handlers

import (
	"net/http"

	"github.com/xneverov/todo-list/internal/config"
)

func StartServer() error {
	port := config.Get("TODO_PORT")

	if port[0] != ':' {
		port = ":" + port
	}

	fileServer := http.FileServer(http.Dir("./web"))
	http.Handle("/", fileServer)

	return http.ListenAndServe(port, nil)
}
