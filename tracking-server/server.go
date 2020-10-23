package tracking_server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httprouter.Router
}

func (s *Server) Serve(port int, database DatabaseOptions) {
	portSpecification := fmt.Sprintf(":%d", port)
	db, err := ConnectToDatabase(database)
	if err != nil {
		panic(err)
	}
	s.setupRoutes(db)

	if err := http.ListenAndServe(portSpecification, s); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) setupRoutes(db *pg.DB) {
	if s.router == nil {
		s.router = httprouter.New()
	}
	s.router.ServeFiles("/js/*filepath", http.Dir("./static/js"))
	s.router.HandlerFunc("POST", "/events/initiate", handleInitiateEventStream(db))
	s.router.HandlerFunc("POST", "/events", handleEvent(db))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
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
