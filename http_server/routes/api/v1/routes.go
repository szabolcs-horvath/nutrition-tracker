package v1

import (
	"net/http"
	"shorvath/nutrition-tracker/http_server/routes"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/items"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/nutritions"
)

const Prefix = "/v1"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		items.Prefix:      routes.SubRouteHandlers(items.Routes()),
		nutritions.Prefix: routes.SubRouteHandlers(nutritions.Routes()),
	}
}
