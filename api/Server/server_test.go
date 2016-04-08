/*
  Test API -- Integration test
*/
package server

/*
  Endpoints:
  (31 of them)
                                                                    STRUCT INPUT
  <Auth> ------------------------------------------------------
  POST -- register user ("/auth/register", register)              - models.CreateUser
  POST -- recover user pw ("/auth/recovery", recovery)            - models.Recovery
  GET -- activate user ("/auth/activate", activate)               - [id required]
  DELETE -- deactivate user ("/auth/remove", deactivate)          - [id required]
  POST -- login user ("/auth/login", login)                       - models.AuthUser
  POST -- update password ("/auth/changepass", updatePassword)    - models.ChangePW
  </Auth>

  <User> ------------------------------------------------------
  GET get user info ("/user/", getUser)                           - [id required]
  GET get saved ("/user/saved", saved)                            - [id required]
  POST add name ("/user/name", name)                              - models.Name
  POST add saved ("/user/saved", saved)                           - models.Saved
  PUT removed saved ("/user/saved", saved)                        - models.Saved
  POST get user's threads ("/user/threads", threads)              - models.GetUserFeed
  POST add a username ("/user/username", username)                - models.Username
  PUT to change it ("/user/username", username)                   - models.Username
  DELETE to remove it ("/user/username", username)                - models.Username
  GET get all notifications ("/user/notifications", notifications)- [id required]
  GET -- gets friends list ("/user/friends", friends)             - [id required]
  POST add a friend -- creates request ("/user/friends", friends) - models.Friend
  PUT accept it ("/user/friends", friends)                        - models.Friend
  DELETE to remove it ("/user/friends", friends)                  - models.Friend
  </User>

  <Groups> ----------------------------------------------------
  POST -- to get group -- paginated ("/group/", getGroup)         - models.GetGroup
  POST -- check if admin or author ("/group/auth", getPermission) - models.Grp
  POST -- for groups -- creating them ("/group/modify", grp)      - models.CreateGroup
  DELETE -- for groups -- deleting ("/group/modify", grp)         - models.Grp
  POST -- set admins for group ("/group/admin", admn)             - models.GroupAdmin
  PUT -- delete Admins in groups ("/group/admin", admn)           - models.GroupAdmin
  </Groups>

  <Threads> ---------------------------------------------------
  POST -- Get Thread ("/thread/", getThread)                      - models.GetThread
  POST -- Create Thread ("/thread/modify", thrd)                  - models.NewThread
  DELETE -- Delete thread ("/thread/modify", thrd)                - models.RemoveThread
  POST -- Create Post ("/thread/post", pst)                       - models.NewPost
  DELETE -- Delete Post ("/thread/post", pst)                     - models.DeletePost
  PUT -- Edit Post ("/thread/post", pst)                          - models.EditPost
  </Threads>

*/
import (
  "net/http"
  "github.com/ahermida/dartboardAPI/api/Util"
  "github.com/ahermida/dartboardAPI/api/DB"
  "net/http/httptest"
  "testing"
  "gopkg.in/mgo.v2/bson"
  "strings"
  "fmt"
)

//setup server so we can test endpoints
var server *httptest.Server

func init() {

  //grab mux from server.go and run it
  server = httptest.NewServer(Server)
  fmt.Printf(server.URL)
}

/*
  Let it be known that at this point I gave and decided to write a function
  to do this json thing automatically -- feeling stupid for not doing this
  earlier...
*/
func DoTest(method, url, json string, expected int) bool {
  id := db.GetIdFromUsername("test")
  if id == "" {
    fmt.Printf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    fmt.Printf("Couldn't make token for group.")
    return false
  }
  reader := strings.NewReader(json)
  request, err := http.NewRequest(method, url, reader)
  request.Header.Set("access_token", token)

  res, errSending := http.DefaultClient.Do(request)
  if errSending != nil {
    fmt.Printf("Couldn't send request.")
    return false
  }
  if res.StatusCode != expected {
    fmt.Printf("StatusCode wasn't correct.")
    return false
  }
  return true
}

func DoTest1(method, url, json string, expected int) bool {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    fmt.Printf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    fmt.Printf("Couldn't make token for group.")
    return false
  }
  reader := strings.NewReader(json)
  request, err := http.NewRequest(method, url, reader)
  request.Header.Set("access_token", token)

  res, errSending := http.DefaultClient.Do(request)
  if errSending != nil {
    fmt.Printf("Couldn't send request.")
    return false
  }
  if res.StatusCode != expected {
    fmt.Printf("StatusCode wasn't correct.")
    return false
  }
  return true
}

func DoTest2(method, url, json, username string, expected int) bool {
  id := db.GetIdFromUsername(username)
  if id == "" {
    fmt.Printf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    fmt.Printf("Couldn't make token for group.")
    return false
  }
  reader := strings.NewReader(json)
  request, err := http.NewRequest(method, url, reader)
  request.Header.Set("access_token", token)

  res, errSending := http.DefaultClient.Do(request)
  if errSending != nil {
    fmt.Printf("Couldn't send request.")
    return false
  }
  if res.StatusCode != expected {
    fmt.Printf("StatusCode wasn't correct.")
    return false
  }
  return true
}

func DoSimpleTest(method, url string, expected int) bool {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    fmt.Printf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    fmt.Printf("Couldn't make token for group.")
    return false
  }
  request, err := http.NewRequest(method, url, nil)
  request.Header.Set("access_token", token)

  res, errSending := http.DefaultClient.Do(request)
  if errSending != nil {
    fmt.Printf("Couldn't send request.")
    return false
  }

  if res.StatusCode != expected {
    fmt.Printf("StatusCode wasn't correct.")
    return false
  }
  return true
}

func TestCreateUser(t *testing.T) {
  json := `{"username":"test","email":"dkraken@thekrakenisgongetu.com","password":"testtest1"}`
  reader := strings.NewReader(json)
  request, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/register", server.URL), reader)
  if err != nil {
    t.Errorf("Problem setting up request.")
  }
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Error(err1)
  }
  //if user was already made, don't worry about it.
  if res.StatusCode != 204 && res.StatusCode != 500 {
      t.Errorf("Success expected: %d", res.StatusCode)
  }
}

func TestActivateUser(t *testing.T) {
  id := db.GetIdFromUsername("test")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/auth/activate", server.URL)
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token", token)

  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Activating Users")
  }
}

func TestCreateUser2(t *testing.T) {
  json := `{"username":"test1","email":"dkraken@thekrakenisgongetu1.com","password":"testtest2"}`
  reader := strings.NewReader(json)
  request, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/register", server.URL), reader)
  if err != nil {
    t.Errorf("Problem setting up request.")
  }
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 204 && res.StatusCode != 500 {
       t.Errorf("Success expected: %d", res.StatusCode)
  }
}

func TestActivateUser2(t *testing.T) {

  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/auth/activate", server.URL)
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token", token)

  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Activating Users")
  }
}

func TestLogin(t *testing.T) {
  json := `{"email":"dkraken@thekrakenisgongetu.com","password":"testtest1"}`
  reader := strings.NewReader(json)
  request, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/login", server.URL), reader)
  if err != nil {
    t.Errorf("Problem setting up request for login.")
  }
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 200 {
       t.Errorf("Couldn't do login")
  }
}

func TestGetUser(t *testing.T) {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/user/", server.URL)
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token", token)

  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 200 {
    t.Errorf("Unfortunate Statuscode for Getting Users")
  }
}

func TestGetUserSaved(t *testing.T) {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/user/saved", server.URL)
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token", token)

  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 200 {
    t.Errorf("Unfortunate Statuscode for Getting Users")
  }
}

func TestMakeGroup(t *testing.T) {
  id := db.GetIdFromUsername("test")
  if id == "" {
    t.Errorf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token for group.")
  }
  json := `{"anonymous":false,"group":"test"}`
  reader := strings.NewReader(json)
  request, err2 := http.NewRequest("POST", fmt.Sprintf("%s/group/modify", server.URL), reader)
  if err2 != nil {
    t.Errorf("Problem setting up request for login.")
  }
  request.Header.Set("access_token", token)
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 204 {
       t.Errorf("Couldn't make group...")
  }
}

func TestAdmin(t *testing.T) {
  id := db.GetIdFromUsername("test")
  if id == "" {
    t.Errorf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token for group.")
  }
  //get admin level /group/auth - models.Grp
  json := `{"group":"test"}`
  reader := strings.NewReader(json)
  request, err2 := http.NewRequest("POST", fmt.Sprintf("%s/group/auth", server.URL), reader)
  if err2 != nil {
    t.Errorf("Problem setting up permissions request.")
  }
  request.Header.Set("access_token", token)
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 200 {
       t.Errorf("Couldn't make group...")
  }

  otherId := db.GetIdFromUsername("test1")

  //add member to group
  json2 := fmt.Sprintf(`{"group":"test", "user":"%s"}`, otherId)
  reader2 := strings.NewReader(json2)
  request1, err3 := http.NewRequest("POST", fmt.Sprintf("%s/group/admin", server.URL), reader2)
  if err3 != nil {
    t.Errorf("Problem setting up request for add member")
  }
  request1.Header.Set("access_token", token)
  res2, err4 := http.DefaultClient.Do(request1)
  if err4 != nil {
      t.Errorf("Couldn't add member -- send request I mean.")
  }
  if res2.StatusCode != 204 {
       t.Errorf("Couldn't add member to group")
  }
  reader3 := strings.NewReader(json2)
  //remove member from group
  requestAgain, errAgain := http.NewRequest("PUT", fmt.Sprintf("%s/group/admin", server.URL), reader3)
  if errAgain != nil {
    t.Errorf("Problem setting up request for add member")
  }
  requestAgain.Header.Set("access_token", token)

  resAgain, errAgain := http.DefaultClient.Do(requestAgain)
  if errAgain != nil {
    t.Errorf("Couldn't add member -- send request I mean.")
  }
  if resAgain.StatusCode != 204 {
    t.Errorf("Couldn't remove member from group")
  }

  //get group -- good stuff
  sumRareJson := `{"group":"test","page":0}`
  readerNew := strings.NewReader(sumRareJson)
  rerere, erere := http.NewRequest("POST", fmt.Sprintf("%s/group/", server.URL), readerNew)
  if erere != nil {
    t.Errorf("Problem getting group.")
  }
  rerere.Header.Set("access_token", token)

  newRes, newErr := http.DefaultClient.Do(rerere)
  if newErr != nil {
    t.Errorf("Couldn't make request getting group")
  }
  if newRes.StatusCode != 200 {
    t.Errorf("Problem getting group -- statuscode not 200")
  }
}

func TestGetUserThreads(t *testing.T) {
  json := `{"page": 0}`
  if !DoTest("POST", fmt.Sprintf("%s/user/threads", server.URL), json, 200) {
    t.Errorf("Couldn't get user threads")
  }
}

func TestFriends(t *testing.T) {
  id := db.GetIdFromUsername("test")
  id1 := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  if id1 == "" {
    t.Errorf("Couldn't find username")
  }
  json := `{"username":"test", "friend":"test1"}`
  if !DoTest("POST", fmt.Sprintf("%s/user/friends", server.URL), json, 204) {
    t.Errorf("Couldn't add user friend")
  }
  json1 := `{"username":"test1", "friend":"test"}`
  if !DoTest1("PUT", fmt.Sprintf("%s/user/friends", server.URL), json1, 204) {
    t.Errorf("Couldn't accept friends")
  }
  json2 := `{"username":"test", "friend":"test1"}`
  if !DoTest("DELETE", fmt.Sprintf("%s/user/friends", server.URL), json2, 204) {
    t.Errorf("Couldn't delete friend")
  }
}

func TestMakeThread(t *testing.T) {
  json := `{"group":"test","body":"hello","author":"test","content":"linkhere","anonymous":false}`
  if !DoTest("POST", fmt.Sprintf("%s/thread/modify", server.URL), json, 200) {
    t.Errorf("Make thread is messed up.")
  }
}

//can't actually get thread because we don't have ID of thread -- so let's get it
func TestGetThread(t *testing.T) {

  //get thread ID because we need it to move on
  threads, err := db.GetGroup("test", 0)
  if err != nil {
    t.Errorf("couldn't actually get the group")
  }

  //shouldn't throw an error but it might if something went wrong beforehand
  threadId := threads[0].Thread

  //format string to get the thread
  json := fmt.Sprintf(`{"thread":"%s"}`, threadId.Hex())
  if !DoTest("POST", fmt.Sprintf("%s/thread/", server.URL), json, 200) {
    t.Error("Get thread is messed up")
  }

  if !DoTest("POST", fmt.Sprintf("%s/user/saved", server.URL), json, 204) {
    t.Error("Couldn't save thread")
  }

  if !DoTest("PUT", fmt.Sprintf("%s/user/saved", server.URL), json, 204) {
    t.Error("Couldn't unsave thread")
  }
}

//test name endpoint
func TestName(t *testing.T) {
  json := `{"name": "Octopus Is Animal"}`
  if !DoTest("POST", fmt.Sprintf("%s/user/name", server.URL), json, 204) {
    t.Errorf("Couldn't change name")
  }
}

//test usernames
func TestUsername(t *testing.T) {
  json := `{"username": "dingo"}`

  if !DoTest1("POST", fmt.Sprintf("%s/user/username", server.URL), json, 204) {
    t.Errorf("Couldn't add username")
  }

  if !DoTest1("PUT", fmt.Sprintf("%s/user/username", server.URL), json, 204) {
    t.Errorf("Couldn't change username")
  }

  json1 := `{"username": "test1"}`
  if !DoTest2("PUT", fmt.Sprintf("%s/user/username", server.URL), json1, "dingo", 204) {
    t.Errorf("Couldn't change username")
  }

  if !DoTest1("DELETE", fmt.Sprintf("%s/user/username", server.URL), json, 204) {
    t.Errorf("Couldn't rm username")
  }
}

//can't actually get thread because we don't have ID of thread -- so let's get it
func TestPost(t *testing.T) {

  //get thread ID because we need it to move on
  threads, err := db.GetGroup("test", 0)
  if err != nil {
    t.Errorf("couldn't actually get the group")
  }

  //shouldn't throw an error but it might if something went wrong beforehand
  threadId := threads[0].Thread

  //format string to get the thread
  json := fmt.Sprintf(`{"thread":"%s", "body":"This is just a test", "content": "link",
    "responseTo":[],"anonymous":false}`, threadId.Hex())
  if !DoTest("POST", fmt.Sprintf("%s/thread/post", server.URL), json, 200) {
    t.Error("Posting is messed up.")
  }

  postId := threads[0].Head.Id
  //another one -- edits posts
  morejson := fmt.Sprintf(`{"post":"%s", "body":"This is just a test"}`, postId.Hex())
  if !DoTest("PUT", fmt.Sprintf("%s/thread/post", server.URL), morejson, 204) {
    t.Error("Edit thread is messed up")
  }

  //delete that very post
  somejson := fmt.Sprintf(`{"post":"%s"}`, postId.Hex())
  if !DoTest("DELETE", fmt.Sprintf("%s/thread/post", server.URL), somejson, 204) {
    t.Error("Delete thread is messed up")
  }
}

//test search threads
func TestSearch(t *testing.T) {
  //format string to get the thread
  json :=`{"text":"e", "page":0}`
  if !DoTest("POST", fmt.Sprintf("%s/thread/search", server.URL), json, 200) {
    t.Error("Search is messed up.")
  }
}

func TestSearchUser(t *testing.T) {
  //format string to get the thread
  json :=`{"text":"test"}`
  if !DoTest("POST", fmt.Sprintf("%s/user/search", server.URL), json, 200) {
    t.Error("Search User is messed up.")
  }
}

func TestSearchGroup(t *testing.T) {
  //format string to get the thread
  json :=`{"text":"test"}`
  if !DoTest("POST", fmt.Sprintf("%s/group/search", server.URL), json, 200) {
    t.Error("Search Group is messed up.")
  }
}

//simple gets
func TestGets(t *testing.T) {
  if !DoSimpleTest("GET", fmt.Sprintf("%s/user/notifications", server.URL), 200) {
    t.Error("broken notifictions!")
  }
  if !DoSimpleTest("GET", fmt.Sprintf("%s/user/friends", server.URL), 200) {
    t.Error("We couldn't even get friends")
  }
}

func TestRemoveThread(t *testing.T) {
  //get thread ID because we need it to move on
  threads, err := db.GetGroup("test", 0)
  if err != nil {
    t.Errorf("couldn't actually get the group")
  }

  //shouldn't throw an error but it might if something went wrong beforehand
  threadId := threads[0].Thread
  json := fmt.Sprintf(`{"thread":"%s"}`, threadId.Hex())
  if !DoTest("DELETE", fmt.Sprintf("%s/thread/modify", server.URL), json, 204) {
    t.Error("Delete thread is messed up")
  }
}

func TestRemoveGroup(t *testing.T) {
  json := `{"group":"test"}`
  if !DoTest("DELETE", fmt.Sprintf("%s/group/modify", server.URL), json, 204) {
    t.Error("Delete thread is messed up")
  }
}


func TestDectivateUser(t *testing.T) {
  id := db.GetIdFromUsername("test")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/auth/remove", server.URL)
  request, err := http.NewRequest("DELETE", url, nil)
  request.Header.Set("access_token", token)

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Deactivating Users")
  }
}

func TestDectivateUser2(t *testing.T) {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token")
  }
  url := fmt.Sprintf("%s/auth/remove", server.URL)
  request, err := http.NewRequest("DELETE", url, nil)
  request.Header.Set("access_token", token)

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Decativating Users")
  }
}

func TestDeleteUsers(t *testing.T){
  //cleanup groups
  id := db.GetIdFromUsername("test")
  id1 := db.GetIdFromUsername("test1")
  db.RemoveUser(bson.ObjectIdHex(id))
  db.RemoveUser(bson.ObjectIdHex(id1))
}
