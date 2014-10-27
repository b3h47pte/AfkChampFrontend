/*
 * 'event' provides basic functionality to perform basic SQL operations on the 'events' database.
 */
package event

import (
	"AfkChampFrontend/model"
	_ "github.com/go-sql-driver/mysql"
)

type EventRow struct {
	EventId        int64
	OwnerId        int64
	EventName      string
	CurrentGameId  int64
	StreamUrl      string
	EventShorthand string
}

type EventRowJoined struct {
	EventId         int64
	EventOwner      string
	EventName       string
	CurrentGameName string
	StreamUrl       string
	EventShorthand  string
}

// GetEventsJoined returns an array of the 'count' events starting at 'offset' assuming the
// events are listed by eventID. This returns an array of 'EventRowJoined' which contains
// the full name of the correspodning owner and game.
func GetEventsJoined(offset int, count int, gameName string) ([]EventRowJoined, error) {
	requestedEvents := make([]EventRowJoined, 0, 0)

	// Find Events
	rows, err := model.Database.Queryx(`SELECT events.eventid AS EventId,
																						 users.username AS EventOwner,
																						 events.eventname AS EventName,
																						 games.gamename AS CurrentGameName,
																						 events.streamurl AS StreamUrl,
																						 events.eventshorthand AS EventShorthand
																			FROM events
																			INNER JOIN users ON events.ownerid = users.userid
																			INNER JOIN games on events.currentgameid = games.gameid
																			ORDER BY events.eventid ASC LIMIT ?, ?`, offset, count)
	if err != nil {
		return nil, err
	}

	// Create structs
	for rows.Next() {
		newObj := EventRowJoined{}
		err = rows.StructScan(&newObj)
		if err != nil {
			continue
		}
		requestedEvents = append(requestedEvents, newObj)
	}

	return requestedEvents, nil
}

// GetEvents returns an array of the 'count' events starting at 'offset' assuming the
// events are listed by eventID. This returns an array of 'EventRow' and we will not
// perform a join here.
func GetEvents(offset int, count int) ([]EventRow, error) {
	requestedEvents := make([]EventRow, 0, 0)

	// Find Events
	rows, err := model.Database.Queryx("SELECT * FROM events ORDER BY eventid ASC LIMIT ?, ?",
		offset, count)
	if err != nil {
		return nil, err
	}

	// Create structs
	for rows.Next() {
		newObj := EventRow{}
		err = rows.StructScan(&newObj)
		if err != nil {
			continue
		}
		requestedEvents = append(requestedEvents, newObj)
	}

	return requestedEvents, nil
}

// AddEvent takes in the event data and adds it to the database
func AddEvent(newEvent *EventRow) error {
	_, nerr := model.Database.Exec("INSERT INTO events (ownerid, eventname, currentgameid, streamurl, eventshorthand) VALUES (?, ?, ?, ?, ?)",
		newEvent.OwnerId, newEvent.EventName, newEvent.CurrentGameId, newEvent.StreamUrl, newEvent.EventShorthand)
	if nerr != nil {
		return nerr
	}
	return nil
}

// RemoveEvent removes an event given an event ID.
func RemoveEvent(eventId int64) error {
	_, err := model.Database.Exec("DELETE FROM events WHERE eventid = ?", eventId)
	if err != nil {
		return err
	}
	return nil
}

// ModifyEvent takes in the event ID which we wish to modify and a structure containing the modified properties
func ModifyEvent(eventId int64, newProperties *EventRow) error {
	_, nerr := model.Database.Exec("UPDATE events SET ownerid = ?, eventname = ?, currentgameid = ?, streamurl = ?, eventshorthand = ? WHERE eventid = ?",
		newProperties.OwnerId, newProperties.EventName, newProperties.CurrentGameId, newProperties.StreamUrl, newProperties.EventShorthand, newProperties.EventId)
	if nerr != nil {
		return nerr
	}
	return nil
}

// GetEventById returns an event.
func GetEventById(eventId int64) (*EventRow, error) {
	row := model.Database.QueryRowx("SELECT * FROM events WHERE eventid = ?", eventId)
	newEvent := EventRow{}
	err := row.StructScan(&newEvent)
	if err != nil {
		return nil, err
	}
	return &newEvent, nil
}
