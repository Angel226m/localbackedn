package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EmbarcacionController maneja los endpoints de embarcaciones
type EmbarcacionController struct {
	embarcacionService *servicios.EmbarcacionService
}

// NewEmbarcacionController crea una nueva instancia de EmbarcacionController
func NewEmbarcacionController(embarcacionService *servicios.EmbarcacionService) *EmbarcacionController {
	return &EmbarcacionController{
		embarcacionService: embarcacionService,
	}
}

// Create crea una nueva embarcación
func (c *EmbarcacionController) Create(ctx *gin.Context) {
	var embarcacionReq entidades.NuevaEmbarcacionRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&embarcacionReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(embarcacionReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear embarcación
	id, err := c.embarcacionService.Create(&embarcacionReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear embarcación", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Embarcación creada exitosamente", gin.H{"id": id}))
}

// GetByID obtiene una embarcación por su ID
func (c *EmbarcacionController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener embarcación
	embarcacion, err := c.embarcacionService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Embarcación no encontrada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Embarcación obtenida", embarcacion))
}

// Update actualiza una embarcación
func (c *EmbarcacionController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var embarcacionReq entidades.ActualizarEmbarcacionRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&embarcacionReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(embarcacionReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar embarcación
	err = c.embarcacionService.Update(id, &embarcacionReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar embarcación", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Embarcación actualizada exitosamente", nil))
}

// Delete elimina una embarcación (borrado lógico)
func (c *EmbarcacionController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar embarcación
	err = c.embarcacionService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar embarcación", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Embarcación eliminada exitosamente", nil))
}

// List lista todas las embarcaciones
func (c *EmbarcacionController) List(ctx *gin.Context) {
	// Listar embarcaciones
	embarcaciones, err := c.embarcacionService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar embarcaciones", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Embarcaciones listadas exitosamente", embarcaciones))
}

// ListByChofer lista todas las embarcaciones de un chofer específico
func (c *EmbarcacionController) ListByChofer(ctx *gin.Context) {
	// Parsear ID del chofer de la URL
	idChofer, err := strconv.Atoi(ctx.Param("idChofer"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Listar embarcaciones del chofer
	embarcaciones, err := c.embarcacionService.ListByChofer(idChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar embarcaciones del chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Embarcaciones del chofer listadas exitosamente", embarcaciones))
}
