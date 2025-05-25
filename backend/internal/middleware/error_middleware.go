package middleware

import (
	"net/http"
	"sistema-toursseft/internal/utils"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware maneja errores globales en la aplicación
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Solo procesar si hay errores y no se ha enviado una respuesta
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err

			// Verificar si es un error de validación
			switch v := err.(type) {
			case *utils.ValidationError:
				c.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", v))
			default:
				c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error interno del servidor", err))
			}
		}
	}
}
