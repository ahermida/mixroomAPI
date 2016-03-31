
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
  "fmt"
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

  //get token from header
  userToken := req.Header.Get("access_token")

  var id string

  //if user has no token
  if userToken == "" {

    //fail
    id = ""
  } else {

    //check token that we got
    token, err := util.CheckToken(userToken)

    //if nothing went wrong
    if err != nil || !token.Valid {

      //fail
      id = ""
    } else {

      //success
      id = token.Claims["id"].(string)
    }
  }

  //check if we're a member of the group
  if !db.IsMember(filledThread.Group, id) {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //onward
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(res).Encode(filledThread); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }
}

/* handle POST /thread/modify
  --requires new post
    -group name -- just check id IsMember
    -author 'name' -- stomp this out
    -body
    -content
    -anonymous y / n
*/
func createThread(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Admin Test Passed!")
}

//handle DELETE /thread/modify
func removeThread(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Admin Test Passed!")
}

//handle POST /thread/post
func createPost(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Admin Test Passed!")
}

//handle PUT /thread/post
func editPost(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Admin Test Passed!")
}

//handle DELETE /thread/post
func removePost(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Admin Test Passed!")
}
