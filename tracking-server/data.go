package tracking_server

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type Stream struct {
	Id         uuid.UUID `pg:"type:uuid,pk"`
	Events     []Event   `pg:"rel:has-many"`
	InsertedAt time.Time `json:"insertedAt" pg:",notnull"`
	UpdatedAt  time.Time `json:"updatedAt" pg:",notnull"`
}

type Event struct {
	Name       string                 `json:"name" pg:",notnull"`
	Ip         string                 `json:"ip" pg:",notnull"`
	Uuid       uuid.UUID              `json:"uuid" pg:"type:uuid,unique,notnull"`
	Parameters map[string]interface{} `json:"parameters"`
	StreamID   uuid.UUID              `json:"streamId" pg:"type:uuid,notnull"`
	Stream     *Stream                `json:"stream" pg:"rel:has-one"`
	InsertedAt time.Time              `json:"insertedAt" pg:",notnull"`
	UpdatedAt  time.Time              `json:"updatedAt" pg:",notnull"`
}

type DatabaseOptions struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func ConnectToDatabase(database DatabaseOptions) (*pg.DB, error) {
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

func InsertEvent(db *pg.DB, event *Event) error {
	if _, err := db.Model(event).Insert(); err != nil {
		return err
	}

	return nil
}

func InsertStream(db *pg.DB, stream *Stream) error {
	if _, err := db.Model(stream).Insert(); err != nil {
		return err
	}

	return nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{(*Stream)(nil), (*Event)(nil)}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			return err
		}
	}

	return nil
}
