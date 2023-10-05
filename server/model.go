package main

import (
	"database/sql"
)

type Task struct {
  ID string `json:"id"`
  Name string  `json:"name"`
  IsComplete bool  `json:"isComplete"`
}

func getTasks(db *sql.DB) ([]Task, error) {
  queryString := "SELECT id, name, is_complete FROM tasks"
  rows, err := db.Query(queryString)
  if err != nil {
    return nil, err
  }

  tasks := []Task{}

  for rows.Next() {
    var task Task
    err := rows.Scan(&task.ID, &task.Name, &task.IsComplete)
    if err != nil {
      return nil, err
    }
    tasks = append(tasks, task)
  }

  return tasks, nil
}
