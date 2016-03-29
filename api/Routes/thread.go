/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package

   handles: creating a thread, creating a post, removing a post,
            removing a thread, changing a post
            -- could be written in this file

*/
package routes

import (
  //  "fmt"
  //  "log"
    "net/http"
  //  "encoding/json"
//    "github.com/dgrijalva/jwt-go"
)

// Routes with /api/ prefix
var ApiMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST for creating thread, DELETE removing thread
  ApiMux.HandleFunc("/thread/", test)

  //POST for creating a post, DELETE for removing it, PUT for changing it
  ApiMux.HandleFunc("/thread/post", user)
}

/*
   Route handlers for Thread  Routes
*/
