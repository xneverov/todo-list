package db

import (
	"github.com/xneverov/todo-list/internal/models"
	"strconv"
)

func CreateTask(task *models.Task) (string, error) {
	const query = `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Description, task.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func GetTask(id string) (models.Task, error) {
	task := models.Task{}

	const query = `SELECT * FROM scheduler WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Description, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}
