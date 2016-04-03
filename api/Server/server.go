/**
 *  server.go -- Starts Dartboard Server
 */
package server

import (
  "net/http"
  "github.com/ahermida/dartboardAPI/api/Routes"
  "github.com/ahermida/dartboardAPI/api/Util"
  "os"
  "github.com/ahermida/dartboardAPI/api/DB"
)

//Use this to handle routes
var server = http.NewServeMux()

func init() {
  //handle group manipulation routes
  server.Handle("/group/", util.Log(routes.GroupMux))

  //handle auth routes
  server.Handle("/auth/", util.Log(routes.AuthMux))

  //handle thread manipulation routes
  server.Handle("/thread/", util.Log(routes.ThreadMux))

  //handle user routes
  server.Handle("/user/", util.Log(routes.UserMux))

}

func Start(port string) {
  //Start Server
  go http.ListenAndServe(port, server) // set listen port
}

func Close() {
  //Exit App -- but close connection to DB first
  db.Connection.Close()
  os.Exit(0)
}
