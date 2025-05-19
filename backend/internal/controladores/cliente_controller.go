package controladores

import (
	"net/http"
	"sistema-tours/internal/config"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"
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
