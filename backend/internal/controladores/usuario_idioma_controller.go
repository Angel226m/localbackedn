package controladores

import (
	"net/http"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UsuarioIdiomaController maneja los endpoints de la relación usuario-idioma
type UsuarioIdiomaController struct {
	usuarioIdiomaService *servicios.UsuarioIdiomaService
}

// NewUsuarioIdiomaController crea una nueva instancia de UsuarioIdiomaController
func NewUsuarioIdiomaController(usuarioIdiomaService *servicios.UsuarioIdiomaService) *UsuarioIdiomaController {
	return &UsuarioIdiomaController{
		usuarioIdiomaService: usuarioIdiomaService,
	}
}

// GetIdiomasByUsuarioID obtiene todos los idiomas de un usuario
func (c *UsuarioIdiomaController) GetIdiomasByUsuarioID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener idiomas del usuario
	idiomas, err := c.usuarioIdiomaService.GetIdiomasByUsuarioID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener idiomas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idiomas obtenidos", idiomas))
}

// AsignarIdioma asigna un idioma a un usuario
func (c *UsuarioIdiomaController) AsignarIdioma(ctx *gin.Context) {
	// Parsear ID de usuario de la URL
	usuarioID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de usuario inválido", err))
		return
	}

	// Parsear request
	var request struct {
		IdiomaID int    `json:"id_idioma" validate:"required"`
		Nivel    string `json:"nivel"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Asignar idioma
	err = c.usuarioIdiomaService.AsignarIdioma(usuarioID, request.IdiomaID, request.Nivel)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al asignar idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma asignado exitosamente", nil))
}

// DesasignarIdioma elimina la asignación de un idioma a un usuario
func (c *UsuarioIdiomaController) DesasignarIdioma(ctx *gin.Context) {
	// Parsear ID de usuario de la URL
	usuarioID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de usuario inválido", err))
		return
	}

	// Parsear ID de idioma de la URL
	idiomaID, err := strconv.Atoi(ctx.Param("idioma_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de idioma inválido", err))
		return
	}

	// Desasignar idioma
	err = c.usuarioIdiomaService.DesasignarIdioma(usuarioID, idiomaID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al desasignar idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma desasignado exitosamente", nil))
}

// ActualizarIdiomasUsuario actualiza todos los idiomas de un usuario
func (c *UsuarioIdiomaController) ActualizarIdiomasUsuario(ctx *gin.Context) {
	// Parsear ID de usuario de la URL
	usuarioID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de usuario inválido", err))
		return
	}

	// Parsear request
	var request struct {
		IdiomasIDs []int `json:"idiomas_ids" validate:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar idiomas
	err = c.usuarioIdiomaService.ActualizarIdiomasUsuario(usuarioID, request.IdiomasIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar idiomas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idiomas actualizados exitosamente", nil))
}

// GetUsuariosByIdiomaID obtiene todos los usuarios que hablan un idioma específico
func (c *UsuarioIdiomaController) GetUsuariosByIdiomaID(ctx *gin.Context) {
	// Parsear ID de idioma de la URL
	idiomaID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de idioma inválido", err))
		return
	}

	// Obtener usuarios con este idioma
	usuarios, err := c.usuarioIdiomaService.GetUsuariosByIdiomaID(idiomaID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener usuarios", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Usuarios obtenidos", usuarios))
}
