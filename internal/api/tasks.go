package api

import (
	"encoding/json"
	"net/http"

	"github.com/xneverov/todo-list/internal/db"
	"github.com/xneverov/todo-list/internal/models"
)

type successResponse struct {
	Tasks []models.Task `json:"tasks"`
}

type errorResponse struct {
	Error string `json:"error"`
}

const tasksLimit = 30

func HandleTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	search := req.URL.Query().Get("search")

	tasks, err := db.GetTasks(tasksLimit, search)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(errorResponse{Error: err.Error()})
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	_ = json.NewEncoder(res).Encode(successResponse{Tasks: tasks})
}
