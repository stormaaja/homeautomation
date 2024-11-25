package middleware

import (
	"net/http"
	"stormaaja/go-ha/security/requestvalidators"

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
