package portions

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/portions"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /item/{itemId}": listForItemHandler,
		"GET /{id}":          findByIdHandler,
		"POST /{$}":          createHandler,
		"PUT /{$}":           updateHandler,
		"DELETE /{id}":       deleteHandler,
	}
}

func listForItemHandler(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("itemId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerId, err := repository.GetOwnerIdByItemId(r.Context(), itemId)
	list, err := repository.ListPortionsForItemAndUser(r.Context(), itemId, ownerId)
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := repository.FindPortionById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestPortion repository.CreatePortionRequest
	if err := util.ReadJson(r, requestPortion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	portion, err := repository.CreatePortion(r.Context(), requestPortion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, portion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestPortion repository.UpdatePortionRequest
	if err := util.ReadJson(r, requestPortion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	portion, err := repository.UpdatePortion(r.Context(), requestPortion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, portion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.DeletePortion(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
