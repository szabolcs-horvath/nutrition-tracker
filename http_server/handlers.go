package http_server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
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

func HandleItemGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idParam := r.PathValue("id")
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
		if err = helpers.WriteJson(w, http.StatusOK, item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}

func HandleItemPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
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
	if err = helpers.WriteJson(w, http.StatusCreated, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
