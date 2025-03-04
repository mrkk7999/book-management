package http

import (
	"book-management/controller"
	"book-management/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetUpRouter(ctrl *controller.Controller, log *logrus.Logger) *gin.Engine {
	var (
		router = gin.New()
	)
	router.Use(middleware.LoggingMiddleware(log))
	router.GET("/books", ctrl.GetAllBooks)
	router.GET("/books/:id", ctrl.GetBookByID)
	router.POST("/books", ctrl.CreateBook)
	router.PUT("/books/:id", ctrl.UpdateBook)
	router.DELETE("/books/:id", ctrl.DeleteBook)

	return router

}
