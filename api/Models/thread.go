/*
  thread.go -- models for thread operations
  thread model:
    {
      "_id" : "12312313",
      "created" : timestamp,
      "head" : 123131321, //post by mongo _id
      "posts": [123121312, 12312131, 12313213] //mongo id's (resolved aferwards)
      "group" : 123132131, //group id (group that it originally belongs to)
      "alive" : true
    }
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Thread struct {
	Id      bson.ObjectId   `bson:"_id,omitempty"`
  Created time.Time       `bson:"created"`
  Posts   []bson.ObjectId `bson:"posts"`
  Alive   bool            `bson:"alive"`
  Group   string          `bson:"group"`
}

type ResThread struct {
  Id      bson.ObjectId   `json:"-" bson:"_id,omitempty"`
  SId     string          `json:"id,omitempty" bson:"-"`
  Created time.Time       `json:"created" bson:"created"`
  Posts   []Post          `json:"posts" bson:"posts"`
  Alive   bool            `json:"-" bson:"alive"`
  Group   string          `json:"group" bson:"group"`
}

/*
  JSON Handling
*/

type GetThread struct {
  Thread string `json:"page"`
}

type NewThread struct {
  Group     string `json:"group"`
  Body      string `json:"body"`
  Username  string `json:"username"`
  Content   string `json:"content"`
  Anonymous bool   `json:"anonymous"`
}

type RemoveThread struct {
  Thread string `json:"thread"`
}
