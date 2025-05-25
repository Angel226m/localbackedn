package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MetodoPagoController maneja los endpoints de métodos de pago
type MetodoPagoController struct {
	metodoPagoService *servicios.MetodoPagoService
}

// NewMetodoPagoController crea una nueva instancia de MetodoPagoController
func NewMetodoPagoController(metodoPagoService *servicios.MetodoPagoService) *MetodoPagoController {
	return &MetodoPagoController{
		metodoPagoService: metodoPagoService,
	}
}

// Create crea un nuevo método de pago
func (c *MetodoPagoController) Create(ctx *gin.Context) {
	var metodoPagoReq entidades.NuevoMetodoPagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&metodoPagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(metodoPagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear método de pago
	id, err := c.metodoPagoService.Create(&metodoPagoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear método de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Método de pago creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un método de pago por su ID
func (c *MetodoPagoController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener método de pago
	metodoPago, err := c.metodoPagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Método de pago no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Método de pago obtenido", metodoPago))
}

// Update actualiza un método de pago
func (c *MetodoPagoController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var metodoPagoReq entidades.ActualizarMetodoPagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&metodoPagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(metodoPagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar método de pago
	err = c.metodoPagoService.Update(id, &metodoPagoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar método de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Método de pago actualizado exitosamente", nil))
}

// Delete elimina un método de pago
func (c *MetodoPagoController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar método de pago
	err = c.metodoPagoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar método de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Método de pago eliminado exitosamente", nil))
}

// List lista todos los métodos de pago
func (c *MetodoPagoController) List(ctx *gin.Context) {
	// Listar métodos de pago
	metodosPago, err := c.metodoPagoService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar métodos de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Métodos de pago listados exitosamente", metodosPago))
}

// ListBySede lista todos los métodos de pago de una sede específica
func (c *MetodoPagoController) ListBySede(ctx *gin.Context) {
	// Parsear ID de la sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar métodos de pago de la sede
	metodosPago, err := c.metodoPagoService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar métodos de pago de la sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Métodos de pago de la sede listados exitosamente", metodosPago))
}
