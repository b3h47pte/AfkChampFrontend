/*
 * 'Login' Provides basic login functionality that can be used across the site.
 */
package controller

import(
  "net/http"
  "code.google.com/p/gcfg"
  "encoding/hex"
  "log"
  "github.com/gorilla/sessions"
  "AfkChampFrontend/model/user"
  "code.google.com/p/go-uuid/uuid"
  "time"
  "AfkChampFrontend/utility"
)

type LoginRegisterPostData struct {
  Username string
  Password string
  Email string
}

type LoginRegisterErrorCode int
const (
  ErrorNoError LoginRegisterErrorCode = iota
  ErrorInvalidRegisterUserName
  ErrorInvalidRegisterPassword
  ErrorInvalidRegisterEmail
  ErrorInvalidLoginCredentials
  ErrorAlreadyLoggedIn
  ErrorUnspecifiedError
)

type LoginRegisterResponse struct {
  ErrorCode LoginRegisterErrorCode
  RedirectUrl string
}

type LoginConfig struct {
  AuthSection struct {
    AuthKey string
  }
}

type LoginTemplateData struct {
  Data BaseTemplateData
  MinimumPasswordLength uint8
  MaximumEmailLength uint8
}

var AuthSessionKey []byte
var LoginStore *sessions.CookieStore

// Generate default login template data
func CreateLoginTemplateData() *LoginTemplateData {
  t := LoginTemplateData{Data: CreateTemplateData(),
    MinimumPasswordLength: user.MinPasswordLength,
    MaximumEmailLength: user.MaxEmailLength}
  return &t
}

// Setup the secure session for the a user's login session.
func init() {
  var config LoginConfig
  err := gcfg.ReadFileInto(&config, "config/login.config")
  if err != nil {
    log.Fatal(err)
  }
  hexAuthKey := config.AuthSection.AuthKey
  AuthSessionKey, err = hex.DecodeString(hexAuthKey)
  if err != nil {
    log.Fatal(err)
  }
  LoginStore = sessions.NewCookieStore(AuthSessionKey)
}

// HandleLoginPageRoute displays a login page if the user is not logged in. Otherwise,
// the user is redirected to the home page.
func HandleLoginPageRoute(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }

  t := CreateLoginTemplateData()
  TemplateMapping["login/login.html"].ExecuteTemplate(w, "tbase", t)
}

// HandleLogoutPageRoute logs out the current user and directs him/her back to the home page.
func HandleLogoutPageRoute(w http.ResponseWriter, r *http.Request) {
  if sessionKey, err := GetUserSession(w, r); err == nil {
    RemoveUserSession(sessionKey, w, r)
  }
  http.Redirect(w, r, "/", http.StatusFound)
}


// HandlRegisterPageRoute displays the registration page if the user is not logged in. Otherwise, the 
// user is redirected to the home page.
func HandleRegisterPageRoute(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }
  t := CreateLoginTemplateData()
  TemplateMapping["login/register.html"].ExecuteTemplate(w, "tbase", t)
}

// HandleLoginAction takes in the user's name and password and checks whether or not they are registered. Sets 
// relevant information in the cookie store to remember the user's session. This fails if the user is logged in already.
func HandleLoginAction(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    LoginRegisterRespondJsonError(ErrorAlreadyLoggedIn, "/", w)
    return
  }
  
  userData := LoginRegisterPostData{}
  err := utility.ReadJsonFromRequestBodyStruct(r, &userData)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
 
  newUser, err := user.VerifyUser(userData.Username, userData.Password)
  if err != nil {
    LoginRegisterRespondJsonError(ErrorInvalidLoginCredentials, "", w)
    return
  }
  
  err = CreateUserSession(newUser, w, r)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "/login", w)
    return
  }
  // Success!
  LoginRegisterRespondJsonError(ErrorNoError, "/", w)
}

// HandleRegisterAction allows you to register a new user given a username and password. If registration is successful, also
// set relevant information in the user's cookies to remember their session. This fails if the user is logged in already.
func HandleRegisterAction(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    LoginRegisterRespondJsonError(ErrorAlreadyLoggedIn, "/", w)
    return
  }
  
  userData := LoginRegisterPostData{}
  err := utility.ReadJsonFromRequestBodyStruct(r, &userData)
  if err != nil {
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
  
  newUser, err := user.CreateUser(userData.Username, userData.Password, userData.Email)
  switch {
  case err == user.UserExistsError:
    // Case where we need to inform user
    LoginRegisterRespondJsonError(ErrorInvalidRegisterUserName, "", w)
    return
  case err != nil:
    // Just keep a mental note to ourself but display another error to the user
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
  err = CreateUserSession(newUser, w, r)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
  // Assume success and redirect to front page
  LoginRegisterRespondJsonError(ErrorNoError, "/", w)
}

// 'GetCurrentUser' will retrieve the currently logged in user.
func GetCurrentUser(w http.ResponseWriter, r *http.Request) (*user.UserEntry, error) {
  sessionKey, err := GetUserSession(w, r)
  if err != nil {
    return nil, err 
  }
  userId, err := user.VerifySession(sessionKey)
  if err != nil {
    return nil, err
  }
  return user.RetrieveUser(userId)
}

// 'GetUserSession' gets the current user's session (if the user has logged in before hand and we stored the session cookie).
func GetUserSession(w http.ResponseWriter, r *http.Request) (string, error) {
  session, err := LoginStore.Get(r, "user-session")
  if err != nil {
    return "", err
  }
  sessionKey, ok := session.Values["key"]
  if !ok {
    return "", user.NoSessionError
  }
  return sessionKey.(string), nil
}

// 'CreateUserSession' takes in a given user and creates a new session cookie for the user. This function
// will clear any session cookie already set.
func CreateUserSession(newUser *user.UserEntry, w http.ResponseWriter, r *http.Request) error {
  session, err := LoginStore.Get(r, "user-session")
  if err != nil {
    return err
  }
  sessionKey, ok := session.Values["key"]
  // Remove session key if one exists already
  if ok {
    sessionKeyStr := sessionKey.(string)
    err = RemoveUserSession(sessionKeyStr, w, r)
    if err != nil {
      return err
    }
  }
  
  // Generate a new session key.
  newSessionKey := uuid.New()
  session.Values["key"] = newSessionKey
  err = session.Save(r, w)
  if err != nil {
    return err
  }
  
  // Associate key with the user for 3 months.
  expirationTime := time.Now().AddDate(0,3,0)
  err = user.AddSessionForUser(newUser, newSessionKey, &expirationTime)
  if err != nil {
    return err 
  }
  return nil
}

// 'RemoveUserSession' deletes a user's session key and causes us to no longer accept it as a valid session.
// Also removes it as a user cookie.
func RemoveUserSession(key string, w http.ResponseWriter, r *http.Request) error {
  session, err := LoginStore.Get(r, "user-session")
  if err != nil {
    return err
  }
  delete(session.Values, "key")
  return user.RemoveSessionForUser(key)
}

// LoginRegisterRespondJsonError takes in an error code and an appropriate redirectURL and sends it to the client in JSON form.
func LoginRegisterRespondJsonError(errorCode LoginRegisterErrorCode, redirectUrl string, w http.ResponseWriter) {
  response := LoginRegisterResponse{ErrorCode: errorCode, RedirectUrl: redirectUrl}
  if errorCode != ErrorNoError {
    htmlErrCode, errString := getErrorCodeFromLoginError(errorCode)
    http.Error(w, errString, htmlErrCode)
  }
  
  // If any error happens here, then the only thing we can redirect the user to is an error page.
  err := utility.WriteJsonToResponse(w, response)
  if err != nil {
    log.Print(err)
    return
  }
}

func getErrorCodeFromLoginError(errorCode LoginRegisterErrorCode) (int, string) {
  switch errorCode {
  case ErrorInvalidRegisterUserName, ErrorInvalidRegisterPassword, ErrorInvalidRegisterEmail, ErrorInvalidLoginCredentials:
    return 401, ""
  case ErrorAlreadyLoggedIn:
    return 403, ""
  case ErrorUnspecifiedError:
    return 500, ""
  default:
    return 200, ""
  }
  return 200, ""
}