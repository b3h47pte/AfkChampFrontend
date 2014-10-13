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
)
type LoginConfig struct {
  AuthSection struct {
    AuthKey string
  }
}

type LoginTemplateData struct {
  Data BaseTemplateData
}

var AuthSessionKey []byte
var LoginStore *sessions.CookieStore

// Setup the secure session for the a user's login session.
func InitializeLogin() {
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
  t := LoginTemplateData{Data: CreateTemplateData()}
  TemplateMapping["login/login.html"].ExecuteTemplate(w, "tbase", t)
}

// HandlRegisterPageRoute displays the registration page if the user is not logged in. Otherwise, the 
// user is redirected to the home page.
func HandleRegisterPageRoute(w http.ResponseWriter, r *http.Request) {
  t := LoginTemplateData{Data: CreateTemplateData()}
  TemplateMapping["login/register.html"].ExecuteTemplate(w, "tbase", t)
}

// HandleLoginAction takes in the user's name and password and checks whether or not they are registered. Sets 
// relevant information in the cookie store to remember the user's session.
func HandleLoginAction(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  err := user.VerifyUser(r.PostFormValue("username"), r.PostFormValue("password"))
  if err != nil {
    w.Write([]byte(user.UserDoesNotExist.Error()))
    return
  }
  
  // Success!
  http.Redirect(w, r, "/", http.StatusFound)
}

// HandleRegisterAction allows you to register a new user given a username and password. If registration is successful, also
// set relevant information in the user's cookies to remember their session.
func HandleRegisterAction(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  err := user.CreateUser(r.PostFormValue("username"), r.PostFormValue("password"), r.PostFormValue("email"))
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
  
  // Assume success and redirect to front page
  http.Redirect(w, r, "/", http.StatusFound)
}