package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	// bcrypt is used by the handler, so it's part of the test scope
)

// Helper to create a new App with a mock DB and router for testing
func newTestApp(t *testing.T) (*App, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	app := &App{
		DB:     db,
		Router: mux.NewRouter(),
	}
	app.handleRoutes() // Important to set up the routes
	return app, mock
}

func TestRegisterUserHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful registration", func(t *testing.T) {
		userInput := User{Name: "testuser", Email: "test@example.com", Password: "password123"}
		jsonBody, _ := json.Marshal(userInput)
		req, _ := http.NewRequest("POST", "/registerUser", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		// Expect the INSERT query from model.go's registerUser
		// Note: bcrypt.GenerateFromPassword makes the exact password arg unpredictable here if we were to mock it.
		// We are testing that the handler calls the model's registerUser, which then makes this query.
		// The password in WithArgs should be the hashed one. Since bcrypt is non-deterministic without seeding,
		// we use sqlmock.AnyArg() for the password field for robustness in this test.
		// Or, more accurately, we know the model's registerUser will be called, and we tested that separately.
		// Here, we focus on the handler's behavior around it.

		query := regexp.QuoteMeta("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email")
		mock.ExpectQuery(query).
			WithArgs(userInput.Name, userInput.Email, sqlmock.AnyArg()). // Password will be hashed
			WillReturnRows(sqlmock.NewRows([]string{"username", "password", "id", "email"}).
				AddRow(userInput.Name, "hashedPassword", "mockUserID", userInput.Email))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var returnedUser User
		err := json.NewDecoder(rr.Body).Decode(&returnedUser)
		if err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if returnedUser.Name != userInput.Name {
			t.Errorf("handler returned unexpected body for Name: got %s want %s", returnedUser.Name, userInput.Name)
		}
		if returnedUser.Id != "mockUserID" {
			t.Errorf("handler returned unexpected body for ID: got %s want %s", returnedUser.Id, "mockUserID")
		}
		// Password should not be returned or should be the hashed one; the model test handles this.
		// Here, we check that the handler marshals what the model returns.

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/registerUser", strings.NewReader("{malformed"))
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError { // Based on SendError in app.go
			t.Errorf("handler returned wrong status code for malformed JSON: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("database error on registration", func(t *testing.T) {
		userInput := User{Name: "dbfailuser", Email: "dbfail@example.com", Password: "password123"}
		jsonBody, _ := json.Marshal(userInput)
		req, _ := http.NewRequest("POST", "/registerUser", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email")
		mock.ExpectQuery(query).
			WithArgs(userInput.Name, userInput.Email, sqlmock.AnyArg()).
			WillReturnError(errors.New("db insert error"))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest { // Based on SendError in app.go
			t.Errorf("handler returned wrong status code for db error: got %v want %v", status, http.StatusBadRequest)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetTasksHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful retrieval", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/getTasks", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		rows := sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
			AddRow(1, "Task 1", false).
			AddRow(2, "Task 2", true)
		mock.ExpectQuery(query).WillReturnRows(rows)

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var tasks []Task
		if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/getTasks", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
		mock.ExpectQuery(query).WillReturnError(errors.New("db query error"))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetTaskHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful retrieval", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/getTask/1", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		rows := sqlmock.NewRows([]string{"task_name", "is_complete"}).AddRow("Task 1", false)
		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}
		if task.Name != "Task 1" {
			t.Errorf("expected task name 'Task 1', got '%s'", task.Name)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("invalid task ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/getTask/abc", nil)
		rr := httptest.NewRecorder()

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		// No DB interaction expected
	})

	t.Run("task not found (db error)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/getTask/404", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		mock.ExpectQuery(query).WithArgs(404).WillReturnError(errors.New("sql: no rows in result set")) // Simulate model error

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError { // As per app.getTask
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCreateTaskHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful creation", func(t *testing.T) {
		taskInput := Task{Name: "New Task"}
		jsonBody, _ := json.Marshal(taskInput)
		req, _ := http.NewRequest("POST", "/createTask", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("INSERT INTO tasks (task_name) VALUES ($1) RETURNING id")
		mock.ExpectQuery(query).WithArgs(taskInput.Name).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}
		if task.ID != 123 {
			t.Errorf("expected task ID 123, got %d", task.ID)
		}
		if task.Name != taskInput.Name { // Model populates ID, name should persist from input
			t.Errorf("expected task name '%s', got '%s'", taskInput.Name, task.Name)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		taskInput := Task{Name: "DB Fail Task"}
		jsonBody, _ := json.Marshal(taskInput)
		req, _ := http.NewRequest("POST", "/createTask", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("INSERT INTO tasks (task_name) VALUES ($1) RETURNING id")
		mock.ExpectQuery(query).WithArgs(taskInput.Name).WillReturnError(errors.New("db insert error"))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestUpdateTaskHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful update", func(t *testing.T) {
		taskInput := Task{Name: "Updated Name", IsComplete: true} // ID comes from URL
		jsonBody, _ := json.Marshal(taskInput)
		req, _ := http.NewRequest("PUT", "/updateTask/1", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
		mock.ExpectQuery(query).
			WithArgs(taskInput.Name, taskInput.IsComplete, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
				AddRow(1, taskInput.Name, taskInput.IsComplete))

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}
		if task.ID != 1 || task.Name != taskInput.Name || task.IsComplete != taskInput.IsComplete {
			t.Errorf("handler returned unexpected body: got %+v want ID=1, Name=%s, IsComplete=%t", task, taskInput.Name, taskInput.IsComplete)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("invalid task ID in URL", func(t *testing.T) {
		taskInput := Task{Name: "Updated Name", IsComplete: true}
		jsonBody, _ := json.Marshal(taskInput)
		req, _ := http.NewRequest("PUT", "/updateTask/abc", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		app.Router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("database error on update", func(t *testing.T) {
		taskInput := Task{Name: "Updated Name Fail", IsComplete: true}
		jsonBody, _ := json.Marshal(taskInput)
		req, _ := http.NewRequest("PUT", "/updateTask/2", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
		mock.ExpectQuery(query).
			WithArgs(taskInput.Name, taskInput.IsComplete, 2).
			WillReturnError(errors.New("db update error"))

		app.Router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDeleteTaskHandler(t *testing.T) {
	app, mock := newTestApp(t)
	defer app.DB.Close()

	t.Run("successful deletion", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/deleteTask/1", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		app.Router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var resp map[string]bool
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}
		if val, ok := resp["retval"]; !ok || !val {
			t.Errorf("expected response map[string]bool{\"retval\": true}, got %v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("invalid task ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/deleteTask/xyz", nil)
		rr := httptest.NewRecorder()

		app.Router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("database error on delete", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/deleteTask/2", nil)
		rr := httptest.NewRecorder()

		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).WithArgs(2).WillReturnError(errors.New("db delete error"))

		app.Router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

// Note: The SendResponse and SendError functions are implicitly tested by checking
// the status code and body of the responses in these handler tests.
// If they had more complex logic, they might warrant their own direct tests.
// For bcrypt, a more thorough test would involve an interface for password hashing
// to allow mocking it, but that's beyond the scope of just testing the existing app.go.
// The current tests for registerUser handler assume bcrypt works and focuses on the DB interaction part.
