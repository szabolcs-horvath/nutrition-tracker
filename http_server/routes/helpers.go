package routes

import (
	"net/http"
	"strings"
)

func SubRouteHandlerFuncs(routes map[string]http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	for pattern, handlerFunc := range routes {
		mux.HandleFunc(pattern, handlerFunc)
	}
	return mux
}

func SubRoute(routes map[string]*http.ServeMux) *http.ServeMux {
	mux := http.NewServeMux()
	for subRoute, serveMux := range routes {
		mux.Handle(subRoute+"/", http.StripPrefix(subRoute, serveMux))
	}
	return mux
}

func ServeRoute(router *http.ServeMux, prefix string, routes map[string]*http.ServeMux) {
	mux := SubRoute(routes)
	router.Handle(prefix+"/", http.StripPrefix(prefix, mux))
}

func ServeRouteHandlers(router *http.ServeMux, prefix string, routes map[string]http.HandlerFunc) {
	for pattern, handlerFunc := range routes {
		split := strings.SplitN(pattern, " ", 2)
		path := split[0] + " " + prefix + split[1]
		router.HandleFunc(path, handlerFunc)
	}
}

func ServeFS(router *http.ServeMux, prefix, dir string) {
	router.Handle(prefix+"/", http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
}
