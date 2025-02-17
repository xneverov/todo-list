package api

import (
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("/api/nextdate", HandleNextDate)
	mux.HandleFunc("/api/tasks", HandleTasks)
	mux.HandleFunc("/api/task", HandleTask)

	return mux
}
