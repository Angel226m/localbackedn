package rutas

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ClienteHandlers contiene funciones de manejo específicas para clientes
type ClienteHandlers struct {
	reservaService     *servicios.ReservaService
	clienteService     *servicios.ClienteService
	mercadoPagoService *servicios.MercadoPagoService
	baseURLProduccion  string
	baseURLDesarrollo  string
}

// NewClienteHandlers crea una nueva instancia de manejadores para clientes
func NewClienteHandlers(
	reservaService *servicios.ReservaService,
	clienteService *servicios.ClienteService,
	mercadoPagoService *servicios.MercadoPagoService,
) *ClienteHandlers {
	return &ClienteHandlers{
		reservaService:     reservaService,
		clienteService:     clienteService,
		mercadoPagoService: mercadoPagoService,
		baseURLProduccion:  "https://reservas.angelproyect.com",
		baseURLDesarrollo:  "https://localhost:5174",
	}
}

// GetReservaDetalle obtiene el detalle de una reserva para un cliente
func (h *ClienteHandlers) GetReservaDetalle(ctx *gin.Context) {
	reservaID := ctx.Param("id")
	clienteID := ctx.GetInt("userID")

	// Obtener la reserva
	id, _ := strconv.Atoi(reservaID)
	reserva, err := h.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
		return
	}

	// Verificar que la reserva pertenece al cliente
	if reserva.IDCliente != clienteID {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
		return
	}

	// Mostrar la reserva
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva obtenida exitosamente", reserva))
}

// determinarBaseURL determina qué URL base usar basado en el encabezado Origin
func (h *ClienteHandlers) determinarBaseURL(ctx *gin.Context) string {
	origin := ctx.GetHeader("Origin")

	if origin == "" {
		// Si no hay Origin, intentar con Referer
		referer := ctx.GetHeader("Referer")
		if referer != "" {
			if strings.Contains(referer, "localhost") {
				return h.baseURLDesarrollo
			} else {
				return h.baseURLProduccion
			}
		}
		// Si no hay Referer, usar producción por defecto
		return h.baseURLProduccion
	}

	// Si hay Origin, verificar si es localhost
	if strings.Contains(origin, "localhost") {
		return h.baseURLDesarrollo
	}

	return h.baseURLProduccion
}

// PagarReserva crea una preferencia de pago para una reserva existente
func (h *ClienteHandlers) PagarReserva(ctx *gin.Context) {
	reservaID := ctx.Param("id")
	clienteID := ctx.GetInt("userID")

	// Obtener la reserva
	id, _ := strconv.Atoi(reservaID)
	reserva, err := h.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
		return
	}

	// Verificar que la reserva pertenece al cliente
	if reserva.IDCliente != clienteID {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
		return
	}

	// Obtener datos del cliente
	cliente, err := h.clienteService.GetByID(clienteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener datos del cliente", err))
		return
	}

	// Determinar la URL base según el entorno
	baseURL := h.determinarBaseURL(ctx)

	// Usar la URL específica para el proceso de pago
	frontendURL := baseURL + "/proceso-pago"

	// Crear la preferencia de pago para esta reserva
	response, err := h.mercadoPagoService.GeneratePreferenceForExistingReserva(
		id, reserva.TotalPagar, cliente, frontendURL)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al generar preferencia de pago", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Preferencia de pago generada exitosamente", response))
}

// CancelarReserva cancela una reserva de un cliente
func (h *ClienteHandlers) CancelarReserva(ctx *gin.Context) {
	reservaID := ctx.Param("id")
	clienteID := ctx.GetInt("userID")

	// Obtener la reserva
	id, _ := strconv.Atoi(reservaID)
	reserva, err := h.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
		return
	}

	// Verificar que la reserva pertenece al cliente
	if reserva.IDCliente != clienteID {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
		return
	}

	// Crear el request para cambiar estado
	estadoReq := entidades.CambiarEstadoReservaRequest{
		Estado: "CANCELADA",
	}

	// Actualizar estado directamente con el servicio
	err = h.reservaService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cancelar la reserva", err))
		return
	}

	// Obtener la reserva actualizada
	reservaActualizada, err := h.reservaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener la reserva actualizada", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva cancelada exitosamente", reservaActualizada))
}

// ReservarConMercadoPago es un wrapper para el controlador de reservas con Mercado Pago
func (h *ClienteHandlers) ReservarConMercadoPago(ctx *gin.Context, reservaController interface{}) {
	// Determinar la URL base según el entorno
	baseURL := h.determinarBaseURL(ctx)

	// Usar la URL específica para el proceso de pago
	frontendURL := baseURL + "/proceso-pago"

	// Guardar la URL en el contexto para que el controlador pueda acceder a ella
	ctx.Set("frontendURL", frontendURL)

	// Llamar al controlador original
	rc, ok := reservaController.(interface{ ReservarConMercadoPago(*gin.Context) })
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error interno del servidor", nil))
		return
	}

	rc.ReservarConMercadoPago(ctx)
}
