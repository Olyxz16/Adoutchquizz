package server

import (
	"Adoutchquizz/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	e.GET("/", handler.Index)

	return e
}
