package http_server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"shorvath/nutrition-tracker/repository"
	"strconv"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"test\":\"value\"}"))
	if err != nil {
		log.Fatal(err)
	}
}

func HandleItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		idParam := r.Form.Get("id")
		if idParam != "" {
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			item, err := repository.FindItemByIdWithNutrition(ctx, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			jsonResponse, err := json.Marshal(item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err = w.Write(jsonResponse); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			return
		}
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
		return
	case "POST":
		var requestItem repository.Item
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := repository.CreateItem(ctx, requestItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(jsonResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
	}
}
