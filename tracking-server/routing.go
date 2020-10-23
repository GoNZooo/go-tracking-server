package tracking_server

import "net/http"

type post struct {
	path    string
	handler http.HandlerFunc
}

func (p post) toRouteInfo() routeInfo {
	return routeInfo{
		method:  "POST",
		path:    p.path,
		handler: p.handler,
	}
}

type routeSpecification interface {
	toRouteInfo() routeInfo
}

type routeInfo struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func (s *Server) addHandlerFunctions(specifications []routeSpecification) {
	for _, spec := range specifications {
		info := spec.toRouteInfo()
		s.router.HandlerFunc(info.method, info.path, info.handler)
	}
}
