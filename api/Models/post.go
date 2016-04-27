/*
  post.go -- models for post operations
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Post struct {
	Id          bson.ObjectId   `bson:"_id,omitempty" json:"-"`
  SId         string          `bson:"id,omitempty" json:"id"`
  Thread      bson.ObjectId   `bson:"thread,omitempty" json:"-"`
  Created     time.Time       `bson:"created" json:"created"`
	Author      string          `bson:"author" json:"author"`
  AuthorId    bson.ObjectId   `bson:"authorId,omitempty" json:"-"`
  Replies     []bson.ObjectId `bson:"replies" json:"replies"`
  ResponseTo  []bson.ObjectId `bson:"responseTo" json:"responseTo"`
  Content     string          `bson:"content,omitempty" json:"content"`
  ContentType string          `bson:"contentType,omitempty" json:"contentType"`
  Body        string          `bson:"body" json:"body"`
}

//pretty much a post, just with a group and size field attached
type PopularPost struct {
  Id          bson.ObjectId   `bson:"_id,omitempty" json:"-"`
  SId         string          `bson:"id,omitempty" json:"id"`
  Thread      bson.ObjectId   `bson:"thread,omitempty" json:"-"`
  ThreadId    string          `bson:"-" json:"thread"`
  Created     time.Time       `bson:"created" json:"created"`
	Author      string          `bson:"author" json:"author"`
  AuthorId    bson.ObjectId   `bson:"authorId,omitempty" json:"-"`
  Replies     []bson.ObjectId `bson:"replies" json:"replies"`
  ResponseTo  []bson.ObjectId `bson:"responseTo" json:"responseTo"`
  Content     string          `bson:"content,omitempty" json:"content"`
  ContentType string          `bson:"contentType,omitempty" json:"contentType"`
  Body        string          `bson:"body" json:"body"`
  Group       string          `bson:"-" json:"group"`
  Size        int             `bson:"-" json:"size"`

}

/*
  JSON Handling
*/

type NewPost struct {
  Body        string   `json:"body"`
  Content     string   `json:"content"`
  ContentType string   `json:"contentType"`
  Author      string   `json:"author"`
  ResponseTo  []string `json:"responseTo"`
  Anonymous   bool     `json:"anonymous"`
  Thread      string   `json:"thread"`
}

type PopularPosts struct {
  Posts []PopularPost `json:"posts"`
}

type EditPost struct {
  Body string `json:"body"`
  Post string `json:"post"`
  Id   string `json:"id,omitempty"`
}

type GetPop struct {
  Skip int `json:"skip"`
}

type DeletePost struct {
  Post string `json:"post"`
  Id   string `json:"id,omitempty"`
}

type SendPost struct {
  Id string `json:"id"`
  PostId string `json:"postId"`
}
