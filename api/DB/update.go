/*
   update.go -- DB queries that edit objects
*/
package db

//[UPDATE] activates account (when email is verified)
func ActivateAccount(uid string) {
  //call DB
}

//[UPDATE] changes password for a given uid
func ChangePassword(newPassword, uid string) {
  //call DB
}

//[UPDATE] change the text for a given post (id)
func EditPost(text, id string) {
  //call DB
}

//[UPDATE] deletes a user's account (in reality, updates 'deleted' to 'true')
func DeleteUser(uid string) {
  //call DB
}

//[UPDATE] saves a thread to a user's watchlist
func WatchThread(threadID, userID string) {
  //call DB
}

//[UPDATE] removes a thread from a user's watchlist
func UnwatchThread(threadID, userID string) {
  //call DB
}

//[UPDATE] saves a friend
func AddFriend(friendID, userID string) {
  //this happens when somebody accepts a request
}

//[UPDATE] removes a friend
func RemoveFriend(friendID, userID string) {
  //this happens when somebody accepts
}
