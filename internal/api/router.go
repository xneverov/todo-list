package api

import (
	"github.com/xneverov/todo-list/internal/auth"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("/api/signin", auth.HandleAuth)
	mux.HandleFunc("/api/nextdate", HandleNextDate)
	mux.HandleFunc("/api/task", auth.Middleware(HandleTask))
	mux.HandleFunc("/api/tasks", auth.Middleware(HandleTasks))
	mux.HandleFunc("/api/task/done", auth.Middleware(HandleTaskComplete))

	return mux
}
