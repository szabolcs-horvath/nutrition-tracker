package api

import (
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1"
	"net/http"
)

const Prefix = "/api"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		v1.Prefix: routes.SubRoute(v1.Routes()),
	}
}
