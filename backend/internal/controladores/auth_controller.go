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
)

// AuthController maneja las peticiones de autenticación
type AuthController struct {
	authService *servicios.AuthService
}

// NewAuthController crea una nueva instancia de AuthController
func NewAuthController(authService *servicios.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login maneja la petición de inicio de sesión
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq entidades.LoginRequest

	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos de entrada inválidos", err))
		return
	}

	// Obtener el parámetro remember_me del query
	rememberMe, _ := strconv.ParseBool(ctx.DefaultQuery("remember_me", "false"))

	// Pasar remember_me al servicio
	loginResp, err := c.authService.Login(&loginReq, rememberMe)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Error de autenticación", err))
		return
	}

	// Configurar cookie para el token JWT
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token",  // Nombre
		loginResp.Token, // Valor
		60*15,           // Tiempo de vida en segundos (15 minutos)
		"/",             // Path
		"",              // Domain (vacío = dominio actual)
		true,            // Secure (false para desarrollo local)
		false,           // HttpOnly (false para permitir acceso desde JS en desarrollo)
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
		"refresh_token",        // Nombre
		loginResp.RefreshToken, // Valor
		refreshExpiry,          // Tiempo de vida en segundos (variable)
		"/",                    // Path
		"",                     // Domain
		true,                   // Secure (false para desarrollo local)
		false,                  // HttpOnly (false para permitir acceso desde JS en desarrollo)
	)

	// Para administradores, no incluir sede en la respuesta
	// Para otros roles, incluir la sede asignada si tiene sede
	var sede *entidades.Sede = nil
	if loginResp.Usuario.Rol != "ADMIN" && loginResp.Usuario.IdSede != nil {
		sede, _ = c.authService.GetSedeByID(*loginResp.Usuario.IdSede)
	}

	// Para depuración: incluir el token en respuesta durante desarrollo
	responseData := gin.H{
		"usuario": loginResp.Usuario,
		"sede":    sede,
	}

	// Solo incluir token en desarrollo para facilitar depuración
	if gin.Mode() != gin.ReleaseMode {
		responseData["token"] = loginResp.Token
	}

	// No devolver los tokens en la respuesta JSON para mayor seguridad en producción
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", responseData))
}

// RefreshToken maneja la regeneración de token usando refresh token
// RefreshToken maneja la regeneración de token usando refresh token
/*func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// Añadir logs para depuración
	fmt.Println("RefreshToken: Iniciando regeneración de token")
	fmt.Printf("Headers recibidos: %v\n", ctx.Request.Header)

	// Obtener refresh token de la cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	fmt.Printf("RefreshToken: Token obtenido de cookie: %v, Error: %v\n", refreshToken != "", err)

	// Si no hay cookie, intentar obtenerlo del Header Authorization
	if err != nil || refreshToken == "" {
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			fmt.Printf("RefreshToken: Token obtenido del header Authorization: %v\n", refreshToken != "")
		}
	}

	// Si aún no tenemos el token, intentar obtenerlo del cuerpo de la solicitud
	if refreshToken == "" {
		var refreshReq struct {
			RefreshToken string `json:"refresh_token"`
		}
		if ctx.ShouldBindJSON(&refreshReq) == nil && refreshReq.RefreshToken != "" {
			refreshToken = refreshReq.RefreshToken
			fmt.Printf("RefreshToken: Token obtenido del body JSON: %v\n", refreshToken != "")
		} else {
			// Si aún no se encuentra el token, probar obtenerlo del form data
			refreshToken = ctx.PostForm("refresh_token")
			fmt.Printf("RefreshToken: Token obtenido de form data: %v\n", refreshToken != "")
		}
	}

	// Verificar si encontramos el token
	if refreshToken == "" {
		fmt.Println("RefreshToken: No se encontró el refresh token en ninguna fuente")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token no proporcionado", nil))
		return
	}

	// Regenerar tokens
	loginResp, err := c.authService.RefreshToken(refreshToken)
	if err != nil {
		fmt.Printf("RefreshToken: Error al regenerar token: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Error al actualizar token", err))
		return
	}

	// Verificar si el refresh token original tenía remember_me activo
	// Esto se puede obtener verificando su tiempo de expiración
	claims, _ := utils.GetRefreshTokenClaims(refreshToken, ctx.MustGet("config").(*config.Config))
	var isRememberMe bool
	if claims != nil {
		// Si el token expira en más de 24 horas desde su emisión, consideramos que tiene remember_me activo
		issuedAt := claims.IssuedAt.Time
		expiresAt := claims.ExpiresAt.Time
		isRememberMe = expiresAt.Sub(issuedAt) > 24*time.Hour
	}

	// Configurar cookie para el nuevo token JWT - CONFIGURACIÓN CONSISTENTE
	ctx.SetSameSite(http.SameSiteNoneMode) // Para permitir CORS con HTTPS/HTTP
	ctx.SetCookie(
		"access_token",  // Nombre
		loginResp.Token, // Valor
		60*15,           // Tiempo de vida en segundos (15 minutos)
		"/",             // Path
		"",              // Domain (vacío = dominio actual)
		true,            // Secure (false para desarrollo local con HTTP/HTTPS mixto)
		false,           // HttpOnly (false para permitir acceso desde JS - necesario para tu flujo)
	)

	// Configurar cookie para el nuevo refresh token manteniendo la misma duración
	var refreshExpiry int
	if isRememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
	}

	ctx.SetSameSite(http.SameSiteNoneMode) // Para permitir CORS con HTTPS/HTTP
	ctx.SetCookie(
		"refresh_token",        // Nombre
		loginResp.RefreshToken, // Valor
		refreshExpiry,          // Tiempo de vida en segundos (variable)
		"/",                    // Path
		"",                     // Domain
		true,                   // Secure (false para desarrollo local con HTTP/HTTPS mixto)
		false,                  // HttpOnly (false para permitir acceso desde JS - necesario para tu flujo)
	)

	fmt.Println("RefreshToken: Cookies establecidas exitosamente")

	// Para administradores, no incluir sede en la respuesta
	// Para otros roles, incluir la sede asignada
	var sede *entidades.Sede = nil
	if loginResp.Usuario.Rol != "ADMIN" && loginResp.Usuario.IdSede != nil {
		sede, _ = c.authService.GetSedeByID(*loginResp.Usuario.IdSede)
	}

	// Retornar la respuesta incluyendo los tokens en desarrollo para facilitar debug
	responseData := gin.H{
		"usuario": loginResp.Usuario,
		"sede":    sede,
	}

	// Solo incluir los tokens en desarrollo
	if gin.Mode() != gin.ReleaseMode {
		responseData["token"] = loginResp.Token
		responseData["refresh_token"] = loginResp.RefreshToken
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Token actualizado exitosamente", responseData))
}
*/

// RefreshToken maneja la regeneración de token usando refresh token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	fmt.Println("RefreshToken: Iniciando regeneración de token")
	fmt.Printf("Headers recibidos: %v\n", ctx.Request.Header)

	// Obtener refresh token de la cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	fmt.Printf("RefreshToken: Token obtenido de cookie: %v, Error: %v\n", refreshToken != "", err)

	// Si no hay cookie, intentar obtenerlo del Header Authorization
	if err != nil || refreshToken == "" {
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			fmt.Printf("RefreshToken: Token obtenido del header Authorization: %v\n", refreshToken != "")
		}
	}

	// Si aún no tenemos el token, intentar obtenerlo del cuerpo de la solicitud
	if refreshToken == "" {
		var refreshReq struct {
			RefreshToken string `json:"refresh_token"`
		}
		if ctx.ShouldBindJSON(&refreshReq) == nil && refreshReq.RefreshToken != "" {
			refreshToken = refreshReq.RefreshToken
			fmt.Printf("RefreshToken: Token obtenido del body JSON: %v\n", refreshToken != "")
		} else {
			// Si aún no se encuentra el token, probar obtenerlo del form data
			refreshToken = ctx.PostForm("refresh_token")
			fmt.Printf("RefreshToken: Token obtenido de form data: %v\n", refreshToken != "")
		}
	}

	// Verificar si encontramos el token
	if refreshToken == "" {
		fmt.Println("RefreshToken: No se encontró el refresh token en ninguna fuente")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token no proporcionado", nil))
		return
	}

	// Validar el refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken, ctx.MustGet("config").(*config.Config))
	if err != nil {
		fmt.Printf("RefreshToken: Error al validar refresh token: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Refresh token inválido", err))
		return
	}

	// Obtener las propiedades del token
	userID := claims.UserID
	userRole := claims.Role
	sedeID := claims.SedeID

	// Determinar si tiene remember_me
	issuedAt := claims.IssuedAt.Time
	expiresAt := claims.ExpiresAt.Time
	isRememberMe := expiresAt.Sub(issuedAt) > 24*time.Hour

	// Generar nuevos tokens
	var token, newRefreshToken string

	// Si tiene sede seleccionada
	if sedeID > 0 {
		token, newRefreshToken, err = c.authService.GenerateTokensForAdminWithSede(userID, sedeID, isRememberMe)
	} else {
		// Sin sede seleccionada
		token, newRefreshToken, err = c.authService.GenerateTokensWithoutDb(userID, userRole, isRememberMe)
	}

	if err != nil {
		fmt.Printf("RefreshToken: Error al generar nuevos tokens: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al generar nuevos tokens", err))
		return
	}

	// Configurar cookie para el nuevo token JWT
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token", // Nombre
		token,          // Valor
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

	fmt.Println("RefreshToken: Cookies establecidas exitosamente")

	// Obtener sede si es necesario
	var sede *entidades.Sede = nil
	if sedeID > 0 {
		sede, _ = c.authService.GetSedeByID(sedeID)
	}

	// Crear usuario simplificado para la respuesta
	usuario := gin.H{
		"id_usuario": userID,
		"rol":        userRole,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Token actualizado exitosamente", gin.H{
		"usuario": usuario,
		"sede":    sede,
	}))
}

// CheckStatus verifica si el usuario tiene una sesión válida
func (c *AuthController) CheckStatus(ctx *gin.Context) {
	// Logs para depuración
	fmt.Println("Ejecutando CheckStatus")

	// Obtener usuario del contexto (puesto por el middleware de autenticación)
	userID, exists := ctx.Get("userID")
	fmt.Printf("CheckStatus: userID en contexto = %v (exists: %v)\n", userID, exists)

	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	// SOLUCIÓN TEMPORAL: Usar respuesta simplificada para desarrollo
	// Esto es para evitar el error en DB mientras lo solucionas
	if gin.Mode() != gin.ReleaseMode {
		// En desarrollo, generar respuesta simplificada
		userRole, _ := ctx.Get("userRole")
		fmt.Printf("CheckStatus: usando respuesta simplificada para desarrollo (userID=%v, userRole=%v)\n",
			userID, userRole)

		ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuario autenticado (modo desarrollo)", gin.H{
			"usuario": gin.H{
				"id_usuario": userID.(int),
				"rol":        userRole.(string),
				"nombres":    "Usuario",
				"apellidos":  "Temporal",
				"correo":     "demo@ejemplo.com",
			},
			"sede": nil,
		}))
		return
	}

	// En producción, usar el flujo normal
	// Obtener datos del usuario
	usuario, err := c.authService.GetUserByID(userID.(int))
	if err != nil {
		fmt.Printf("CheckStatus: Error al obtener usuario: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no encontrado", err))
		return
	}

	// Para administradores, no incluir sede en la respuesta a menos que tenga una sede seleccionada temporalmente
	// Para otros roles, incluir la sede asignada
	var sede *entidades.Sede = nil

	// Verificar si hay una sede seleccionada en la sesión (para administradores)
	sedeID, sedeExists := ctx.Get("sedeID")
	if usuario.Rol == "ADMIN" && sedeExists && sedeID.(int) > 0 {
		sede, _ = c.authService.GetSedeByID(sedeID.(int))
	} else if usuario.Rol != "ADMIN" && usuario.IdSede != nil {
		sede, _ = c.authService.GetSedeByID(*usuario.IdSede)
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuario autenticado", gin.H{
		"usuario": usuario,
		"sede":    sede,
	}))
}

// GetUserSedes obtiene todas las sedes disponibles para el usuario administrador
// GetUserSedes obtiene todas las sedes disponibles para el usuario administrador
//
/*
func (c *AuthController) GetUserSedes(ctx *gin.Context) {
	// Obtener usuario del contexto (puesto por el middleware de autenticación)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	// Obtener datos del usuario
	usuario, err := c.authService.GetUserByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no encontrado", err))
		return
	}

	// Verificar que sea administrador
	if usuario.Rol != "ADMIN" {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("Acceso denegado. Solo administradores pueden ver todas las sedes", nil))
		return
	}

	// Obtener todas las sedes disponibles
	sedes, err := c.authService.GetAllSedes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener sedes", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sedes obtenidas exitosamente", gin.H{
		"sedes": sedes,
	}))
}*/
// GetUserSedes modificado para menor dependencia en la base de datos
func (c *AuthController) GetUserSedes(ctx *gin.Context) {
	// Obtener rol directamente del contexto (ya validado por el middleware)
	userRole, exists := ctx.Get("userRole")
	if !exists || userRole != "ADMIN" {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("Acceso denegado. Solo administradores pueden ver todas las sedes", nil))
		return
	}

	// Obtener todas las sedes disponibles
	sedes, err := c.authService.GetAllSedes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener sedes", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sedes obtenidas exitosamente", gin.H{
		"sedes": sedes,
	}))
}

// SelectSede permite a un administrador seleccionar una sede para la sesión
/*func (c *AuthController) SelectSede(ctx *gin.Context) {
	var selectSedeReq struct {
		IdSede int `json:"id_sede" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&selectSedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos de entrada inválidos", err))
		return
	}

	// Obtener usuario del contexto (puesto por el middleware de autenticación)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	// Obtener datos del usuario
	usuario, err := c.authService.GetUserByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no encontrado", err))
		return
	}

	// Verificar que sea administrador
	if usuario.Rol != "ADMIN" {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("Acceso denegado. Solo administradores pueden seleccionar sede", nil))
		return
	}

	// Verificar que la sede exista
	sede, err := c.authService.GetSedeByID(selectSedeReq.IdSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Sede no encontrada", err))
		return
	}

	// Verificar el estado actual del refresh token para mantener su configuración
	refreshToken, _ := ctx.Cookie("refresh_token")
	var isRememberMe bool
	if refreshToken != "" {
		claims, _ := utils.GetRefreshTokenClaims(refreshToken, ctx.MustGet("config").(*config.Config))
		if claims != nil {
			issuedAt := claims.IssuedAt.Time
			expiresAt := claims.ExpiresAt.Time
			isRememberMe = expiresAt.Sub(issuedAt) > 24*time.Hour
		}
	}

	// Actualizar la sesión con la sede seleccionada
	token, newRefreshToken, err := c.authService.GenerateTokensWithSede(usuario.ID, selectSedeReq.IdSede, isRememberMe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al actualizar sesión", err))
		return
	}

	// Actualizar cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token", // Nombre
		token,          // Valor
		60*15,          // Tiempo de vida en segundos (15 minutos)
		"/",            // Path
		"",             // Domain
		true,           // Secure
		false,          // HttpOnly
	)

	// Configurar cookie para el refresh token manteniendo la misma duración
	var refreshExpiry int
	if isRememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
	}

	// Actualizar refresh token
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede seleccionada exitosamente", gin.H{
		"sede": sede,
	}))
}*/

// SelectSede permite a un administrador seleccionar una sede para la sesión
func (c *AuthController) SelectSede(ctx *gin.Context) {
	var selectSedeReq struct {
		IdSede int `json:"id_sede" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&selectSedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos de entrada inválidos", err))
		return
	}

	// Obtener información directamente del contexto (ya validado por el middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole != "ADMIN" {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("Acceso denegado. Solo administradores pueden seleccionar sede", nil))
		return
	}

	// Verificar que la sede exista
	sede, err := c.authService.GetSedeByID(selectSedeReq.IdSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Sede no encontrada", err))
		return
	}

	// Verificar el estado actual del refresh token para mantener su configuración
	refreshToken, _ := ctx.Cookie("refresh_token")
	var isRememberMe bool
	if refreshToken != "" {
		claims, _ := utils.GetRefreshTokenClaims(refreshToken, ctx.MustGet("config").(*config.Config))
		if claims != nil {
			issuedAt := claims.IssuedAt.Time
			expiresAt := claims.ExpiresAt.Time
			isRememberMe = expiresAt.Sub(issuedAt) > 24*time.Hour
		}
	}

	// Actualizar la sesión con la sede seleccionada - usando userID del contexto
	// La clave del cambio está aquí: ya no busca el usuario en la base de datos
	token, newRefreshToken, err := c.authService.GenerateTokensForAdminWithSede(userID.(int), selectSedeReq.IdSede, isRememberMe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al actualizar sesión", err))
		return
	}

	// Actualizar cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"access_token", // Nombre
		token,          // Valor
		60*15,          // Tiempo de vida en segundos (15 minutos)
		"/",            // Path
		"",             // Domain
		true,           // Secure
		false,          // HttpOnly
	)

	// Configurar cookie para el refresh token manteniendo la misma duración
	var refreshExpiry int
	if isRememberMe {
		refreshExpiry = 60 * 60 * 24 * 7 // 7 días si remember_me estaba activo
	} else {
		refreshExpiry = 60 * 60 // 1 hora si remember_me no estaba activo
	}

	// Actualizar refresh token
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede seleccionada exitosamente", gin.H{
		"sede": sede,
	}))
}

// ChangePassword maneja el cambio de contraseña
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos de entrada inválidos", err))
		return
	}

	// Obtener el ID del usuario del contexto (establecido por el middleware de autenticación)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	// Llamar al servicio para cambiar la contraseña
	err := c.authService.ChangePassword(userID.(int), changePassReq.CurrentPassword, changePassReq.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar contraseña", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Contraseña actualizada exitosamente", nil))
}

// Logout maneja el cierre de sesión
func (c *AuthController) Logout(ctx *gin.Context) {
	// Eliminar cookies estableciendo tiempo de expiración en el pasado
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", "", -1, "/", "", true, false)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, false)

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sesión cerrada exitosamente", nil))
}
