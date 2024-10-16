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

    /***************/
    /*    INDEX    */
    /***************/
	e.GET("/", handler.Index)
    
    /***************/
    /*    VIDEO    */
    /***************/
    e.GET("/video/:uid", handler.Video)
    e.GET("/video/new", handler.Video)
    e.GET("/video/update/:uid", handler.Video)
    e.POST("/video/:action", handler.VideoPost) 
    e.DELETE("/video/clip/:uid", handler.VideoDelete)

	return e
}
