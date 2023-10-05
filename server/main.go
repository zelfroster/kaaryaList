package main

import (
	"fmt"
	"log"
)

func main() {
  app := App{}

  err := app.Initialise(USER, DBNAME, PORT)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connnected DB")
  fmt.Println("Serving on http://localhost:9000")

  app.Run("localhost:9000")
}
