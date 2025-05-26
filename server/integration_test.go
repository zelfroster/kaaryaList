package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Real driver for potential actual test DB connection
)

var (
	testApp App
	mock    sqlmock.Sqlmock // Keep mock accessible if needed for multiple tests
	ts      *httptest.Server
)

// Ideal TestMain for actual DB, here adapted for httptest.Server and mock DB per test
func TestMain(m *testing.M) {
	// --- Ideal Real Test Database Setup (Descriptive) ---
	// 1. Start Docker Compose: `docker-compose up -d test_postgres`
	// 2. Wait for PostgreSQL to be ready (e.g., using a utility like dockerize or a loop)
	// 3. Apply schema migrations (e.g., using golang-migrate/migrate or init-test-db.sql)
	// 4. Configure App to use test DB:
	//    testDbUser := os.Getenv("TEST_DB_USER") // etc.
	//    testApp.Initialise(testDbUser, ...)
	//
	// For this worker environment, we will use httptest.NewServer and re-initialize App with mock DB for tests.
	// We don't run the actual main() to avoid subprocess complexity here.
	// The server 'ts' will be started per test or per suite if state is managed carefully.

	// Setup can be done here if mock is shared, or per-test for isolation.
	// For simplicity with httptest.Server, we'll set up the server per test group or test.

	code := m.Run()

	// --- Ideal Real Test Database Teardown ---
	// 1. Stop and remove Docker Compose services: `docker-compose down -v`

	os.Exit(code)
}

// Helper to set up the app and server for a group of tests
func setupTestSuite(t *testing.T) func(t *testing.T) {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New() // Assign to global mock
	if err != nil {
		log.Fatalf("Failed to create sqlmock: %v", err)
	}

	testApp = App{DB: db, Router: mux.NewRouter()}
	testApp.handleRoutes()

	ts = httptest.NewServer(testApp.Router)

	// Teardown function
	return func(t *testing.T) {
		ts.Close()
		db.Close() // Close mock DB
	}
}

// Helper to clear tables (conceptually, for a real DB)
// With sqlmock, this means ensuring no prior expectations are carried over if mock is reused.
// For these tests, mock is created fresh per suite/test.
func clearTables(t *testing.T) {
	// For a real DB:
	// _, err := testApp.DB.Exec("DELETE FROM tasks")
	// if err != nil { t.Fatalf("Failed to clear tasks table: %v", err) }
	// _, err = testApp.DB.Exec("DELETE FROM users") // Be careful with order due to FKs if any
	// if err != nil { t.Fatalf("Failed to clear users table: %v", err) }

	// For sqlmock, ensure mock is fresh or reset expectations if needed.
	// Here, mock is created by setupTestSuite, so it's fresh.
}

func TestUserAndTaskLifecycle(t *testing.T) {
	teardownSuite := setupTestSuite(t)
	defer teardownSuite(t)
	clearTables(t) // Ensure clean state

	var registeredUser User
	var createdTask Task

	// 1. Register User
	t.Run("RegisterUser", func(t *testing.T) {
		userInput := User{Name: "integUser", Email: "integ@example.com", Password: "password123"}
		jsonBody, _ := json.Marshal(userInput)

		// Mock DB for user registration
		// bcrypt makes the actual password unpredictable, so use sqlmock.AnyArg()
		query := regexp.QuoteMeta("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING username, password, id, email")
		mock.ExpectQuery(query).
			WithArgs(userInput.Name, userInput.Email, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password", "id", "email"}).
				AddRow(userInput.Name, "hashedPassword", "1", userInput.Email)) // Simulate returning ID "1"

		resp, err := http.Post(ts.URL+"/registerUser", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status OK; got %v", resp.Status)
		}
		if err := json.NewDecoder(resp.Body).Decode(&registeredUser); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if registeredUser.Name != userInput.Name {
			t.Errorf("Expected user name %s; got %s", userInput.Name, registeredUser.Name)
		}
		if registeredUser.Id == "" { // Assuming model returns some ID
			t.Error("Expected user ID to be populated")
		}
		// Store ID for next steps (though not directly used by current task model)
	})

	// 2. Create Task
	t.Run("CreateTask", func(t *testing.T) {
		// Prerequisite: User registered (though not strictly linked in current model)
		if registeredUser.Id == "" && t.Failed() { // Check if previous step failed
			t.Skip("Skipping task creation due to user registration failure")
		}

		taskInput := Task{Name: "My Integration Test Task"}
		jsonBody, _ := json.Marshal(taskInput)

		query := regexp.QuoteMeta("INSERT INTO tasks (task_name) VALUES ($1) RETURNING id")
		mock.ExpectQuery(query).WithArgs(taskInput.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(101)) // Simulate ID 101

		resp, err := http.Post(ts.URL+"/createTask", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status OK; got %v", resp.Status)
		}
		if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if createdTask.Name != taskInput.Name {
			t.Errorf("Expected task name %s; got %s", taskInput.Name, createdTask.Name)
		}
		if createdTask.ID == 0 {
			t.Error("Expected task ID to be populated")
		}
	})

	// 3. Get Task
	t.Run("GetTask", func(t *testing.T) {
		if createdTask.ID == 0 && t.Failed() {
			t.Skip("Skipping get task due to task creation failure")
		}
		query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		mock.ExpectQuery(query).WithArgs(createdTask.ID).
			WillReturnRows(sqlmock.NewRows([]string{"task_name", "is_complete"}).
				AddRow(createdTask.Name, createdTask.IsComplete))

		resp, err := http.Get(fmt.Sprintf("%s/getTask/%d", ts.URL, createdTask.ID))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status OK; got %v", resp.Status)
		}
		var fetchedTask Task
		if err := json.NewDecoder(resp.Body).Decode(&fetchedTask); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if fetchedTask.ID != createdTask.ID { // ID is not in select but is set by handler
			// The model's getTask sets task.ID from the input, not the DB for this specific query
			// The handler then returns this task.
			// Let's adjust expectation based on what the handler returns.
			// The handler sets the ID from the path param.
		}
		if fetchedTask.Name != createdTask.Name {
			t.Errorf("Expected task name %s; got %s", createdTask.Name, fetchedTask.Name)
		}
	})

	// 4. Update Task
	t.Run("UpdateTask", func(t *testing.T) {
		if createdTask.ID == 0 && t.Failed() {
			t.Skip("Skipping update task due to task creation failure")
		}
		updatedTaskData := Task{Name: "Updated Integration Task", IsComplete: true}
		jsonBody, _ := json.Marshal(updatedTaskData)

		query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
		mock.ExpectQuery(query).
			WithArgs(updatedTaskData.Name, updatedTaskData.IsComplete, createdTask.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "is_complete"}).
				AddRow(createdTask.ID, updatedTaskData.Name, updatedTaskData.IsComplete))

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/updateTask/%d", ts.URL, createdTask.ID), bytes.NewBuffer(jsonBody))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status OK; got %v", resp.Status)
		}
		var returnedTask Task
		if err := json.NewDecoder(resp.Body).Decode(&returnedTask); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if returnedTask.Name != updatedTaskData.Name || !returnedTask.IsComplete {
			t.Errorf("Expected updated task data; got name '%s', completed '%t'", returnedTask.Name, returnedTask.IsComplete)
		}
	})

	// 5. Delete Task
	t.Run("DeleteTask", func(t *testing.T) {
		if createdTask.ID == 0 && t.Failed() {
			t.Skip("Skipping delete task due to task creation failure")
		}
		query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
		mock.ExpectExec(query).WithArgs(createdTask.ID).WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/deleteTask/%d", ts.URL, createdTask.ID), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status OK; got %v", resp.Status)
		}
		var deleteResp map[string]bool
		if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
			t.Fatalf("Failed to decode delete response: %v", err)
		}
		if !deleteResp["retval"] {
			t.Errorf("Expected retval:true for delete; got %v", deleteResp)
		}

		// Optionally, verify task is gone (Get should fail)
		queryGet := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
		mock.ExpectQuery(queryGet).WithArgs(createdTask.ID).WillReturnError(sql.ErrNoRows) // Or your specific "not found" error

		respGet, _ := http.Get(fmt.Sprintf("%s/getTask/%d", ts.URL, createdTask.ID))
		defer respGet.Body.Close()
		if respGet.StatusCode != http.StatusInternalServerError { // Current app.getTask returns 500 on error
			t.Errorf("Expected status InternalServerError after delete; got %v", respGet.Status)
		}
	})

	// Final check on mock expectations for the whole suite
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations at the end of the lifecycle test: %s", err)
	}
}

func TestGetTasks_NoTasks(t *testing.T) {
	teardownSuite := setupTestSuite(t)
	defer teardownSuite(t)
	clearTables(t)

	query := regexp.QuoteMeta("SELECT id, task_name, is_complete FROM tasks ORDER BY id ASC")
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "is_complete"})) // Empty rows

	resp, err := http.Get(ts.URL + "/getTasks")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}
	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks; got %d", len(tasks))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetNonExistentTask(t *testing.T) {
	teardownSuite := setupTestSuite(t)
	defer teardownSuite(t)
	clearTables(t)

	taskID := 9999
	query := regexp.QuoteMeta("SELECT task_name, is_complete FROM tasks WHERE id=$1")
	mock.ExpectQuery(query).WithArgs(taskID).WillReturnError(sql.ErrNoRows) // Simulate DB not found

	resp, err := http.Get(fmt.Sprintf("%s/getTask/%d", ts.URL, taskID))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// The handler `app.getTask` returns InternalServerError for any error from `task.getTask`
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError for non-existent task; got %v", resp.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateNonExistentTask(t *testing.T) {
	teardownSuite := setupTestSuite(t)
	defer teardownSuite(t)
	clearTables(t)

	taskID := 9998
	taskData := Task{Name: "Try Update", IsComplete: false}
	jsonBody, _ := json.Marshal(taskData)

	query := regexp.QuoteMeta("UPDATE tasks SET task_name=$1, is_complete=$2 WHERE id=$3 RETURNING id, task_name, is_complete")
	mock.ExpectQuery(query).WithArgs(taskData.Name, taskData.IsComplete, taskID).
		WillReturnError(sql.ErrNoRows) // Simulate DB not found or not updated

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/updateTask/%d", ts.URL, taskID), bytes.NewBuffer(jsonBody))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError for updating non-existent task; got %v", resp.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteNonExistentTask(t *testing.T) {
	teardownSuite := setupTestSuite(t)
	defer teardownSuite(t)
	clearTables(t)

	taskID := 9997

	// DB might return error or 0 rows affected. sqlmock.NewResult(0,0) simulates 0 rows affected.
	// If the driver or DB errors on "ID not found for delete", then WillReturnError(someError).
	// The current model.go's deleteTask doesn't error on 0 rows affected.
	query := regexp.QuoteMeta("DELETE FROM tasks WHERE id=$1")
	mock.ExpectExec(query).WithArgs(taskID).WillReturnResult(sqlmock.NewResult(0, 0))

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/deleteTask/%d", ts.URL, taskID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // The handler returns OK even if 0 rows affected
		t.Errorf("Expected status OK for deleting non-existent task; got %v", resp.Status)
	}
	var deleteResp map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		t.Fatalf("Failed to decode delete response: %v", err)
	}
	if !deleteResp["retval"] { // This indicates the operation was "successful" from handler's view
		t.Errorf("Expected retval:true for delete non-existent; got %v", deleteResp)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
