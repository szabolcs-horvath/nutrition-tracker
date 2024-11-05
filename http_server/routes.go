package routes

import (
	"net/http"
	"shorvath/nutrition-tracker/http_server/routes/v1/items"
)

func ServeApiV1Routes(router *http.ServeMux) *http.ServeMux {
	itemsApi := addRoutes(router, items.Prefix, items.Routes())

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", itemsApi))
	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", v1))
	return api
}

func addRoutes(serveMux *http.ServeMux, prefix string, routes map[string]func(w http.ResponseWriter, r *http.Request)) *http.ServeMux {
	for route, handler := range routes {
		serveMux.HandleFunc(route, handler)
	}
	servedApi := http.NewServeMux()
	servedApi.Handle(prefix+"/", http.StripPrefix(prefix, serveMux))
	return servedApi
}
