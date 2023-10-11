package middleware

import (
	"github.com/labstack/echo/v4"
)

type (
	JsonResponse struct {
		Message string `json:"message"`
	}
)

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if true {
			res := JsonResponse{
				Message: "Unauthorized",
			}

			return c.JSON(401, &res)
		}

		return next(c)
	}
}
