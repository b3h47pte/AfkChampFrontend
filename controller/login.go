/*
 * Provides basic login functionality that can be used across the site.
 */
package controller

import(
  "net/http"
  "code.google.com/p/gcfg"
  "encoding/hex"
  "log"
  "github.com/gorilla/sessions"
)
type LoginConfig struct {
  AuthSection struct {
    AuthKey string
  }
}

type LoginTemplateData struct {
  Data BaseTemplateData
}

var authSessionKey []byte
var loginStore *sessions.CookieStore

// Setup the secure session for the a user's login session.
func InitializeLogin() {
  var config LoginConfig
  err := gcfg.ReadFileInto(&config, "config/login.config")
  if err != nil {
    log.Fatal(err)
  }
  hexAuthKey := config.AuthSection.AuthKey
  authSessionKey, err := hex.DecodeString(hexAuthKey)
  if err != nil {
    log.Fatal(err)
  }
  loginStore = sessions.NewCookieStore(authSessionKey)
}

func HandleLoginPageRoute(w http.ResponseWriter, r *http.Request) {
  t := LoginTemplateData{Data: CreateTemplateData()}
  TemplateMapping["login/login.html"].ExecuteTemplate(w, "tbase", t)
}