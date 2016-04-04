/*
  post.go -- models for post operations
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Post struct {
	Id         bson.ObjectId   `bson:"_id,omitempty" json:"id"`
  Thread     bson.ObjectId   `bson:"thread,omitempty" json:"-"`
  Created    time.Time       `bson:"created" json:"created"`
	Author     string          `bson:"author" json:"author"`
  AuthorId   bson.ObjectId   `bson:"authorId,omitempty" json:"-"`
  Replies    []bson.ObjectId `bson:"replies" json:"replies"`
  ResponseTo []bson.ObjectId `bson:"responseTo" json:"responseTo"`
  Content    string          `bson:"content,omitempty" json:"content"`
  Body       string          `bson:"body" json:"body"`
}

/*
  JSON Handling
*/

type NewPost struct {
  Body       string   `json:"body"`
  Content    string   `json:"content"`
  Author     string   `json:"author"`
  ResponseTo []string `json:"responseTo"`
  Anonymous  bool     `json:"anonymous"`
  Thread     string   `json:"thread"`
}

type EditPost struct {
  Body string `json:"body"`
  Post string `json:"post"`
  Id   string `json:"id,omitempty"`
}

type DeletePost struct {
  Post string `json:"post"`
  Id   string `json:"id,omitempty"`
}

type SendPost struct {
  Id string `json:"id"`
}
