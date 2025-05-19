package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware registra información sobre cada solicitud HTTP
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		startTime := time.Now()

		// Procesar solicitud
		c.Next()

		// Tiempo de finalización
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// Obtener datos de la solicitud
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Registrar información
		log.Printf("[%s] %s | %d | %s | %s | %s",
			reqMethod,
			reqURI,
			statusCode,
			clientIP,
			latencyTime,
			c.Errors.String(),
		)
	}
}
