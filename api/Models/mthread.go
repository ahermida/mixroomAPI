/*
  mthread.go -- models for meta-thread operations -- each group gets its own collection
  meta-thread model (meta threads):
    {
      "_id" : "12312313",
      "created" : timestamp,
      "thread" : 1231245, // mongo _id representing thread
      "head" : post // -- embedded post
    }
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

//For the Group View
type Mthread struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
  Created time.Time     `bson:"created"`
  Thread  bson.ObjectId `bson:"thread"`
  Head    *Post         `bson:"head"`
  Group   string        `bson:"group"`
}
