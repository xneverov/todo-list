package db

import (
	"fmt"
	"github.com/xneverov/todo-list/internal/models"
	"time"
)

const (
	queryTasks = `
		SELECT *
		FROM scheduler 
		ORDER BY date DESC 
		LIMIT ?;`

	queryTasksSearch = `
		SELECT *
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ? 
		ORDER BY date DESC
		LIMIT ?`

	queryTasksSearchDate = `
		SELECT *
		FROM scheduler
		WHERE date = ?
		LIMIT ?`
)

func GetTasks(limit int, search string) ([]models.Task, error) {
	var query string
	var args []interface{}
	var date string

	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		date = t.Format("20060102")
	}

	switch {
	case date != "":
		query = queryTasksSearchDate
		args = append(args, date, limit)
	case search != "":
		query = queryTasksSearch
		search = "%" + search + "%"
		args = append(args, search, search, limit)
	default:
		query = queryTasks
		args = append(args, limit)
	}

	return fetchTasks(query, args...)
}

func fetchTasks(query string, args ...interface{}) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Description, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}
