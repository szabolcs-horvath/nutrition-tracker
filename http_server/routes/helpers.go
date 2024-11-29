package routes

import "net/http"

func SubRouteHandlers(routes map[string]http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	for pattern, handler := range routes {
		mux.HandleFunc(pattern, handler)
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
