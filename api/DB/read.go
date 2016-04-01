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

//[READ] checks if user is a member of a group
func IsMember(group string, user string) bool {
  db := Connection.DB("dartboard")
  var g models.Group
  err := db.C("groups").Find(bson.M{"name": group}).One(&g)

  //check if something went wrong finding the group
  if err != nil {

    //if something goes wrong, return false to be safe
    return false
  }
  if !g.Private {
    return true
  }
  if user == g.Author.Hex() {
    return true
  }
  //bring back group metadata
  for _, member := range g.Admins {
    if member.Hex() == user {
      return true
    }
  }
  return false
}

//[READ] returns group that the given thead (hex string) belongs to
func GetThreadParent(thread string) string {
  db := Connection.DB("dartboard")
  var thrd struct {
    Group string
  }
  err := db.C("threads").Find(bson.M{"_id": bson.ObjectIdHex(thread)}).Select(bson.M{"group": 1}).One(&thrd)
  if err != nil {
    return ""
  }

  return thrd.Group
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

//[READ]
func GetSaved(userId bson.ObjectId) ([]models.Mthread, error){
  db := Connection.DB("dartboard")
  var user struct {
    Saved []bson.ObjectId `bson:"saved"`
  }
  if err := db.C("users").Find(bson.M{"_id": userId}).Select(bson.M{"saved": 1}).One(&user); err != nil {
    return nil, err
  }
  //make slice
  threads := make([]models.Mthread, 0)
  for _, save := range user.Saved {
    var item models.Mthread
    err := db.C("mthreads").Find(bson.M{"_id": save}).One(&item)
    if err != nil {
      return nil, err
    }
    //add thread to grouping
    threads = append(threads, item)
  }
  return threads, nil
}

//[READ]
func GetNotifications(userId bson.ObjectId) ([]models.Notification, error){
  db := Connection.DB("dartboard")
  var user struct {
    Notifications []bson.ObjectId `bson:"notifications"`
  }
  if err := db.C("users").Find(bson.M{"_id": userId}).Select(bson.M{"notifications": 1}).One(&user); err != nil {
    return nil, err
  }

  //make slice
  notifications := make([]models.Notification, 0)
  for _, note := range user.Notifications {
    var noted models.Notification
    err := db.C("notifications").Find(bson.M{"_id": note}).One(&noted)
    if err != nil {
      return nil, err
    }
    //add notification to slice
    notifications = append(notifications, noted)
  }
  return notifications, nil
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

//[READ] gets user data -- posts, watchlist, thread likes, friends -- just returns what we need
func GetUser(user string) (*models.GetUser, error) {
  db := Connection.DB("dartboard")
  usr := bson.ObjectIdHex(user)
  var userData models.GetUser
  fields := bson.M{"email": 1, "username": 1, "unread": 1}
  if err := db.C("users").Find(bson.M{"_id": usr}).Select(fields).One(&userData); err != nil {
    return nil, err
  }

  return &userData, nil
}

//[READ] gets user data -- posts, watchlist, thread likes, friends -- just returns what we need
func GetIdFromUsername(username string) string {
  db := Connection.DB("dartboard")

  var userData struct {
    Id bson.ObjectId
  }
  fields := bson.M{"_id": 1}
  if err := db.C("users").Find(bson.M{"username": username}).Select(fields).One(&userData); err != nil {
    return ""
  }

  return userData.Id.Hex()
}

//[READ] compares the user & hashed password given to the one in the DB
func LoginCheck(email, hashword string) (string, bool) {

  //get proper DB
  db := Connection.DB("dartboard")

  //make struct for query
  var usr struct {
    Password string `bson:"password"`
    Id bson.ObjectId `bson:"_id"`
  }

  //run query, only getting password
  err := db.C("users").Find(bson.M{"email": email}).Select(bson.M{"password": 1, "_id": 1}).One(&usr);

  //if there's a problem, or unmatching passwords, return false
  if err != nil || usr.Password != hashword {
    return "", false
  }

  //it's a match!
  return usr.Id.Hex(), true
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

//[READ] gets a user's friends -- in a string slice format
func GetFriendsJoined(id bson.ObjectId) ([]string, error) {
  //grab proper db
  db := Connection.DB("dartboard").C("users")

  //anonymous struct just so type doesn't fail us on unmarshalling to []bson.ObjectId
  var people struct{
    Friends []bson.ObjectId
  }

  //get friends
  if err := db.Find(bson.M{"_id": id}).Select(bson.M{"friends" : 1}).One(&people); err != nil {
    return nil, err
  }

  //create slice that we'll return
  friends := make([]string,0)

  //go through friends
  for _, friend := range people.Friends {
    var fr struct {
      Username string
    }
    err := db.Find(bson.M{"_id": friend}).Select(bson.M{"username": 1}).One(&fr)
    if err != nil {
      return nil, err
    }
    friends = append(friends, fr.Username)
  }

  //if all goes well, send back nil error and friends list
  return friends, nil
}

//[READ] gets user's username
func GetUsername(id bson.ObjectId) string {

  //grab proper db
  db := Connection.DB("dartboard").C("users")

  //anonymous struct just so type doesn't fail us on unmarshalling to []bson.ObjectId
  var person struct{
    Username string
  }

  //get friends
  err := db.Find(bson.M{"_id": id}).Select(bson.M{"username" : 1}).One(&person)
  if err != nil {
    return ""
  }

  //all is well, so return nil error and friends
  return person.Username
}
