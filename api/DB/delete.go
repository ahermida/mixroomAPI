/*
   delete.go -- DB queries that delete objects
*/
package db

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/dartboardAPI/api/Models"
  "github.com/ahermida/dartboardAPI/api/Config"

)

/**
    GROUPS -------------------------------------------------------------
 */

//[DELETE] remove a thread
func DeleteGroup(group string, user bson.ObjectId) error {

  //call DB
  db := Connection.DB(config.DBName)

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

/**
    THREADS -------------------------------------------------------------
 */

//[DELETE] remove a thread
func DeleteThread(threadID, userID bson.ObjectId) error {

  //call DB
  db := Connection.DB(config.DBName)

  //post that will be populated by thread head
  var author struct {
    Id bson.ObjectId `bson:"author"`
    Mthread bson.ObjectId `bson:"mthread"`
    Group string `bson:"group"`
  }
  q := db.C("threads").Find(bson.M{"_id": threadID}).Select(bson.M{"author" : 1, "mthread": 1, "group": 1})
  //verify author by head post for threads
  if err := q.One(&author); err != nil {
    return err
  }

  //compare user ids -- if this passes, then go on to do great things
  if author.Id != userID {
    return errors.New("User can't delete the thread as he's not the author.")
  }

  //set mthread to its id
  mthread := author.Mthread

  //get friends of author by mongo _id
  friends, _ := GetFriendsById(userID)

  //check if we have friends -- not anonymous -- in this case, remove the mthreads from feeds
  if friends != nil {
    //get friends of author for post & deliver thread to their feeds
    for _, friend := range friends {

      //insert hex representation of MongoId
      err :=  db.C(friend.Hex()).Remove(bson.M{"_id": mthread})
      if err != nil {
        //return err if anything fails -- no reason why it should though
        return err
      }
    }
  }

  //remove thread from everything
  if err := db.C(author.Group).Remove(bson.M{"_id": mthread}); err != nil {
    return err
  }

  //remove thread
  if err := db.C("threads").Remove(bson.M{"_id": threadID}); err != nil {
    return err
  }

  //nothing wrong -- so we return nil error (all good)
  return nil
}

/**
    POSTS -------------------------------------------------------------
 */

//[DELETE] removes a post
func DeletePost(postID, userID bson.ObjectId) error {
  //call DB
  db := Connection.DB(config.DBName)

  //post that will be populated by thread head
  var post struct {
    AuthorId bson.ObjectId `bson:"authorId"`
    Thread bson.ObjectId `bson:"thread"`
  }

  //verify author by head post for threads
  if err := db.C("posts").Find(bson.M{"_id": postID}).Select(bson.M{"authorId": 1, "thread": 1}).One(&post); err != nil {
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

/**
    USER -------------------------------------------------------------
 */

//[DELETE] actually deletes a user (for dev use only!)
func RemoveUser(user bson.ObjectId) error {
  //call DB
  err := Connection.DB(config.DBName).C("users").Remove(bson.M{"_id": user})
  if err != nil {
    return err
  }
  errDropping := Connection.DB(config.DBName).C(user.Hex()).DropCollection()
  if errDropping != nil {
    return err
  }
  return nil
}
