/*
   User Routes that will be used to manage basic user management utils
*/
package routes

import (
    "net/http"
    "github.com/ahermida/sudopostAPI/api/DB"
    "github.com/ahermida/sudopostAPI/api/Models"
    "github.com/ahermida/sudopostAPI/api/Config"
    "github.com/ahermida/sudopostAPI/api/Util"
    "gopkg.in/mgo.v2/bson"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
)

// Routes with /users/ prefix
var AuthMux = http.NewServeMux()

// Setup Routes with Mux
func init() {

  //POST register user
  AuthMux.HandleFunc("/auth/register", register)

  //POST activate user
  AuthMux.HandleFunc("/auth/activate", activate)

  //DELETE deactivate user
  AuthMux.HandleFunc("/auth/remove", deactivate)

  //POST login user
  AuthMux.HandleFunc("/auth/login", login)

  //POST password recovery
  AuthMux.HandleFunc("/auth/recovery", recovery)

  //POST update password
  AuthMux.HandleFunc("/auth/changepass", updatePassword)

}

/*
   Route handlers for User Routes
*/

// Handle /user/register
func register(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  //make room for user
  var usr models.CreateUser

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&usr); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //validate forms
  isValid := validateEmail(usr.Email) && validateGeneral(usr.Username) && validatePassword(usr.Password)
  if !isValid {
    http.Error(res, http.StatusText(422), 422)
    return
  }

  //hash pw
  key := []byte(config.Secret)
  hasher := hmac.New(sha256.New, key)
  hasher.Write([]byte(usr.Password))
  hashword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

  //create user since forms are valid
  id, err := db.CreateUser(usr.Email, usr.Username, hashword);

  //check if something went wrong
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //makes token with id hex baked in
  token, errToken := util.MakeToken(id)
  if errToken != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //send user validation email -- this should be done in a goroutine
  go setupEmail(usr.Email, token)

  //statuscode 204
  res.WriteHeader(http.StatusNoContent)
}

// Handle /user/register
func recovery(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }
  //make room for user
  var rec models.Recovery

  //POST request handling
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&rec); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //create user since forms are valid
  id := db.GetIdFromEmail(rec.Email);

  //don't continue if fake email
  if id == "" {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //makes token with id hex baked in
  token, err := util.MakeToken(id)
  if err != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //send user validation email
  go recoverEmail(rec.Email, token)

  //statuscode 204
  res.WriteHeader(http.StatusCreated)
}

// Handle /user/activate
func activate(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    http.Error(res, http.StatusText(405), 405)
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
      errAuth := db.ActivateAccount(bson.ObjectIdHex(token.Claims["id"].(string)))

      //check if something went wrong
      if errAuth != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 202 response code -- No Content
      res.WriteHeader(http.StatusAccepted)
    } else {

      //return err
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

// Handle /user/deactivate
func deactivate(res http.ResponseWriter, req *http.Request) {
  if req.Method != "DELETE" {
    http.Error(res, http.StatusText(405), 405)
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
      errAuth := db.DeleteUser(bson.ObjectIdHex(token.Claims["id"].(string)))

      //check if something went wrong
      if errAuth != nil {

        //return err
        http.Error(res, http.StatusText(401), 401)
        return
      }

      //send the headers with a 202 response code -- No Content
      res.WriteHeader(http.StatusAccepted)
    } else {

      //return err
      http.Error(res, http.StatusText(401), 401)
    }
  }
}

// Handle /user/login
func login(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }

  //setup struct to recieve json
  var user models.AuthUser
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&user); err != nil {
    http.Error(res, http.StatusText(400), 400)
    return
  }

  //hash pw
  key := []byte(config.Secret)
  hasher := hmac.New(sha256.New, key)
  hasher.Write([]byte(user.Password))
  hashword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

  id, match := db.LoginCheck(user.Email, hashword)

  //if login failed, return error
  if !match {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //get token -- hex
  token, errToken := util.MakeToken(id)
  //check if something went wrong with token
  if errToken != nil {
    http.Error(res, http.StatusText(500), 500)
    return
  }

  //in the other case that everything went as planned
  result := &models.AuthedUser{
    Token: token,
  }

  //send back no error response
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")
  res.WriteHeader(http.StatusOK)

  //send back Token in JSON
  if err := json.NewEncoder(res).Encode(result); err != nil {
    panic(err)
  }
}

// Handle /user/updatepassword
func updatePassword(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.Error(res, http.StatusText(405), 405)
    return
  }

  //get token from header
  id := util.GetId(req)

  //check if the token is legit
  if id == "" {
    http.Error(res, http.StatusText(401), 401)
  }

  //check if password is the same
  var usr models.ChangePW

  //decode
  decoder := json.NewDecoder(req.Body)
  if err := decoder.Decode(&usr); err != nil {
    http.Error(res, http.StatusText(400), 400)
  }

  //setup key for hashing
  key := []byte(config.Secret)

  //hash old pw
  hasher := hmac.New(sha256.New, key)
  hasher.Write([]byte(usr.Password))
  hashword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

  //hash new password
  newHasher := hmac.New(sha256.New, key)
  newHasher.Write([]byte(usr.NewPassword))
  newHashword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

  //call change password
  errChange := db.ChangePassword(newHashword, hashword, bson.ObjectIdHex(id))
  if errChange != nil {
    http.Error(res, http.StatusText(401), 401)
    return
  }

  //no probs, so sent back no content
  res.WriteHeader(http.StatusNoContent)
}
