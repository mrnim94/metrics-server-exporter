package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	Echo *echo.Context
}

func (m *MetricsHandler) HandlerMetrics(c echo.Context) error {
	handler := promhttp.Handler()
	handler.ServeHTTP(c.Response(), c.Request())
	return nil
}
