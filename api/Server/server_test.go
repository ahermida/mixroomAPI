/*
  Test API -- Integration test
*/
package server

/*
  Endpoints:
  (31 of them)
                                                                    STRUCT INPUT
  <Auth> ------------------------------------------------------
  xPOST -- register user ("/auth/register", register)              - models.CreateUser
  xPOST -- recover user pw ("/auth/recovery", recovery)            - models.Recovery
  xGET -- activate user ("/auth/activate", activate)               - [id required]
  xDELETE -- deactivate user ("/auth/remove", deactivate)          - [id required]
  xPOST -- login user ("/auth/login", login)                       - models.AuthUser
  xPOST -- update password ("/auth/changepass", updatePassword)    - models.ChangePW
  </Auth>

  <User> ------------------------------------------------------
  xGET get user info ("/user/", getUser)                           - [id required]
  xGET get saved ("/user/saved", saved)                            - [id required]
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
  xPOST -- to get group -- paginated ("/group/", getGroup)         - models.GetGroup
  xPOST -- check if admin or author ("/group/auth", getPermission) - models.Grp
  xPOST -- for groups -- creating them ("/group/modify", grp)      - models.CreateGroup
  DELETE -- for groups -- deleting ("/group/modify", grp)         - models.Grp
  xPOST -- set admins for group ("/group/admin", admn)             - models.GroupAdmin
  xPUT -- delete Admins in groups ("/group/admin", admn)           - models.GroupAdmin
  </Groups>

  <Threads> ---------------------------------------------------
  xPOST -- Get Thread ("/thread/", getThread)                      - models.GetThread
  xPOST -- Create Thread ("/thread/modify", thrd)                  - models.NewThread
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
  "strings"
  "fmt"
)

//setup server so we can test endpoints
var server *httptest.Server

func init() {

  //grab mux from server.go and run it
  server = httptest.NewServer(Server)
}

func TestCreateUser(t *testing.T) {
  json := `{"username":"test","email":"dkraken@thekrakenisgongetu.com","password":"testtest1"}`
  reader := strings.NewReader(json)
  request, err := http.NewRequest("POST", "http://localhost:8000/auth/register", reader)
  if err != nil {
    t.Errorf("Problem setting up request.")
  }
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 201 {
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
  url := "http://localhost:8000/auth/activate"
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token": token)

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
  request, err := http.NewRequest("POST", "http://localhost:8000/auth/register", reader)
  if err != nil {
    t.Errorf("Problem setting up request.")
  }
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 201 {
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
  url := "http://localhost:8000/auth/activate"
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token": token)

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
  request, err := http.NewRequest("POST", "http://localhost:8000/auth/login", reader)
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
  url := "http://localhost:8000/user/"
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token": token)

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
  url := "http://localhost:8000/user/saved"
  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Set("access_token": token)

  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 200 {
    t.Errorf("Unfortunate Statuscode for Getting Users")
  }
}

func TestMakeGroup(t *testing.T) {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token for group.")
  }
  json := `{"anonymous":false,"group":"test"}`
  reader := strings.NewReader(json)
  request, err2 := http.NewRequest("POST", "http://localhost:8000/group/modify", reader)
  if err2 != nil {
    t.Errorf("Problem setting up request for login.")
  }
  request.Header.Set("access_token": token)
  res, err1 := http.DefaultClient.Do(request)
  if err1 != nil {
      t.Errorf("Couldn't send request")
  }
  if res.StatusCode != 204 {
       t.Errorf("Couldn't make group...")
  }
}

func TestAdmin(t *testing.T) {
  id := db.GetIdFromUsername("test1")
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
  request, err2 := http.NewRequest("POST", "http://localhost:8000/group/auth", reader)
  if err2 != nil {
    t.Errorf("Problem setting up permissions request.")
  }
  request.Header.Set("access_token": token)
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
  request1, err3 := http.NewRequest("POST", "http://localhost:8000/group/modify", reader2)
  if err3 != nil {
    t.Errorf("Problem setting up request for add member")
  }
  request1.Header.Set("access_token": token)
  res2, err4 := http.DefaultClient.Do(request1)
  if err4 != nil {
      t.Errorf("Couldn't add member -- send request I mean.")
  }
  if res2.StatusCode != 204 {
       t.Errorf("Couldn't add member to group")
  }

  //remove member from group
  requestAgain, errAgain := http.NewRequest("PUT", "http://localhost:8000/group/modify", reader2)
  if errAgain != nil {
    t.Errorf("Problem setting up request for add member")
  }
  requestAgain.Header.Set("access_token": token)

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
  rerere, erere := http.NewRequest("POST", "http://localhost:8000/group/", readerNew)
  if erere != nil {
    t.Errorf("Problem getting group.")
  }
  rerere.Header.Set("access_token": token)

  newRes, newErr := http.DefaultClient.Do(rerere)
  if newErr != nil {
    t.Errorf("Couldn't make request getting group")
  }
  if newRes.StatusCode != 200 {
    t.Errorf("Problem getting group -- statuscode not 200")
  }
}
/*
  Let it be known that at this point I gave and decided to write a function
  to do this json thing automatically -- feeling stupid af for not doing this
  earlier...
*/
func DoTest(method, url, json string, expected int) bool {
  id := db.GetIdFromUsername("test1")
  if id == "" {
    t.Errorf("Couldn't find username for Group Creation")
  }
  token, err := util.MakeToken(id)
  if err != nil {
    t.Errorf("Couldn't make token for group.")
    return false
  }
  reader := strings.NewReader(json)
  request, err := http.NewRequest(method, url, reader)
  request.Header.Set("access_token", token)

  res, errSending := http.DefaultClient.Do(request)
  if errSending != nil {
    t.Error("Couldn't send request.")
    return false
  }
  if res.StatusCode != expected {
    t.Errorf("StatusCode wasn't correct.")
    return false
  }
  return true
}
func TestMakeThread(t *testing.T) {
  // type NewThread struct {
  //   Group     string `json:"group"`
  //   Body      string `json:"body"`
  //   Author    string `json:"author"`
  //   Content   string `json:"content"`
  //   Anonymous bool   `json:"anonymous"`
  // }
  json := `{"group":"test","body":"hello","author":"test","content":"linkhere",anonymous:false}`
  if !DoTest("POST", "http://localhost:8000/thread/modify", json, 200) {
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
  if !DoTest("POST", "http://localhost:8000/thread/", json, 200) {
    t.Error("Get thread is messed up")
  }
}

// type NewPost struct {
//   Body       string   `json:"body"`
//   Content    string   `json:"content"`
//   Author     string   `json:"author"`
//   ResponseTo []string `json:"responseTo"`
//   Anonymous  bool     `json:"anonymous"`
//   Thread     string   `json:"thread"`
// }
//
// type EditPost struct {
//   Body string `json:"body"`
//   Post string `json:"post"`
//   Id   string `json:"id,omitempty"`
// }
//
// type DeletePost struct {
//   Post string `json:"post"`
//   Id   string `json:"id,omitempty"`
// }

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
  if !DoTest("POST", "http://localhost:8000/thread/post", json, 200) {
    t.Error("Posting is messed up.")
  }

  postId := threads[0].Head.Id
  //another one -- edits posts
  morejson := fmt.Sprintf(`{"post":"%s", "body":"This is just a test", }`, threadId.Hex())
  if !DoTest("PUT", "http://localhost:8000/thread/post", morejson, 200) {
    t.Error("Get thread is messed up")
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
  url := "http://localhost:8000/auth/remove"
  request, err := http.NewRequest("DELETE", url, nil)
  request.Header.Set("access_token": token)

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Activating Users")
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
  url := "http://localhost:8000/auth/remove"
  request, err := http.NewRequest("DELETE", url, nil)
  request.Header.Set("access_token": token)

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error("Couldn't send request") //Something is wrong while sending request
  }
  if res.StatusCode != 202 {
    t.Errorf("Unfortunate Statuscode for Activating Users")
  }
}
