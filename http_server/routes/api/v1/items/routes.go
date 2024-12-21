package items

import (
	"encoding/json"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/items"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /":          listHandler,
		"GET /{id}":      findByIdHandler,
		"POST /":         createHandler,
		"POST /multiple": createMultipleHandler,
		"PUT /":          updateHandler,
		"DELETE /{id}":   deleteHandler,
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	list, err := repository.ListItems(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam != "" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := repository.FindItemById(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err = util.WriteJson(w, http.StatusOK, item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestItem *repository.Item
	if err := json.NewDecoder(r.Body).Decode(requestItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := repository.CreateItem(r.Context(), requestItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func createMultipleHandler(w http.ResponseWriter, r *http.Request) {
	var requestItems []*repository.Item
	if err := json.NewDecoder(r.Body).Decode(&requestItems); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := repository.CreateMultipleItems(r.Context(), requestItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestItem *repository.Item
	if err := json.NewDecoder(r.Body).Decode(requestItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := repository.UpdateItem(r.Context(), requestItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam != "" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = repository.DeleteItem(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}
