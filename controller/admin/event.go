package admin

/*
 * event.go in the 'admin' package handles the displays and actions related to handling the
 * 'events' which are what we associate with the streams to get the live stats.
 */
import (
	"AfkChampFrontend/controller"
	"AfkChampFrontend/model/event"
	"AfkChampFrontend/utility"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const MaxEventShorthandLength = 20

type AdminEventTemplateData struct {
	Data                  controller.BaseTemplateData
	Events                []event.EventRowJoined
	IsNewEvent            bool
	SelectedEvent         event.EventRowJoined
	CurrentGameShorthand  string
	CurrentEventShorthand string

	EventShorthandCharLimit int
}

type AdminEventErrorCode int

const (
	errorEventNoError AdminEventErrorCode = iota
	errorEventUnspecifiedError
)

type AdminEventPostData struct {
	IsNew                  bool
	Event                  event.EventRowJoined
	OriginalEventShorthand string
	OriginalGameShorthand  string
}

type AdminEventResponseData struct {
	ErrorCode AdminEventErrorCode
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
		log.Print(err)
		allEvents = make([]event.EventRowJoined, 0, 0)
	}
	t.Events = allEvents
	t.CurrentGameShorthand = gameName
	controller.TemplateMapping["admin/event/index.html"].ExecuteTemplate(w, "tbase", t)
}

// 'HandleAdminEventNewRoute' displays the page to create a new event.
func HandleAdminEventNewRoute(w http.ResponseWriter, r *http.Request) {
	handleAdminEventNewEditRoute(w, r, true)
}

// 'HandleAdminEventEditRoute' displays the page to edit an event.
func HandleAdminEventEditRoute(w http.ResponseWriter, r *http.Request) {
	handleAdminEventNewEditRoute(w, r, false)
}

// handleAdminEventNewEditRoute handles the inner workings of the new route and edit route for
// events.
func handleAdminEventNewEditRoute(w http.ResponseWriter, r *http.Request, isNew bool) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}
	t := CreateBaseEventAdminTemplateData()
	t.IsNewEvent = isNew

	if !isNew {
		eventVars := mux.Vars(r)
		// We are guaranteed to get a gamename and event shorthand because it's in the route.
		gameName, _ := eventVars["gameName"]
		eventShorthand, _ := eventVars["eventShorthand"]

		// Make sure this game exists...if it doesn't redirect to a new event page.
		currentEvent, err := event.GetEventByShorthandAndGameJoined(eventShorthand, gameName)
		if err != nil {
			http.Redirect(w, r, "/admin/event/new", http.StatusFound)
			return
		}
		t.SelectedEvent = *currentEvent
		t.CurrentGameShorthand = gameName
		t.CurrentEventShorthand = eventShorthand
	}

	controller.TemplateMapping["admin/event/newedit.html"].ExecuteTemplate(w, "tbase", t)
}

// HandleAdminEventNewEditPost handles a new/edit request for an event.
func HandleAdminEventNewEditPost(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}

	eventData := AdminEventPostData{}
	err := utility.ReadJsonFromRequestBodyStruct(r, &eventData)
	if err != nil {
		log.Print(err)
		adminEventRespondJsonError(errorEventUnspecifiedError, w)
		return
	}

	if eventData.IsNew {
		err = event.AddEventJoined(&eventData.Event)
	} else {
		err = event.ModifyEventByShorthandJoined(eventData.OriginalEventShorthand, eventData.OriginalGameShorthand, &eventData.Event)
	}

	if err != nil {
		log.Print(err)
		adminEventRespondJsonError(errorEventUnspecifiedError, w)
		return
	}

	adminEventRespondJsonError(errorEventNoError, w)
}

// 'CreateBaseEventAdminTemplateData' creates the template data for rendering.
func CreateBaseEventAdminTemplateData() *AdminEventTemplateData {
	t := AdminEventTemplateData{Data: controller.CreateTemplateData(),
		EventShorthandCharLimit: MaxEventShorthandLength}
	return &t
}

// adminEventRespondJsonError takes in an error code and passes it back to the client in the form of a JSON response.
func adminEventRespondJsonError(errorCode AdminEventErrorCode, w http.ResponseWriter) {
	response := AdminEventResponseData{ErrorCode: errorCode}
	if errorCode != errorEventNoError {
		htmlErrCode := getErrorCodeFromEventError(errorCode)
		http.Error(w, "", htmlErrCode)
	}

	// If any error happens here, then the only thing we can redirect the user to is an error page.
	err := utility.WriteJsonToResponse(w, response)
	if err != nil {
		log.Print(err)
		return
	}
}

// 'getErrorCodeFromEventError' takes in a game error code and returns a HTTP error code along with it.
func getErrorCodeFromEventError(errorCode AdminEventErrorCode) int {
	switch errorCode {
	case errorEventUnspecifiedError:
		return 500
	default:
		return 200
	}
	return 200
}
