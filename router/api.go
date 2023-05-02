package router

import (
	"github.com/labstack/echo/v4"
	"metrics-server-exporter/handler"
	"net/http"
)

type API struct {
	Echo                  *echo.Echo
	UsageResourcesHandler handler.UsageResourcesHandler
	PromHandler           http.Handler
}

func (api *API) SetupRouter() {

	api.Echo.GET("/", handler.Welcome)
	api.Echo.GET("/metrics", echo.WrapHandler(api.PromHandler))
	api.Echo.GET("/pod-usage", api.UsageResourcesHandler.HandlerPodUsage)
}
