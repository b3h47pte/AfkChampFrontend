/*
 * 'live' handles the live stats functionality of the site. How the client actually receives
 * live data is handled by the API server as well as the client-side Javascript. This package just
 * lets the client know what they requested sets up that information in the template data.
 */
package live

import (
	"AfkChampFrontend/controller"
	"github.com/gorilla/mux"
	"net/http"
)

type LiveStatsTemplateData struct {
	Data controller.BaseTemplateData
}

// HandleLiveStatsPageRoute shows the live stats page given a specific event name.
func HandleLiveStatsPageRoute(w http.ResponseWriter, r *http.Request) {
	liveVars := mux.Vars(r)
	eventShorthand, _ := liveVars["eventShorthand"]

	t := CreateLiveStatsTemplateData()
	controller.TemplateMapping["live/live.html"].ExecuteTemplate(w, "tbase", t)
}

// 'CreateLiveStatsTemplateData' creates the template data for rendering.
func CreateLiveStatsTemplateData() *LiveStatsTemplateData {
	t := LiveStatsTemplateData{Data: controller.CreateTemplateData()}
	return &t
}
