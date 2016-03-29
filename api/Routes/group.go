/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package

   handles: creating a group, removing a group, getting a group,
            adding admins, removing admins -- could logically be
            handled in this file
*/
package routes

import (
    "fmt"
  //  "log"
    "net/http"
  //  "encoding/json"
//    "github.com/dgrijalva/jwt-go"
)

// Routes with /api/ prefix
var GroupMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //Handles GET, POST, and DELETE for groups
  GroupMux.HandleFunc("/group/", test)

  //Handles GET, PUT for Admins in groups
  GroupMux.HandleFunc("/group/admin", user)
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
