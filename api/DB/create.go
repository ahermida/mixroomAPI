/*
   create.go -- DB queries that create objects
*/
package db

import (
  "fmt"
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/sudopostAPI/api/Models"
  "github.com/ahermida/sudopostAPI/api/Config"

)
/**
    GROUPS -------------------------------------------------------------
 */

//[CREATE] creates a group, rather, reserves a namespace for a group
func CreateGroup(group string, user bson.ObjectId, private bool) error {
  //get proper db
  db := Connection.DB(config.DBName)

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
  db := Connection.DB(config.DBName)

  //check DB
  count, err := db.C("groups").Find(bson.M{"name": group}).Count()

  //something went wrong, so play it safe and say group doesn't exist
  if err != nil {
    return true
  }

  //if there are any groups with that name, it should be 1
  return count == 1
}

/**
    THREADS -------------------------------------------------------------
 */

//[CREATE] creates a thread from a struct of a given JSON -- also creates an mthread for each location
func CreateThread(group string, anonymous bool, post *models.Post) error {

  //connect to DB
  db := Connection.DB(config.DBName)

  //create thread, then attach to mthread
  thread := &models.Thread{
    Id: bson.NewObjectId(),
    Created: bson.Now(),
    Posts: []bson.ObjectId{post.Id},
    Alive: true,
    Group: group,
    Author: post.AuthorId,
  }

  //id for mthread
  mId := bson.NewObjectId()

  //set reference to mthread (for easy saves)
  thread.Mthread = mId

  //create mthread -- for view
  mthread := &models.Mthread{
    Id: mId,
    SId: mId.Hex(),
    Created: thread.Created,
    Thread: thread.Id,
    ThreadId: thread.Id.Hex(),
    Head: post,
    Group: group,
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

  //insert mthread to their own collection, this will allow for easy browsing of old saved threads
  if err := db.C("mthreads").Insert(mthread); err != nil {
    //if there's an error on inserting, return the error
    return err
  }

  //if anonymous poster, don't share apart from thread
  if !anonymous {

    //get friends of author
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

  //set thread in post
  post.Thread = thread.Id

  //save post
  if err := db.C("posts").Insert(post); err != nil {
    //insert it, shouldn't result in error
    return err
  }

  //return nil error -- nothing went wrong
  return nil
}

//[CREATE] creates a head post for a given thread -- leave thread id unset
func CreateHeadPost(author, body, content, contentType string, authorId bson.ObjectId) *models.Post {

  //post's id
  id := bson.NewObjectId()

  //make post
  post := &models.Post{
    Id: id,
    SId: id.Hex(),
    Created: bson.Now(),
    Author: author,
    AuthorId: authorId,
    Replies: make([]bson.ObjectId,0),
    ResponseTo: make([]bson.ObjectId,0),
    Content: content,
    ContentType: contentType,
    Body: body,
  }

  //return post
  return post
}

/**
    POSTS -------------------------------------------------------------
 */

//[CREATE] creates a post for a given thread
func CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId,
                author, body, content, contentType string) (*models.Post, error) {
  //connect to DB
  db := Connection.DB(config.DBName)

  id :=  bson.NewObjectId()

  //add post to DB, notify people that it exists, add replies to responses
  post := &models.Post{
    Id: id,
    SId: id.Hex(),
    Created: bson.Now(),
    Author: author,
    AuthorId: authorId,
    Replies: make([]bson.ObjectId,0),
    ResponseTo: responseTo,
    Content: content,
    ContentType: contentType,
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

/**
    User & User Ops -------------------------------------------------------------
 */

//[CREATE] creates user in the Database with given email, username, and hashed password
func CreateUser(email, username, password string) (string, error) {
  //connect to appropriate DB
  db := Connection.DB(config.DBName)

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
    Unread:        0,
    Requests:      make([]bson.ObjectId, 0),
    Activated:     false,
    Suspended:     false,
    Saved:         make([]bson.ObjectId, 0),
  }

  //insert user -- before feed so we don't create collections needlessly
  if err := db.C("users").Insert(usr); err != nil {
    return "", err
  }

  //create content feed for given user
  if err := CreateGroup(usr.Id.Hex(), usr.Id, true); err != nil {
    return "", err
  }

  //all went well, so return nil err
  return usr.Id.Hex(), nil
}

//[CREATE] creates a notification for the Id
func CreateNotification(id bson.ObjectId, link, text string) error {
  db := Connection.DB(config.DBName)
  note := &models.Notification{
    Id: bson.NewObjectId(),
    Recipient: id,
    Link: link,
    Text: text,
  }
  change := bson.M{"$addToSet": bson.M{"notifications": note.Id}}
  if err := db.C("users").Update(bson.M{"_id": id}, change); err != nil {
    return err
  }
  incUnread := bson.M{"$inc": bson.M{"unread": 1}}
  if err := db.C("users").Update(bson.M{"_id": id}, incUnread); err != nil {
    return err
  }
  if err := db.C("notifications").Insert(note); err != nil {
    return err
  }

  //everything went as planned
  return nil
}

func RequestFriend(user bson.ObjectId, username, friend string) error {
  db := Connection.DB(config.DBName)

  var newFriend struct {
    Id bson.ObjectId `bson:"_id"`
  }
  q := db.C("users").Find(bson.M{"username": friend})
  if err := q.Select(bson.M{"_id": 1}).One(&newFriend); err != nil {
    return err
  }

  //setup change
  change := bson.M{"$addToSet": bson.M{"requests": user}}

  //execute change
  if err := db.C("users").Update(bson.M{"_id": newFriend.Id}, change); err != nil {
    return err
  }

  //note that we made -- USE OUR OWN MARKUP HERE so we can style it appropriately
  note := fmt.Sprintf("You have a new friend request from %s", username)

  //create link that they can click on to view our profile
  link := fmt.Sprintf("localhost:8080/user/%s", username)

  //add notification to new friend
  if err := CreateNotification(newFriend.Id, link, note); err != nil {
    return err
  }

  //everything went as planned
  return nil
}
