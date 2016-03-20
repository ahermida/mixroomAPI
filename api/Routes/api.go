/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package
*/
package routes

import (
    "fmt"
  //  "log"
    "net/http"
  //  "encoding/json"
//    "github.com/dgrijalva/jwt-go"
//    "github.com/ahermida/Writer/resourceGo/DB"
//    "github.com/ahermida/Writer/resourceGo/Config"
)

// Routes with /api/ prefix
var ApiMux = http.NewServeMux()

// Setup Routes with Mux
func init() {
  ApiMux.HandleFunc("/api/test", test)
  ApiMux.HandleFunc("/api/user", user)
  ApiMux.HandleFunc("/api/other", other)
}

/*
   Route handlers for API Routes
*/

// Handle /api/test
func test(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

//handle /api/user (gets information about key via protected route)
func user(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /api/other
func other(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}
