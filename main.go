package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"math/rand"
	"os"
	"time"
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

	// setup tracer
	tracer, closer := initJaeger("hands-on-app")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	rand.Seed(time.Now().UnixNano())

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

	e.GET("/slow", func(c echo.Context) error {
		// Start a span
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "http_request")
		defer span.Finish()

		// your slow code here with random sleep to simulate
		slowFunction(ctx)

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

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(err.Error())
	}
	return tracer, closer
}

func slowFunction(ctx context.Context) {
	// Start a span
	span, _ := opentracing.StartSpanFromContext(ctx, "slow_function")
	defer span.Finish()

	// Simulating unpredictability, just like your coding habits
	randomSleep := time.Duration(rand.Intn(5)) * time.Second // Random sleep between 0-5 seconds
	time.Sleep(randomSleep)
}
