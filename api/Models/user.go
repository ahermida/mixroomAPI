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

/*
  JSON Handling
*/

type CreateUser struct {
  Username string `json:"username"`
  Email    string `json:"email"`
  Password string `json:"password"`
}

type GetUser struct {
  Username      string `json:"username"`
  Email         string `json:"email"`
  Notifications int    `json:"notifications"`
}

type GetUserFeed struct {
  Page int `json:"page"`
}
type RemoveFriend struct {
  friend string `json:"username"`
}

type AuthUser struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

type AuthedUser struct {
  Token string `json:"token"`
}

type ChangePW struct {
  Password string `json:"password"`
  NewPassword string `json:"newPassword"`
}
