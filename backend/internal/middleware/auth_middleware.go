package middleware

import (
	"net/http"
	"sistema-tours/internal/config"
	"sistema-tours/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware crea un middleware para autenticación JWT
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener token de autorización
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token no proporcionado", nil))
			ctx.Abort()
			return
		}

		// Verificar formato Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Formato de token inválido", nil))
			ctx.Abort()
			return
		}

		// Extraer token
		tokenString := tokenParts[1]

		// Validar token
		claims, err := utils.ValidateToken(tokenString, config)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token inválido", err))
			ctx.Abort()
			return
		}

		// Guardar claims en el contexto
		ctx.Set("userID", claims.UserID)
		ctx.Set("userEmail", claims.Email)
		ctx.Set("userRole", claims.Role)

		ctx.Next()
	}
}

// RoleMiddleware crea un middleware para restricción por rol
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener rol del contexto (establecido por AuthMiddleware)
		userRole, exists := ctx.Get("userRole")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
			ctx.Abort()
			return
		}

		// Verificar si tiene acceso
		hasAccess := false
		for _, role := range roles {
			if userRole.(string) == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permisos para acceder a este recurso", nil))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
