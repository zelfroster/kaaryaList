package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise(USER string, DBNAME string, PORT int, PASSWORD string) error {
	psqlconn := fmt.Sprintf("user=%s dbname=%s sslmode=disable port=%v password=%s", USER, DBNAME, PORT, PASSWORD)

	var err error
	app.DB, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	// router to handle routes
	router := mux.NewRouter().StrictSlash(true)

	// add cors headers
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.Use(handlers.CORS(headers, methods, origins))

	app.Router = router

	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(rw http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.WriteHeader(statusCode)
	rw.Write(response)
}

func sendError(rw http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(rw, statusCode, errorMessage)
}

func (app *App) registerUser(rw http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
	}

	// create a hash of password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		sendError(rw, http.StatusBadRequest, err.Error())
	}

	user.Password = string(password)

	err = user.registerUser(app.DB)
	if err != nil {
		sendError(rw, http.StatusBadRequest, err.Error())
	}
	sendResponse(rw, http.StatusOK, user)
}

func (app *App) getTasks(rw http.ResponseWriter, req *http.Request) {
	tasks, err := getTasks(app.DB)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(rw, http.StatusOK, tasks)
}

func (app *App) getTask(rw http.ResponseWriter, req *http.Request) {
	// get the id from the route accessed e.g. for route /getTask/1 , id will be 1
	varMap := mux.Vars(req)
	key, err := strconv.Atoi(varMap["id"])
	if err != nil {
		sendError(rw, http.StatusBadRequest, "invalid task id")
		return
	}

	task := Task{ID: key}

	err = task.getTask(app.DB)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(rw, http.StatusOK, task)
}

func (app *App) createTask(rw http.ResponseWriter, req *http.Request) {
	var task Task
	// decode data coming in request and assign to task
	json.NewDecoder(req.Body).Decode(&task)

	err := task.createTask(app.DB)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(rw, http.StatusOK, task)
}

func (app *App) updateTask(rw http.ResponseWriter, req *http.Request) {
	varMap := mux.Vars(req)
	key, err := strconv.Atoi(varMap["id"])
	if err != nil {
		sendError(rw, http.StatusBadRequest, "invalid task id")
		return
	}

	task := Task{ID: key}

	json.NewDecoder(req.Body).Decode(&task)

	err = task.updateTask(app.DB)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(rw, http.StatusOK, task)
}

func (app *App) deleteTask(rw http.ResponseWriter, req *http.Request) {
	varMap := mux.Vars(req)
	key, err := strconv.Atoi(varMap["id"])
	if err != nil {
		sendError(rw, http.StatusBadRequest, "invalid task id")
		return
	}

	task := Task{ID: key}

	err = task.deleteTask(app.DB)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(rw, http.StatusOK, map[string]bool{"retval": true})
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/registerUser", app.registerUser).Methods("POST", "OPTIONS")
	app.Router.HandleFunc("/getTasks", app.getTasks).Methods("GET", "OPTIONS")
	app.Router.HandleFunc("/getTask/{id}", app.getTask).Methods("GET", "OPTIONS")
	app.Router.HandleFunc("/createTask", app.createTask).Methods("POST", "OPTIONS")
	app.Router.HandleFunc("/updateTask/{id}", app.updateTask).Methods("PUT", "OPTIONS")
	app.Router.HandleFunc("/deleteTask/{id}", app.deleteTask).Methods("DELETE", "OPTIONS")
}
