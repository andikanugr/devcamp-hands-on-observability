package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()
	requestCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests, labeled by success or error",
		},
		[]string{"method", "endpoint", "status", "is_error"},
	)

	e.GET("/", func(c echo.Context) error {
		requestCounter.WithLabelValues(c.Request().Method, c.Path(), "200", "false").Inc()
		return c.JSON(200, echo.Map{
			"message": "Hello World",
		})
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Logger.Fatal(e.Start(":8080"))
}
