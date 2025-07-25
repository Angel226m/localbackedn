/*
package controladores

import (

	"fmt" // Añadido para logs de depuración
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"strings"
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
		var loginReq struct {
			Correo     string `json:"correo" validate:"required,email"`
			Contrasena string `json:"contrasena" validate:"required"`
		}

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
		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie(
			"access_token", // Nombre
			token,          // Valor
			60*15,          // Tiempo de vida en segundos (15 minutos)
			"/",            // Path
			"",             // Domain (vacío = dominio actual)
			true,           // Secure (true para producción)
			false,          // HttpOnly (false para permitir acceso desde JS)
		)

		// Configurar cookie para el refresh token con duración variable
		var refreshExpiry int
		if rememberMe {
			refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me es true
		} else {
			refreshExpiry = 60 * 60 // 1 hora si remember_me es false
		}

		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie(
			"refresh_token", // Nombre
			refreshToken,    // Valor
			refreshExpiry,   // Tiempo de vida en segundos (variable)
			"/",             // Path
			"",              // Domain
			true,            // Secure (true para producción)
			false,           // HttpOnly (false para permitir acceso desde JS)
		)

		// Para depuración: incluir el token en respuesta durante desarrollo
		responseData := gin.H{
			"usuario": gin.H{
				"id_cliente":       cliente.ID,
				"nombres":          cliente.Nombres,
				"apellidos":        cliente.Apellidos,
				"nombre_completo":  cliente.Nombres + " " + cliente.Apellidos,
				"tipo_documento":   cliente.TipoDocumento,
				"numero_documento": cliente.NumeroDocumento,
				"correo":           cliente.Correo,
				"numero_celular":   cliente.NumeroCelular, // Incluir el número de celular
				"rol":              "CLIENTE",
			},
		}

		// Solo incluir token en desarrollo para facilitar depuración
		if gin.Mode() != gin.ReleaseMode {
			responseData["token"] = token
			responseData["refresh_token"] = refreshToken
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", responseData))
	}

// RefreshToken renueva los tokens de un cliente

	func (c *ClienteController) RefreshToken(ctx *gin.Context) {
		fmt.Println("RefreshToken Cliente: Iniciando regeneración de token")
		fmt.Printf("Headers recibidos: %v\n", ctx.Request.Header)

		// Obtener refresh token de la cookie
		refreshToken, err := ctx.Cookie("refresh_token")
		fmt.Printf("RefreshToken Cliente: Token obtenido de cookie: %v, Error: %v\n", refreshToken != "", err)

		// Si no hay cookie, intentar obtenerlo del Header Authorization
		if err != nil || refreshToken == "" {
			authHeader := ctx.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
				fmt.Printf("RefreshToken Cliente: Token obtenido del header Authorization: %v\n", refreshToken != "")
			}
		}

		// Si aún no tenemos el token, intentar obtenerlo del cuerpo de la solicitud
		if refreshToken == "" {
			var refreshReq struct {
				RefreshToken string `json:"refresh_token"`
			}
			if ctx.ShouldBindJSON(&refreshReq) == nil && refreshReq.RefreshToken != "" {
				refreshToken = refreshReq.RefreshToken
				fmt.Printf("RefreshToken Cliente: Token obtenido del body JSON: %v\n", refreshToken != "")
			} else {
				// Si aún no se encuentra el token, probar obtenerlo del form data
				refreshToken = ctx.PostForm("refresh_token")
				fmt.Printf("RefreshToken Cliente: Token obtenido de form data: %v\n", refreshToken != "")
			}
		}

		// Verificar si encontramos el token
		if refreshToken == "" {
			fmt.Println("RefreshToken Cliente: No se encontró el refresh token en ninguna fuente")
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token no proporcionado", nil))
			return
		}

		// Renovar tokens
		newToken, newRefreshToken, cliente, err := c.clienteService.RefreshClienteToken(refreshToken)
		if err != nil {
			fmt.Printf("RefreshToken Cliente: Error al actualizar token: %v\n", err)
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
		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie(
			"access_token", // Nombre
			newToken,       // Valor
			60*15,          // Tiempo de vida en segundos (15 minutos)
			"/",            // Path
			"",             // Domain (vacío = dominio actual)
			true,           // Secure
			false,          // HttpOnly
		)

		// Configurar cookie para el nuevo refresh token manteniendo la misma duración
		var refreshExpiry int
		if isRememberMe {
			refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
		} else {
			refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
		}

		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie(
			"refresh_token", // Nombre
			newRefreshToken, // Valor
			refreshExpiry,   // Tiempo de vida en segundos (variable)
			"/",             // Path
			"",              // Domain
			true,            // Secure
			false,           // HttpOnly
		)

		fmt.Println("RefreshToken Cliente: Cookies establecidas exitosamente")

		// Crear respuesta para el cliente
		responseData := gin.H{
			"usuario": gin.H{
				"id_cliente":       cliente.ID,
				"nombres":          cliente.Nombres,
				"apellidos":        cliente.Apellidos,
				"nombre_completo":  cliente.Nombres + " " + cliente.Apellidos,
				"tipo_documento":   cliente.TipoDocumento,
				"numero_documento": cliente.NumeroDocumento,
				"correo":           cliente.Correo,
				"numero_celular":   cliente.NumeroCelular, // Incluir el número de celular
				"rol":              "CLIENTE",
			},
		}

		// Solo incluir token en desarrollo para facilitar depuración
		if gin.Mode() != gin.ReleaseMode {
			responseData["token"] = newToken
			responseData["refresh_token"] = newRefreshToken
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse("Token actualizado exitosamente", responseData))
	}

// ChangePassword cambia la contraseña de un cliente

	func (c *ClienteController) ChangePassword(ctx *gin.Context) {
		fmt.Println("ChangePassword Cliente: Iniciando cambio de contraseña")

		// Parsear ID del cliente del contexto (establecido por el middleware de autenticación)
		clienteIDValue, exists := ctx.Get("userID")
		if !exists {
			fmt.Println("ChangePassword Cliente: Usuario no autenticado")
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
			return
		}

		clienteID, ok := clienteIDValue.(int)
		if !ok {
			fmt.Println("ChangePassword Cliente: Error en identificación de usuario")
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error en identificación de usuario", nil))
			return
		}

		var changePassReq struct {
			CurrentPassword string `json:"current_password" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required,min=6"`
		}

		// Parsear request
		if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
			fmt.Printf("ChangePassword Cliente: Datos inválidos: %v\n", err)
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
			return
		}

		// Cambiar contraseña
		err := c.clienteService.ChangePassword(clienteID, changePassReq.CurrentPassword, changePassReq.NewPassword)
		if err != nil {
			fmt.Printf("ChangePassword Cliente: Error al cambiar contraseña: %v\n", err)
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar contraseña", err))
			return
		}

		fmt.Println("ChangePassword Cliente: Contraseña actualizada exitosamente")
		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Contraseña actualizada exitosamente", nil))
	}

// Logout cierra la sesión de un cliente

	func (c *ClienteController) Logout(ctx *gin.Context) {
		fmt.Println("Logout Cliente: Cerrando sesión")

		// Eliminar cookies estableciendo tiempo de expiración en el pasado
		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie("access_token", "", -1, "/", "", true, false)
		ctx.SetCookie("refresh_token", "", -1, "/", "", true, false)

		fmt.Println("Logout Cliente: Cookies eliminadas exitosamente")
		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Sesión cerrada exitosamente", nil))
	}
*/
package controladores

import (
	"fmt" // Añadido para logs de depuración
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"strings"
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

// GetByDocumento obtiene un cliente por su tipo y número de documento
func (c *ClienteController) GetByDocumento(ctx *gin.Context) {
	tipoDocumento := ctx.Query("tipo")
	numeroDocumento := ctx.Query("numero")

	if tipoDocumento == "" || numeroDocumento == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tipo y número de documento son requeridos", nil))
		return
	}

	cliente, err := c.clienteService.GetByDocumento(tipoDocumento, numeroDocumento)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Cliente no encontrado", err))
		return
	}

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

// UpdateDatosEmpresa actualiza solo los datos de empresa de un cliente
func (c *ClienteController) UpdateDatosEmpresa(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var datosReq entidades.ActualizarDatosEmpresaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&datosReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(datosReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar datos de empresa
	err = c.clienteService.UpdateDatosEmpresa(id, &datosReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar datos de empresa", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Datos de empresa actualizados exitosamente", nil))
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
	searchType := ctx.Query("type") // Nuevo: tipo de búsqueda (name, doc)

	var clientes []*entidades.Cliente
	var err error

	// Determinar tipo de búsqueda
	if query != "" {
		if searchType == "doc" {
			// Buscar por documento
			clientes, err = c.clienteService.SearchByDocumento(query)
		} else {
			// Por defecto buscar por nombre
			clientes, err = c.clienteService.SearchByName(query)
		}
	} else {
		// Listar todos
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
	var loginReq struct {
		Correo     string `json:"correo" validate:"required,email"`
		Contrasena string `json:"contrasena" validate:"required"`
	}

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
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token", // Nombre
		token,          // Valor
		60*15,          // Tiempo de vida en segundos (15 minutos)
		"/",            // Path
		"",             // Domain (vacío = dominio actual)
		true,           // Secure (true para producción)
		false,          // HttpOnly (false para permitir acceso desde JS)
	)

	// Configurar cookie para el refresh token con duración variable
	var refreshExpiry int
	if rememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me es true
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me es false
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"refresh_token", // Nombre
		refreshToken,    // Valor
		refreshExpiry,   // Tiempo de vida en segundos (variable)
		"/",             // Path
		"",              // Domain
		true,            // Secure (true para producción)
		false,           // HttpOnly (false para permitir acceso desde JS)
	)

	// Preparar respuesta según tipo de cliente (persona natural o empresa)
	responseData := gin.H{
		"usuario": gin.H{
			"id_cliente":       cliente.ID,
			"tipo_documento":   cliente.TipoDocumento,
			"numero_documento": cliente.NumeroDocumento,
			"correo":           cliente.Correo,
			"numero_celular":   cliente.NumeroCelular,
			"rol":              "CLIENTE",
		},
	}

	// Añadir campos específicos según tipo de cliente
	if cliente.TipoDocumento == "RUC" {
		// Para empresas
		responseData["usuario"].(gin.H)["razon_social"] = cliente.RazonSocial
		responseData["usuario"].(gin.H)["direccion_fiscal"] = cliente.DireccionFiscal
		responseData["usuario"].(gin.H)["nombre_completo"] = cliente.RazonSocial
	} else {
		// Para personas naturales
		responseData["usuario"].(gin.H)["nombres"] = cliente.Nombres
		responseData["usuario"].(gin.H)["apellidos"] = cliente.Apellidos
		responseData["usuario"].(gin.H)["nombre_completo"] = cliente.Nombres + " " + cliente.Apellidos
	}

	// Solo incluir token en desarrollo para facilitar depuración
	if gin.Mode() != gin.ReleaseMode {
		responseData["token"] = token
		responseData["refresh_token"] = refreshToken
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", responseData))
}

// RefreshToken renueva los tokens de un cliente
func (c *ClienteController) RefreshToken(ctx *gin.Context) {
	fmt.Println("RefreshToken Cliente: Iniciando regeneración de token")
	fmt.Printf("Headers recibidos: %v\n", ctx.Request.Header)

	// Obtener refresh token de la cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	fmt.Printf("RefreshToken Cliente: Token obtenido de cookie: %v, Error: %v\n", refreshToken != "", err)

	// Si no hay cookie, intentar obtenerlo del Header Authorization
	if err != nil || refreshToken == "" {
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			fmt.Printf("RefreshToken Cliente: Token obtenido del header Authorization: %v\n", refreshToken != "")
		}
	}

	// Si aún no tenemos el token, intentar obtenerlo del cuerpo de la solicitud
	if refreshToken == "" {
		var refreshReq struct {
			RefreshToken string `json:"refresh_token"`
		}
		if ctx.ShouldBindJSON(&refreshReq) == nil && refreshReq.RefreshToken != "" {
			refreshToken = refreshReq.RefreshToken
			fmt.Printf("RefreshToken Cliente: Token obtenido del body JSON: %v\n", refreshToken != "")
		} else {
			// Si aún no se encuentra el token, probar obtenerlo del form data
			refreshToken = ctx.PostForm("refresh_token")
			fmt.Printf("RefreshToken Cliente: Token obtenido de form data: %v\n", refreshToken != "")
		}
	}

	// Verificar si encontramos el token
	if refreshToken == "" {
		fmt.Println("RefreshToken Cliente: No se encontró el refresh token en ninguna fuente")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token no proporcionado", nil))
		return
	}

	// Renovar tokens
	newToken, newRefreshToken, cliente, err := c.clienteService.RefreshClienteToken(refreshToken)
	if err != nil {
		fmt.Printf("RefreshToken Cliente: Error al actualizar token: %v\n", err)
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
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token", // Nombre
		newToken,       // Valor
		60*15,          // Tiempo de vida en segundos (15 minutos)
		"/",            // Path
		"",             // Domain (vacío = dominio actual)
		true,           // Secure
		false,          // HttpOnly
	)

	// Configurar cookie para el nuevo refresh token manteniendo la misma duración
	var refreshExpiry int
	if isRememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"refresh_token", // Nombre
		newRefreshToken, // Valor
		refreshExpiry,   // Tiempo de vida en segundos (variable)
		"/",             // Path
		"",              // Domain
		true,            // Secure
		false,           // HttpOnly
	)

	fmt.Println("RefreshToken Cliente: Cookies establecidas exitosamente")

	// Preparar respuesta según tipo de cliente (persona natural o empresa)
	responseData := gin.H{
		"usuario": gin.H{
			"id_cliente":       cliente.ID,
			"tipo_documento":   cliente.TipoDocumento,
			"numero_documento": cliente.NumeroDocumento,
			"correo":           cliente.Correo,
			"numero_celular":   cliente.NumeroCelular,
			"rol":              "CLIENTE",
		},
	}

	// Añadir campos específicos según tipo de cliente
	if cliente.TipoDocumento == "RUC" {
		// Para empresas
		responseData["usuario"].(gin.H)["razon_social"] = cliente.RazonSocial
		responseData["usuario"].(gin.H)["direccion_fiscal"] = cliente.DireccionFiscal
		responseData["usuario"].(gin.H)["nombre_completo"] = cliente.RazonSocial
	} else {
		// Para personas naturales
		responseData["usuario"].(gin.H)["nombres"] = cliente.Nombres
		responseData["usuario"].(gin.H)["apellidos"] = cliente.Apellidos
		responseData["usuario"].(gin.H)["nombre_completo"] = cliente.Nombres + " " + cliente.Apellidos
	}

	// Solo incluir token en desarrollo para facilitar depuración
	if gin.Mode() != gin.ReleaseMode {
		responseData["token"] = newToken
		responseData["refresh_token"] = newRefreshToken
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Token actualizado exitosamente", responseData))
}

// ChangePassword cambia la contraseña de un cliente
func (c *ClienteController) ChangePassword(ctx *gin.Context) {
	fmt.Println("ChangePassword Cliente: Iniciando cambio de contraseña")

	// Parsear ID del cliente del contexto (establecido por el middleware de autenticación)
	clienteIDValue, exists := ctx.Get("userID")
	if !exists {
		fmt.Println("ChangePassword Cliente: Usuario no autenticado")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	clienteID, ok := clienteIDValue.(int)
	if !ok {
		fmt.Println("ChangePassword Cliente: Error en identificación de usuario")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error en identificación de usuario", nil))
		return
	}

	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	// Parsear request
	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		fmt.Printf("ChangePassword Cliente: Datos inválidos: %v\n", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Cambiar contraseña
	err := c.clienteService.ChangePassword(clienteID, changePassReq.CurrentPassword, changePassReq.NewPassword)
	if err != nil {
		fmt.Printf("ChangePassword Cliente: Error al cambiar contraseña: %v\n", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar contraseña", err))
		return
	}

	fmt.Println("ChangePassword Cliente: Contraseña actualizada exitosamente")
	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Contraseña actualizada exitosamente", nil))
}

// Logout cierra la sesión de un cliente
func (c *ClienteController) Logout(ctx *gin.Context) {
	fmt.Println("Logout Cliente: Cerrando sesión")

	// Eliminar cookies estableciendo tiempo de expiración en el pasado
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", "", -1, "/", "", true, false)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, false)

	fmt.Println("Logout Cliente: Cookies eliminadas exitosamente")
	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sesión cerrada exitosamente", nil))
}

// GetPerfilCliente obtiene el perfil del cliente autenticado
func (c *ClienteController) GetPerfilCliente(ctx *gin.Context) {
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

	// Obtener cliente
	cliente, err := c.clienteService.GetByID(clienteID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Cliente no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Perfil obtenido", cliente))
}
