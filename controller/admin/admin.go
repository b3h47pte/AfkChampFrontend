/*
 * 'admin' handles the main admin page and general admin control related activities here.
 */
package admin

import(
  "net/http"
  "AfkChampFrontend/controller"
  "AfkChampFrontend/model/game"
//  "github.com/gorilla/mux"
  "strconv"
)

type AdminTemplateData struct {
  Data controller.BaseTemplateData
  SelectedUser interface{}
  SelectedEvent interface{}
  Events []interface{}
  Users []interface{}
  Games []game.GameRow
}

const DefaultPageSize = 15

// 'RequireAdminRelogin' requires the user to re-login when visiting the admin page. Returns nil if user is already logged in as an admin.
func RequireAdminRelogin(w http.ResponseWriter, r *http.Request) error {
  if _, err := controller.GetCurrentUser(w, r, true); err != nil {
    http.Redirect(w, r, "login?admin=true", 302)
    return err
  }
  return nil
}

// 'HandleAdminRoute' Handles the home page for the admin section. 
func HandleAdminRoute(w http.ResponseWriter, r *http.Request) {
  if err := RequireAdminRelogin(w, r); err != nil {
    return
  }
  
  // If we get here then we know that the user is an admin.
  t := AdminTemplateData{Data: controller.CreateTemplateData()}
  controller.TemplateMapping["admin/admin.html"].ExecuteTemplate(w, "tbase", t)
}

// 'HandleAdminGamePageRoute' displays the page of all games that we support.
func HandleAdminGamePageRoute(w http.ResponseWriter, r *http.Request) {
  t := AdminTemplateData{Data: controller.CreateTemplateData()}
  // Show games in pages. So figure out which page we want
  pageIdx, err := strconv.Atoi(r.FormValue("p"))
  if err != nil {
    pageIdx = 0
  }
  
  // Then figure out how many entries we want
  entryCount, err := strconv.Atoi(r.FormValue("c"))
  if err != nil {
    entryCount = DefaultPageSize
  }
  
  // Get a portion of the games. 
  allGames, err := game.GetGames(pageIdx, entryCount)
  if err != nil {
    allGames = make([]game.GameRow,0,0)
  }
  t.Games = allGames  
  controller.TemplateMapping["admin/game/index.html"].ExecuteTemplate(w, "tbase", t)
}

// 'HandleAdminGameRoute' takes a game and returns all the events related to that game. This page allows us to add/delete/modify anyh
// events.
func HandleAdminGameRoute(w http.ResponseWriter, r *http.Request) {
//  gameVars := mux.Vars(r)
  // We are guaranteed to get a gamename because it's in the route.
//  gameName, _ := gameVars["gameName"]
  
//  t := AdminTemplateData{Data: controller.CreateTemplateData()}
  
}

// 'HandleAdminEventRoute' takes the given game and event name and returns related details.
func HandleAdminEventRoute(w http.ResponseWriter, r *http.Request) {
//  eventVars := mux.Vars(r)
//  gameName, ok := eventVars["gameName"]
}

// 'HandleAdminUserRoute' presents a list of all users and allows us to view/modify any user properties as necessary.
func HandleAdminUserPageRoute(w http.ResponseWriter, r *http.Request) {
}

func HandleAdminUserRoute(w http.ResponseWriter, r *http.Request) {
//  userVars := mux.Vars(r)
//  username, ok := userVars["username"]
}
