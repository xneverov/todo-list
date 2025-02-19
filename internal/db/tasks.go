package db

import (
	"fmt"
	"github.com/xneverov/todo-list/internal/models"
)

func GetTasks(n int) ([]models.Task, error) {
	var tasks []models.Task

	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date DESC 
		LIMIT ?;`

	rows, err := db.Query(query, n)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Description, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}
