/*
   delete.go -- DB queries that delete objects
*/
package db

import (
  //"github.com/ahermida/dartboardAPI/api/Config"
//  "gopkg.in/mgo.v2"
  "errors"
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/dartboardAPI/api/Models"
)

//[DELETE] remove a thread
func DeleteGroup(group string, user bson.ObjectId) error {

  //call DB
  db := Connection.DB("dartboard")

  //make struct for group
  var grp models.Group

  //verify author by head post for threads
  if err := db.C("groups").Find(bson.M{"name": group}).One(&grp); err != nil {
    return err
  }

  //compare user ids
  if grp.Author != user {
    return errors.New("User can't delete the group as he's not the author.")
  }

  //remove thread
  if err := db.C("groups").Remove(bson.M{"name": group}); err != nil {
    return err
  }

  //remove collection named after group
  if err := db.C(group).DropCollection(); err != nil {
    return err
  }

  //nothing wrong -- so we return nil error (all good)
  return nil
}

//[DELETE] remove a thread
func DeleteThread(threadID, userID bson.ObjectId) error {

  //call DB
  db := Connection.DB("dartboard")

  //post that will be populated by thread head
  var post models.Post

  //verify author by head post for threads
  if err := db.C("threads").Find(bson.M{"_id": threadID}).Select(bson.M{"posts.0" : 1}).One(&post); err != nil {
    return err
  }

  //compare user ids
  if post.AuthorId != userID {
    return errors.New("User can't delete the thread as he's not the author.")
  }

  //remove thread
  if err := db.C("threads").Remove(bson.M{"_id": threadID}); err != nil {
    return err
  }

  //nothing wrong -- so we return nil error (all good)
  return nil
}

//[DELETE] removes a post
func DeletePost(postID, userID bson.ObjectId) error {
  //call DB
  db := Connection.DB("dartboard")

  //post that will be populated by thread head
  var post struct {
    AuthorId bson.ObjectId `bson:"authorId"`
  }

  //verify author by head post for threads
  if err := db.C("posts").Find(bson.M{"_id": postID}).Select(bson.M{"authorId": 1}).One(&post); err != nil {
    return err
  }

  //compare user ids
  if post.AuthorId != userID {
    return errors.New("User can't delete the post as he's not the author.")
  }

  //remove post
  if err := db.C("posts").Remove(bson.M{"_id": postID}); err != nil {
    return err
  }

  //remove post from thread
  err := db.C("threads").UpdateId(post.Thread, bson.M{"$pull": bson.M{"posts": postID}})
  if err != nil {
    return err
  }

  //nothing wrong -- so we return nil error (all good)
  return nil
}

//[DELETE] actually deletes a user (for dev use only!)
func RemoveUser(user bson.ObjectId) error {
  //call DB
  err := Connection.DB("dartboard").C("users").Remove(bson.M{"_id": user})
  if err != nil {
    return err
  }
  return nil
}
