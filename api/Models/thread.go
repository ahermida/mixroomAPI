/*
  thread.go -- models for thread operations
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Thread struct {
	Id      bson.ObjectId   `bson:"_id,omitempty"`
  Created time.Time       `bson:"created"`
  Author  bson.ObjectId   `bson:"author"`
  Posts   []bson.ObjectId `bson:"posts"`
  Alive   bool            `bson:"alive"`
  Group   string          `bson:"group"`
  Mthread bson.ObjectId   `bson:"mthread"`
}

type ResThread struct {
  Id      bson.ObjectId   `json:"-" bson:"_id,omitempty"`
  SId     string          `json:"id,omitempty" bson:"-"`
  Created time.Time       `json:"created" bson:"created"`
  Posts   []Post          `json:"posts" bson:"posts"`
  Alive   bool            `json:"-" bson:"alive"`
  Group   string          `json:"group" bson:"group"`
  Mthread string          `json:"mthread" bson:"-"`
}

/*
  JSON Handling
*/

type GetThread struct {
  Thread string `json:"thread"`
}

type NewThread struct {
  Group       string `json:"group"`
  Body        string `json:"body"`
  Author      string `json:"author"`
  Content     string `json:"content"`
  ContentType string `json:"contentType"`
  Anonymous   bool   `json:"anonymous"`
}

type RemoveThread struct {
  Thread string `json:"thread"`
  Id     string `json:"id,omitempty"`
}

type ThreadLen struct {
  Length int    `json:"size,omitempty"`
}
