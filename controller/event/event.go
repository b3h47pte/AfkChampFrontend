/*
 * controller/event handles displaying the event pages. Mainly just feeds any initialization data
 * to the client and lets the client-side Javascript do the rest.
 */
package event

import (
	"AfkChampFrontend/controller"
	"AfkChampFrontend/model/event"
	"github.com/gorilla/mux"
	"net/http"
)

type EventTemplateData struct {
	Data           controller.BaseTemplateData
	EventShorthand string
}

/*
 * HandleEventPageRoute displays the landing page for this particular event. This page shows
 * any relevant news/match data to this specific event.
 *
 * Just render the page and the javascript will handle the rest.
 */
func HandleEventPageRoute(w http.ResponseWriter, r *http.Request) {
	t := createEventMainTemplateData(w, r)

	// Modify the header to have the game name.
	eventVars := mux.Vars(r)
	// We are guaranteed to an event shorthand because it's in the route.
	eventShorthand, _ := eventVars["eventShorthand"]
	currentEvent, err := event.GetEventByShorthandJoined(eventShorthand)
	if err != nil {
		controller.Handle404Page(w, r)
		return
	}
	t.EventShorthand = eventShorthand
	t.Data.WebsiteName = currentEvent.EventName
	controller.TemplateMapping["event/main.html"].ExecuteTemplate(w, "tbase", t)
}

// 'createEventMainTemplateData' creates the template data for rendering.
func createEventMainTemplateData(w http.ResponseWriter, r *http.Request) *EventTemplateData {
	t := EventTemplateData{Data: controller.CreateTemplateData(w, r)}
	return &t
}
