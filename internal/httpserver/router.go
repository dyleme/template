package httpserver

import (
	"github.com/labstack/echo/v4"

	"github.com/dyleme/template/internal/handler/example"
)

func Route(exmplHandler example.Handler) *echo.Echo {
	e := echo.New()
	e.Group("/example")

	e.PUT("/:id", exmplHandler.Update)
	e.GET("/:id", exmplHandler.Get)

	return e
}
