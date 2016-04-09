/*
   API Routes that will be used by the source of the application.
   Creates a ServeMux from the default http package
*/
package routes

import (
    "github.com/ahermida/sudopostAPI/api/Util"
    "github.com/ahermida/sudopostAPI/api/Models"
    "net/http"
    "encoding/json"
    "gopkg.in/mgo.v2/bson"
    "github.com/ahermida/sudopostAPI/api/DB"
)

// Routes with /group/ prefix
var GroupMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //Handles POST to get group -- paginated
  GroupMux.HandleFunc("/group/", getGroup)

  //Handles POST to get permission
  GroupMux.HandleFunc("/group/auth", getPermission)

  //Handles POST, and DELETE for groups -- administration (creating and deleting)
  GroupMux.HandleFunc("/group/modify", grp)

  //Handles POST, PUT for Admins in groups
  GroupMux.HandleFunc("/group/admin", admn)

  //Handles POST for searching group names
  GroupMux.HandleFunc("/group/search", searchGroups)
}

// Handle /group/modify
func grp(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST":
    createGroup(res, req)
  case "DELETE":
    removeGroup(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

// Handle /group/admin
func admn(res http.ResponseWriter, req *http.Request) {
  switch req.Method {
  case "POST" :
    addAdmin(res, req)
  case "PUT":
    removeAdmin(res, req)
  default:
    http.Error(res, http.StatusText(405), 405)
  }
}

//handle POST to /group/
func getGroup(res http.ResponseWriter, req *http.Request) {

  //only handle POST
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }

  //data returned by request
  var reqGroup models.GetGroup

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&reqGroup); err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //user _id in hex
  id := util.GetId(req)

  //check if private group, if so, we need to do some auth
  if !db.IsMember(reqGroup.Group, id) {

    //if it's private and we don't have permission, throw back 401 -- unauthorized
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //else, resolve thread
  group, errFromQuery := db.GetGroup(reqGroup.Group, reqGroup.Page)

  //check for query issues
  if errFromQuery != nil {

    //if there's an error in getting the group, return a 500
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //get json struct that we're gonna send over
  grp := &models.SendGroup{
    Threads: group,
  }

  //send back no error response
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  //send over data
  if err := json.NewEncoder(res).Encode(grp); err != nil {
    panic(err)
  }
}

//handle POST /group/modify
func createGroup(res http.ResponseWriter, req *http.Request) {

  //data returned by request
  var reqGroup models.CreateGroup

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&reqGroup); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //validate input from post -- particularly the group name -- if fails, return err code
  if !validateGeneral(reqGroup.Group){
    http.Error(res, http.StatusText(422), 422)
    return
  }

  //get token from header
  userToken := req.Header.Get("access_token")

  //if user has no token
  if userToken == "" {

    //return err
    http.Error(res, http.StatusText(401), 401)
    return
  } else {

    //check token that we got
    token, err := util.CheckToken(userToken)

    //if nothing went wrong
    if err == nil && token.Valid {

      //success -- validate account & continue
      errGrp := db.CreateGroup(reqGroup.Group, bson.ObjectIdHex(token.Claims["id"].(string)), reqGroup.Private)

      //check if something went wrong
      if errGrp != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 204 response code -- No Content
      res.WriteHeader(http.StatusNoContent)
    } else {

      //return err
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

//handle DELETE /group/
func removeGroup(res http.ResponseWriter, req *http.Request) {
  //data returned by request
  var reqGroup models.Grp

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&reqGroup); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get token from header
  userToken := req.Header.Get("access_token")

  //if user has no token
  if userToken == "" {

    //return err
    http.Error(res, http.StatusText(401), 401)
    return
  } else {

    //check token that we got
    token, err := util.CheckToken(userToken)

    //if nothing went wrong
    if err == nil && token.Valid {

      //success -- validate account & continue
      errGrp := db.DeleteGroup(reqGroup.Group, bson.ObjectIdHex(token.Claims["id"].(string)))

      //check if something went wrong
      if errGrp != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 204 response code -- No Content
      res.WriteHeader(http.StatusNoContent)
    } else {

      //return err if invalid token
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

//handle POST /group/auth
func getPermission(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //manage post
  var post models.Grp

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&post); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get token from header
  id := util.GetId(req)

  //get permission
  permission := db.GetPermission(post.Group, id)

  if permission == nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //set headers
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)

  //send over data
  if err := json.NewEncoder(res).Encode(permission); err != nil {
    panic(err)
  }
}

//handle POST /group/admin
func addAdmin(res http.ResponseWriter, req *http.Request) {

  //data returned by request
  var reqAdmin models.GroupAdmin

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&reqAdmin); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get token from header
  userToken := req.Header.Get("access_token")

  //if user has no token
  if userToken == "" {

    //return err
    http.Error(res, http.StatusText(401), 401)
    return
  } else {

    //check token that we got
    token, err := util.CheckToken(userToken)

    //if nothing went wrong
    if err == nil && token.Valid {

      //get Ids for both users
      newbie := bson.ObjectIdHex(reqAdmin.User)
      user := bson.ObjectIdHex(token.Claims["id"].(string))

      //success -- validate account & continue
      errGrp := db.AddAdmin(user, newbie, reqAdmin.Group)

      //check if something went wrong
      if errGrp != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 204 response code -- No Content
      res.WriteHeader(http.StatusNoContent)
    } else {

      //return err if invalid token
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

//handle PUT /group/admin
func removeAdmin(res http.ResponseWriter, req *http.Request) {
  //data returned by request
  var reqAdmin models.GroupAdmin

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&reqAdmin); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //get token from header
  userToken := req.Header.Get("access_token")

  //if user has no token
  if userToken == "" {

    //return err
    http.Error(res, http.StatusText(401), 401)
    return
  } else {

    //check token that we got
    token, err := util.CheckToken(userToken)

    //if nothing went wrong
    if err == nil && token.Valid {

      //get Ids for both users
      newbie := bson.ObjectIdHex(reqAdmin.User)
      user := bson.ObjectIdHex(token.Claims["id"].(string))

      //success -- validate account & continue
      errGrp := db.RemoveAdmin(user, newbie, reqAdmin.Group)

      //check if something went wrong
      if errGrp != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 204 response code -- No Content
      res.WriteHeader(http.StatusNoContent)
    } else {

      //return err if invalid token
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

//handle POST to /group/search
func searchGroups(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
  }

  //user _id in hex
  id := util.GetId(req)

  //users only searchable if we're logged in
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //data returned by request
  var request models.Search

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&request); err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  groups, err := db.SearchGroups(request.Text)
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
  }

  send := &models.SendGroupSearch{
    Groups: groups,
  }

  //else send back user which is already json formatted
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(res).Encode(send); err != nil {
    panic(err)
  }
}
