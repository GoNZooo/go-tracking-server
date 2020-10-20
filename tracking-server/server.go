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

func eventHandler(database *pg.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody := request.Body
		requestIp := request.Header.Get("X-Real-Ip")
		if requestIp == "" {
			requestIp = strings.Split(request.RemoteAddr, ":")[0]
		}

		if requestBody == nil {
			http.Error(writer, "Need event data", http.StatusBadRequest)

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

		if err := InsertEvent(database, &event); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Printf("Error writing event: %#v: %s", event, err.Error())

			return
		}

		_, err := fmt.Fprintf(writer, "{\"status\":\"ok\"}")
		if err != nil {
			log.Println("Error writing for request:", err.Error())
		}
	}
}

func Serve(port int, database DatabaseOptions) {
	portSpecification := fmt.Sprintf(":%d", port)
	fileServer := http.FileServer(http.Dir("./static"))
	db, err := ConnectToDatabase(database)
	if err != nil {
		panic(err)
	}

	http.Handle("/js/", fileServer)

	http.HandleFunc("/events", eventHandler(db))

	if err := http.ListenAndServe(portSpecification, nil); err != nil {
		log.Fatal(err)
	}
}
