/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package

   handles: creating a thread, creating a post, removing a post,
            removing a thread, changing a post
            -- could be written in this file

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
var ThreadMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST for creating thread, DELETE removing thread, GET thread
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
  fmt.Fprintf(res, "Admin Test Passed!")
}

//handle POST /thread/modify
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
