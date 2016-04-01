/*
  Test DB Operations
*/
package server

import (
  //  "net/http"
  //  "net/http/httptest"
  //  "testing"
)

/*
 TEST these:
  CreateUser(email, username, password string) (string, error)
  CreateThread(group string, anonymous bool, post *models.Post) error
  CreateHeadPost(author, body, content string, authorId bson.ObjectId) *models.Post
  CreatePost(authorId, thread bson.ObjectId, responseTo []bson.ObjectId, author, body, content string) (*models.Post, error)
  CreateGroup(group string, user bson.ObjectId, private bool) error
  GroupExists(group string) bool
  CreateNotification(id bson.ObjectId, link, text string) error
  RequestFriend(user bson.ObjectId, username, friend string) error
  DeleteGroup(group string, user bson.ObjectId) error
  DeleteThread(threadID, userID bson.ObjectId)
  DeletePost(postID, userID bson.ObjectId) error
  RemoveUser(user bson.ObjectId) error
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
  ChangePassword(newPassword string, oldPassword string, user bson.ObjectId) error
  ChangeUsername(username string, user bson.ObjectId) error
  AddUsername(username string, user bson.ObjectId) error
  RemoveUsername(username string, user bson.ObjectId) error
  EditPost(text string, post, user bson.ObjectId) error
  DeleteUser(user bson.ObjectId) error
  AddAdmin(oldAdmin, user bson.ObjectId, group string) error
  RemoveAdmin(oldAdmin, user bson.ObjectId, group string) error
  SaveThread(thread, user bson.ObjectId) error
  UnsaveThread(thread, user bson.ObjectId) error
  AddFriend(friend, user bson.ObjectId) error
  RemoveFriend(friend, user bson.ObjectId) error
*/
