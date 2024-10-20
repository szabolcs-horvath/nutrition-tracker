package main

import (
	"log"
	"net/http"
	"shorvath/nutrition-tracker/http_server"
)

func main() {
	for k, v := range http_server.GetRoutes() {
		http.HandleFunc(k, v)
	}
	log.Fatalln(http.ListenAndServe(":80", nil))
}
