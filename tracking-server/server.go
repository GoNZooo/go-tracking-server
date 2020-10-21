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
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	database *pg.DB
	router   *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	s.router.ServeFiles("/js/*filepath", http.Dir("./static/js"))
	s.router.HandlerFunc("POST", "/events/initiate", s.handleInitiateEventStream())
	s.router.HandlerFunc("POST", "/events", s.handleEvent())
}

func InitializeServer() *Server {
	server := &Server{router: httprouter.New()}
	server.setupRoutes()

	return server
}

func (s *Server) handleInitiateEventStream() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Println("Error creating stream ID:", err.Error())
			http.Error(writer, "ERROR_NO_STREAM_ID", http.StatusInternalServerError)

			return
		}

		timeNow := time.Now()
		stream := Stream{Id: id, InsertedAt: timeNow, UpdatedAt: timeNow}
		if err := InsertStream(s.database, &stream); err != nil {
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

func (s *Server) handleEvent() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody := request.Body
		requestIp := headerOrDefault(request.Header, "X-Real-Ip", strings.Split(request.RemoteAddr, ":")[0])

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

		if err := InsertEvent(s.database, &event); err != nil {
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

func (s *Server) Serve(port int, database DatabaseOptions) {
	portSpecification := fmt.Sprintf(":%d", port)
	db, err := ConnectToDatabase(database)
	if err != nil {
		panic(err)
	}
	s.database = db

	if err := http.ListenAndServe(portSpecification, s); err != nil {
		log.Fatal(err)
	}
}

// Returns the existing header value if it exists, otherwise the default value. `headerName` is not case-sensitive,
// since this uses `Header.Get()` internally.
func headerOrDefault(header http.Header, headerName string, defaultValue string) string {
	if value := header.Get(headerName); value != "" {
		return value
	} else {
		return defaultValue
	}
}
