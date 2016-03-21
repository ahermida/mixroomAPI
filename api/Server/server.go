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

  //handle api routes -- Protected
  http.Handle("/api/", util.Log(util.Protect(routes.ApiMux)))

  //handle user routes -- Unprotected
  http.Handle("/user/", util.Log(routes.UserMux))

  //Start Server
  go http.ListenAndServe(port, nil) // set listen port
}
