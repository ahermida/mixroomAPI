/*
  user.go -- models for user operations
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
  Usernames     []string        `bson:"usernames"`
  Name          string          `bson:"name",omitempty`
	Email         string          `bson:"email"`
  Password      string          `bson:"password"`
  Friends       []bson.ObjectId `bson:"friends"`
  Notifications []bson.ObjectId `bson:"notifications"`
  Unread        int             `bson:"unread"`
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
  Username      string   `bson:"username" json:"username"`
  Email         string   `bson:"email" json:"email"`
  Unread        int      `bson:"unread" json:"unread"`
  Usernames     []string `bson:"usernames" json:"usernames"`
}

type GetUserFeed struct {
  Page int `json:"page"`
}

type GetFriends struct {
  Friends []string `json:"friends"`
}

type Friend struct {
  Username string `json:"username"`
  Friend string `json:"friend"`
}

type AuthUser struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

type AuthedUser struct {
  Token string `json:"token"`
}

type Recovery struct {
  Email string `json:"email"`
}

type ChangePW struct {
  Password string `json:"password"`
  NewPassword string `json:"newPassword"`
}

type Saved struct {
  Thread string `json:"thread"`
}

type Username struct {
  Username string `json:"username"`
}

type Name struct {
  Name string `json:"name"`
}
