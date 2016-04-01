/*
  This is the API for dartboard, an anonymous social app in the works.
*/
package main

import (
  "fmt"
  "github.com/ahermida/dartboardAPI/api/Server"
  "github.com/ahermida/dartboardAPI/api/Config"
  "github.com/ahermida/dartboardAPI/api/DB"
  "os"
  "runtime"
  "syscall"
  "time"
)

/**
 * Initialize the API, all go programs start with 'main' function
 */
func main() {

  //Log Started Server & provide hint for interface
  fmt.Printf("API Server started at: %s\n", time.Now().Format(time.RFC822))
  fmt.Println("\nYou can type: \"exit\" or \"quit\" to shut down the app.\nType \"stats\" to show info about the server.")

  //Get start time
  var startTime int64 = time.Now().Unix()

  //Automatically set to NumCPU -- the number of CPUs available
  runtime.GOMAXPROCS(0)

  //Start the server
  server.Start(config.Port)

  /**
   * Interface for Communicating With Server
   */
  var input string
  for input != "exit" {
    _, _ = fmt.Scanf("%v", &input)
    if input != "exit" {
      switch input {
        case "", "help":
          fmt.Println("You can type: \"exit\" or \"quit\" to quit and \"stats\" to show stats about the server")
        case "exit", "quit":
          fmt.Println("\nShutting down...")
          input = "exit"
        case "stats":
          fmt.Printf("\nCPU cores: %d\nGo routines: %d\nProcess ID: %v\n", runtime.NumCPU(), runtime.NumGoroutine(), syscall.Getpid())
          fmt.Printf("The Application has been running for %d minutes\n", (time.Now().Unix() - startTime) / 60)
      }
    }
  }

  //Exit App -- but close connection to DB first
  db.Connection.Close()
  os.Exit(0)
}
