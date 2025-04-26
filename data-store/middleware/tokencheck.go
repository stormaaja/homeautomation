package middleware

import (
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/requestvalidators"

	"github.com/gin-gonic/gin"
)

func TokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requestvalidators.IsApiTokenValid(c.Request.Header) {
			c.String(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
	}
}

func Deprecated() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Deprecated endpoint: %s", c.Request.URL.Path)
		c.Next()
	}
}
