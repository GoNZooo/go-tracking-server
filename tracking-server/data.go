package tracking_server

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

// A `stream` identifies a series of events that all happen in the same session on a page.
type stream struct {
	Id         uuid.UUID `pg:"type:uuid,pk"`
	Events     []event   `pg:"rel:has-many"`
	InsertedAt time.Time `json:"insertedAt" pg:",notnull"`
	UpdatedAt  time.Time `json:"updatedAt" pg:",notnull"`
}

// An `event` represents an action on a page, such as the loading of the page, or hovering/mouseover on an element.
type event struct {
	// Identifies the type of the event.
	Name string `json:"name" pg:",notnull"`

	// Identifies the IP address that spawned the event.
	Ip string `json:"ip" pg:",notnull"`

	// A unique identifier for the event, such that it can be identified. This is meant only for uniqueness.
	Uuid uuid.UUID `json:"uuid" pg:"type:uuid,unique,notnull"`

	// Extra parameters; this can be different for each event and the only commonality that definitely should exist
	// is that it is a map from strings to values.
	Parameters map[string]interface{} `json:"parameters"`

	// Identifies the stream that the event belongs to. This is useful for grouping the events into a session.
	StreamID uuid.UUID `json:"streamId" pg:"type:uuid,notnull"`
	Stream   *stream   `json:"stream" pg:"rel:has-one"`

	InsertedAt time.Time `json:"insertedAt" pg:",notnull"`
	UpdatedAt  time.Time `json:"updatedAt" pg:",notnull"`
}

type DatabaseOptions struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func connectToDatabase(database DatabaseOptions) (*pg.DB, error) {
	connection := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", database.Host, database.Port),
		User:     database.User,
		Password: database.Password,
		Database: database.Database,
	})

	if err := createSchema(connection); err != nil {
		panic(err)
	}

	return connection, nil
}

func insertEvent(db *pg.DB, event *event) error {
	if _, err := db.Model(event).Insert(); err != nil {
		return err
	}

	return nil
}

func insertStream(db *pg.DB, stream *stream) error {
	if _, err := db.Model(stream).Insert(); err != nil {
		return err
	}

	return nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{(*stream)(nil), (*event)(nil)}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			return err
		}
	}

	return nil
}
