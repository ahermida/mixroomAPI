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
  Id      bson.ObjectId   `bson:"_id,omitempty"`
  Created time.Time       `bson:"created"`
  Posts   []Post          `bson:"posts"`
  Alive   bool            `bson:"alive"`
  Group   string          `bson:"group"`
}
