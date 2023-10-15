package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
  // Load Environment Variables
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  user := os.Getenv("USER")
  dbname := os.Getenv("DBNAME")
  port, err := strconv.Atoi(os.Getenv("PORT"))
  if err != nil {
    log.Fatal("Error converting string to integer: ", err)
  }

  // Initialise App
  app := App{}

  err = app.Initialise(user, dbname, port)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connnected DB")
  fmt.Println("Serving on http://localhost:9000")

  app.Run("localhost:9000")
}
