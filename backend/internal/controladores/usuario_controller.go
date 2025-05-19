package controladores

import (
	"net/http"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UsuarioController maneja los endpoints de usuarios
type UsuarioController struct {
	usuarioService *servicios.UsuarioService
}

// NewUsuarioController crea una nueva instancia de UsuarioController
func NewUsuarioController(usuarioService *servicios.UsuarioService) *UsuarioController {
	return &UsuarioController{
		usuarioService: usuarioService,
	}
}

// Create crea un nuevo usuario
func (c *UsuarioController) Create(ctx *gin.Context) {
	var usuarioReq entidades.NuevoUsuarioRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&usuarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(usuarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear usuario
	id, err := c.usuarioService.Create(&usuarioReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear usuario", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Usuario creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un usuario por su ID
func (c *UsuarioController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener usuario
	usuario, err := c.usuarioService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Usuario no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuario obtenido", usuario))
}

// Update actualiza un usuario
func (c *UsuarioController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var usuario entidades.Usuario

	// Parsear request
	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(usuario); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar usuario
	err = c.usuarioService.Update(id, &usuario)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar usuario", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuario actualizado exitosamente", nil))
}

// Delete elimina un usuario (borrado lógico)
func (c *UsuarioController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar usuario
	err = c.usuarioService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar usuario", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuario eliminado exitosamente", nil))
}

// ListByRol lista usuarios por rol
func (c *UsuarioController) ListByRol(ctx *gin.Context) {
	// Obtener rol de la URL
	rol := ctx.Param("rol")

	// Validar rol
	validRoles := map[string]bool{
		"ADMIN":    true,
		"VENDEDOR": true,
		"CHOFER":   true,
		"CLIENTE":  true,
	}

	if !validRoles[rol] {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Rol inválido", nil))
		return
	}

	// Listar usuarios
	usuarios, err := c.usuarioService.ListByRol(rol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar usuarios", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuarios listados exitosamente", usuarios))
}

// List lista todos los usuarios
func (c *UsuarioController) List(ctx *gin.Context) {
	// Listar usuarios
	usuarios, err := c.usuarioService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar usuarios", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuarios listados exitosamente", usuarios))
}
