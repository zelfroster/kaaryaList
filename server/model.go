package main

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsComplete bool   `json:"isComplete"`
}

func getTasks(db *sql.DB) ([]Task, error) {
	queryString := "SELECT id, name, is_complete FROM tasks ORDER BY id ASC"
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}

	// array of Tasks
	tasks := []Task{}

	// Loop over the rows returned after querying the database
	for rows.Next() {
		var task Task
		// store the values in task
		err := rows.Scan(&task.ID, &task.Name, &task.IsComplete)
		if err != nil {
			return nil, err
		}
		// append to array
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (task *Task) getTask(db *sql.DB) error {
	queryString := fmt.Sprintf("SELECT name, is_complete FROM tasks WHERE id=%v", task.ID)
	row := db.QueryRow(queryString)

	err := row.Scan(&task.Name, &task.IsComplete)
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) createTask(db *sql.DB) error {
	queryString := fmt.Sprintf("INSERT INTO tasks (name) VALUES ('%v') RETURNING id", task.Name)
	err := db.QueryRow(queryString).Scan(&task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) updateTask(db *sql.DB) error {
  queryString := fmt.Sprintf("UPDATE tasks SET name='%v', is_complete=%v WHERE id=%v RETURNING *", task.Name, task.IsComplete, task.ID)
	row := db.QueryRow(queryString)
	err := row.Scan(&task.ID, &task.Name, &task.IsComplete)
	if err != nil {
		return err
	}
	return nil
}
