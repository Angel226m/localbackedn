package controladores

import (
	"net/http"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PagoController maneja los endpoints de pagos
type PagoController struct {
	pagoService *servicios.PagoService
}

// NewPagoController crea una nueva instancia de PagoController
func NewPagoController(pagoService *servicios.PagoService) *PagoController {
	return &PagoController{
		pagoService: pagoService,
	}
}

// Create crea un nuevo pago
func (c *PagoController) Create(ctx *gin.Context) {
	var pagoReq entidades.NuevoPagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&pagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(pagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear pago
	id, err := c.pagoService.Create(&pagoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Pago creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un pago por su ID
func (c *PagoController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener pago
	pago, err := c.pagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Pago no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pago obtenido", pago))
}

// Update actualiza un pago
func (c *PagoController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var pagoReq entidades.ActualizarPagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&pagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(pagoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar pago
	err = c.pagoService.Update(id, &pagoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pago actualizado exitosamente", nil))
}

// CambiarEstado cambia el estado de un pago
func (c *PagoController) CambiarEstado(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var estadoReq entidades.CambiarEstadoPagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&estadoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(estadoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Cambiar estado
	err = c.pagoService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado del pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado del pago actualizado exitosamente", nil))
}

// Delete elimina un pago
func (c *PagoController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar pago
	err = c.pagoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pago eliminado exitosamente", nil))
}

// List lista todos los pagos
func (c *PagoController) List(ctx *gin.Context) {
	// Listar pagos
	pagos, err := c.pagoService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar pagos", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pagos listados exitosamente", pagos))
}

// ListByReserva lista todos los pagos de una reserva específica
func (c *PagoController) ListByReserva(ctx *gin.Context) {
	// Parsear ID de reserva de la URL
	idReserva, err := strconv.Atoi(ctx.Param("idReserva"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de reserva inválido", err))
		return
	}

	// Listar pagos por reserva
	pagos, err := c.pagoService.ListByReserva(idReserva)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar pagos por reserva", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pagos listados exitosamente", pagos))
}

// ListByFecha lista todos los pagos de una fecha específica
func (c *PagoController) ListByFecha(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Listar pagos por fecha
	pagos, err := c.pagoService.ListByFecha(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar pagos por fecha", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pagos listados exitosamente", pagos))
}

// GetTotalPagadoByReserva obtiene el total pagado de una reserva específica
func (c *PagoController) GetTotalPagadoByReserva(ctx *gin.Context) {
	// Parsear ID de reserva de la URL
	idReserva, err := strconv.Atoi(ctx.Param("idReserva"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de reserva inválido", err))
		return
	}

	// Obtener total pagado
	totalPagado, err := c.pagoService.GetTotalPagadoByReserva(idReserva)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al obtener total pagado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Total pagado obtenido exitosamente", gin.H{"total_pagado": totalPagado}))
}

// ListByEstado lista todos los pagos con un estado específico
func (c *PagoController) ListByEstado(ctx *gin.Context) {
	// Parsear estado de la URL
	estado := ctx.Param("estado")

	// Listar pagos por estado
	pagos, err := c.pagoService.ListByEstado(estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar pagos por estado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pagos listados exitosamente", pagos))
}

// ListByCliente lista todos los pagos de un cliente específico
func (c *PagoController) ListByCliente(ctx *gin.Context) {
	// Parsear ID de cliente de la URL
	idCliente, err := strconv.Atoi(ctx.Param("idCliente"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de cliente inválido", err))
		return
	}

	// Listar pagos por cliente
	pagos, err := c.pagoService.ListByCliente(idCliente)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar pagos por cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pagos listados exitosamente", pagos))
}
