/*
  Test DB Operations
*/
package db

import (
  "testing"
  "gopkg.in/mgo.v2/bson"
)

/*
 TEST these funcs:
  CreateUser(email, username, password string) (string, error)
  CreateThread(group string, anonymous bool, post *models.Post) error
  CreateHeadPost(author, body, content string, authorId bson.ObjectId) *models.Post
  CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId, author, body, content string) (*models.Post, error)
  CreateGroup(group string, user bson.ObjectId, private bool) error
  GroupExists(group string) bool
  CreateNotification(id bson.ObjectId, link, text string) error
  RequestFriend(user bson.ObjectId, username, friend string) error
  CheckGroup(group string) (*models.Group, error)
  IsMember(group string, user string) bool
  GetThreadParent(thread string) string
  GetGroup(group string, page int) ([]models.Mthread, error)
  GetSaved(userId bson.ObjectId) ([]models.Mthread, error)
  GetNotifications(userId bson.ObjectId) ([]models.Notification, error)
  GetThread(threadID bson.ObjectId) (*models.ResThread, error)
  GetUser(user string) (*models.GetUser, error)
  GetIdFromUsername(username string) string
  LoginCheck(email, hashword string) (string, bool)
  GetFriends(author string) ([]bson.ObjectId, error)
  GetFriendsJoined(id bson.ObjectId) ([]string, error)
  GetUsername(id bson.ObjectId) string
  ActivateAccount(user bson.ObjectId) error
  EditPost(text string, post, user bson.ObjectId) error
  ChangePassword(newPassword string, oldPassword string, user bson.ObjectId) error
  ChangeUsername(username string, user bson.ObjectId) error
  AddUsername(username string, user bson.ObjectId) error
  RemoveUsername(username string, user bson.ObjectId) error
  AddAdmin(oldAdmin, user bson.ObjectId, group string) error
  RemoveAdmin(oldAdmin, user bson.ObjectId, group string) error
  SaveThread(thread, user bson.ObjectId) error
  UnsaveThread(thread, user bson.ObjectId) error
  AddFriend(friend, user bson.ObjectId) error
  DeleteUser(user bson.ObjectId) error
  DeleteGroup(group string, user bson.ObjectId) error
  DeleteThread(threadID, userID bson.ObjectId)
  DeletePost(postID, userID bson.ObjectId) error
  RemoveFriend(friend, user bson.ObjectId) error
  RemoveUser(user bson.ObjectId) error

  sequence: -> CreateUser -> ActivateAccount -> LoginCheck ->
           -> GetUser -> CreateUser2 -> ActivateAccount2 -> LoginCheck2 ->
           -> GetUser2 -> CreateGroup -> IsMember -> RequestFriend -> AddFriend -> GetFriends ->
           -> GetFriendsJoined-> Add Admin -> RemoveAdmin -> CreateHeadPost ->
           -> CreateThread -> GetThreadParent -> Create Post -> Edit Post -> Delete Post ->
           -> SaveThread -> GetSaved -> RemoveSaved -> DeleteUser2 -> RemoveUser2 ->
           -> GetIdFromUsername -> CreateNotification -> GetNotifications -> DeleteThread ->
           -> DeleteGroup -> RemoveUser
*/

//create store to hold user Ids, thread Ids, and other string data
var store = make(map[string]string)

func TestUsers(t *testing.T) {
  //t.Errorf("something went wrong!")
  id, _ := CreateUser("test@gmail.com","test","testtest1")
  if id != "" {
    errActivate := ActivateAccount(bson.ObjectIdHex(id))
    if errActivate != nil {
      t.Errorf("problem activating account")
    }
  }
  id2, _ := CreateUser("test2@gmail.com","test2","testtest2")
  if id2 != "" {
    errActivate := ActivateAccount(bson.ObjectIdHex(id2))
    if errActivate != nil {
      t.Errorf("problem activating account")
    }
  }

  //get Id from Username
  id = GetIdFromUsername("test")
  id2 = GetIdFromUsername("test2")

  //get Id from Username
  usrname := GetUsername(bson.ObjectIdHex(id))
  if usrname != "test" {
    t.Errorf("GetUsername (from id) failed!")
  }
  //get Id from Username
  usrname2 := GetUsername(bson.ObjectIdHex(id2))
  if usrname2 != "test2" {
    t.Errorf("GetUsername (from id) failed (2)!")
  }

  //get user for check
  usr, err := GetUser(id)
  if err != nil {
    t.Errorf("Couldn't get user via GetUser(id)")
  }
  if usr.Username != "test" {
    t.Errorf("Username that we got from GetUser wasn't that which we made")
  }

  //LoginCheck(email, hashword string) (string, bool)
  _, success := LoginCheck("test@gmail.com","testtest1")
  if !success {
    t.Errorf("Failed login - password didn't match for test@gmail.com & testtest1")
  }

  //LoginCheck(email, hashword string) (string, bool)
  _, success2 := LoginCheck("test2@gmail.com","testtest2")
  if !success2 {
    t.Errorf("Failed login - password didn't match for test2@gmail.com & testtest2")
  }

  store["id1"] = id
  store["id2"] = id2

  //AddName(id bson.ObjectId, name string)
  errName := AddName(bson.ObjectIdHex(id), "Octodan is Here")

  if errName != nil {
    t.Errorf("Couldn't add name")
  }
}

func TestFriends(t *testing.T) {
  //RequestFriend(user bson.ObjectId, username, friend string) error -- should refactor this one
  errFriend := RequestFriend(bson.ObjectIdHex(store["id1"]), "test", "test2")
  if errFriend != nil {
    t.Errorf("Couldn't add friend.")
  }

  // AddFriend(friend, user bson.ObjectId) error
  errAddFriend := AddFriend(bson.ObjectIdHex(store["id1"]), bson.ObjectIdHex(store["id2"]))
  if errAddFriend != nil {
    t.Errorf("Couldn't accept friend request.")
  }

  checked := false
  //check friends
  friends, errGetFriends := GetFriends("test")
  if errGetFriends != nil {
    t.Errorf("Couldn't get friends properly for username test")
  }
  for _, el := range friends {
    if el == bson.ObjectIdHex(store["id2"]) {
      checked = true
    }
  }
  //see if user is actually in there
  if !checked {
    t.Errorf("couldn't find friend in friendslist")
  }

  checkedJoined := false
  friendsJoined, errGetFriendsJoined := GetFriendsJoined(bson.ObjectIdHex(store["id1"]))
  if errGetFriendsJoined != nil {
    t.Errorf("Couldn't get friends (joined) properly for username: test")
  }
  for _, el := range friendsJoined {
    if el == "test2" {
      checkedJoined = true
    }
  }
  //see if user is actually in there
  if !checkedJoined {
    t.Errorf("couldn't find friend in friendslist")
  }

  //GetNotifications(userId bson.ObjectId) ([]models.Notification, error)
  notes, errNotes := GetNotifications(bson.ObjectIdHex(store["id1"]))
  if errNotes != nil {
    t.Errorf("Couldn't get notifications")
  }
  if len(notes) == 0 {
    t.Errorf("Notifications are broken!")
  }
}

func TestOtherUserOps(t *testing.T) {
  // ChangePassword(newPassword string, oldPassword string, user bson.ObjectId) error
  errPass := ChangePassword("testtest1", "testtest1", bson.ObjectIdHex(store["id1"]))
  if errPass != nil {
    t.Errorf("Problem with ChangePassword!")
  }
  // AddUsername(username string, user bson.ObjectId) error
  errAddingUserName := AddUsername("octodan", bson.ObjectIdHex(store["id1"]))
  if errAddingUserName != nil {
    t.Errorf("Problem adding username!")
  }
  // ChangeUsername(username string, user bson.ObjectId) error
  errChangeUsername := ChangeUsername("octodan", bson.ObjectIdHex(store["id1"]))
  if errChangeUsername != nil {
    t.Errorf("Problem changing username!")
  }
  // RemoveUsername(username string, user bson.ObjectId) error
  errChangeUsernameback := ChangeUsername("test", bson.ObjectIdHex(store["id1"]))
  if errChangeUsernameback != nil {
    t.Errorf("Problem changing username back!")
  }
  //RemoveUsername(username string, user bson.ObjectId) error
  errRemoveUsername := RemoveUsername("octodan", bson.ObjectIdHex(store["id1"]))
  if errRemoveUsername != nil {
    t.Errorf("Problem removing username: octodan")
  }
}

func TestGroup(t *testing.T) {
  //CreateGroup(group string, user bson.ObjectId, private bool) error
  CreateGroup("test", bson.ObjectIdHex(store["id1"]), true)

  if !GroupExists("test") {
    t.Errorf("Says that group (test) doesn't exist.")
  }

  if !IsMember("test", store["id1"]) {
    t.Errorf("IsMember didn't work as intended. test@gmail.com's Id failed.")
  }

  //Add admin
  errAdmin := AddAdmin(bson.ObjectIdHex(store["id1"]), bson.ObjectIdHex(store["id2"]), "test")
  if errAdmin != nil {
    t.Errorf("Problem adding an admin to: test")
  }

  //Remove admin
  errRmAdmin := RemoveAdmin(bson.ObjectIdHex(store["id1"]), bson.ObjectIdHex(store["id2"]), "test")
  if errRmAdmin != nil {
    t.Errorf("Problem adding an admin to: test")
  }
}
func TestThreads(t *testing.T) {

  //CreateHeadPost(author, body, content string, authorId bson.ObjectId) *models.Post
  post := CreateHeadPost("test", "Hello, this is a test head post", "youtube!", bson.ObjectIdHex(store["id1"]))

  store["post"] = post.Id.Hex()
  //CreateThread(group string, anonymous bool, post *models.Post) error
  errThrd := CreateThread("test", false, post)
  if errThrd != nil {
    t.Errorf("Problem creating thread in: test")
  }

  //GetGroup(group string, page int) ([]models.Mthread, error)
  mthreads, errMthreads := GetGroup("test", 0)
  if errMthreads != nil {
    t.Errorf("Problem getting group: test")
  }

  threadId := mthreads[0].Thread
  store["thread"] = threadId.Hex()
  parentgrp := GetThreadParent(threadId.Hex())
  if parentgrp != "test" {
    t.Errorf("Problem with GetThreadParent.")
  }

  //SaveThread(thread, user bson.ObjectId) error
  errSavingThread := SaveThread(mthreads[0].Id, bson.ObjectIdHex(store["id1"]))
  if errSavingThread != nil {
    t.Errorf("Problem saving thread!")
  }

  //GetSaved(userId bson.ObjectId) ([]models.Mthread, error)
  thrds, errGetSaved := GetSaved(bson.ObjectIdHex(store["id1"]))
  if errGetSaved != nil {
    t.Errorf("Problem getting saved threads!")
  }

  if thrds[0].Head.Author != "test" {
    t.Errorf("Author isn't correct in GetSaved -- [0th] thread head post")
  }

  //UnsaveThread(thread, user bson.ObjectId) error
  errUnsavingThread := UnsaveThread(threadId, bson.ObjectIdHex(store["id1"]))
  if errUnsavingThread != nil {
    t.Errorf("Problem unsaving thread!")
  }

  //CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId, author, body, content string) (*models.Post, error)
  pst, errPost := CreatePost(bson.ObjectIdHex(store["id1"]), threadId, make([]bson.ObjectId, 0), "test", "hello", "youtube!")
  if errPost != nil {
    t.Errorf("Couldn't make the post!")
  }

  //EditPost(text string, post, user bson.ObjectId) error
  editErr := EditPost("DINGO", pst.Id, bson.ObjectIdHex(store["id1"]))
  if editErr != nil {
    t.Errorf("Problem with ID of post.")
  }

  //GetThread(threadID bson.ObjectId) (*models.ResThread, error)
  resThread, errGetThread := GetThread(threadId)
  if errGetThread != nil {
    t.Errorf("Problem getting thread")
  }

  if len(resThread.Posts) == 0 {
    t.Errorf("GetThread had empty posts! That's a problem")
  }

  //Search function
  _, err := SearchThreads(store["id1"], "hello", 0)
  if err != nil {
    t.Errorf("Problem searching threads!")
  }

  _, errUserSearch := SearchUsers("Octodan")
  if errUserSearch != nil {
    t.Errorf("Couldn't search users")
  }

  _, errGrpSearch := SearchGroups("test")
  if errGrpSearch != nil {
    t.Errorf("Couldn't search groups")
  }

  //DeletePost(postID, userID bson.ObjectId) error
  errDeletePost := DeletePost(pst.Id, bson.ObjectIdHex(store["id1"]))
  if errDeletePost != nil {
    t.Errorf("Problem deleting post!")
  }

  //DeletePost(postID, userID bson.ObjectId) error
  errDeletePost2 := DeletePost(bson.ObjectIdHex(store["post"]), bson.ObjectIdHex(store["id1"]))
  if errDeletePost2 != nil {
    t.Errorf("Problem deleting post!")
  }
}

func TestRemoval(t *testing.T) {
  // RemoveFriend(friend, user bson.ObjectId) error
  errRemoveFriend := RemoveFriend(bson.ObjectIdHex(store["id2"]), bson.ObjectIdHex(store["id1"]))
  if errRemoveFriend != nil {
    t.Errorf("Problem removing test2 as a friend")
  }
  //DeleteThread(threadID, userID bson.ObjectId)
  errDeleteThread := DeleteThread(bson.ObjectIdHex(store["thread"]), bson.ObjectIdHex(store["id1"]))
  if errDeleteThread != nil {
    t.Errorf("Problem deleting thread")
  }
  // DeleteGroup(group string, user bson.ObjectId) error
  errDeleteGroup := DeleteGroup("test", bson.ObjectIdHex(store["id1"]))
  if errDeleteGroup != nil {
    t.Errorf("Problem deleting group")
  }
  // DeleteUser(user bson.ObjectId) error
  errDeleteUser1 := DeleteUser(bson.ObjectIdHex(store["id1"]))
  if errDeleteUser1 != nil {
    t.Errorf("Problem deleting test as a user")
  }
  errDeleteUser2 := DeleteUser(bson.ObjectIdHex(store["id2"]))
  if errDeleteUser2 != nil {
    t.Errorf("Problem deleting test2 as a user")
  }
  //RemoveUser(user bson.ObjectId) error
  err := RemoveUser(bson.ObjectIdHex(store["id1"]))
  if err != nil {
    t.Errorf("Problem removing test1 as a user completely")
  }
  err1 := RemoveUser(bson.ObjectIdHex(store["id2"]))
  if err1 != nil {
    t.Errorf("Problem removing test2 as a user completely")
  }
}
