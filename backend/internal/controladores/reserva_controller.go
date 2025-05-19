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

// ReservaController maneja los endpoints de reservas
type ReservaController struct {
	reservaService *servicios.ReservaService
}

// NewReservaController crea una nueva instancia de ReservaController
func NewReservaController(reservaService *servicios.ReservaService) *ReservaController {
	return &ReservaController{
		reservaService: reservaService,
	}
}

// Create crea una nueva reserva
func (c *ReservaController) Create(ctx *gin.Context) {
	var reservaReq entidades.NuevaReservaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&reservaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(reservaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Si es una reserva de vendedor, obtener el ID del vendedor del contexto
	if ctx.GetString("rol") == "VENDEDOR" {
		vendedorID := ctx.GetInt("user_id")
		reservaReq.IDVendedor = &vendedorID
	}

	// Crear reserva
	id, err := c.reservaService.Create(&reservaReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear reserva", err))
		return
	}

	// Obtener la reserva creada
	reserva, err := c.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener la reserva creada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Reserva creada exitosamente", reserva))
}

// GetByID obtiene una reserva por su ID
func (c *ReservaController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener reserva
	reserva, err := c.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva obtenida", reserva))
}

// Update actualiza una reserva
func (c *ReservaController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var reservaReq entidades.ActualizarReservaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&reservaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(reservaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Si es una reserva de vendedor, obtener el ID del vendedor del contexto
	if ctx.GetString("rol") == "VENDEDOR" {
		vendedorID := ctx.GetInt("user_id")
		reservaReq.IDVendedor = &vendedorID
	}

	// Actualizar reserva
	err = c.reservaService.Update(id, &reservaReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar reserva", err))
		return
	}

	// Obtener la reserva actualizada
	reserva, err := c.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener la reserva actualizada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva actualizada exitosamente", reserva))
}

// CambiarEstado cambia el estado de una reserva
func (c *ReservaController) CambiarEstado(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var estadoReq entidades.CambiarEstadoReservaRequest

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
	err = c.reservaService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado de la reserva", err))
		return
	}

	// Obtener la reserva actualizada
	reserva, err := c.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener la reserva actualizada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado de la reserva actualizado exitosamente", reserva))
}

// Delete elimina una reserva
func (c *ReservaController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar reserva
	err = c.reservaService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar reserva", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva eliminada exitosamente", nil))
}

// List lista todas las reservas
func (c *ReservaController) List(ctx *gin.Context) {
	// Listar reservas
	reservas, err := c.reservaService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar reservas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reservas listadas exitosamente", reservas))
}

// ListByCliente lista todas las reservas de un cliente
func (c *ReservaController) ListByCliente(ctx *gin.Context) {
	// Parsear ID del cliente de la URL
	idCliente, err := strconv.Atoi(ctx.Param("idCliente"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de cliente inválido", err))
		return
	}

	// Listar reservas del cliente
	reservas, err := c.reservaService.ListByCliente(idCliente)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar reservas del cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reservas del cliente listadas exitosamente", reservas))
}

// ListByTourProgramado lista todas las reservas para un tour programado
func (c *ReservaController) ListByTourProgramado(ctx *gin.Context) {
	// Parsear ID del tour programado de la URL
	idTourProgramado, err := strconv.Atoi(ctx.Param("idTourProgramado"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour programado inválido", err))
		return
	}

	// Listar reservas del tour programado
	reservas, err := c.reservaService.ListByTourProgramado(idTourProgramado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar reservas del tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reservas del tour programado listadas exitosamente", reservas))
}

// ListByFecha lista todas las reservas para una fecha específica
func (c *ReservaController) ListByFecha(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Listar reservas por fecha
	reservas, err := c.reservaService.ListByFecha(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar reservas por fecha", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reservas por fecha listadas exitosamente", reservas))
}

// ListByEstado lista todas las reservas por estado
func (c *ReservaController) ListByEstado(ctx *gin.Context) {
	// Parsear estado de la URL
	estado := ctx.Param("estado")

	// Listar reservas por estado
	reservas, err := c.reservaService.ListByEstado(estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar reservas por estado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reservas por estado listadas exitosamente", reservas))
}

// ListMyReservas lista todas las reservas del cliente autenticado
// ListMyReservas lista todas las reservas del cliente autenticado
// ListMyReservas lista todas las reservas del cliente autenticado
func (c *ReservaController) ListMyReservas(ctx *gin.Context) {
	// Obtener ID del cliente autenticado
	idCliente, exists := ctx.Get("userID") // CORREGIDO: user_id -> userID
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Cliente no autenticado", nil))
		return
	}

	// Listar reservas del cliente
	reservas, err := c.reservaService.ListByCliente(idCliente.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar reservas del cliente", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Mis reservas listadas exitosamente", reservas))
}
