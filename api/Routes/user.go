/*
   User Routes that will be used to Create, Read, and Update user
   Creates a ServeMux from the default http package

   handles: updating user's username, getting a user's friends, adding a saved thread,
            getting a user's saved threads, removing a saved thread, adding a friend,
            removing a friend, getting notifications, add username, remove username
*/
package routes

import (
    "github.com/ahermida/dartboardAPI/api/Util"
    "fmt"
    "github.com/ahermida/dartboardAPI/api/Models"
    "net/http"
    "encoding/json"
    "gopkg.in/mgo.v2/bson"
    "github.com/ahermida/dartboardAPI/api/DB"
)

// Routes with /users/ prefix
var UserMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST get user info
  UserMux.HandleFunc("/user/info", getUser)

  //POST & PUT add and removed saved
  UserMux.HandleFunc("/user/saved", saved)

  //POST get user's saved threads
  UserMux.HandleFunc("/user/threads", threads)

  //POST add a username, PUT to change it, DELETE to remove it
  UserMux.HandleFunc("/user/username", username)

  //POST add a friend -- creates request, PUT accept it, DELETE to remove it,
  UserMux.HandleFunc("/user/friends", friends)
}

/*
   Route handlers for User Routes
*/

// POST Handle /user/info
func getUser(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/saved
func saved(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    addSaved(res, req)
  case "PUT":
    removeSaved(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

// POST Handle /user/threads -- get user's threads
func threads(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  fmt.Fprintf(res, "Test Passed!")
}

// Handle /user/username
func username(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    addUsername(res, req)
  case "PUT":
    changeUsername(res, req)
  case "DELETE":
    removeUsername(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

// Handle /user/userID
func friends(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    addFriend(res, req)
  case "PUT":
    acceptFriend(res, req)
  case "DELETE":
    removeFriend(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

// POST Handle adding saved thread
func addSaved(res http.ResponseWriter, req *http.Request) {
}
// PUT Handle
func removeSaved(res http.ResponseWriter, req *http.Request) {
}
// POST Handle
func addUsername(res http.ResponseWriter, req *http.Request) {
}
// PUT Handle
func changeUsername(res http.ResponseWriter, req *http.Request) {
}
// DELETE Handle
func removeUsername(res http.ResponseWriter, req *http.Request) {
}
// POST Handle
func addFriend(res http.ResponseWriter, req *http.Request) {
}
// PUT Handle
func acceptFriend(res http.ResponseWriter, req *http.Request) {
}
// DELETE Handle
func removeFriend(res http.ResponseWriter, req *http.Request) {
}
