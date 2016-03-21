/*
   read.go -- DB queries that read objects
*/
package db

//[READ] gets all posts for a given thread
func GetThread(threadID string) {
  //call DB
}

//[READ] gets available threads for a given group on a particular page
func GetGroup(group string, page int) {
  //call DB
}

//[READ] gets user data -- posts, watchlist, thread likes, friends
func GetUser(userID string) {
  //call DB
}

//[READ] gets feed for person (integrate friends first)
func GetFeed(userID string) {
  //call DB
}
