package controladores

import (
	"net/http"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthController maneja los endpoints de autenticación
type AuthController struct {
	authService *servicios.AuthService
}

// NewAuthController crea una nueva instancia de AuthController
func NewAuthController(authService *servicios.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login maneja la autenticación de usuarios
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq entidades.LoginRequest

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
	loginResp, err := c.authService.Login(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autenticación fallida", err))
		return
	}

	// Devolver tokens y datos de usuario
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Login exitoso", loginResp))
}

// RefreshToken renueva el token de acceso
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var refreshReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Parsear request
	if err := ctx.ShouldBindJSON(&refreshReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Renovar token
	loginResp, err := c.authService.RefreshToken(refreshReq.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token de actualización inválido", err))
		return
	}

	// Devolver nuevos tokens
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Token renovado exitosamente", loginResp))
}

// ChangePassword cambia la contraseña del usuario
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// Obtener ID de usuario del contexto (establecido por el middleware de autenticación)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	// Parsear request
	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Cambiar contraseña
	err := c.authService.ChangePassword(userID.(int), changePassReq.CurrentPassword, changePassReq.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar contraseña", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Contraseña cambiada exitosamente", nil))
}
