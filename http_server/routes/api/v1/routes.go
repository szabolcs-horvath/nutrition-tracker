package v1

import (
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/daily_quotas"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/items"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/meallogs"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/meals"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/notifications"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/portions"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/users"
	"net/http"
)

const Prefix = "/v1"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		daily_quotas.Prefix:  routes.SubRouteHandlerFuncs(daily_quotas.HandlerFuncs()),
		items.Prefix:         routes.SubRouteHandlerFuncs(items.HandlerFuncs()),
		meallogs.Prefix:      routes.SubRouteHandlerFuncs(meallogs.HandlerFuncs()),
		meals.Prefix:         routes.SubRouteHandlerFuncs(meals.HandlerFuncs()),
		notifications.Prefix: routes.SubRouteHandlerFuncs(notifications.HandlerFuncs()),
		portions.Prefix:      routes.SubRouteHandlerFuncs(portions.HandlerFuncs()),
		users.Prefix:         routes.SubRouteHandlerFuncs(users.HandlerFuncs()),
	}
}
