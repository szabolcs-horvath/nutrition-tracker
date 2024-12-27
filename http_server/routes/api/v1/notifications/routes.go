package notifications

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/notifications"

func HandlerFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /{id}":       findByIdHandler,
		"GET /owner/{id}": listByOwnerHandler,
		"POST /{$}":       createHandler,
		"PUT /{$}":        updateHandler,
		"DELETE /{id}":    deleteHandler,
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notification, err := repository.FindNotificationById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func listByOwnerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list, err := repository.ListNotificationsByUserId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestNotification repository.CreateNotificationRequest
	if err := util.ReadJson(r, &requestNotification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notification, err := repository.CreateNotification(r.Context(), requestNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestNotification repository.UpdateNotificationRequest
	if err := util.ReadJson(r, &requestNotification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notification, err := repository.UpdateNotification(r.Context(), requestNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.DeleteNotification(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
