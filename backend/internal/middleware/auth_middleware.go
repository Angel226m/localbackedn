/*
package middleware

import (
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/utils"
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

		// Guardar también con los nombres que usa el controlador
		ctx.Set("user_id", claims.UserID)
		ctx.Set("rol", claims.Role)

		// Para el rol ADMIN, establecer sede_id con un valor especial
		if claims.Role == "ADMIN" {
			// Usar -1 para indicar "todas las sedes" (o cualquier valor distinto de 0)
			ctx.Set("sede_id", -1)
		} else {
			// Para otros roles, buscar la sede del usuario
			// NOTA: Aquí necesitarías implementar la lógica para obtener la sede del usuario
			// Esto dependerá de cómo manejas la asignación de sedes en tu sistema
			// Por ahora, establecemos un valor temporal
			// TODO: Implementar lógica para obtener sede_id del usuario
			sedeID := obtenerSedeIDDelUsuario(claims.UserID)
			ctx.Set("sede_id", sedeID)
		}

		ctx.Next()
	}
}

// Función para obtener la sede_id del usuario
// Debes implementar esta función según tu lógica de negocio
func obtenerSedeIDDelUsuario(userID int) int {
	// Aquí deberías consultar la base de datos o algún servicio
	// para obtener la sede del usuario
	// Por ejemplo:
	// 1. Consultar la tabla de usuarios para obtener su sede_id
	// 2. Si no tiene sede asignada, retornar 0 (o un valor por defecto)

	// Por ahora, retornamos un valor temporal
	// IMPLEMENTAR: Reemplazar esto con lógica real
	return 1 // Valor temporal para no bloquear el flujo
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
*/

// auth_middleware.go
package middleware

import (
	"fmt"
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware crea un middleware para autenticación JWT
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Agregar logs para depuración
		fmt.Println("AuthMiddleware: Iniciando verificación de autenticación")
		fmt.Printf("Headers recibidos: %v\n", ctx.Request.Header)
		fmt.Printf("Cookies recibidas: %v\n", ctx.Request.Cookies())

		// Obtener token de la cookie
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			fmt.Printf("AuthMiddleware: Cookie access_token no encontrada: %v\n", err)
		} else if len(tokenString) > 20 {
			fmt.Printf("Token de cookie (primeros 20 caracteres): %s...\n", tokenString[:20]+"...")
		} else {
			fmt.Printf("Token de cookie: %v\n", tokenString)
		}

		// Si no hay cookie, intentar obtenerlo del header Authorization
		if err != nil || tokenString == "" {
			authHeader := ctx.GetHeader("Authorization")
			fmt.Printf("AuthMiddleware: Authorization header: %v\n", authHeader != "")

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
			tokenString = tokenParts[1]
			fmt.Printf("AuthMiddleware: Token obtenido del header Authorization\n")
		}

		// Parsear y validar el token
		token, err := jwt.ParseWithClaims(tokenString, &utils.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verificar que el algoritmo de firma es el esperado
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Algoritmo de firma no esperado", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(config.JWTSecret), nil
		})

		// Mejorar el manejo de errores para mostrar más detalles
		if err != nil {
			fmt.Printf("AuthMiddleware: Error validando token: %v, tipo: %T\n", err, err)

			// Verificar causas comunes de error
			if validationErr, ok := err.(*jwt.ValidationError); ok {
				if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
					fmt.Println("AuthMiddleware: El token ha expirado")
					ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token expirado", nil))
				} else if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					fmt.Println("AuthMiddleware: La firma del token es inválida")
					ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Firma del token inválida", nil))
				} else if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("AuthMiddleware: El token está malformado")
					ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token malformado", nil))
				} else if validationErr.Errors&jwt.ValidationErrorClaimsInvalid != 0 {
					fmt.Println("AuthMiddleware: Los claims del token son inválidos")
					ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Claims del token inválidos", nil))
				} else {
					ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token inválido", nil))
				}
			} else {
				ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token inválido", nil))
			}

			ctx.Abort()
			return
		}

		// Verificar que el token es válido y obtener los claims
		if claims, ok := token.Claims.(*utils.TokenClaims); ok && token.Valid {
			fmt.Printf("AuthMiddleware: Token válido para usuario ID: %d, Rol: %s\n", claims.UserID, claims.Role)

			// Guardar claims en el contexto
			ctx.Set("userID", claims.UserID)
			ctx.Set("userRole", claims.Role)

			// Para administradores con sede seleccionada
			if claims.Role == "ADMIN" && claims.SedeID > 0 {
				// El admin ha seleccionado una sede específica
				ctx.Set("sedeID", claims.SedeID)
				ctx.Set("adminConSede", true) // Bandera para indicar admin con sede específica
				fmt.Printf("AuthMiddleware: Admin con sede seleccionada: %d\n", claims.SedeID)
			} else if claims.Role == "ADMIN" {
				// Admin sin sede específica (puede ver todas)
				ctx.Set("adminSinSede", true) // Bandera para indicar admin sin sede específica
				fmt.Println("AuthMiddleware: Admin sin sede específica")
			} else if claims.SedeID > 0 {
				// Otros roles solo pueden ver su sede asignada
				ctx.Set("sedeID", claims.SedeID)
				fmt.Printf("AuthMiddleware: Usuario con sede asignada: %d\n", claims.SedeID)
			}

			ctx.Next()

			fmt.Printf("AuthMiddleware: Contexto establecido - userID: %d, userRole: %s\n",
				claims.UserID, claims.Role)
		} else {
			fmt.Println("AuthMiddleware: Error: token inválido o claims no coinciden con el tipo esperado")
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token inválido", nil))
			ctx.Abort()
			return
		}
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

// SedeMiddleware verifica que el usuario tenga acceso a la sede solicitada
func SedeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener la sede solicitada del parámetro de la URL
		sedeIDParam := ctx.Param("idSede")
		if sedeIDParam == "" {
			// Si no hay parámetro de sede, continuar (algunas rutas no necesitan este parámetro)
			ctx.Next()
			return
		}

		// Verificar si es admin sin sede específica (puede acceder a cualquier sede)
		adminSinSede, existsFlag := ctx.Get("adminSinSede")
		if existsFlag && adminSinSede.(bool) {
			// Los admins sin sede específica pueden acceder a cualquier sede
			ctx.Next()
			return
		}

		// Obtener sedeID del token
		sedeIDToken, existsSedeToken := ctx.Get("sedeID")
		if !existsSedeToken {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene una sede asignada", nil))
			ctx.Abort()
			return
		}

		// Si es admin con sede, debe coincidir con la sede solicitada
		// Si es otro rol, debe ser la sede asignada
		if sedeIDParam != "" && sedeIDToken.(int) != 0 {
			// Convertir sedeIDParam a entero
			sedeIDSolicitada, err := strconv.Atoi(sedeIDParam)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
				ctx.Abort()
				return
			}

			// Verificar si coincide con la sede del token
			if sedeIDSolicitada != sedeIDToken.(int) {
				ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta sede", nil))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
