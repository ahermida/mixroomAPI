/*
   read.go -- DB queries that read objects
*/
package db

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "github.com/ahermida/sudopostAPI/api/Models"
  "github.com/ahermida/sudopostAPI/api/Config"

)

/**
    GROUPS -------------------------------------------------------------
 */

//[READ] gets available threads for a given group on a particular page
func GetGroup(group string, page int) ([]models.Mthread, error) {

  //check if group exists
  if !GroupExists(group) {

    //if it doesn't, return an error
    return nil, errors.New("Can't get a group that doens't exist.")
  }

  //get DB
  db := Connection.DB(config.DBName)

  //Sort by Timestamp --> Get The Range from (page * 30) -- 30 items -- project all fields
  pipeline := []bson.M{bson.M{"$sort": bson.M{"created": -1 }},
                       bson.M{"$skip": page * 30},
                       bson.M{"$limit": 30},
                       bson.M{"$project": bson.M{
                         "_id": 1,
                         "id": 1,
                         "created": 1,
                         "thread": 1,
                         "threadId": 1,
                         "head": 1,
                         "group": 1}}}

  //set up Pipe to actually run query
  pipe := db.C(group).Pipe(pipeline)

  //slice of threads that will be populated by query
  var threads []models.Mthread

  //run it
  if err := pipe.All(&threads); err != nil {

    //if something is wrong, return err
    return nil, err
  }

  //get size for each thread
  for i := 0; i < len(threads); i++ {
    threads[i].Size = GetThreadSize(threads[i].ThreadId)
  }


  //else return the threads, and a nil error
  return threads, nil
}


//[READ] gets group info -- meta info about group from group collection
func CheckGroup(group string) (*models.Group, error) {
  db := Connection.DB(config.DBName)
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
  db := Connection.DB(config.DBName)
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

//[READ] checks if user is a member of a group
func GetPermission(group string, user string) *models.Permission {
  db := Connection.DB(config.DBName)
  var g models.Group
  err := db.C("groups").Find(bson.M{"name": group}).One(&g)

  //check if something went wrong finding the group
  if err != nil {
    return &models.Permission{
      Allowed: false,
      Author: false,
      Admin: false,
      Mod: false,
    }
  }

  if user == g.Author.Hex() {
    return &models.Permission{
      Allowed: true,
      Author: true,
      Admin: true,
      Mod: true,
    }
  }

  //bring back group metadata
  for _, member := range g.Admins {
    if member.Hex() == user {
      if g.Private {
        return &models.Permission{
              Allowed: true,
              Author: false,
              Admin: true,
              Mod: false,
        }
      } else {
       return &models.Permission{
              Allowed: true,
              Author: false,
              Admin: true,
              Mod: true,
        }
      }
    }
  }

  if g.Private {
    return &models.Permission{
        Allowed: false,
        Author: false,
        Admin: false,
        Mod: false,
      }
  }

  //this shouldn't happen if we're authors or admins
  return &models.Permission{
    Allowed: true,
    Author: false,
    Admin: false,
    Mod: false,
  }
}

//[READ] searches groups for match
func SearchGroups(str string) ([]string, error) {
  db := Connection.DB(config.DBName)

  //slice of usernames that we'll return to the user
  groups := make([]string, 0)

  //groups that we'll get from query
  var getGroup []struct{
    Name string `bson:"name"`
  }

  if err := db.C("groups").Find(bson.M{"$text": bson.M{"$search": str}}).Select(bson.M{"name": 1}).All(&getGroup); err != nil {
    return nil, err
  }

  //go through users and add to usernames
  for _, grp := range getGroup {
    groups = append(groups, grp.Name)
  }

  //send it back & nil error
  return groups, nil
}


//[READ] get metadata about group
func GetGroupInfo(grp string) (*models.GroupInfo, error) {

  //slice of usernames that we'll return to the user
  check, err := CheckGroup(grp)

  // type GroupInfo struct {
  //   Created time.Time `json:"created"`
  //   Name    string    `json:"name"`
  //   Author  string    `json:"author"`
  //   Private bool      `bson:"private"`
  // }

  if err != nil {
    return nil, err
  }

  username := GetUsername(check.Author)

  if username == "" {
    username = "Anonymous"
  }

  group := &models.GroupInfo{
    Created: check.Created,
    Name: check.Name,
    Author: username,
    Private: check.Private,
  }

  //send it back & nil error
  return group, nil
}


/**
    THREADS -------------------------------------------------------------
 */

//[READ] gets all posts for a given thread
func GetThread(threadID bson.ObjectId) (*models.ResThread, error) {

 //call DB
 db := Connection.DB(config.DBName)

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
   Mthread: thread.Mthread.Hex(),
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

//[READ] kinda like getting a thread, but it's just the most popular posts (8 of them)
func GetPopularPosts(id string, skip int) (*models.PopularPosts, error) {
  db := Connection.DB(config.DBName)

  //Get only those posted in the past 24 hours, sort them by # of replies, and limit them
  pipeline := []bson.M{bson.M{"$match": bson.M{"created": bson.M{"$gt": bson.Now().AddDate(0, 0, -1)}}},
                       bson.M{"$project": bson.M{"replies_size": bson.M{"$size": bson.M{"$ifNull": []string{"$replies", "[]"}}},
                       "_id": 1,
                       "id": 1,
                       "created": 1,
                       "thread": 1,
                       "replies": 1,
                       "author": 1,
                       "content": 1,
                       "contentType": 1,
                       "body": 1,
                       }},
                       bson.M{"$sort": bson.M{"replies_size": 1}},
                       bson.M{"$skip": skip * 8},
                       bson.M{"$limit": 8}}

  /* NOTE: We get 20 because the top 8 might be blocked */

  //set up Pipe to actually run query
  pipe := db.C("posts").Pipe(pipeline)

  //for population later
  pop := &models.PopularPosts{
    Posts: make([]models.PopularPost,0),
  }

  //gonna populate this slice with the query
  var pops []models.PopularPost

  //run it
  if err := pipe.All(&pops); err != nil {

    //if something is wrong, return err
    return nil, err
  }

  //filter out the posts that we're not allowed to see
  for _, el := range pops {
    grp := GetThreadParent(el.Thread.Hex())
    if IsMember(grp, id){
      el.Size = GetThreadSize(el.Thread.Hex())
      el.ThreadId = el.Thread.Hex()
      el.Group = grp
      pop.Posts = append(pop.Posts, el)
    }
  }

  //else return the threads, and a nil error
  return pop, nil
}

//[READ] search posts text index for a given string
func SearchThreads(id string, str string, page int) ([]models.Mthread, error) {
  db := Connection.DB(config.DBName)

  //Sort by Timestamp --> Get The Range from (page * 30) -- 30 items -- project all fields
  pipeline := []bson.M{bson.M{"$match": bson.M{"$text": bson.M{"$search": str}}},
                       bson.M{"$sort": bson.M{"created": -1}},
                       bson.M{"$skip": page * 30},
                       bson.M{"$limit": 30},
                       bson.M{"$project": bson.M{
                         "_id": 1,
                         "id": 1,
                         "created": 1,
                         "thread": 1,
                         "threadId": 1,
                         "head": 1,
                         "group": 1}}}

  //set up Pipe to actually run query
  pipe := db.C("mthreads").Pipe(pipeline)

  //for population later
  var threads []models.Mthread

  //run it
  if err := pipe.All(&threads); err != nil {

    //if something is wrong, return err
    return nil, err
  }

  //could do it in place, but it's much easier to just have a slice that we can append to
  var goodThreads = make([]models.Mthread, 0)

  //filter out the threads that we're not allowed to see
  for _, el := range threads {
    if IsMember(el.Group, id){
      el.Size = GetThreadSize(el.ThreadId)
      goodThreads = append(goodThreads, el)
    }
  }

  //else return the threads, and a nil error
  return goodThreads, nil
}

//[READ] searches usernames and names for match -- EXACT
func SearchUsers(str string) ([]string, error) {
  db := Connection.DB(config.DBName)

  //slice of usernames that we'll return to the user
  usernames := make([]string, 0)

  //search name & usernames
  pipeline := []bson.M{bson.M{"$match": bson.M{"$or": []interface{}{
                         bson.M{"usernames": str},
                         bson.M{"name": str}}}},
                         bson.M{"$limit": 10},
                         bson.M{"$project": bson.M{
                           "username": 1}}}

  //pipe it up pipe it up pipe it up
  pipe := db.C("users").Pipe(pipeline)

  var users []models.SearchUsers

  //run it
  if err := pipe.All(&users); err != nil {

    //if something is wrong, return err
    return nil, err
  }

  //go through users and add to usernames
  for _, usr := range users {
    usernames = append(usernames, usr.Username)
  }

  //send it back & nil error
  return usernames, nil
}

//[READ] returns group that the given thead (hex string) belongs to
func GetThreadParent(thread string) string {
  db := Connection.DB(config.DBName)
  var thrd struct {
    Group string `bson:"group"`
  }
  err := db.C("threads").Find(bson.M{"_id": bson.ObjectIdHex(thread)}).Select(bson.M{"group": 1}).One(&thrd)
  if err != nil {
    return ""
  }

  return thrd.Group
}

//[READ] get a thread's length (how many posts are in a thread)
func GetThreadSize(thread string) int {
  db := Connection.DB(config.DBName)
  var length []struct {
    Size int `bson:"size"`
  }
  pipeline := []bson.M{
    bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(thread)}},
    bson.M{"$project": bson.M{"size": bson.M{"$size": "$posts"}}},
  }

  //pipe it up pipe it up pipe it up
  pipe := db.C("threads").Pipe(pipeline)

  //run it
  if err := pipe.All(&length); err != nil {

    //if something is wrong, return err
    return 0
  }

  //should only be one match, and its size should be this. Kinda hacky method -- should be revisited
  return length[0].Size - 1
}

/**
    USER  -------------------------------------------------------------
 */

//[READ] gets user data -- posts, watchlist, thread likes, friends -- just returns what we need
func GetUser(user string) (*models.GetUser, error) {
  db := Connection.DB(config.DBName)
  usr := bson.ObjectIdHex(user)
  var userData models.GetUser
  fields := bson.M{"email": 1, "username": 1, "unread": 1, "usernames": 1, "saved": 1}
  if err := db.C("users").Find(bson.M{"_id": usr}).Select(fields).One(&userData); err != nil {
    return nil, err
  }

  //populate saved strings
  for _, el := range userData.Saved {
    userData.SavedStr = append(userData.SavedStr, el.Hex())
  }

  return &userData, nil
}

/**
    SAVED THREADS -------------------------------------------------------------
 */

//[READ]
func GetSaved(userId bson.ObjectId) ([]models.Mthread, error){
  db := Connection.DB(config.DBName)
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
    if err == nil {
      item.Size = GetThreadSize(item.ThreadId)
      //add thread to grouping
      threads = append(threads, item)
    }
  }
  return threads, nil
}

/**
    USER NOTIFICATIONS -------------------------------------------
 */

//[READ]
func GetNotifications(userId bson.ObjectId) ([]models.Notification, error){
  db := Connection.DB(config.DBName)
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

/**
    USER AUTH -------------------------------------------------------
 */

//[READ] compares the user & hashed password given to the one in the DB
func LoginCheck(email, hashword string) (string, bool) {

  //get proper DB
  db := Connection.DB(config.DBName)

  //make struct for query
  var usr struct {
    Password string `bson:"password"`
    Id bson.ObjectId `bson:"_id"`
    Activated bool `bson:"activated"`
  }

  //run query, only getting password
  err := db.C("users").Find(bson.M{"email": email}).Select(bson.M{"password": 1, "_id": 1, "activated": 1}).One(&usr);

  //if there's a problem, or unmatching passwords, return false
  if err != nil || usr.Password != hashword || !usr.Activated {
    return "", false
  }

  //it's a match!
  return usr.Id.Hex(), true
}

/**
    USER FRIENDS -------------------------------------------------------
 */

//[READ] gets user's friends -- resolving "joins"
func GetFriends(author string) ([]bson.ObjectId, error) {

  //grab proper db
  db := Connection.DB(config.DBName).C("users")

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

//[READ] gets user's friends -- resolving "joins"
func GetFriendsById(author bson.ObjectId) ([]bson.ObjectId, error) {

  //grab proper db
  db := Connection.DB(config.DBName).C("users")

  //anonymous struct just so type doesn't fail us on unmarshalling to []bson.ObjectId
  var people struct{
    Friends []bson.ObjectId `bson:"friends"`
  }

  //get friends
  err := db.Find(bson.M{"_id": author}).Select(bson.M{"friends" : 1}).One(&people)
  if err != nil {
    return nil, err
  }

  //all is well, so return nil error and friends
  return people.Friends, nil
}

//[READ] gets a user's friends -- in a string slice format
func GetFriendsJoined(id bson.ObjectId) ([]string, error) {
  //grab proper db
  db := Connection.DB(config.DBName).C("users")

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

/**
    USER UTIL -------------------------------------------------------
 */

//[READ] gets user's username
func GetUsername(id bson.ObjectId) string {

  //grab proper db
  db := Connection.DB(config.DBName).C("users")

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

//[READ] gets user data -- posts, watchlist, thread likes, friends -- just returns what we need
func GetIdFromUsername(username string) string {
  db := Connection.DB(config.DBName)

  var userData struct {
    Id bson.ObjectId `bson:"_id"`
  }
  fields := bson.M{"_id": 1}
  if err := db.C("users").Find(bson.M{"username": username}).Select(fields).One(&userData); err != nil {
    return ""
  }

  return userData.Id.Hex()
}

//[READ] gets user data -- posts, watchlist, thread likes, friends -- just returns what we need
func GetIdFromEmail(email string) string {
  db := Connection.DB(config.DBName)

  var userData struct {
    Id bson.ObjectId `bson:"_id"`
  }
  fields := bson.M{"_id": 1}
  if err := db.C("users").Find(bson.M{"email": email}).Select(fields).One(&userData); err != nil {
    return ""
  }

  return userData.Id.Hex()
}
