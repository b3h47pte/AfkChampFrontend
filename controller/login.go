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
  // Admin determines whether or not the login request is for an admin login
  Admin string
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
  forceAdmin := r.FormValue("admin") == "true"
  if _, err := GetCurrentUser(w, r, forceAdmin); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }

  t := CreateLoginTemplateData()
  TemplateMapping["login/login.html"].ExecuteTemplate(w, "tbase", t)
}

// HandleLogoutPageRoute logs out the current user and directs him/her back to the home page. Removes any admin sessions as well
func HandleLogoutPageRoute(w http.ResponseWriter, r *http.Request) {
  if sessionKey, err := GetUserSession(w, r, false); err == nil {
    RemoveUserSession(sessionKey, w, r, false)
  }
  
  if sessionKey, err := GetUserSession(w, r, true); err == nil {
    RemoveUserSession(sessionKey, w, r, true)
  }
  
  http.Redirect(w, r, "/", http.StatusFound)
}


// HandlRegisterPageRoute displays the registration page if the user is not logged in. Otherwise, the 
// user is redirected to the home page.
func HandleRegisterPageRoute(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r, false); err == nil {
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
  userData := LoginRegisterPostData{}
  err := utility.ReadJsonFromRequestBodyStruct(r, &userData)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
  
  isAdminRequest := (userData.Admin == "true")
  if _, err := GetCurrentUser(w, r, isAdminRequest); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    LoginRegisterRespondJsonError(ErrorAlreadyLoggedIn, "/", w)
    return
  }
 
  newUser, err := user.VerifyUser(userData.Username, userData.Password)
  if err != nil {
    LoginRegisterRespondJsonError(ErrorInvalidLoginCredentials, "", w)
    return
  }
  
  err = CreateUserSession(newUser, w, r, isAdminRequest)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "/login", w)
    return
  }
  // Success!
  if isAdminRequest {
    LoginRegisterRespondJsonError(ErrorNoError, "/admin", w)
  } else {
    LoginRegisterRespondJsonError(ErrorNoError, "/", w)
  }
}

// HandleRegisterAction allows you to register a new user given a username and password. If registration is successful, also
// set relevant information in the user's cookies to remember their session. This fails if the user is logged in already.
func HandleRegisterAction(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r, false); err == nil {
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
  err = CreateUserSession(newUser, w, r, false)
  if err != nil {
    log.Print(err)
    LoginRegisterRespondJsonError(ErrorUnspecifiedError, "", w)
    return
  }
  // Assume success and redirect to front page
  LoginRegisterRespondJsonError(ErrorNoError, "/", w)
}

// 'GetCurrentUser' will retrieve the currently logged in user.
func GetCurrentUser(w http.ResponseWriter, r *http.Request, forceAdmin bool) (*user.UserEntry, error) {
  sessionKey, err := GetUserSession(w, r, forceAdmin)
  if err != nil {
    return nil, err 
  }
  userId, err := user.VerifySession(sessionKey, forceAdmin)
  if err != nil {
    return nil, err
  }
  return user.RetrieveUser(userId)
}

// 'GetUserSession' gets the current user's session (if the user has logged in before hand and we stored the session cookie).
func GetUserSession(w http.ResponseWriter, r *http.Request, forceAdmin bool) (string, error) {
  session, err := LoginStore.Get(r, "user-session")
  if forceAdmin {
    session, err = LoginStore.Get(r, "admin-session")
  }
  
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
func CreateUserSession(newUser *user.UserEntry, w http.ResponseWriter, r *http.Request, forceAdmin bool) error {
  session, err := LoginStore.Get(r, "user-session")
  if forceAdmin && newUser.IsAdmin {
    session, err = LoginStore.Get(r, "admin-session")
  } else {
    forceAdmin = false
  }
  
  // AFTER THIS POINT: forceAdmin is a boolean that determines whether or not we are currently trying to create an admin session.
  if err != nil {
    return err
  }
  // Keep session key if one exists already
  sessionKey, ok := session.Values["key"].(string)
  
  // If this is a valid key, then we can keep it otherwise we want to make a new one.
  // If we have a valid key already, then we can just ignore the request.
  _, err = user.VerifySession(sessionKey, forceAdmin)
  if err != nil {
    ok = false
  }
  
  // Generate a new session key if necessary
  if !ok {
    sessionKey = uuid.New()
  } else {
    return nil
  }
  
  session.Values["key"] = sessionKey
  err = session.Save(r, w)
  if err != nil {
    return err
  }
  
  // Associate key with the user for 3 months. Admin keys only last one day.
  expirationTime := time.Now().AddDate(0,3,0)
  if forceAdmin {
    expirationTime = time.Now().AddDate(0,0,1)
  }
  
  err = user.AddSessionForUser(newUser, sessionKey, &expirationTime, forceAdmin)
  if err != nil {
    return err 
  }
  
  // Successfully created the user key. However, if we just made an admin key, we also want to make sure they get a user key too.
  if forceAdmin {
    return CreateUserSession(newUser, w, r, false)
  }
  
  return nil
}

// 'RemoveUserSession' deletes a user's session key and causes us to no longer accept it as a valid session.
// Also removes it as a user cookie.
func RemoveUserSession(key string, w http.ResponseWriter, r *http.Request, forceAdmin bool) error {
  session, err := LoginStore.Get(r, "user-session")
  if forceAdmin {
    session, err = LoginStore.Get(r, "admin-session")
  }
  
  if err != nil {
    return err
  }
  delete(session.Values, "key")
  return user.RemoveSessionForUser(key, forceAdmin)
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