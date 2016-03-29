/*
  user.go -- models for user operations
  user model:
    {
      "_id" : "12312313",
      "created" : timestamp,
      "username" : "example",
      "email" : "example@gmail.com",
      "password": "thisishashed",
      "friends" : [20130123, 12301031, 1231021], //by mongo _id <-- for users
      "notifications": [20123023, 12030132, 10230044], //by mongo _id <-- for posts
      "requests": [123012031], //by mongo _id <-- for users
      "activated": true,
      "suspended" : false,
      "saved" : [1023121, 10230103, 12312301], // by mongo _id for threads
      "description" : "I'm the dude.",
      "feed" : "username",//group by username
      "likes": [123123, 12312123, 1231231321] //post ids
    }
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId   `bson:"_id,omitempty"`
  Created       time.Time       `bson:"created"`
	Username      string          `bson:"username"`
  Usernames     []string        `bson:"usernames,omitempty"`
	Email         string          `bson:"email"`
  Password      string          `bson:"password"`
  Friends       []bson.ObjectId `bson:"friends"`
  Notifications []bson.ObjectId `bson:"notifications"`
  Requests      []bson.ObjectId `bson:"requests"`
  Activated     bool            `bson:"activated"`
  Suspended     bool            `bson:"suspended"`
  Saved         []bson.ObjectId `bson:"saved"`
  Description   string          `bson:"description,omitempty"`
}
