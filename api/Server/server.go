/**
 *  Server.go -- Starts Dartboard Server
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
var Server = http.NewServeMux()

func init() {
  //handle group manipulation routes
  Server.Handle("/group/", util.Log(routes.GroupMux))

  //handle auth routes
  Server.Handle("/auth/", util.Log(routes.AuthMux))

  //handle thread manipulation routes
  Server.Handle("/thread/", util.Log(routes.ThreadMux))

  //handle user routes
  Server.Handle("/user/", util.Log(routes.UserMux))

}

func Start(port string) {
  //Start Server
  go http.ListenAndServe(port, Server) // set listen port
}

func Close() {
  //Exit App -- but close connection to DB first
  db.Connection.Close()
  os.Exit(0)
}
