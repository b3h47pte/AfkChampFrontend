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
)
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
    w.Write([]byte(user.UserLoggedInError.Error()))
    return
  }
  r.ParseForm()
  newUser, err := user.VerifyUser(r.PostFormValue("username"), r.PostFormValue("password"))
  if err != nil {
    w.Write([]byte(user.UserDoesNotExist.Error()))
    return
  }
  err = CreateUserSession(newUser, w, r)
  if err != nil {
    log.Print(err)
    w.Write([]byte(user.UserUnspecifiedError.Error()))
    return
  }
  // Success!
  http.Redirect(w, r, "/", http.StatusFound)
}

// HandleRegisterAction allows you to register a new user given a username and password. If registration is successful, also
// set relevant information in the user's cookies to remember their session. This fails if the user is logged in already.
func HandleRegisterAction(w http.ResponseWriter, r *http.Request) {
  if _, err := GetCurrentUser(w, r); err == nil {
    // If we succeed, for whatever reason, it means we're logged in already.
    w.Write([]byte(user.UserLoggedInError.Error()))
    return
  }
  r.ParseForm()
  newUser, err := user.CreateUser(r.PostFormValue("username"), r.PostFormValue("password"), r.PostFormValue("email"))
  switch {
  case err == user.UserExistsError:
    // Case where we need to inform user
    w.Write([]byte(err.Error()))
    return
  case err != nil:
    // Just keep a mental note to ourself but display another error to the user
    log.Print(err)
    w.Write([]byte(user.UserUnspecifiedError.Error()))
    return
  }
  err = CreateUserSession(newUser, w, r)
  if err != nil {
    log.Print(err)
    w.Write([]byte(user.UserUnspecifiedError.Error()))
    return
  }
  // Assume success and redirect to front page
  http.Redirect(w, r, "/", http.StatusFound)
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