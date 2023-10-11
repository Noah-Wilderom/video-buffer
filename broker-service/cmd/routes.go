package main

import (
	"github.com/Noah-Wilderom/video-buffer/broker-service/handlers"
	"github.com/Noah-Wilderom/video-buffer/broker-service/middleware"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Config) routes() {

	app.Echo.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	app.Echo.POST("api/auth/login", handlers.Login)
	app.Echo.POST("api/auth/register", handlers.Register)

	app.Echo.GET("api/auth/user", handlers.User, middleware.Authenticated)
}
