package api

import (
	"fmt"
	"net/http"

	"github.com/xneverov/todo-list/internal/tasks"
)

func HandleNextDate(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	result, err := tasks.NextDate(now, date, repeat)
	if err != nil {
		_, _ = fmt.Fprint(res, err)
		return
	}

	_, _ = fmt.Fprint(res, result)
}
