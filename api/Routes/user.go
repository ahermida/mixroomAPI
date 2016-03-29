/*
   User Routes that will be used to Create, Read, and Update user
   Creates a ServeMux from the default http package

   handles: registering user, authorizing user, logging in user,
            updating user's password, updating user's username,
            getting a user's friends, adding a saved thread, getting a user's saved threads,
            removing a saved thread, adding a friend, removing a friend, getting notifications,
            add username, remove username
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

  //POST register user
  UserMux.HandleFunc("/user/register", register)

  //POST activate user
  UserMux.HandleFunc("/user/activate", activate)

  //DELETE fake delete user
  UserMux.HandleFunc("/user/remove", activate)

  //POST login user
  UserMux.HandleFunc("/user/login", login)

  //POST update password
  UserMux.HandleFunc("/user/changepass", updatePassword)

  //POST get user info
  UserMux.HandleFunc("/user/info", getUser)

  //POST & PUT add and removed saved
  UserMux.HandleFunc("/user/saved", register)

  //POST get user's saved threads
  UserMux.HandleFunc("/user/threads", register)

  //POST add a username, DELETE to remove it
  UserMux.HandleFunc("/user/username", register)

  //POST change username
  UserMux.HandleFunc("/user/changeusername", register)
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
