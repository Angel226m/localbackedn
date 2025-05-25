package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TipoPasajeController maneja los endpoints de tipos de pasaje
type TipoPasajeController struct {
	tipoPasajeService *servicios.TipoPasajeService
}

// NewTipoPasajeController crea una nueva instancia de TipoPasajeController
func NewTipoPasajeController(tipoPasajeService *servicios.TipoPasajeService) *TipoPasajeController {
	return &TipoPasajeController{
		tipoPasajeService: tipoPasajeService,
	}
}

// Create crea un nuevo tipo de pasaje
func (c *TipoPasajeController) Create(ctx *gin.Context) {
	var tipoPasajeReq entidades.NuevoTipoPasajeRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tipoPasajeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tipoPasajeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear tipo de pasaje
	id, err := c.tipoPasajeService.Create(&tipoPasajeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tipo de pasaje", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tipo de pasaje creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un tipo de pasaje por su ID
func (c *TipoPasajeController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener tipo de pasaje
	tipoPasaje, err := c.tipoPasajeService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tipo de pasaje no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de pasaje obtenido", tipoPasaje))
}

// Update actualiza un tipo de pasaje
func (c *TipoPasajeController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var tipoPasajeReq entidades.ActualizarTipoPasajeRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tipoPasajeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tipoPasajeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar tipo de pasaje
	err = c.tipoPasajeService.Update(id, &tipoPasajeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tipo de pasaje", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de pasaje actualizado exitosamente", nil))
}

// Delete elimina un tipo de pasaje
func (c *TipoPasajeController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar tipo de pasaje
	err = c.tipoPasajeService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar tipo de pasaje", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipo de pasaje eliminado exitosamente", nil))
}

// ListBySede lista todos los tipos de pasaje de una sede específica
func (c *TipoPasajeController) ListBySede(ctx *gin.Context) {
	// Parsear ID de sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("id_sede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar tipos de pasaje por sede
	tiposPasaje, err := c.tipoPasajeService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tipos de pasaje", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipos de pasaje listados exitosamente", tiposPasaje))
}

// List lista todos los tipos de pasaje
func (c *TipoPasajeController) List(ctx *gin.Context) {
	// Listar tipos de pasaje
	tiposPasaje, err := c.tipoPasajeService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tipos de pasaje", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tipos de pasaje listados exitosamente", tiposPasaje))
}
