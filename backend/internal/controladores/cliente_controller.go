/*
package controladores

import (

	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"

)

// ClienteController maneja los endpoints de clientes

	type ClienteController struct {
		clienteService *servicios.ClienteService
		config         *config.Config
	}

// NewClienteController crea una nueva instancia de ClienteController

	func NewClienteController(clienteService *servicios.ClienteService, config *config.Config) *ClienteController {
		return &ClienteController{
			clienteService: clienteService,
			config:         config,
		}
	}

// Create crea un nuevo cliente

	func (c *ClienteController) Create(ctx *gin.Context) {
		var clienteReq entidades.NuevoClienteRequest

		// Parsear request
		if err := ctx.ShouldBindJSON(&clienteReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
			return
		}

		// Validar datos
		if err := utils.ValidateStruct(clienteReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
			return
		}

		// Crear cliente
		id, err := c.clienteService.Create(&clienteReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear cliente", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusCreated, utils.SuccessResponse("Cliente creado exitosamente", gin.H{"id": id}))
	}

// GetByID obtiene un cliente por su ID

	func (c *ClienteController) GetByID(ctx *gin.Context) {
		// Parsear ID de la URL
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
			return
		}

		// Obtener cliente
		cliente, err := c.clienteService.GetByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Cliente no encontrado", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente obtenido", cliente))
	}

// Update actualiza un cliente

	func (c *ClienteController) Update(ctx *gin.Context) {
		// Parsear ID de la URL
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
			return
		}

		var clienteReq entidades.ActualizarClienteRequest

		// Parsear request
		if err := ctx.ShouldBindJSON(&clienteReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
			return
		}

		// Validar datos
		if err := utils.ValidateStruct(clienteReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
			return
		}

		// Actualizar cliente
		err = c.clienteService.Update(id, &clienteReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar cliente", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente actualizado exitosamente", nil))
	}

// Delete elimina un cliente

	func (c *ClienteController) Delete(ctx *gin.Context) {
		// Parsear ID de la URL
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
			return
		}

		// Eliminar cliente
		err = c.clienteService.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar cliente", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente eliminado exitosamente", nil))
	}

// List lista todos los clientes

	func (c *ClienteController) List(ctx *gin.Context) {
		// Obtener parámetro de búsqueda
		query := ctx.Query("search")

		var clientes []*entidades.Cliente
		var err error

		// Buscar por nombre o listar todos
		if query != "" {
			clientes, err = c.clienteService.SearchByName(query)
		} else {
			clientes, err = c.clienteService.List()
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar clientes", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Clientes listados exitosamente", clientes))
	}

// Login maneja el inicio de sesión de un cliente

	func (c *ClienteController) Login(ctx *gin.Context) {
		var loginReq entidades.LoginClienteRequest

		// Parsear request
		if err := ctx.ShouldBindJSON(&loginReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
			return
		}

		// Validar datos
		if err := utils.ValidateStruct(loginReq); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
			return
		}

		// Intentar login
		cliente, err := c.clienteService.Login(loginReq.Correo, loginReq.Contrasena)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Credenciales incorrectas", err))
			return
		}

		// CAMBIO IMPORTANTE: Asegurar que el rol sea exactamente como lo espera el middleware RoleMiddleware
		usuarioEquivalente := &entidades.Usuario{
			ID:     cliente.ID,
			Correo: cliente.Correo,
			Rol:    "CLIENTE", // Asegurar que coincida exactamente con lo que espera el RoleMiddleware
		}

		// Generar token JWT
		token, err := utils.GenerateJWT(usuarioEquivalente, c.config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al generar token", err))
			return
		}

		// Generar refresh token
		refreshToken, err := utils.GenerateRefreshToken(usuarioEquivalente, c.config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al generar refresh token", err))
			return
		}

		// Respuesta exitosa - NO cambiar el formato de esta respuesta
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"usuario": gin.H{
				"id_cliente":       cliente.ID,
				"nombres":          cliente.Nombres,
				"apellidos":        cliente.Apellidos,
				"nombre_completo":  cliente.Nombres + " " + cliente.Apellidos,
				"tipo_documento":   cliente.TipoDocumento,
				"numero_documento": cliente.NumeroDocumento,
				"correo":           cliente.Correo,
				"rol":              "CLIENTE",
			},
		}))
	}
*/

package controladores

import (
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// ClienteController maneja los endpoints de clientes
type ClienteController struct {
	clienteService *servicios.ClienteService
	config         *config.Config
}

// NewClienteController crea una nueva instancia de ClienteController
func NewClienteController(clienteService *servicios.ClienteService, config *config.Config) *ClienteController {
	return &ClienteController{
		clienteService: clienteService,
		config:         config,
	}
}

// Create crea un nuevo cliente
func (c *ClienteController) Create(ctx *gin.Context) {
	var clienteReq entidades.NuevoClienteRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&clienteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(clienteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear cliente
	id, err := c.clienteService.Create(&clienteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Cliente creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un cliente por su ID
func (c *ClienteController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener cliente
	cliente, err := c.clienteService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Cliente no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente obtenido", cliente))
}

// Update actualiza un cliente
func (c *ClienteController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var clienteReq entidades.ActualizarClienteRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&clienteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(clienteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar cliente
	err = c.clienteService.Update(id, &clienteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente actualizado exitosamente", nil))
}

// Delete elimina un cliente
func (c *ClienteController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar cliente
	err = c.clienteService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cliente eliminado exitosamente", nil))
}

// List lista todos los clientes
func (c *ClienteController) List(ctx *gin.Context) {
	// Obtener parámetro de búsqueda
	query := ctx.Query("search")

	var clientes []*entidades.Cliente
	var err error

	// Buscar por nombre o listar todos
	if query != "" {
		clientes, err = c.clienteService.SearchByName(query)
	} else {
		clientes, err = c.clienteService.List()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar clientes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Clientes listados exitosamente", clientes))
}

// Login maneja el inicio de sesión de un cliente
func (c *ClienteController) Login(ctx *gin.Context) {
	var loginReq entidades.LoginClienteRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Obtener el parámetro remember_me del query
	rememberMe, _ := strconv.ParseBool(ctx.DefaultQuery("remember_me", "false"))

	// Intentar login
	cliente, token, refreshToken, err := c.clienteService.Login(loginReq.Correo, loginReq.Contrasena, rememberMe)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Credenciales incorrectas", err))
		return
	}

	// Configurar cookie para el token JWT
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",                    // Nombre
		token,                             // Valor
		60*15,                             // Tiempo de vida en segundos (15 minutos)
		"/",                               // Path
		"",                                // Domain (vacío = dominio actual)
		ctx.Request.URL.Scheme == "https", // Secure (true en HTTPS)
		true,                              // HttpOnly
	)

	// Configurar cookie para el refresh token con duración variable
	var refreshExpiry int
	if rememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me es true
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me es false
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",                   // Nombre
		refreshToken,                      // Valor
		refreshExpiry,                     // Tiempo de vida en segundos (variable)
		"/",                               // Path
		"",                                // Domain
		ctx.Request.URL.Scheme == "https", // Secure (true en HTTPS)
		true,                              // HttpOnly
	)

	// Respuesta exitosa - NO cambiar el formato de esta respuesta
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"usuario": gin.H{
			"id_cliente":       cliente.ID,
			"nombres":          cliente.Nombres,
			"apellidos":        cliente.Apellidos,
			"nombre_completo":  cliente.Nombres + " " + cliente.Apellidos,
			"tipo_documento":   cliente.TipoDocumento,
			"numero_documento": cliente.NumeroDocumento,
			"correo":           cliente.Correo,
			"rol":              "CLIENTE",
		},
	}))
}

// RefreshToken renueva los tokens de un cliente
func (c *ClienteController) RefreshToken(ctx *gin.Context) {
	// Obtener refresh token de la cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		// Si no hay cookie, intentar obtenerlo del cuerpo de la solicitud para compatibilidad
		var refreshReq struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := ctx.ShouldBindJSON(&refreshReq); err != nil || refreshReq.RefreshToken == "" {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token no proporcionado", err))
			return
		}
		refreshToken = refreshReq.RefreshToken
	}

	// Renovar tokens
	newToken, newRefreshToken, cliente, err := c.clienteService.RefreshClienteToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Error al actualizar token", err))
		return
	}

	// Verificar si el refresh token original tenía remember_me activo
	// Analizamos el token original para obtener sus claims
	token, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.config.JWTRefreshSecret), nil
	})

	var isRememberMe bool
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		// Convertir Unix timestamps a time.Time para poder compararlos
		issuedAt := time.Unix(claims.IssuedAt, 0)
		expiresAt := time.Unix(claims.ExpiresAt, 0)

		// Si el token expira en más de 24 horas desde su emisión, consideramos que tiene remember_me activo
		isRememberMe = expiresAt.Sub(issuedAt) > 24*time.Hour
	}

	// Configurar cookie para el nuevo token JWT
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",                    // Nombre
		newToken,                          // Valor
		60*15,                             // Tiempo de vida en segundos (15 minutos)
		"/",                               // Path
		"",                                // Domain
		ctx.Request.URL.Scheme == "https", // Secure (true en HTTPS)
		true,                              // HttpOnly
	)

	// Configurar cookie para el nuevo refresh token manteniendo la misma duración
	var refreshExpiry int
	if isRememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",                   // Nombre
		newRefreshToken,                   // Valor
		refreshExpiry,                     // Tiempo de vida en segundos (variable)
		"/",                               // Path
		"",                                // Domain
		ctx.Request.URL.Scheme == "https", // Secure (true en HTTPS)
		true,                              // HttpOnly
	)

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Token actualizado exitosamente", gin.H{
		"token":         newToken,
		"refresh_token": newRefreshToken,
		"usuario": gin.H{
			"id_cliente":       cliente.ID,
			"nombres":          cliente.Nombres,
			"apellidos":        cliente.Apellidos,
			"nombre_completo":  cliente.Nombres + " " + cliente.Apellidos,
			"tipo_documento":   cliente.TipoDocumento,
			"numero_documento": cliente.NumeroDocumento,
			"correo":           cliente.Correo,
			"rol":              "CLIENTE",
		},
	}))
}

// ChangePassword cambia la contraseña de un cliente
func (c *ClienteController) ChangePassword(ctx *gin.Context) {
	// Parsear ID del cliente del contexto (establecido por el middleware de autenticación)
	clienteIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}
	clienteID, ok := clienteIDValue.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error en identificación de usuario", nil))
		return
	}

	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	// Parsear request
	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Cambiar contraseña
	err := c.clienteService.ChangePassword(clienteID, changePassReq.CurrentPassword, changePassReq.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar contraseña", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Contraseña actualizada exitosamente", nil))
}

// Logout cierra la sesión de un cliente
func (c *ClienteController) Logout(ctx *gin.Context) {
	// Eliminar cookies estableciendo tiempo de expiración en el pasado
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("access_token", "", -1, "/", "", ctx.Request.URL.Scheme == "https", true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", ctx.Request.URL.Scheme == "https", true)

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sesión cerrada exitosamente", nil))
}
