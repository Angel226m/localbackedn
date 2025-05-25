package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TipoTourController maneja los endpoints de tipos de tour
type TipoTourController struct {
	tipoTourService *servicios.TipoTourService
}

// NewTipoTourController crea una nueva instancia de TipoTourController
func NewTipoTourController(tipoTourService *servicios.TipoTourService) *TipoTourController {
	return &TipoTourController{
		tipoTourService: tipoTourService,
	}
}

// Create crea un nuevo tipo de tour
func (c *TipoTourController) Create(ctx *gin.Context) {
	var tipoTourReq entidades.NuevoTipoTourRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tipoTourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tipoTourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear tipo de tour
	id, err := c.tipoTourService.Create(&tipoTourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tipo de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tipo de tour creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un tipo de tour por su ID
func (c *TipoTourController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener tipo de tour
	tipoTour, err := c.tipoTourService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tipo de tour no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de tour obtenido", tipoTour))
}

// Update actualiza un tipo de tour
func (c *TipoTourController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var tipoTourReq entidades.ActualizarTipoTourRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tipoTourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tipoTourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar tipo de tour
	err = c.tipoTourService.Update(id, &tipoTourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tipo de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de tour actualizado exitosamente", nil))
}

// Delete elimina un tipo de tour (borrado lógico)
func (c *TipoTourController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar tipo de tour
	err = c.tipoTourService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar tipo de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de tour eliminado exitosamente", nil))
}

// List lista todos los tipos de tour
func (c *TipoTourController) List(ctx *gin.Context) {
	// Listar tipos de tour
	tiposTour, err := c.tipoTourService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tipos de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipos de tour listados exitosamente", tiposTour))
}

// ListBySede lista todos los tipos de tour de una sede específica
func (c *TipoTourController) ListBySede(ctx *gin.Context) {
	// Parsear ID de la sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar tipos de tour de la sede
	tiposTour, err := c.tipoTourService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tipos de tour de la sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipos de tour de la sede listados exitosamente", tiposTour))
}
