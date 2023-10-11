package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}

func Register(c echo.Context) error {
	return c.String(http.StatusOK, "Register")
}

func User(c echo.Context) error {
	return c.String(http.StatusOK, "User")
}
