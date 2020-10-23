package tracking_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

func handleInitiateEventStream(db *pg.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Println("Error creating stream ID:", err.Error())
			http.Error(writer, "stream id could not be created", http.StatusInternalServerError)

			return
		}

		timeNow := time.Now()
		stream := Stream{Id: id, InsertedAt: timeNow, UpdatedAt: timeNow}
		if err := InsertStream(db, &stream); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Printf("Error making new stream: %#v: %s", stream, err.Error())

			return
		}

		if _, err := fmt.Fprintf(writer, "%s", id.String()); err != nil {
			log.Println("Unable to send stream ID response:", err.Error())

			return
		}
	}
}

func handleEvent(db *pg.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody := request.Body
		ipFromRemoteAddress := strings.Split(request.RemoteAddr, ":")[0]
		requestIp := headerOrDefault(request.Header, "X-Real-Ip", ipFromRemoteAddress)

		if requestBody == nil {
			http.Error(writer, "need event data", http.StatusBadRequest)

			return
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

		if err := InsertEvent(db, &event); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Printf("Error writing event: %#v: %s", event, err.Error())

			return
		}

		_, err := fmt.Fprintf(writer, "{\"status\":\"ok\"}")
		if err != nil {
			log.Println("Error writing to socket for request:", err.Error())
		}
	}
}
