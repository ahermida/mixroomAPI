/*
   User Routes that will be used to Create, Read, and Update user
*/
package routes

import (
    "github.com/ahermida/sudopostAPI/api/Util"
    "github.com/ahermida/sudopostAPI/api/Models"
    "net/http"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "github.com/ahermida/sudopostAPI/api/DB"
)

// Routes with /users/ prefix
var UserMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //GET get user info
  UserMux.HandleFunc("/user/", getUser)

  //GET, POST & PUT add and removed saved
  UserMux.HandleFunc("/user/saved", saved)

  //POST change user's name
  UserMux.HandleFunc("/user/name", name)

  //POST search usernames and users' names
  UserMux.HandleFunc("/user/search", searchUsers)

  //POST get user's threads
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
    panic(err)
  }
}

//Handle /user/search
func searchUsers(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
  }

  //user _id in hex
  id := util.GetId(req)

  //users only searchable if we're logged in
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //data returned by request
  var request models.Search

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&request); err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  usernames, err := db.SearchUsers(request.Text)
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
  }

  send := &models.SendUserSearch{
    Usernames: usernames,
  }
  //else send back user which is already json formatted
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(res).Encode(send); err != nil {
    panic(err)
  }
}

// Handle /user/saved
func saved(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "GET":
    getSaved(res, req)
  case "POST":
    savedToggle(res, req)
  case "PUT":
    savedToggle(res, req)
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
  if err := json.NewEncoder(res).Encode(grp); err != nil {
    panic(err)
  }
}

// Handle /user/username
func username(res http.ResponseWriter, req *http.Request) {
  if req.Method == "GET" {
    http.Error(res, http.StatusText(405), 405)
    return
  }

  var user models.Username
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&user); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //if it's a bad input, return 400
  if !validateGeneral(user.Username) {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  var err error
  if req.Method == "POST" {
    err = db.AddUsername(user.Username, bson.ObjectIdHex(id))
  }
  if req.Method == "PUT" {
    err = db.ChangeUsername(user.Username, bson.ObjectIdHex(id))
  }
  if req.Method == "DELETE" {
    err = db.RemoveUsername(user.Username, bson.ObjectIdHex(id))
  }

  if err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //else we're all good and give a no content code
  res.WriteHeader(http.StatusNoContent)
}

// Handle /user/userID
func friends(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "GET":
    getFriends(res, req)
  default:
    friend(res, req)
  }
}

// PUT Handle
func savedToggle(res http.ResponseWriter, req *http.Request) {
  var save models.Saved
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&save); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //actually add it -- ex. of the 'un- ification of language'
  var err error
  if req.Method == "POST" {
    err = db.SaveThread(bson.ObjectIdHex(save.Thread), bson.ObjectIdHex(id));
  }
  if req.Method == "PUT" {
    err = db.UnsaveThread(bson.ObjectIdHex(save.Thread), bson.ObjectIdHex(id));
  }

  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //return err if there's a problem
  res.WriteHeader(http.StatusNoContent)
}

// GET Handle getting saved threads -- kinda like group
func getSaved(res http.ResponseWriter, req *http.Request) {
  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //get saved -- GetSaved(userId bson.ObjectId)
  saved, err := db.GetSaved(bson.ObjectIdHex(id))
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //set up response struct
  response := &models.SendGroup{
    Threads: saved,
  }

  //send back json formatted threads
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  //send over data
  if err := json.NewEncoder(res).Encode(response); err != nil {
      panic(err)
  }
}

//GET Handle -- getting friend's names
func getFriends(res http.ResponseWriter, req *http.Request) {
  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //GetFriendsJoined(id bson.ObjectId) ([]string, error)
  friends, err := db.GetFriendsJoined(bson.ObjectIdHex(id))
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //send friends over
  sendFriends := &models.GetFriends{
    Friends: friends,
  }

  //send back json formatted threads
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)

  //send over data
  if err := json.NewEncoder(res).Encode(sendFriends); err != nil {
      panic(err)
  }
}

// POST Handle -- notifications are going to be a link with text
func friend(res http.ResponseWriter, req *http.Request) {

  //handle post with friends name
  var request models.Friend
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&request); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }
  var err error
  if req.Method == "POST" {
    //RequestFriend(user bson.ObjectId, username, friend string)
    err = db.RequestFriend(bson.ObjectIdHex(id), request.Username, request.Friend)
  }
  if req.Method == "PUT" {
    //AddFriend(friend, user bson.ObjectId)
    friendId := db.GetIdFromUsername(request.Friend)
    if friendId == "" {
      http.Error(res, http.StatusText(400), 400)
      return
    }
    err = db.AddFriend(bson.ObjectIdHex(id), bson.ObjectIdHex(friendId))
  }
  if req.Method == "DELETE" {
    //RemoveFriend(friend, user bson.ObjectId)
    friendId := db.GetIdFromUsername(request.Friend)
    if friendId == "" {
      http.Error(res, http.StatusText(400), 400)
      return
    }
    err = db.RemoveFriend(bson.ObjectIdHex(friendId), bson.ObjectIdHex(id))
  }

  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //if all goes well, we should end here
  res.WriteHeader(http.StatusNoContent)
}

// Handle GET
func notifications(res http.ResponseWriter, req *http.Request) {

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //get notifications
  notes, err := db.GetNotifications(bson.ObjectIdHex(id))

  //let ourselves know if something went wrong
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  sendNotes := &models.GetNotifications{
    Notifications: notes,
  }

  //send back json formatted threads
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)

  //send over data
  if err := json.NewEncoder(res).Encode(sendNotes); err != nil {
      panic(err)
  }
}


// Handle POST to /user/name
func name(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
  }

  //user _id in hex
  id := util.GetId(req)
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  var post models.Name
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&post); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //check if name is valid
  if !validateName(post.Name){
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get notifications
  err := db.AddName(bson.ObjectIdHex(id), post.Name)

  //let ourselves know if something went wrong
  if err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  res.WriteHeader(http.StatusNoContent)
}
