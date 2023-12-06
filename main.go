package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// setup logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

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

	e.GET("/error", func(c echo.Context) error {
		requestCounter.WithLabelValues(c.Request().Method, c.Path(), "500", "true").Inc()

		log.WithFields(logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Path(),
			"status": "500",
		}).Error("error")

		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	})

	e.GET("/success", func(c echo.Context) error {

		log.WithFields(logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Path(),
			"status": "200",
		}).Info("success")

		requestCounter.WithLabelValues(c.Request().Method, c.Path(), "200", "false").Inc()
		return c.JSON(200, echo.Map{
			"message": "success",
		})
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Logger.Fatal(e.Start(":8080"))
}
