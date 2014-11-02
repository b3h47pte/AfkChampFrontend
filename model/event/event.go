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
	rows, err := model.Database.Queryx(`SELECT events.eventid AS eventid,
																						 users.username AS eventowner,
																						 events.eventname AS eventname,
																						 games.gamename AS currentgamename,
																						 events.streamurl AS streamurl,
																						 events.eventshorthand AS eventshorthand 
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
			return nil, err
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
			return nil, err
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

// AddEventJoined takes in the event data that has the string owner and game name and adds it
// to the database.
func AddEventJoined(newEvent *EventRowJoined) error {
	_, nerr := model.Database.Exec(
		`INSERT INTO events (ownerid, eventname, currentgameid, streamurl, eventshorthand)
			 SELECT users.userid, ?, games.gameid, ?, ? 
			 FROM games, users
			 WHERE games.gameshorthand = ?
			 AND users.username = ?`, newEvent.EventName, newEvent.StreamUrl, newEvent.EventShorthand,
		newEvent.CurrentGameName, newEvent.EventOwner)
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

// ModifyEventByShorthandJoined takes in the event shorthand and modifies its properties.
func ModifyEventByShorthandJoined(eventShorthand string, gameShorthand string, newProperties *EventRowJoined) error {
	_, nerr := model.Database.Exec(`UPDATE events 
																	INNER JOIN users ON users.username = ?
																	INNER JOIN games ON games.gameshorthand = ?
																	SET events.ownerid = users.userid, events.eventname = ?, events.currentgameid = games.gameid, events.streamurl = ?, events.eventshorthand = ? 
																	WHERE events.eventshorthand = ? AND games.gameshorthand = ?`,
		newProperties.EventOwner, gameShorthand,
		newProperties.EventName, newProperties.StreamUrl,
		newProperties.EventShorthand, eventShorthand, gameShorthand)
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

// GetEventByShorthand returns an event given an event shorthand.
func GetEventByShorthand(shorthand string) (*EventRow, error) {
	row := model.Database.QueryRowx("SELECT * FROM events WHERE eventshorthand = ?", shorthand)
	newEvent := EventRow{}
	err := row.StructScan(&newEvent)
	if err != nil {
		return nil, err
	}
	return &newEvent, nil
}

// GetEventByShorthandAndGameJoined returns an event given an event shorthand with the game id and the owner id joined to get the actual game/owner name.
func GetEventByShorthandAndGameJoined(shorthand string, gamename string) (*EventRowJoined, error) {
	row := model.Database.QueryRowx(`SELECT events.eventid AS eventid,
																						 users.username AS eventowner,
																						 events.eventname AS eventname,
																						 games.gamename AS currentgamename,
																						 events.streamurl AS streamurl,
																						 events.eventshorthand AS eventshorthand 
																			FROM events
																			INNER JOIN users ON events.ownerid = users.userid
																			INNER JOIN games on events.currentgameid = games.gameid
																			WHERE events.eventshorthand = ? AND games.gameshorthand = ?`, shorthand, gamename)
	newEvent := EventRowJoined{}
	err := row.StructScan(&newEvent)
	if err != nil {
		return nil, err
	}
	return &newEvent, nil
}
