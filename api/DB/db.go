/*
   db.go -- Database Connection for Application
   (moved everything into CRUD letter files for simplicity)
   >Keeping all DB access functions in this package to keep things clean
*/
package db

import (
  "gopkg.in/mgo.v2/bson"
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
  Connection.SetMode(mgo.Monotonic, true)
  ensureUserIndex()
  ensureGroupIndex()
  //reserveNamespaces()
}

//Makes Sure that Users cannot be duplicates
func ensureUserIndex() {
  index := mgo.Index{
    Key:        []string{"email", "usernames"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }

  //ensure indices are unique
  err := Connection.DB(config.DBName).C("users").EnsureIndex(index)

  //index needs to be unique
  if err != nil {
    log.Panic(err)
  }
}

//Makes Sure that groups cannot be duplicates
func ensureGroupIndex() {
  index := mgo.Index{
    Key:        []string{"name"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }

  //ensure indices are unique
  err := Connection.DB(config.DBName).C("groups").EnsureIndex(index)

  //index needs to be unique
  if err != nil {
    log.Panic(err)
  }
}

//make sure that no groups can be made for "users", "threads", "posts", "groups", or "mthreads"
func reserveNamespaces() {
  reserved := []string{"users", "threads", "posts", "groups", "mthreads"}
  for _, namespace := range reserved {
    CreateGroup(namespace, bson.NewObjectId(), true)
  }
}
