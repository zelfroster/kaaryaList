package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
  Router *mux.Router
  DB *sql.DB
}

func (app *App) Initialise(USER string, DBNAME string, PORT int) error {
  psqlconn := fmt.Sprintf("user=%s dbname=%s sslmode=disable port=%v",USER,DBNAME,PORT)

  var err error
  app.DB, err = sql.Open("postgres", psqlconn)
  if err != nil {
    return err
  }

  // Ping the db to check connection (not necessary)
  err = app.DB.Ping()
  if err != nil {
    return err
  }

  // router to handle routes
  app.Router = mux.NewRouter().StrictSlash(true)

  app.handleRoutes()
  return nil
}

func (app *App) Run(address string) {
  log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(rw http.ResponseWriter, statusCode int, payload interface{}) {
  response, _ := json.Marshal(payload)
  rw.Header().Set("Content-type", "application/json")
  rw.Header().Set("Access-Control-Allow-Origin", "*")
  rw.WriteHeader(statusCode)
  rw.Write(response)
}

func sendError(rw http.ResponseWriter, statusCode int, err string) {
  errorMessage := map[string]string{"error": err}
  sendResponse(rw, statusCode, errorMessage)
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

func (app *App) handleRoutes() {
  app.Router.HandleFunc("/getTasks", app.getTasks).Methods("GET")
  app.Router.HandleFunc("/getTask/{id}", app.getTask).Methods("GET")
  app.Router.HandleFunc("/createTask", app.createTask).Methods("POST")
  app.Router.HandleFunc("/updateTask/{id}", app.updateTask).Methods("PUT")
}
