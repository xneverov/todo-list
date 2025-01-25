package handlers

import (
	"fmt"
	"net/http"

	"github.com/xneverov/todo-list/internal/config"
	"github.com/xneverov/todo-list/internal/tasks"
)

var port = config.Get("TODO_PORT")

func StartServer() error {
	if port[0] != ':' {
		port = ":" + port
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("/api/nextdate", HandleNextDate)

	return http.ListenAndServe(port, mux)
}

func HandleNextDate(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	result, err := tasks.NextDate(now, date, repeat)
	if err != nil {
		fmt.Fprint(res, err)
		return
	}

	fmt.Fprint(res, result)
}
