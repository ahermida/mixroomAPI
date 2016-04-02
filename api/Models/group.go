/*
  group.go -- each group gets its own collection, this is just how we resolve permissions
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Group struct {
  Id      bson.ObjectId   `bson:"_id,omitempty"`
  Created time.Time       `bson:"created"`
  Name    string          `bson:"name"`
  Author  bson.ObjectId   `bson:"author"`
  Admins  []bson.ObjectId `bson:"admins"`
  Private bool            `bson:"private"`
}

//For the Group View
type Mthread struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"-"`
  SId      string        `bson:"id" json:"id"`
  Created  time.Time     `bson:"created" json:"created"`
  Thread   bson.ObjectId `bson:"thread" json:"-"`
  ThreadId string        `bson:"threadId" json:"thread"`
  Head     *Post         `bson:"head" json:"head"`
  Group    string        `bson:"group" json:"group"`
}


/*
  JSON Handling
*/

type GetGroup struct {
  Page  int    `json:"page"`
  Group string `json:"group"`
}

type CreateGroup struct {
  Private   bool   `json:"anonymous"`
  Group     string `json:"group"`
}

type SendGroup struct {
  Threads []Mthread `json:"threads"`
}

type DeleteGroup struct {
  Group string `json:"group"`
}

type GroupAdmin struct {
  User  string `json:"user"`
  Group string `json:"group"`
}
