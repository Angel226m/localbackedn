package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CanalVentaController maneja los endpoints de canales de venta
type CanalVentaController struct {
	canalVentaService *servicios.CanalVentaService
}

// NewCanalVentaController crea una nueva instancia de CanalVentaController
func NewCanalVentaController(canalVentaService *servicios.CanalVentaService) *CanalVentaController {
	return &CanalVentaController{
		canalVentaService: canalVentaService,
	}
}

// Create crea un nuevo canal de venta
func (c *CanalVentaController) Create(ctx *gin.Context) {
	var canalReq entidades.NuevoCanalVentaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&canalReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(canalReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear canal de venta
	id, err := c.canalVentaService.Create(&canalReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear canal de venta", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Canal de venta creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un canal de venta por su ID
func (c *CanalVentaController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener canal de venta
	canal, err := c.canalVentaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Canal de venta no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Canal de venta obtenido", canal))
}

// Update actualiza un canal de venta
func (c *CanalVentaController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var canalReq entidades.ActualizarCanalVentaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&canalReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(canalReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar canal de venta
	err = c.canalVentaService.Update(id, &canalReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar canal de venta", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Canal de venta actualizado exitosamente", nil))
}

// Delete elimina un canal de venta
func (c *CanalVentaController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar canal de venta
	err = c.canalVentaService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar canal de venta", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Canal de venta eliminado exitosamente", nil))
}

// List lista todos los canales de venta
func (c *CanalVentaController) List(ctx *gin.Context) {
	// Listar canales de venta
	canales, err := c.canalVentaService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar canales de venta", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Canales de venta listados exitosamente", canales))
}

// ListBySede lista todos los canales de venta de una sede específica
func (c *CanalVentaController) ListBySede(ctx *gin.Context) {
	// Parsear ID de la sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar canales de venta de la sede
	canales, err := c.canalVentaService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar canales de venta de la sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Canales de venta de la sede listados exitosamente", canales))
}
