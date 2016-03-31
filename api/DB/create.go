/*
   create.go -- DB queries that create objects
*/
package db

import (
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/dartboardAPI/api/Models"
)

//[CREATE] creates user in the Database with given email, username, and hashed password
func CreateUser(email, username, password string) (string, error) {
  //connect to appropriate DB
  db := Connection.DB("dartboard")

  //create user struct
  usr := &models.User{
    Id:            bson.NewObjectId(),
    Created:       bson.Now(),
    Username:      username,
    Usernames:     []string{username},
  	Email:         email,
    Password:      password,
    Friends:       make([]bson.ObjectId, 0),
    Notifications: make([]bson.ObjectId, 0),
    Requests:      make([]bson.ObjectId, 0),
    Activated:     false,
    Suspended:     false,
    Saved:         make([]bson.ObjectId, 0),
  }

  //create content feed for given user
  if err := CreateGroup(usr.Id.Hex(), usr.Id, true); err != nil {
    return "", err
  }

  //insert user
  if err := db.C("users").Insert(usr); err != nil {
    return "", err
  }

  //all went well, so return nil err
  return usr.Id.Hex(), nil
}

//[CREATE] creates a thread from a struct of a given JSON -- also creates an mthread for each location
func CreateThread(group string, anonymous bool, post *models.Post) error {

  //connect to DB
  db := Connection.DB("dartboard")

  //create thread, then attach to mthread
  thread := &models.Thread{
    Id: bson.NewObjectId(),
    Created: bson.Now(),
    Posts: []bson.ObjectId{post.Id},
    Alive: true,
    Group: group,
  }

  //create mthread -- for view
  mthread := &models.Mthread{
    Id: bson.NewObjectId(),
    Created: thread.Created,
    Thread: thread.Id,
    Head: post,
    Group: group,
  }

  //set thread in post
  post.Thread = thread.Id

  //save post
  if err := db.C("posts").Insert(post); err != nil {
    //insert it, shouldn't result in error
    return err
  }

  //add thread to threads, this is where we'll query for specific thread views
  if err := db.C("threads").Insert(thread); err != nil {
    //if there's an error on inserting, return the error
    return err
  }

  //add mthread to group, this is where we'll query for mthreads
  if err := db.C(group).Insert(mthread); err != nil {
    //if there's an error on inserting, return the error
    return err
  }

  //if anonymous poster, don't share apart from thread
  if !anonymous {

    //get friends of author by mongo _id
    friends, err := GetFriends(post.Author)

    //check if something went wrong getting friends
    if err != nil {

      //if there's a problem getting friends, let's treat it as not a problem
      return nil
    }

    //get friends of author for post & deliver thread to their feeds
    for _, friend := range friends {

      //insert hex representation of MongoId
      err :=  db.C(friend.Hex()).Insert(mthread)
      if err != nil {
        //return err if anything fails -- no reason why it should though
        return err
      }
    }
  }

  //return nil error -- nothing went wrong
  return nil
}

//[CREATE] creates a head post for a given thread -- leave thread id unset
func CreateHeadPost(author, body, content string, authorId bson.ObjectId) *models.Post {

  //make post
  post := &models.Post{
    Id: bson.NewObjectId(),
    Created: bson.Now(),
    Author: author,
    AuthorId: authorId,
    Replies: make([]bson.ObjectId,0),
    ResponseTo: make([]bson.ObjectId,0),
    Content: content,
    Body: body,
  }

  //return post
  return post
}

//[CREATE] creates a post for a given thread
func CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId, author, body, content string) (*models.Post, error) {
  //connect to DB
  db := Connection.DB("dartboard")

  //add post to DB, notify people that it exists, add replies to responses
  post := &models.Post{
    Id: bson.NewObjectId(),
    Created: bson.Now(),
    Author: author,
    AuthorId: authorId,
    Replies: make([]bson.ObjectId,0),
    ResponseTo: responseTo,
    Content: content,
    Body: body,
    Thread: thread,
  }

  //push the id of the post to each we're responding to
  postChange := bson.M{"$push": bson.M{"replies": post.Id}}

  //add reply to each of the posts that it was to
  for _, person := range responseTo {

    //insert hex representation of MongoId
    if err := db.C("posts").Update(bson.M{"_id": person}, postChange); err != nil {

      //insert it, shouldn't result in error
      return nil, err
    }
  }

  //insert post into db
  if err := db.C("posts").Insert(post); err != nil {
    //insert it, shouldn't result in error
    return nil, err
  }

  //thread changes
  threadChange := bson.M{"$push": bson.M{"posts": post.Id}}

  //insert post into thread
  if err := db.C("threads").Update(bson.M{"_id": thread}, threadChange); err != nil {

    //insert it, shouldn't result in error
    return nil, err
  }

  //return post created
  return post, nil

}

//[CREATE] creates a group, rather, reserves a namespace for a group
func CreateGroup(group string, user bson.ObjectId, private bool) error {
  //get proper db
  db := Connection.DB("dartboard")

  //create group with index options (anonymity allowed, thread lifetimes, etc)
  g := &models.Group{
    Created: bson.Now(),
    Name: group,
    Author: user,
    Admins: make([]bson.ObjectId,0),
    Private: private,
  }

  //insert group into DB
  err := db.C("groups").Insert(g)

  //err should only happen on name duplicates
  if err != nil {
    return err
  }

  //nothing went wrong, so return nil
  return nil
}

//[Misc] Checks if group exists
func GroupExists(group string) bool {
  //get DB
  db := Connection.DB("dartboard")

  //check DB
  count, err := db.C("groups").Find(bson.M{"name": group}).Count()

  //something went wrong, so play it safe and say group doesn't exist
  if err != nil {
    return true
  }

  //if there are any groups with that name, it should be 1
  return count == 1
}
