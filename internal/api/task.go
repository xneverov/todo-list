package api

import (
	"encoding/json"
	"github.com/xneverov/todo-list/internal/db"
	"github.com/xneverov/todo-list/internal/models"
	"github.com/xneverov/todo-list/internal/tasks"
	"net/http"
	"time"
)

type createTaskResponse struct {
	ID    *string `json:"id,omitempty"`
	Error string  `json:"error,omitempty"`
}

func HandleTask(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

	case http.MethodPost:
		createTask(res, req)
	case http.MethodDelete:

	default:
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func createTask(res http.ResponseWriter, req *http.Request) {
	var task models.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(createTaskResponse{Error: "Invalid JSON"})
		return
	}

	if task.Title == "" {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(createTaskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	now := time.Now().Format("20060102")

	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(createTaskResponse{Error: "Некорректный формат даты"})
		return
	}

	var nextDate string

	if task.Repeat != "" {
		nextDate, err = tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(res).Encode(createTaskResponse{Error: err.Error()})
			return
		}
	}

	if task.Date < now {
		if nextDate != "" {
			task.Date = nextDate
		} else {
			task.Date = now
		}
	}

	taskID, err := db.CreateTask(&task)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(createTaskResponse{Error: err.Error()})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(res).Encode(createTaskResponse{ID: &taskID})
}
