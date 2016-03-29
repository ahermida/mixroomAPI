/*
   read.go -- DB queries that read objects
*/
package db

import (
  //"github.com/ahermida/dartboardAPI/api/Config"
//  "gopkg.in/mgo.v2"
  "errors"
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/dartboardAPI/api/Models"
)

//[READ] gets group info -- meta info about group from group collection
func CheckGroup(group string) (*models.Group, error) {
  db := Connection.DB("dartboard")
  var g models.Group
  err := db.C("groups").Find(bson.M{"name": group}).One(&g)
  if err != nil {
    return nil, err
  }
  //bring back group metadata
  return &g, nil
}

//[READ] gets available threads for a given group on a particular page
func GetGroup(group string, page int) ([]models.Mthread, error) {

  //check if group exists
  if !GroupExists(group) {

    //if it doesn't, return an error
    return nil, errors.New("Can't get a group that doens't exist.")
  }
  //get DB
  db := Connection.DB("dartboard")

  //Sort by Timestamp --> Get The Range from (page * 30) -- 30 items
  pipeline := []bson.M{bson.M{"$sort": bson.M{"created": -1 }},
                        bson.M{"$limit": 30},
                        bson.M{"$skip": page * 30}}

  //set up Pipe to actually run query
  pipe := db.C(group).Pipe(pipeline)

  //slice of threads that will be populated by query
  var threads []models.Mthread

  //run it
  if err:= pipe.All(&threads); err != nil {

    //if something is wrong, return err
    return nil, err
  }

  //else return the threads, and a nil error
  return threads, nil
}

//[READ] gets all posts for a given thread
func GetThread(threadID bson.ObjectId) (*models.ResThread, error) {

  //call DB
  db := Connection.DB("dartboard")

  //thread model
  var thread models.Thread

  //get thread
  if err := db.C("threads").FindId(threadID).One(&thread); err != nil {
    return nil, err
  }

  //if thread is dead, don't return it.
  if !thread.Alive {
    return nil, errors.New("Couldn't get thread, it's dead")
  }

  //make new thread to be resolved
  resThread := &models.ResThread{
    Id: thread.Id,
    Created: thread.Created,
    Posts: make([]models.Post,0),
    Alive: thread.Alive,
    Group: thread.Group,
  }

  //add reply to each of the posts that it was to
  for _, postId := range thread.Posts {

    //make room for post
    var post models.Post

    //get post by ID
    if err := db.C("posts").Find(bson.M{"_id": postId}).One(&post); err != nil {
      //insert it, shouldn't result in error
      return nil, err
    }

    //merge slices
    resThread.Posts = append(resThread.Posts, post)
  }

  //return thread and nil error (nothing went wrong)
  return resThread, nil
}

//[READ] gets user data -- posts, watchlist, thread likes, friends
func GetUser(userID string) {
  //call DB
}

//[READ] gets user's friends -- resolving "joins"
func GetFriends(author string) ([]bson.ObjectId, error) {

  //grab proper db
  db := Connection.DB("dartboard").C("users")

  //anonymous struct just so type doesn't fail us on unmarshalling to []bson.ObjectId
  var people struct{
    Friends []bson.ObjectId
  }

  //get friends
  err := db.Find(bson.M{"username": author}).Select(bson.M{"friends" : 1, "_id" : 0}).One(&people)
  if err != nil {
    return nil, err
  }

  //all is well, so return nil error and friends
  return people.Friends, nil
}
