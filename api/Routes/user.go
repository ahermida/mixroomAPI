/*
   User Routes that will be used to Create, Read, and Update user
   Creates a ServeMux from the default http package

   handles: updating user's username, getting a user's friends, adding a saved thread,
            getting a user's saved threads, removing a saved thread, adding a friend,
            removing a friend, getting notifications, add username, remove username
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

  //POST get user info
  UserMux.HandleFunc("/user/info", getUser)

  //POST & PUT add and removed saved
  UserMux.HandleFunc("/user/saved", register)

  //POST get user's saved threads
  UserMux.HandleFunc("/user/threads", getUserThreads)

  //POST add a username, PUT to change it, DELETE to remove it
  UserMux.HandleFunc("/user/username", username)

  //POST add a friend -- creates request, PUT accept it, DELETE to remove it,
  UserMux.HandleFunc("/user/friends", friends)
}

/*
   Route handlers for User Routes
*/

// Handle /user/userID
func getUser(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/userID
func getUserThreads(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/userID
func username(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/userID
func friends(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}
