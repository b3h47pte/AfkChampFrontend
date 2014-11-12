package admin

/*
 * game.go in the 'admin' package handles the game related display of the new/edit/delete pages for everything ame related along with
 * the appropriate REST API.
 */
import (
	"AfkChampFrontend/controller"
	"AfkChampFrontend/model/game"
	"AfkChampFrontend/utility"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const MaxGameNameLength = 100
const MaxGameShorthandLength = 20

type GameNewEditErrorCode int

const (
	errorGameNoError GameNewEditErrorCode = iota
	errorGameInvalidOperation
	errorGameUnspecifiedError
)

type AdminGameTemplateData struct {
	Data                   controller.BaseTemplateData
	SelectedGame           game.GameRow
	Games                  []game.GameRow
	IsNewGame              bool
	GameNameCharLimit      int
	GameShorthandCharLimit int
	OldGameShorthand       string
}

type AdminGameNewEditPostData struct {
	IsNew         bool
	GameName      string
	GameShorthand string
	OldShorthand  string
}

type AdminGameNewEditResponseData struct {
	ErrorCode GameNewEditErrorCode
}

// 'HandleAdminGamePageRoute' displays the page of all games that we support.
func HandleAdminGamePageRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}

	t := CreateBaseGameAdminTemplateData(w, r)
	// Show games in pages. So figure out which page we want
	pageIdx, err := strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageIdx = 0
	}

	// Get a portion of the games.
	allGames, err := game.GetGames(pageIdx*DefaultPageSize, DefaultPageSize)
	if err != nil {
		allGames = make([]game.GameRow, 0, 0)
	}
	t.Games = allGames
	controller.TemplateMapping["admin/game/index.html"].ExecuteTemplate(w, "tbase", t)
}

// 'HandleAdminGameEditRoute' takes a specific game and presents a page with its information and the option to edit any of the values.
func HandleAdminGameEditRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	gameVars := mux.Vars(r)
	// We are guaranteed to get a gamename because it's in the route.
	gameName, _ := gameVars["gameName"]
	handleAdminNewEditGameRoute(w, r, false, gameName)
}

// 'HandleAdminGameNewRoute' takes a given shorthand name (or none at all) and present a page for users to fill in the values.
func HandleAdminGameNewRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	shorthand := r.FormValue("shorthand")
	handleAdminNewEditGameRoute(w, r, true, shorthand)
}

// 'handleAdminNewEditGameRoute' will take the game name and whether or not it is a new game and create the proper response page for it.
// This saves us some code since the new game/edit game page should look the same.
func handleAdminNewEditGameRoute(w http.ResponseWriter, r *http.Request, isNew bool, gameName string) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	gameRow, err := game.GetGame(gameName)
	// If we didn't find the game and we're somehow looking to 'edit' it, then make it so that we're creating doing a 'new' request instead. Change the URL to make it clear to the user.
	if err != nil && !isNew {
		http.Redirect(w, r, "/admin/new/game/shorthand="+gameName, http.StatusFound)
		return
	} else if err != nil {
		gameRow = &game.GameRow{GameShorthand: gameName}
	}

	t := CreateBaseGameAdminTemplateData(w, r)
	t.SelectedGame = *gameRow
	t.IsNewGame = isNew
	t.OldGameShorthand = gameName
	controller.TemplateMapping["admin/game/newedit.html"].ExecuteTemplate(w, "tbase", t)
}

// 'CreateBaseGameAdminTemplateData' creates the template data for rendering.
func CreateBaseGameAdminTemplateData(w http.ResponseWriter, r *http.Request) *AdminGameTemplateData {
	t := AdminGameTemplateData{Data: controller.CreateTemplateData(w, r),
		GameNameCharLimit:      MaxGameNameLength,
		GameShorthandCharLimit: MaxGameShorthandLength}
	return &t
}

// 'HandleNewEditGamePost' handles POST requests for either new and/or edit requests. For both requests, we do a check to see if the game already exists, but in the new case, it it already exists, we fail while in the edit
// case, if it doesn't exist, we fail.
func HandleNewEditGamePost(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	gameData := AdminGameNewEditPostData{}
	err := utility.ReadJsonFromRequestBodyStruct(r, &gameData)
	if err != nil {
		NewEditGameRespondJsonError(errorGameUnspecifiedError, w)
		return
	}

	// For a 'new' entry, search to see if the new entry is a duplicate
	if gameData.IsNew {
		_, err = game.GetGame(gameData.GameShorthand)
	} else {
		// For an 'update' entry, need to search to see if the OLD entry is a duplicate
		_, err = game.GetGame(gameData.OldShorthand)
	}
	// A new game doesn't want the game to exist, an edit game update wants the game to exist
	if gameData.IsNew == (err == nil) {
		log.Print(err)
		NewEditGameRespondJsonError(errorGameInvalidOperation, w)
		return
	}

	// Create the game row object from the game data.
	newGame := game.GameRow{GameName: gameData.GameName, GameShorthand: gameData.GameShorthand}

	// If we get to this point, we can modify the item in the database.
	if gameData.IsNew {
		err = game.CreateGame(&newGame)
	} else {
		err = game.UpdateGame(gameData.OldShorthand, &newGame)
	}

	if err != nil {
		log.Print(err)
		NewEditGameRespondJsonError(errorGameUnspecifiedError, w)
		return
	}

	NewEditGameRespondJsonError(errorGameNoError, w)
}

// NewEditGameRespondJsonError takes in an error code and passes it back to the client in the form of a JSON response.
func NewEditGameRespondJsonError(errorCode GameNewEditErrorCode, w http.ResponseWriter) {
	response := AdminGameNewEditResponseData{ErrorCode: errorCode}
	if errorCode != errorGameNoError {
		htmlErrCode := getErrorCodeFromGameError(errorCode)
		http.Error(w, "", htmlErrCode)
	}

	// If any error happens here, then the only thing we can redirect the user to is an error page.
	err := utility.WriteJsonToResponse(w, response)
	if err != nil {
		log.Print(err)
		return
	}
}

// 'HandleAdminGameDeleteRoute' deletes the specified game.
func HandleAdminGameDeleteRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	gameVars := mux.Vars(r)
	// We are guaranteed to get a gamename because it's in the route.
	gameName, _ := gameVars["gameName"]
	if err := game.DeleteGame(gameName); err != nil {
		log.Print(err)
	}

	http.Redirect(w, r, "/admin/game", http.StatusFound)
}

// 'getErrorCodeFromGameError' takes in a game error code and returns a HTTP error code along with it.
func getErrorCodeFromGameError(errorCode GameNewEditErrorCode) int {
	switch errorCode {
	case errorGameUnspecifiedError, errorGameInvalidOperation:
		return 500
	default:
		return 200
	}
	return 200
}
