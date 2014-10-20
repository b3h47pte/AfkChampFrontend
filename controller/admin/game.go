package admin

/*
 * game.go in the 'admin' package handles the game related display of the new/edit/delete pages for everything ame related along with
 * the appropriate REST API.
 */
 import(
  "net/http"
  "AfkChampFrontend/controller"
  "AfkChampFrontend/model/game"
  "strconv"
  "github.com/gorilla/mux"
)


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

// 'HandleAdminGameEditRoute' takes a specific game and presents a page with its information and the option to edit any of the values.
func HandleAdminGameEditRoute(w http.ResponseWriter, r *http.Request) {
  gameVars := mux.Vars(r)
  // We are guaranteed to get a gamename because it's in the route.
  gameName, _ := gameVars["gameName"]
  handleAdminNewEditGameRoute(w, r, false, gameName)
}

// 'HandleAdminGameNewRoute' takes a given shorthand name (or none at all) and present a page for users to fill in the values.
func HandleAdminGameNewRoute(w http.ResponseWriter, r *http.Request) {
  shorthand := r.FormValue("shorthand") 
  handleAdminNewEditGameRoute(w, r, true, shorthand)
}

// 'handleAdminNewEditGameRoute' will take the game name and whether or not it is a new game and create the proper response page for it.
// This saves us some code since the new game/edit game page should look the same.
func handleAdminNewEditGameRoute(w http.ResponseWriter, r *http.Request, isNew bool, gameName string) {
  gameRow, err := game.GetGame(gameName)
  // If we didn't find the game and we're somehow looking to 'edit' it, then make it so that we're creating doing a 'new' request instead. Change the URL to make it clear to the user.
  if err != nil && !isNew {
    http.Redirect(w, r, "/admin/game/new?shorthand=" + gameName, http.StatusFound)
    return
  } else if err != nil {
    gameRow = &game.GameRow{GameShorthand : gameName}
  }
  
  t := AdminTemplateData{Data: controller.CreateTemplateData(), SelectedGame: *gameRow, IsNewGame: isNew}
  controller.TemplateMapping["admin/game/newedit.html"].ExecuteTemplate(w, "tbase", t)
}