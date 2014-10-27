package admin

/*
 * event.go in the 'admin' package handles the displays and actions related to handling the
 * 'events' which are what we associate with the streams to get the live stats.
 */
import (
	"AfkChampFrontend/controller"
	"AfkChampFrontend/model/event"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const MaxEventShorthandLength = 20

type AdminEventTemplateData struct {
	Data          controller.BaseTemplateData
	Events        []event.EventRowJoined
	IsNewEvent    bool
	SelectedEvent event.EventRowJoined

	EventShorthandCharLimit int
}

// 'HandleAdminEventIndexRoute' displays the page that lists out all the events. It is however
// mandatory to choose a specific game that one wants the events for.
func HandleAdminEventIndexRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	t := CreateBaseEventAdminTemplateData()

	eventVars := mux.Vars(r)
	// We are guaranteed to get a gamename because it's in the route.
	gameName, _ := eventVars["gameName"]

	// Get a list of events based on page
	pageIdx, err := strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageIdx = 0
	}
	allEvents, err := event.GetEventsJoined(pageIdx*DefaultPageSize, DefaultPageSize, gameName)
	if err != nil {
		allEvents = make([]event.EventRowJoined, 0, 0)
	}
	t.Events = allEvents
	controller.TemplateMapping["admin/event/index.html"].ExecuteTemplate(w, "tbase", t)
}

// 'HandleAdminEventNewRoute' displays the page to create a new event.
func HandleAdminEventNewRoute(w http.ResponseWriter, r *http.Request) {
	handleAdminEventNewEditRoute(w, r, true)
}

// handleAdminEventNewEditRoute handles the inner workings of the new route and edit route for
// events.
func handleAdminEventNewEditRoute(w http.ResponseWriter, r *http.Request, isNew bool) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	t := CreateBaseEventAdminTemplateData()
	t.IsNewEvent = isNew
	controller.TemplateMapping["admin/event/newedit.html"].ExecuteTemplate(w, "tbase", t)

}

// 'CreateBaseEventAdminTemplateData' creates the template data for rendering.
func CreateBaseEventAdminTemplateData() *AdminEventTemplateData {
	t := AdminEventTemplateData{Data: controller.CreateTemplateData(),
		EventShorthandCharLimit: MaxEventShorthandLength}
	return &t
}
