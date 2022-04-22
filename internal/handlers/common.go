package handlers

import (
	"net/http"

	_ "github.com/agandreev/crime-app-auth/docs" // register my own docs
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// InitCommonRoutes registers handlers and uses middlewares.
func InitCommonRoutes(e *echo.Echo) http.Handler {
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	return e
}
