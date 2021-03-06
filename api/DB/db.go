/*
   db.go -- Database Connection for Application
   (moved everything into CRUD letter files for simplicity)
   >Keeping all DB access functions in this package to keep things clean
*/
package db

import (
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/sudopostAPI/api/Config"
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

  //Create Indices
  ensureUserIndex()
  ensureGroupIndex()
  ensureMthreadTextIndex()
  ensureUserTextIndex()

  //reserve group names
  reserveNamespaces()
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

func ensureUserTextIndex() {
  index := mgo.Index{
    Key: []string{"name", "usernames"},
  }

  //ensure indices are unique
  err := Connection.DB(config.DBName).C("users").EnsureIndex(index)

  //index needs to be unique
  if err != nil {
    log.Panic(err)
  }
}

func ensureMthreadTextIndex() {
  index := mgo.Index{
    Key: []string{"$text:head.body"},
  }

  //ensure indices are unique
  err := Connection.DB(config.DBName).C("mthreads").EnsureIndex(index)

  //index needs to be unique
  if err != nil {
    log.Panic(err)
  }
}


//Makes Sure that groups cannot be duplicates
func ensureGroupIndex() {
  index := mgo.Index{
    Key:        []string{"$text:name"},
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
  //private groups
  _reserved := []string{"users", "threads", "posts", "groups", "mthreads"}
  for _, namespace := range _reserved {
    CreateGroup(namespace, bson.NewObjectId(), true)
  }

  reserved := []string{"/", "/404", "/cs/", "/music/", "/vid/", "/bored/", "/random/"}
  //create "/" group
  for _, namespace := range reserved {
    CreateGroup(namespace, bson.NewObjectId(), false)
  }
}
