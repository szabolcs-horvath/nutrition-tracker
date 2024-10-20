package http_server

import "net/http"

func GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return map[string]func(w http.ResponseWriter, r *http.Request){
		"/":     HandleRoot,
		"/item": HandleItem,
	}
}
