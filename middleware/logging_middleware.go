package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs request details in Gin
func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Println("------------------------------------------------------------")
		log.Infof("Incoming request: Method=%s URL=%s", c.Request.Method, c.Request.URL.Path)

		// Process the request
		c.Next()

		log.Println("------------------------------------------------------------")
	}
}
