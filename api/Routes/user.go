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
    "github.com/ahermida/dartboardAPI/api/Models"
    "net/http"
    "encoding/json"
  //  "gopkg.in/mgo.v2/bson"
    "github.com/ahermida/dartboardAPI/api/DB"
)

// Routes with /users/ prefix
var UserMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //GET get user info
  UserMux.HandleFunc("/user/", getUser)

  //POST & PUT add and removed saved
  UserMux.HandleFunc("/user/saved", saved)

  //POST get user's saved threads
  UserMux.HandleFunc("/user/threads", threads)

  //POST add a username, PUT to change it, DELETE to remove it
  UserMux.HandleFunc("/user/username", username)

  //GET get all notifications
  UserMux.HandleFunc("/user/notifications", notifications)

  //GET -- gets friends list, POST add a friend -- creates request, PUT accept it, DELETE to remove it,
  UserMux.HandleFunc("/user/friends", friends)
}

/*
   Route handlers for User Routes
*/

// GET Handle /user/
func getUser(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }
  user, err := db.GetUser(id)
  if err != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //else send back user which is already json formatted
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(res).Encode(user); err != nil {
    http.Error(res, http.StatusText(500), 500)
  }
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

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  var thread models.GetUserFeed
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&thread); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get group
  group, err := db.GetGroup(id, thread.Page)
  if err != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //get json struct that we're gonna send over
  grp := &models.SendGroup{
    Threads: group,
  }

  //send back no error response
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  //send over data
  if err := json.NewEncoder(res).Encode(&grp); err != nil {
      http.Error(res, http.StatusText(500), 500)
  }
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
// Handle
func notifications(res http.ResponseWriter, req *http.Request) {
}
