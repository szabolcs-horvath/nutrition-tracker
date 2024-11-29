package api

import (
	"net/http"
	"shorvath/nutrition-tracker/http_server/routes"
	"shorvath/nutrition-tracker/http_server/routes/api/v1"
)

const Prefix = "/api"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		v1.Prefix: routes.SubRoute(v1.Routes()),
	}
}
