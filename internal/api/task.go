package api

import (
	"encoding/json"
	"errors"
	"github.com/xneverov/todo-list/internal/db"
	"github.com/xneverov/todo-list/internal/models"
	"github.com/xneverov/todo-list/internal/tasks"
	"net/http"
	"time"
)

type taskResponse struct {
	ID    *string `json:"id,omitempty"`
	Error string  `json:"error,omitempty"`
}

func HandleTask(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		createTask(res, req)
	case http.MethodGet:
		readTask(res, req)
	case http.MethodPut:
		updateTask(res, req)
	case http.MethodDelete:
		//deleteTask(res, req)
	default:
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func createTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var task models.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: "Invalid JSON"})
		return
	}

	if err := validateTask(&task); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: err.Error()})
		return
	}

	taskID, err := db.CreateTask(&task)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: err.Error()})
		return
	}

	res.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(res).Encode(taskResponse{ID: &taskID})
}

func readTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	query := req.URL.Query()
	taskID := query.Get("id")

	if taskID == "" {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: "Не указан идентификатор задачи"})
		return
	}

	task, err := db.ReadTask(taskID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: "Задача не найдена"})
		return
	}

	_ = json.NewEncoder(res).Encode(task)
}

func updateTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var task models.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: "Invalid JSON"})
		return
	}

	if err := validateTask(&task); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: err.Error()})
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(taskResponse{Error: "Задача не найдена"})
		return
	}

	res.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(res).Encode(struct{}{})
}

func validateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("Не указан заголовок задачи")
	}

	now := time.Now().Format("20060102")

	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("Некорректный формат даты")
	}

	var nextDate string

	if task.Repeat != "" {
		nextDate, err = tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if task.Date < now {
		if nextDate != "" {
			task.Date = nextDate
		} else {
			task.Date = now
		}
	}

	return nil
}
