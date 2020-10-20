package tracking_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

func eventHandler(database *pg.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody := request.Body
		if requestBody == nil {
			http.Error(writer, "Need event data", http.StatusBadRequest)
		}

		var event Event
		if err := json.NewDecoder(requestBody).Decode(&event); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)

			return
		}
		event.Uuid = uuid.New()
		event.Ip = requestIp
		timeNow := time.Now()
		event.InsertedAt = timeNow
		event.UpdatedAt = timeNow

		if err := InsertEvent(database, &event); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}

		_, err := fmt.Fprintf(writer, "{\"status\":\"ok\"}")
		if err != nil {
			log.Println("Error writing for request:", err.Error())
		}
	}
}

func Serve(port int, databaseHost string, databasePort int, databaseDatabase string, databaseUser string, databasePassword string) {
	portSpecification := fmt.Sprintf(":%d", port)
	fileServer := http.FileServer(http.Dir("./static"))
	db, err := ConnectToDatabase(databaseHost, databasePort, databaseUser, databasePassword, databaseDatabase)
	if err != nil {
		panic(err)
	}

	http.Handle("/js/", fileServer)

	http.HandleFunc("/events", eventHandler(db))

	if err := http.ListenAndServe(portSpecification, nil); err != nil {
		log.Fatal(err)
	}
}
