package api

import (
	"encoding/json"
	"github.com/xneverov/todo-list/internal/db"
	"github.com/xneverov/todo-list/internal/models"
	"net/http"
)

type SuccessResponse struct {
	Tasks []models.Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

const tasksLimit = 30

func HandleTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	query := req.URL.Query()
	search := query.Get("search")

	tasks, err := db.GetTasks(tasksLimit, search)
	if err != nil {
		_ = json.NewEncoder(res).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	_ = json.NewEncoder(res).Encode(SuccessResponse{Tasks: tasks})
}
