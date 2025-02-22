package db

import (
	"errors"
	"strconv"
	"time"

	"github.com/xneverov/todo-list/internal/models"
	"github.com/xneverov/todo-list/internal/tasks"
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

func ReadTask(id string) (models.Task, error) {
	const query = `SELECT * FROM scheduler WHERE id = ?`

	var task models.Task

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Description, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

func UpdateTask(task *models.Task) error {
	const query = `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Description, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	const query = `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func CompleteTask(id string) error {
	const query = `UPDATE scheduler SET date = ? WHERE id = ?`

	task, err := ReadTask(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		_ = DeleteTask(id)
		return nil
	}

	now := time.Now().Format("20060102")

	nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, nextDate, id)
	if err != nil {
		return err
	}

	return nil
}
