/*
   db.go -- Database Connection for Application
   (moved everything into CRUD letter files for simplicity)
   >Keeping all DB access functions in this package to keep things clean
*/
package db

import (
  "github.com/ahermida/dartboardAPI/api/Config"
  "log"
  "gopkg.in/mgo.v2"
)

var Connection *mgo.Session

//initialize DB
func init() {
  var err error
  Connection, err = mgo.Dial(config.DB)
  if err != nil {
    log.Panic(err)
  }
}
