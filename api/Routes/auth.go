/*
   User Routes that will be used to manage basic user management utils

   handles: registering user, authorizing user, logging in user,
            updating user's password, removing a user
            -- split up for simplicity
*/
package routes

import (
    "net/http"
    "fmt"
)

// Routes with /users/ prefix
var AuthMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST register user
  AuthMux.HandleFunc("/auth/register", register)

  //POST activate user
  AuthMux.HandleFunc("/auth/activate", activate)

  //DELETE deactivate user
  AuthMux.HandleFunc("/auth/remove", deactivate)

  //POST login user
  AuthMux.HandleFunc("/auth/login", login)

  //POST update password
  AuthMux.HandleFunc("/auth/changepass", updatePassword)

}

/*
   Route handlers for User Routes
*/

// Handle /user/register
func register(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/activate
func activate(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/deactivate
func deactivate(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/login
func login(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/updatepassword
func updatePassword(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  //check if password is the same
  //call change password
  fmt.Fprintf(res, "Test Passed!")
}
