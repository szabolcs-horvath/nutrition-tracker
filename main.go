package main

import (
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"test\":\"value\"}"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":80", nil))
}
