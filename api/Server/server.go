/**
 *  server.go -- Starts Dartboard Server
 */
package server

import (
  "net/http"
  "github.com/ahermida/dartboardAPI/api/Routes" //package with route mux(s)
  "github.com/ahermida/dartboardAPI/api/Util"   //package with utility funcs
)

func Start(port string) {

  //handle group manipulation routes
  http.Handle("/group/", util.Log(routes.GroupMux))

  //handle auth routes
  http.Handle("/auth/", util.Log(routes.AuthMux))

  //handle thread manipulation routes
  http.Handle("/thread/", util.Log(routes.ThreadMux))

  //handle user routes
  http.Handle("/user/", util.Log(routes.UserMux))

  //Start Server
  go http.ListenAndServe(port, nil) // set listen port
}
