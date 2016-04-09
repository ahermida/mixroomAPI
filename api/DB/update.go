/*
   update.go -- DB queries that edit objects
*/
package db

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "fmt"
  "github.com/ahermida/sudopostAPI/api/Config"

)

/**
    GROUPS -------------------------------------------------------------
 */

//[UPDATE] pushes a thread to a user's watchlist
func AddAdmin(oldAdmin, user bson.ObjectId, group string) error {

  //get group info
  grp, err := CheckGroup(group)

  //check if something went wrong
  if err != nil {
    return err
  }

  //check for admin in group
  hasAdmin := false

  //if person is author
  if grp.Author == oldAdmin {
    hasAdmin = true
  }

  //check Admins for old admin (only admins can add admins)
  for _, person := range grp.Admins {
    if oldAdmin == person {
      hasAdmin = true
    }
  }

  //if users have permission, add them
  if !hasAdmin {
    return errors.New("If users don't have admin permissions, they can't add other users as admins.")
  }

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change
  change := bson.M{"$addToSet": bson.M{"admins" : user}}

  //run update to group
  errFromUpdate := db.C("groups").Update(bson.M{"name": group}, change)
  if errFromUpdate != nil {
    return errFromUpdate
  }

  //should be nil if nothing went wrong updating
  return nil
}

//[UPDATE] Removes admins, only the author can do this
func RemoveAdmin(oldAdmin, user bson.ObjectId, group string) error {

  //get group info
  grp, errFromCheck := CheckGroup(group)

  //make sure nothing went wrong getting grp
  if errFromCheck != nil {
    return errFromCheck
  }

  //if person is author
  if grp.Author != oldAdmin {
    return errors.New("Only author can remove admins.")
  }

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change
  change := bson.M{"$pull": bson.M{"admins" : user}}

  //run update
  err := db.C("groups").Update(bson.M{"name": group}, change)
  if err != nil {
    return err
  }

  //should be nil if nothing went wrong
  return nil
}

/**
    POSTS -------------------------------------------------------------
 */

 //[UPDATE] change the text for a given post (id)
func EditPost(text string, post, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //check if user is author of post
  var getPost struct {
    AuthorId bson.ObjectId `bson:"authorId"`
  }

  //we have to make sure that we're the author
  if err := db.C("posts").Find(bson.M{"_id": post}).Select(bson.M{"authorId": 1}).One(&getPost); err != nil {

   //if there's an error return it
   return err
  }

  //check if our user has the capability of removing this post
  if getPost.AuthorId != user {
    return errors.New("User doesn't have the authorization to edit this request.")
  }

  //setup change -- modifying the password
  change := bson.M{"$set": bson.M{"body" : text}}

  //run update to user (found by _id)
  err := db.C("posts").Update(bson.M{"_id": post}, change)
  if err != nil {
    return err
  }

  //should be nil if nothing went wrong
  return nil
}

/**
    USER AUTH -------------------------------------------------------------
 */

//[UPDATE] activates account (when email is verified)
func ActivateAccount(user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change -- modifying activated value in user
  change := bson.M{"$set": bson.M{"activated" : true}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

//[UPDATE] activates account (when email is verified)
func AddName(user bson.ObjectId, name string) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change -- modifying activated value in user
  change := bson.M{"$set": bson.M{"name" : name}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

//[UPDATE] changes password for a given uid
func ChangePassword(newPassword string, oldPassword string, user bson.ObjectId) error {
  //get proper DB
  db := Connection.DB(config.DBName)

  var usr struct {
    Password string `bson:"password"`
  }

  if err := db.C("users").Find(bson.M{"_id": user}).Select(bson.M{"password": 1}).One(&usr); err != nil {
    return err
  }

  //if passwords don't match, return err
  if usr.Password != oldPassword {
    return errors.New("Passwords have to match.")
  }

  //setup change -- modifying the password
  change := bson.M{"$set": bson.M{"password" : newPassword}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

//[UPDATE] deletes a user's account (in reality, updates 'deleted' to 'true')
func DeleteUser(user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //find user's name
  var usr struct {
    Username string `bson:"username"`
  }

  //query username
  if err := db.C("users").Find(bson.M{"_id": user}).Select(bson.M{"username": 1}).One(&usr); err != nil {
    return err
  }

  //get friends so we can remove this user from each of their friends lists
  friends, friendErr := GetFriends(usr.Username)

  //return friendErr if something went wrong getting friends
  if friendErr != nil {
    return friendErr
  }

  //go through and make the changes
  for _, friend := range friends {

    //set up change
    change := bson.M{"$pull": bson.M{"friends" : friend}}

    //remove ourself, check for err
    if err := db.C("users").Update(bson.M{"_id": user}, change); err != nil {
      return err
    }
  }

  //remove our entire friends list
  rmFriends := bson.M{"$set": bson.M{"friends": make([]bson.ObjectId,0)}}

  //run update to user (found by _id)
  if err := db.C("users").Update(bson.M{"_id": user}, rmFriends); err != nil {
    return err
  }

  //setup change -- modifying activated value in user -- reactivate with email
  change := bson.M{"$set": bson.M{"activated" : false}}

  //run update to user (found by _id)
  if err := db.C("users").Update(bson.M{"_id": user}, change); err != nil {
    return err
  }

  //should be nil if nothing went wrong
  return nil
}

/**
    USERNAMES -------------------------------------------------------------
 */

//[UPDATE] changes username for a given uid
func ChangeUsername(username string, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //anonymous struct for simplicity in extracting user's usernames
  var person struct {
    Usernames []string `bson:"usernames"`
  }

  //check if we have username
  if err := db.C("users").Find(bson.M{"_id": user}).Select(bson.M{"usernames": 1}).One(&person); err != nil {
    return err
  }

  //check if user actually has username
  hasUsername := false

  //if it's in the user's list of usernames, set to true
  for _, item := range person.Usernames {
    if item == username {
      hasUsername = true
    }
  }

  //username not owned by user, so don't swap to it
  if !hasUsername {

    //let ourselves know that it failed
    return errors.New("Username isn't owned by user.")
  }

  //setup change -- modifying the username
  change := bson.M{"$set": bson.M{"username" : username}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

//[UPDATE] adds username for a given uid
func AddUsername(username string, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //check if the username is already taken
  count, err := db.C("users").Find(bson.M{"usernames": username}).Count()

  //check if something went wrong with query
  if err != nil {
    return err
  }

  //check if somebody has the username
  if count > 0 {
    return errors.New("That username is already taken.")
  }

  //check if we have fewer than 3 usernames
  queryUser := db.C("users").Find(bson.M{"_id": user})

  //anonymous struct for simplicity in extracting user's usernames
  var person struct {
    Usernames []string `bson:"usernames"`
  }

  //check if we have username
  if errFromCheck := queryUser.Select(bson.M{"usernames": 1}).One(&person); err != nil {
    return errFromCheck
  }

  //check for count
  if len(person.Usernames) > 2 {
    return errors.New("Can't have more than 3 usernames.")
  }

  //we're all good, so go ahead and set up the query
  change := bson.M{"$addToSet": bson.M{"usernames" : username}}

  //run update to user (found by _id)
  errFromChange := db.C("users").Update(bson.M{"_id": user}, change)

  if errFromChange != nil {
    return errFromChange
  }

  //should be nil if nothing went wrong
  return nil
}

//[UPDATE] removes username for a given uid -- musn't be username
func RemoveUsername(username string, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //check if we have fewer than 3 usernames
  queryUser := db.C("users").Find(bson.M{"_id": user})

  //anonymous struct for simplicity in extracting user's usernames
  var person struct {
    Usernames []string `bson:"usernames"`
  }

  //check if we have username
  if err := queryUser.Select(bson.M{"usernames": 1}).One(&person); err != nil {
    return err
  }

  //check if we own the username
  hasUsername := false

  //if it's in the user's list of usernames, set to true
  for _, item := range person.Usernames {
    if item == username {
      hasUsername = true
    }
  }

  //if user doesn't own username, user shouldn't be able to remove it
  if !hasUsername {
    return errors.New("User must own username to remove it.")
  }

  //we're all good, so go ahead and set up the query
  change := bson.M{"$pull": bson.M{"usernames" : username}}

  //run update to user (found by _id)
  errFromChange := db.C("users").Update(bson.M{"_id": user}, change)

  if errFromChange != nil {
    return errFromChange
  }

  //should be nil if nothing went wrong
  return nil
}

/**
    SAVED THREADS -------------------------------------------------
 */

//[UPDATE] pushes a thread to a user's watchlist
func SaveThread(mthread, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change -- push thread ID to saved
  change := bson.M{"$addToSet": bson.M{"saved" : mthread}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

//[UPDATE] removes a thread from a user's watchlist
func UnsaveThread(thread, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change -- modifying activated value in user -- reactivate with email
  change := bson.M{"$pull": bson.M{"saved" : thread}}

  //run update to user (found by _id)
  err := db.C("users").Update(bson.M{"_id": user}, change)

  //should be nil if nothing went wrong
  return err
}

/**
    FRIENDS ----------------------------------------------------
 */

//[UPDATE] saves a friend
func AddFriend(friend, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup changes -- add friend to list -- do same for friend
  change := bson.M{"$addToSet": bson.M{"friends" : friend}}
  rmRequest := bson.M{"$pull": bson.M{"requests": friend}}
  friendChange := bson.M{"$addToSet": bson.M{"friends": user}}

  //run update to user (found by _id)
  if err := db.C("users").Update(bson.M{"_id": user}, change); err != nil {
    return err
  }

  //run update to remove request from list
  if err := db.C("users").Update(bson.M{"_id": user}, rmRequest); err != nil {
    return err
  }

  //run update to friend (found by _id)
  if err := db.C("users").Update(bson.M{"_id": friend}, friendChange); err != nil {
    return err
  }

  //get new friend's new username, let them know that we accepted their request
  name := GetUsername(friend)
  if name == "" {
    return errors.New("Issue getting username for a given _id")
  }

  //note that we made -- USE OUR OWN MARKUP HERE so we can style it appropriately
  note := fmt.Sprintf("%s accepted your friend request", name)

  //create link that they can click on to view our profile
  link := fmt.Sprintf("localhost:8080/user/%s", name)

  //add notification to new friend
  if err := CreateNotification(friend, link, note); err != nil {
    return err
  }

  //should be nil if nothing went wrong
  return nil
}

//[UPDATE] removes a friend
func RemoveFriend(friend, user bson.ObjectId) error {

  //get proper DB
  db := Connection.DB(config.DBName)

  //setup change -- remove friend from list -- do same for friend
  change := bson.M{"$pull": bson.M{"friends" : friend}}

  //setup change for friend -- remove self from list
  friendChange := bson.M{"$pull": bson.M{"friends": user}}

  //run update to user (found by _id)
  if err := db.C("users").Update(bson.M{"_id": user}, change); err != nil {
    return err
  }

  //run update to friend (found by _id)
  if err := db.C("users").Update(bson.M{"_id": friend}, friendChange); err != nil {
    return err
  }

  //should be nil if nothing went wrong
  return nil
}
