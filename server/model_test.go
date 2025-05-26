package main

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := User{Name: "testuser", Email: "test@example.com", Password: "password123"}
	hashedPassword := "hashedPassword" // Assume password was hashed before this stage for the model
	user.Password = hashedPassword     // The model receives the already hashed password

	t.Run("successful registration", func(t *testing.T) {
		query := regexp.QuoteMeta("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email")
		mock.ExpectQuery(query).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password", "id", "email"}).
				AddRow(user.Name, user.Password, "mockUserID", user.Email))

		err := user.registerUser(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if user.Id != "mockUserID" {
			t.Errorf("expected user ID to be 'mockUserID', but got %s", user.Id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		query := regexp.QuoteMeta("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email")
		mock.ExpectQuery(query).
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnError(errors.New("db error"))

		err := user.registerUser(db)
		if err == nil {
			t.Errorf("expected a database error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("successful retrieval of multiple tasks", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		rows := sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
			AddRow(1, "Task 1", false).
			AddRow(2, "Task 2", true)
		mock.ExpectQuery(query).WillReturnRows(rows)

		tasks, err := getTasks(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, but got %d", len(tasks))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("successful retrieval of zero tasks", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		rows := sqlmock.NewRows([]string{"id", "task_name", "is_complete"})
		mock.ExpectQuery(query).WillReturnRows(rows)

		tasks, err := getTasks(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected 0 tasks, but got %d", len(tasks))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database query error", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		mock.ExpectQuery(query).WillReturnError(errors.New("db query error"))

		_, err := getTasks(db)
		if err == nil {
			t.Errorf("expected a database query error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("row scanning error", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		rows := sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
			AddRow(1, "Task 1", false).
			AddRow("invalid_id", "Task 2", true) // Error in this row
		mock.ExpectQuery(query).WillReturnRows(rows)

		_, err := getTasks(db)
		if err == nil {
			t.Errorf("expected a row scanning error, but got none")
		}
		// Note: The exact error message might depend on the driver and how Scan handles it.
		// We're checking that *an* error occurs.

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	task := Task{ID: 1}

	t.Run("successful retrieval", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		rows := sqlmock.NewRows([]string{"task_name", "is_complete"}).AddRow("Test Task", false)
		mock.ExpectQuery(query).WithArgs(task.ID).WillReturnRows(rows)

		err := task.getTask(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if task.Name != "Test Task" {
			t.Errorf("expected task name 'Test Task', but got '%s'", task.Name)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("task not found", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		mock.ExpectQuery(query).WithArgs(task.ID).WillReturnError(errors.New("sql: no rows in result set")) // Simulating sql.ErrNoRows behavior

		err := task.getTask(db)
		if err == nil {
			t.Errorf("expected an error for task not found, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("other database error", func(t *testing.T) {
		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		mock.ExpectQuery(query).WithArgs(task.ID).WillReturnError(errors.New("db error"))

		err := task.getTask(db)
		if err == nil {
			t.Errorf("expected a database error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	task := Task{Name: "New Task"}

	t.Run("successful creation", func(t *testing.T) {
		query := regexp.QuoteMeta("INSERT INTO tasks (task_name) VALUES ($1) RETURNING id")
		mock.ExpectQuery(query).
			WithArgs(task.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

		err := task.createTask(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if task.ID != 123 {
			t.Errorf("expected task ID to be 123, but got %d", task.ID)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		query := regexp.QuoteMeta("INSERT INTO tasks (task_name) VALUES ($1) RETURNING id")
		mock.ExpectQuery(query).
			WithArgs(task.Name).
			WillReturnError(errors.New("db error"))

		err := task.createTask(db)
		if err == nil {
			t.Errorf("expected a database error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	task := Task{ID: 1, Name: "Updated Task", IsComplete: true}

	t.Run("successful update", func(t *testing.T) {
		query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
		mock.ExpectQuery(query).
			WithArgs(task.Name, task.IsComplete, task.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
				AddRow(task.ID, task.Name, task.IsComplete))

		err := task.updateTask(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
		mock.ExpectQuery(query).
			WithArgs(task.Name, task.IsComplete, task.ID).
			WillReturnError(errors.New("db error"))

		err := task.updateTask(db)
		if err == nil {
			t.Errorf("expected a database error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	task := Task{ID: 1}

	t.Run("successful deletion", func(t *testing.T) {
		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).
			WithArgs(task.ID).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		err := task.deleteTask(db)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).
			WithArgs(task.ID).
			WillReturnError(errors.New("db error"))

		err := task.deleteTask(db)
		if err == nil {
			t.Errorf("expected a database error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("no rows affected", func(t *testing.T) {
		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).
			WithArgs(task.ID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		err := task.deleteTask(db)
		if err != nil {
			// Depending on how you want to handle "no rows affected" for a DELETE.
			// Some might consider it an error, others not. sqlmock itself doesn't error here.
			// The current model.go deleteTask doesn't check RowsAffected(), so it won't return an error.
			t.Errorf("expected no error for 0 rows affected, but got %s", err)
		}


		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
