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

type User struct {
	Id       string `json:"id"`
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) registerUser(db *sql.DB) error {
	queryString := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email"

	// store data according to the columns in database
	err := db.QueryRow(queryString, user.Name, user.Email, user.Password).Scan(&user.Name, &user.Password, &user.Id, &user.Email)
	if err != nil {
		return err
	}
	return nil
}

func getTasks(db *sql.DB) ([]Task, error) {
	queryString := "SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC"
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
	queryString := "SELECT task_name, is_complete FROM tasks WHERE id=$1"
	row := db.QueryRow(queryString, task.ID)

	err := row.Scan(&task.Name, &task.IsComplete)
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) createTask(db *sql.DB) error {
	queryString := "INSERT INTO tasks (task_name) VALUES ($1) RETURNING id"
	err := db.QueryRow(queryString, task.Name).Scan(&task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) updateTask(db *sql.DB) error {
	queryString := "UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete"
	row := db.QueryRow(queryString, task.Name, task.IsComplete, task.ID)
	err := row.Scan(&task.ID, &task.Name, &task.IsComplete)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) deleteTask(db *sql.DB) error {
	queryString := "DELETE FROM tasks WHERE id=$1"
	_, err := db.Exec(queryString, task.ID)
	if err != nil {
		return err
	}
	return nil
}
