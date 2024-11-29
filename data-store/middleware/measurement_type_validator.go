package middleware

import "github.com/gin-gonic/gin"

var validMeasurementTypes = []string{"temperature", "electricity_consumption"}

func Includes(value string) bool {
	for _, v := range validMeasurementTypes {
		if v == value {
			return true
		}
	}
	return false
}

func MeasurementTypeValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		measurementType := c.Param("measurement")
		if Includes(measurementType) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid measurement type"})
	}
}
