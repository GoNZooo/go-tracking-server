package tracking_server

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type Event struct {
	Name       string                 `json:"name"`
	Ip         string                 `json:"ip"`
	Uuid       uuid.UUID              `json:"uuid"`
	Parameters map[string]interface{} `json:"parameters"`
	InsertedAt time.Time              `json:"insertedAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

func ConnectToDatabase(host string, port int, user string, password string, database string) (*pg.DB, error) {
	connection := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		User:     user,
		Password: password,
		Database: database,
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

func createSchema(db *pg.DB) error {
	models := []interface{}{(*Event)(nil)}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			return err
		}
	}

	return nil
}
