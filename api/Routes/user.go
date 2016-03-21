/*
   User Routes that will be used to Create, Read, and Update user
   Creates a ServeMux from the default http package

*/
package routes

import (
    "net/http"
    "fmt"
)

// Routes with /users/ prefix
var UserMux = http.NewServeMux()

// Setup Routes with Mux
func init() {
  //make user
  UserMux.HandleFunc("/user/register", register)
  //get user info
  UserMux.HandleFunc("/user/activate", activate)
  //login user
  UserMux.HandleFunc("/user/login", login)
  //update password
  UserMux.HandleFunc("/user/updatepassword", updatePassword)
  //get user
  UserMux.HandleFunc("/user/getuser", getUser)
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

// Handle /user/userID
func getUser(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}
