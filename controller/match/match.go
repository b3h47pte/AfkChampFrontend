package match

import (
	"AfkChampFrontend/controller"
    "AfkChampFrontend/model/match"
    "github.com/gorilla/mux"
    "code.google.com/p/gcfg"
	"net/http"
    "strconv"
    "log"
    "fmt"
)

type ApiConfig struct {
  Api struct {
      Url       string
  }
}

type MatchTemplateData struct {
	Data            controller.BaseTemplateData
    MatchId         int64
    WebsocketUrl    string
}

/*
 * HandleMatchPageRoute displays the page for this particular match. 
 * Most of the interactive ability of the page is done via Javascript.
 */
func HandleMatchPageRoute(w http.ResponseWriter, r *http.Request) {
	t := createMatchMainTemplateData(w, r)
    
	matchVars := mux.Vars(r)
	// We are guaranteed to match id because it's in the route.
	matchId, _ := matchVars["matchId"]
	convertedMatchId, nerr := strconv.ParseInt(matchId, 0, 64)
    if nerr != nil {
        log.Print(nerr);
		controller.Handle404Page(w, r)
		return
    }
    
    t.MatchId = convertedMatchId
    // Get the two teams against each other and display it as the title
    // i.e. TEAM1 vs TEAM2 :: Game 2/3 (DATE)
    basicMatch, nerr := match.QueryBasicInformation(t.MatchId)
    if nerr != nil {
        log.Print(nerr);
		controller.Handle404Page(w, r)
		return
    }
    t.Data.WebsiteName = fmt.Sprintf("%s vs %s :: Game %d/%d (%s)",
                                     basicMatch.TeamOne, basicMatch.TeamTwo,
                                     basicMatch.CurrentGame, basicMatch.TotalGames,
                                     basicMatch.MatchDate.Format("01/02/2006"))
    
    // Pass the API server URL to the client.
    // This URL will be used for a Websocket connection. 
    var config ApiConfig
    err := gcfg.ReadFileInto(&config, "config/api.config")
    if err != nil {
        log.Fatal(err)
        controller.Handle404Page(w, r)
		return
    }
    t.WebsocketUrl = config.Api.Url;
    
	controller.TemplateMapping["match/match.html"].ExecuteTemplate(w, "tbase", t)
}

func createMatchMainTemplateData(w http.ResponseWriter, r *http.Request) *MatchTemplateData {
	t := MatchTemplateData{Data: controller.CreateTemplateData(w, r)}
	return &t
}
