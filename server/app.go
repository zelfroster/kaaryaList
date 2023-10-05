package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

  err = app.DB.Ping()
  if err != nil {
    return err
  }

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
  }
  log.Println(tasks)
  sendResponse(rw, http.StatusOK, tasks)
}

func (app *App) handleRoutes() {
  app.Router.HandleFunc("/getTasks", app.getTasks).Methods("GET")
}
