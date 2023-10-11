package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

const (
	port        = 8080
	metricsPort = 8081
)

type Config struct {
	Echo *echo.Echo
}

func main() {
	app := Config{
		Echo: echo.New(),
	}

	app.Echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))

	app.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://**", "https://**"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	app.routes()

	app.Echo.Use(echoprometheus.NewMiddleware("videobufferapi")) // adds middleware to gather metrics

	go func() {
		metrics := echo.New()                                // this Echo will run on separate port 8081
		metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metrics.Start(fmt.Sprintf(":%d", metricsPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	if err := app.Echo.Start(fmt.Sprintf(":%d", port)); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
