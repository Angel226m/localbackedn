package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ComprobantePagoController maneja los endpoints de comprobantes de pago
type ComprobantePagoController struct {
	comprobantePagoService *servicios.ComprobantePagoService
}

// NewComprobantePagoController crea una nueva instancia de ComprobantePagoController
func NewComprobantePagoController(comprobantePagoService *servicios.ComprobantePagoService) *ComprobantePagoController {
	return &ComprobantePagoController{
		comprobantePagoService: comprobantePagoService,
	}
}

// Create crea un nuevo comprobante de pago
func (c *ComprobantePagoController) Create(ctx *gin.Context) {
	var comprobanteReq entidades.NuevoComprobantePagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&comprobanteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(comprobanteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Si no se especifica la sede, usar la sede del usuario autenticado
	if comprobanteReq.IDSede == 0 && ctx.GetInt("sede_id") > 0 {
		comprobanteReq.IDSede = ctx.GetInt("sede_id")
	}

	// Crear comprobante de pago
	id, err := c.comprobantePagoService.Create(&comprobanteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear comprobante de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Comprobante de pago creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un comprobante de pago por su ID
func (c *ComprobantePagoController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener comprobante de pago
	comprobante, err := c.comprobantePagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Comprobante de pago no encontrado", err))
		return
	}

	// Verificar acceso según el rol
	if ctx.GetString("rol") != "ADMIN" && comprobante.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para ver este comprobante", nil))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobante de pago obtenido", comprobante))
}

// GetByTipoAndNumero obtiene un comprobante de pago por su tipo y número
func (c *ComprobantePagoController) GetByTipoAndNumero(ctx *gin.Context) {
	// Parsear tipo y número de los query params
	tipo := ctx.Query("tipo")
	numero := ctx.Query("numero")

	// Validar parámetros
	if tipo == "" || numero == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tipo y número son requeridos", nil))
		return
	}

	// Obtener comprobante de pago
	comprobante, err := c.comprobantePagoService.GetByTipoAndNumero(tipo, numero)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Comprobante de pago no encontrado", err))
		return
	}

	// Verificar acceso según el rol
	if ctx.GetString("rol") != "ADMIN" && comprobante.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para ver este comprobante", nil))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobante de pago obtenido", comprobante))
}

// Update actualiza un comprobante de pago
func (c *ComprobantePagoController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Verificar que el comprobante existe y el usuario tiene acceso
	comprobante, err := c.comprobantePagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Comprobante de pago no encontrado", err))
		return
	}

	// Verificar acceso según el rol
	if ctx.GetString("rol") != "ADMIN" && comprobante.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para modificar este comprobante", nil))
		return
	}

	var comprobanteReq entidades.ActualizarComprobantePagoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&comprobanteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(comprobanteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Verificar que la sede es la misma del usuario para no-administradores
	if ctx.GetString("rol") != "ADMIN" && comprobanteReq.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para cambiar la sede del comprobante", nil))
		return
	}

	// Actualizar comprobante de pago
	err = c.comprobantePagoService.Update(id, &comprobanteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar comprobante de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobante de pago actualizado exitosamente", nil))
}

// CambiarEstado cambia el estado de un comprobante de pago
func (c *ComprobantePagoController) CambiarEstado(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Verificar que el comprobante existe y el usuario tiene acceso
	comprobante, err := c.comprobantePagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Comprobante de pago no encontrado", err))
		return
	}

	// Verificar acceso según el rol
	if ctx.GetString("rol") != "ADMIN" && comprobante.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para modificar este comprobante", nil))
		return
	}

	var estadoReq entidades.CambiarEstadoComprobanteRequest

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
	err = c.comprobantePagoService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado del comprobante de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado del comprobante de pago actualizado exitosamente", nil))
}

// Delete elimina un comprobante de pago
func (c *ComprobantePagoController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Verificar que el comprobante existe y el usuario tiene acceso
	comprobante, err := c.comprobantePagoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Comprobante de pago no encontrado", err))
		return
	}

	// Verificar acceso según el rol
	if ctx.GetString("rol") != "ADMIN" && comprobante.IDSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para eliminar este comprobante", nil))
		return
	}

	// Eliminar comprobante de pago
	err = c.comprobantePagoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar comprobante de pago", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobante de pago eliminado exitosamente", nil))
}

// List lista todos los comprobantes de pago
func (c *ComprobantePagoController) List(ctx *gin.Context) {
	var comprobantes []*entidades.ComprobantePago
	var err error

	// Si es ADMIN, listar todos los comprobantes
	if ctx.GetString("rol") == "ADMIN" {
		comprobantes, err = c.comprobantePagoService.List()
	} else {
		// Si no es ADMIN, listar solo los comprobantes de su sede
		sedeID := ctx.GetInt("sede_id")
		if sedeID <= 0 {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Usuario no tiene sede asignada", nil))
			return
		}
		comprobantes, err = c.comprobantePagoService.ListBySede(sedeID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar comprobantes de pago", err))
		return
	}

	// Si no hay comprobantes, devolver array vacío
	if comprobantes == nil {
		comprobantes = []*entidades.ComprobantePago{}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de pago listados exitosamente", comprobantes))
}

// ListByReserva lista todos los comprobantes de pago de una reserva específica
func (c *ComprobantePagoController) ListByReserva(ctx *gin.Context) {
	// Parsear ID de reserva de la URL
	idReserva, err := strconv.Atoi(ctx.Param("idReserva"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de reserva inválido", err))
		return
	}

	// Listar comprobantes por reserva
	comprobantes, err := c.comprobantePagoService.ListByReserva(idReserva)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar comprobantes de pago por reserva", err))
		return
	}

	// Verificar acceso según el rol para los comprobantes obtenidos
	if ctx.GetString("rol") != "ADMIN" && len(comprobantes) > 0 {
		sedeUsuario := ctx.GetInt("sede_id")
		comprobanteFiltrados := []*entidades.ComprobantePago{}

		for _, comprobante := range comprobantes {
			if comprobante.IDSede == sedeUsuario {
				comprobanteFiltrados = append(comprobanteFiltrados, comprobante)
			}
		}
		comprobantes = comprobanteFiltrados
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de pago listados exitosamente", comprobantes))
}

// ListByFecha lista todos los comprobantes de pago de una fecha específica
func (c *ComprobantePagoController) ListByFecha(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Listar comprobantes por fecha
	comprobantes, err := c.comprobantePagoService.ListByFecha(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar comprobantes de pago por fecha", err))
		return
	}

	// Verificar acceso según el rol para los comprobantes obtenidos
	if ctx.GetString("rol") != "ADMIN" {
		sedeUsuario := ctx.GetInt("sede_id")
		comprobanteFiltrados := []*entidades.ComprobantePago{}

		for _, comprobante := range comprobantes {
			if comprobante.IDSede == sedeUsuario {
				comprobanteFiltrados = append(comprobanteFiltrados, comprobante)
			}
		}
		comprobantes = comprobanteFiltrados
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de pago listados exitosamente", comprobantes))
}

// ListByTipo lista todos los comprobantes de pago de un tipo específico
func (c *ComprobantePagoController) ListByTipo(ctx *gin.Context) {
	// Parsear tipo de la URL
	tipo := ctx.Param("tipo")

	// Listar comprobantes por tipo
	comprobantes, err := c.comprobantePagoService.ListByTipo(tipo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar comprobantes de pago por tipo", err))
		return
	}

	// Verificar acceso según el rol para los comprobantes obtenidos
	if ctx.GetString("rol") != "ADMIN" {
		sedeUsuario := ctx.GetInt("sede_id")
		comprobanteFiltrados := []*entidades.ComprobantePago{}

		for _, comprobante := range comprobantes {
			if comprobante.IDSede == sedeUsuario {
				comprobanteFiltrados = append(comprobanteFiltrados, comprobante)
			}
		}
		comprobantes = comprobanteFiltrados
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de pago listados exitosamente", comprobantes))
}

// ListByEstado lista todos los comprobantes de pago con un estado específico
func (c *ComprobantePagoController) ListByEstado(ctx *gin.Context) {
	// Parsear estado de la URL
	estado := ctx.Param("estado")

	// Listar comprobantes por estado
	comprobantes, err := c.comprobantePagoService.ListByEstado(estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar comprobantes de pago por estado", err))
		return
	}

	// Verificar acceso según el rol para los comprobantes obtenidos
	if ctx.GetString("rol") != "ADMIN" {
		sedeUsuario := ctx.GetInt("sede_id")
		comprobanteFiltrados := []*entidades.ComprobantePago{}

		for _, comprobante := range comprobantes {
			if comprobante.IDSede == sedeUsuario {
				comprobanteFiltrados = append(comprobanteFiltrados, comprobante)
			}
		}
		comprobantes = comprobanteFiltrados
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de pago listados exitosamente", comprobantes))
}

// ListByCliente lista todos los comprobantes de un cliente específico
func (c *ComprobantePagoController) ListByCliente(ctx *gin.Context) {
	// Parsear ID de cliente de la URL
	idCliente, err := strconv.Atoi(ctx.Param("idCliente"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de cliente inválido", err))
		return
	}

	// Obtener comprobantes del cliente
	comprobantes, err := c.comprobantePagoService.ListByCliente(idCliente)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar comprobantes por cliente", err))
		return
	}

	// Verificar acceso según el rol para los comprobantes obtenidos
	if ctx.GetString("rol") != "ADMIN" {
		sedeUsuario := ctx.GetInt("sede_id")
		comprobanteFiltrados := []*entidades.ComprobantePago{}

		for _, comprobante := range comprobantes {
			if comprobante.IDSede == sedeUsuario {
				comprobanteFiltrados = append(comprobanteFiltrados, comprobante)
			}
		}
		comprobantes = comprobanteFiltrados
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes listados exitosamente", comprobantes))
}

// ListBySede lista todos los comprobantes de pago de una sede específica
func (c *ComprobantePagoController) ListBySede(ctx *gin.Context) {
	// Parsear ID de sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Verificar permisos
	if ctx.GetString("rol") != "ADMIN" && idSede != ctx.GetInt("sede_id") {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene permiso para ver comprobantes de otra sede", nil))
		return
	}

	// Listar comprobantes por sede
	comprobantes, err := c.comprobantePagoService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar comprobantes por sede", err))
		return
	}

	// Si no hay comprobantes, devolver array vacío
	if comprobantes == nil {
		comprobantes = []*entidades.ComprobantePago{}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Comprobantes de la sede listados exitosamente", comprobantes))
}
