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
}
