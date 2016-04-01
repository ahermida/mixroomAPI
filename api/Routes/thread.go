
/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package

   handles: creating a thread, creating a post, removing a post,
            removing a thread, changing a post
            -- could be written in this file

*/
package routes

import (
    "github.com/ahermida/dartboardAPI/api/Util"
    "github.com/ahermida/dartboardAPI/api/Models"
    "net/http"
    "encoding/json"
    "gopkg.in/mgo.v2/bson"
    "github.com/ahermida/dartboardAPI/api/DB"
)

// Routes with /api/ prefix
var ThreadMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST for getting thread
  ThreadMux.HandleFunc("/thread/", getThread)

  //POST for creating thread, DELETE removing thread, GET thread
  ThreadMux.HandleFunc("/thread/modify", thrd)

  //POST for creating a post, DELETE for removing it, PUT for changing it
  ThreadMux.HandleFunc("/thread/post", pst)
}

/*
   Route handlers for Thread  Routes
*/

// Handle /thread/
func thrd(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    createThread(res, req)
  case "DELETE":
    removeThread(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

// Handle /thread/post
func pst(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    createPost(res, req)
  case "PUT" :
    editPost(res, req)
  case "DELETE":
    removePost(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

//handle POST /thread/
func getThread(res http.ResponseWriter, req *http.Request) {
  var thread models.GetThread
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&thread); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //just get thread and check if we're authorized afterwards
  filledThread, threadErr := db.GetThread(bson.ObjectIdHex(thread.Thread))

  //fix up filledThread so it has S(tring)Id -- hex representation of _id
  filledThread.SId = filledThread.Id.Hex()

  //check for error in getting thread
  if threadErr != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  id := util.GetId(req)

  //check if we're a member of the group
  if !db.IsMember(filledThread.Group, id) {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //onward
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(res).Encode(filledThread); err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }
}

// handle POST /thread/modify
func createThread(res http.ResponseWriter, req *http.Request) {
  var reqBody models.NewThread
  if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //Should people anonymously create threads in private groups? yup, just with permissions
  id := util.GetId(req)

  //check if we're allowed to create this thread
  if !db.IsMember(reqBody.Group, id) {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //maybe we should make sure that author is the same as username
  if !validateGeneral(reqBody.Author) {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //if author doesn't have a valid Id...
  if id == "" {

    //make sure author gets their id
    id = bson.NewObjectId().Hex()
  }

  //now we're allowed so let's make it
  post := db.CreateHeadPost(reqBody.Author, reqBody.Body, reqBody.Content, bson.ObjectIdHex(id))
  if err := db.CreateThread(reqBody.Group, reqBody.Anonymous, post); err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //send is JSON to be sent
  send := &models.SendPost{
    Id: id, //send user's id
  }

  //onward
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if errSending := json.NewEncoder(res).Encode(send); errSending != nil {
    http.Error(res, http.StatusText(500), 500)
  }
}

//handle DELETE /thread/modify
func removeThread(res http.ResponseWriter, req *http.Request) {
  var reqBody models.RemoveThread
  if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //grab user id
  id := util.GetId(req)

  //if we're not signed in, then we're deleting an anon post
  if id == "" {
    id = reqBody.Id
  }

  err := db.DeleteThread(bson.ObjectIdHex(reqBody.Thread), bson.ObjectIdHex(id))
  if err != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }
  //else we're all good!
  res.WriteHeader(http.StatusNoContent)
}

//handle POST /thread/post
func createPost(res http.ResponseWriter, req *http.Request) {

  //make room for post
  var reqBody models.NewPost
  if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get user's id
  id := util.GetId(req)

  //check if we're posting anonymously
  if id == "" {

    //for anon posts
    id = bson.NewObjectId().Hex()
  }
  grp := db.GetThreadParent(reqBody.Thread)
  if !db.IsMember(grp, id) {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //convert responseTo list to ObjectIds
  responseTo := make([]bson.ObjectId, 0)
  for _, postId := range reqBody.ResponseTo {
    responseTo = append(responseTo, bson.ObjectIdHex(postId))
  }

  //maybe we should make sure that author is the same as username
  if !validateGeneral(reqBody.Author) {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get ID's
  usrId := bson.ObjectIdHex(id)
  thrdId := bson.ObjectIdHex(reqBody.Thread)

  //CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId, author, body, content string)
  _, err := db.CreatePost(usrId, thrdId, responseTo, reqBody.Author, reqBody.Body, reqBody.Content)
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //send is JSON to be sent
  send := &models.SendPost{
    Id: id, //send user's id -- so it can be removed
  }

  //onward
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if errSending := json.NewEncoder(res).Encode(send); errSending != nil {
    http.Error(res, http.StatusText(500), 500)
  }
}

//handle PUT /thread/post
func editPost(res http.ResponseWriter, req *http.Request) {
  var request models.EditPost
  if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //grab user id
  id := util.GetId(req)

  //if we're not signed in, then we're deleting an anon post
  if id == "" {
    id = request.Id
  }

  //run edit
  err := db.EditPost(request.Body, bson.ObjectIdHex(request.Post), bson.ObjectIdHex(id))
  if err != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //else we're clear to change the post
  res.WriteHeader(http.StatusNoContent)
}

//handle DELETE /thread/post
func removePost(res http.ResponseWriter, req *http.Request) {

  //Handle Post
  var request models.DeletePost

  //decode json into request
  if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //grab user id
  id := util.GetId(req)

  //if we're not signed in, then we're deleting an anon post
  if id == "" {
    id = request.Id
  }

  //run edit
  err := db.DeletePost(bson.ObjectIdHex(request.Post), bson.ObjectIdHex(id))

  //if err is not nil, respond with that
  if err != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //all good, so give no content status
  res.WriteHeader(http.StatusNoContent)
}
