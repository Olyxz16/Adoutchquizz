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
    e.GET("/video", handler.Video)
    e.GET("/video/new", handler.Video)
    e.POST("/video/set", handler.VideoPost) 

    /*e.GET("/get/:group/:value/*", handler.Get)

    e.GET("/set/:value", handler.Set)
    e.POST("/set/:value", handler.SetForm);
    
    e.DELETE("/clip/:id", handler.RemoveClip);*/

	return e
}
