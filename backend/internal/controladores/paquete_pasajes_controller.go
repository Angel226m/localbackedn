package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaquetePasajesController maneja los endpoints de paquetes de pasajes
type PaquetePasajesController struct {
	paquetePasajesService *servicios.PaquetePasajesService
}

// NewPaquetePasajesController crea una nueva instancia de PaquetePasajesController
func NewPaquetePasajesController(paquetePasajesService *servicios.PaquetePasajesService) *PaquetePasajesController {
	return &PaquetePasajesController{
		paquetePasajesService: paquetePasajesService,
	}
}

// Create crea un nuevo paquete de pasajes
func (c *PaquetePasajesController) Create(ctx *gin.Context) {
	var paqueteReq entidades.NuevoPaquetePasajesRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&paqueteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(paqueteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear paquete de pasajes
	id, err := c.paquetePasajesService.Create(&paqueteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear paquete de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Paquete de pasajes creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un paquete de pasajes por su ID
func (c *PaquetePasajesController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener paquete de pasajes
	paquete, err := c.paquetePasajesService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Paquete de pasajes no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquete de pasajes obtenido", paquete))
}

// Update actualiza un paquete de pasajes
func (c *PaquetePasajesController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var paqueteReq entidades.ActualizarPaquetePasajesRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&paqueteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(paqueteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar paquete de pasajes
	err = c.paquetePasajesService.Update(id, &paqueteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar paquete de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquete de pasajes actualizado exitosamente", nil))
}

// Delete elimina un paquete de pasajes
func (c *PaquetePasajesController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar paquete de pasajes
	err = c.paquetePasajesService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar paquete de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquete de pasajes eliminado exitosamente", nil))
}

// ListBySede lista todos los paquetes de pasajes de una sede específica
func (c *PaquetePasajesController) ListBySede(ctx *gin.Context) {
	// Parsear ID de sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("id_sede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar paquetes de pasajes por sede
	paquetes, err := c.paquetePasajesService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar paquetes de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquetes de pasajes listados exitosamente", paquetes))
}

// ListByTipoTour lista todos los paquetes de pasajes de un tipo de tour específico
func (c *PaquetePasajesController) ListByTipoTour(ctx *gin.Context) {
	// Parsear ID del tipo de tour de la URL
	idTipoTour, err := strconv.Atoi(ctx.Param("id_tipo_tour"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Listar paquetes de pasajes por tipo de tour
	paquetes, err := c.paquetePasajesService.ListByTipoTour(idTipoTour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar paquetes de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquetes de pasajes listados exitosamente", paquetes))
}

// List lista todos los paquetes de pasajes
func (c *PaquetePasajesController) List(ctx *gin.Context) {
	// Listar paquetes de pasajes
	paquetes, err := c.paquetePasajesService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar paquetes de pasajes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Paquetes de pasajes listados exitosamente", paquetes))
}
