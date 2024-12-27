package daily_quotas

import (
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/daily_quotas"

func HandlerFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /{id}":                   findByIdHandler,
		"GET /owner/{id}":             listByOwnerHandler,
		"GET /owner/{id}/date/{date}": findByOwnerAndDateHandler,
		"GET /owner/{id}/date/{$}":    findByOwnerAndDateHandler,
		"POST /{$}":                   createHandler,
		"PUT /{$}":                    updateHandler,
		"DELETE /{id}":                deleteHandler,
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dailyQuota, err := repository.FindDailyQuotaById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, dailyQuota); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func listByOwnerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list, err := repository.ListDailyQuotasForUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findByOwnerAndDateHandler(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var result *repository.DailyQuota
	dateParam := r.PathValue("date")
	if dateParam == "" {
		dq, dqErr := repository.FindDailyQuotaByOwnerAndCurrentDay(r.Context(), ownerId)
		if dqErr != nil {
			http.Error(w, dqErr.Error(), http.StatusInternalServerError)
			return
		}
		result = dq
	} else {
		date, parseErr := custom_types.ParseDate(r.PathValue("date"))
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}
		dq, dqErr := repository.FindDailyQuotaByOwnerAndDate(r.Context(), ownerId, date.UnderlyingTime())
		if dqErr != nil {
			http.Error(w, dqErr.Error(), http.StatusInternalServerError)
			return
		}
		result = dq
	}
	if err = util.WriteJson(w, http.StatusOK, result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestDailyQuota repository.CreateDailyQuotaRequest
	if err := util.ReadJson(r, &requestDailyQuota); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dailyQuota, err := repository.CreateDailyQuota(r.Context(), requestDailyQuota)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, dailyQuota); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestDailyQuota repository.UpdateDailyQuotaRequest
	if err := util.ReadJson(r, &requestDailyQuota); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dailyQuota, err := repository.UpdateDailyQuota(r.Context(), requestDailyQuota)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, dailyQuota); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.ArchiveDailyQuota(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
