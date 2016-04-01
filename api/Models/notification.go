/*
  notification.go -- each group gets its own collection, this is just how we resolve permissions
*/
package models

import (
  "gopkg.in/mgo.v2/bson"
)

type Notification struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"-"`
  SId       string        `bson:"-" json:"id"`
  Recipient bson.ObjectId `bson:"recipient" json:"recipient"`
  Link      string        `bson:"link" json:"link"`
  Text      string        `bson:"text" json:"text"`
}

/*
  JSON Handling
*/
type GetNotifications struct {
	Notifications []Notification `json:"notifications"`
}
