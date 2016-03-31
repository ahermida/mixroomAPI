/*
  post.go -- models for post operations
  post model:
    {
      "_id" : "12312313",
      "created" : timestamp,
      "author" : "username",
      "thread" : 123132131, //thread id (group that it originally belongs to)
      "replies" : [123102313, 123902103, 141281445], //mongo _id of posts responding to this one
      "responseTo": [12301233, 12302013], //responseTo are the posts this is a response to
      "content": "src/sadsdas.mp4", //content takes precedence over link as main
      "body": "blablabla this is text"
    }
*/
package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Post struct {
	Id         bson.ObjectId   `bson:"_id,omitempty"`
  Thread     bson.ObjectId   `bson:"thread"`
  Created    time.Time       `bson:"created"`
	Author     string          `bson:"author"`
  AuthorId   bson.ObjectId   `bson:"authorId,omitempty"`
  Replies    []bson.ObjectId `bson:"replies"`
  ResponseTo []bson.ObjectId `bson:"responseTo"`
  Content    string          `bson:"content,omitempty"`
  Body       string          `bson:"body"`
}

/*
  JSON Handling
*/

type NewPost struct {
  Body      string   `json:"body"`
  Username  string   `json:"username"`
  Content   []string `json:"content"`
  Anonymous bool     `json:"anonymous"`
  Thread    string   `json:"thread"`
}
