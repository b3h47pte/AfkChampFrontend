/*
 * 'live' handles the live stats functionality of the site. How the client actually receives
 * live data is handled by the API server as well as the client-side Javascript. This package just
 * lets the client know what they requested sets up that information in the template data.
 */
package live

import (
	"net/http"
)

// HandleLiveStatsPageRoute shows the live stats page given a specific event name.
func HandleLiveStatsPageRoute(w http.ResponseWriter, r *http.Request) {
	//	liveVars := mux.Vars(r)
	//	eventShorthand, _ := liveVars["eventShorthand"]

}
