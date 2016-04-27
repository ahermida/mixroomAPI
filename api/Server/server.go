/**
 *  Server.go -- Starts Dartboard Server
 */
package server

import (
  "net/http"
  "github.com/ahermida/sudopostAPI/api/Routes"
  "github.com/ahermida/sudopostAPI/api/Util"
  "os"
  "github.com/ahermida/sudopostAPI/api/DB"
)

//Use this to handle routes
var Server = http.NewServeMux()

func init() {
  //handle group manipulation routes
  Server.Handle("/group/", routes.GroupMux)

  //handle auth routes
  Server.Handle("/auth/", routes.AuthMux)

  //handle thread manipulation routes
  Server.Handle("/thread/", routes.ThreadMux)

  //handle user routes
  Server.Handle("/user/", routes.UserMux)
}

func Start(port string) {
  //Start Server
  go http.ListenAndServe(port, util.Log(Server)) // set listen port
}

func Close() {
  //Exit App -- but close connection to DB first
  db.Connection.Close()
  os.Exit(0)
}
