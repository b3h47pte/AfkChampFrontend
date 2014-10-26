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
