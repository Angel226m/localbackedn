/*package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TourProgramadoController maneja los endpoints de tours programados
type TourProgramadoController struct {
	tourProgramadoService *servicios.TourProgramadoService
}

// NewTourProgramadoController crea una nueva instancia de TourProgramadoController
func NewTourProgramadoController(tourProgramadoService *servicios.TourProgramadoService) *TourProgramadoController {
	return &TourProgramadoController{
		tourProgramadoService: tourProgramadoService,
	}
}

// Create crea un nuevo tour programado
func (c *TourProgramadoController) Create(ctx *gin.Context) {
	var tourReq entidades.NuevoTourProgramadoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear tour programado
	id, err := c.tourProgramadoService.Create(&tourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tour programado creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un tour programado por su ID
func (c *TourProgramadoController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener tour programado
	tour, err := c.tourProgramadoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tour programado no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado obtenido", tour))
}

// Update actualiza un tour programado
func (c *TourProgramadoController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var tourReq entidades.ActualizarTourProgramadoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar tour programado
	err = c.tourProgramadoService.Update(id, &tourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado actualizado exitosamente", nil))
}

// CambiarEstado cambia el estado de un tour programado
func (c *TourProgramadoController) CambiarEstado(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var estadoReq entidades.CambiarEstadoTourRequest

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
	err = c.tourProgramadoService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado del tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado del tour programado actualizado exitosamente", nil))
}

// Delete elimina un tour programado
func (c *TourProgramadoController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar tour programado
	err = c.tourProgramadoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado eliminado exitosamente", nil))
}

// List lista todos los tours programados
func (c *TourProgramadoController) List(ctx *gin.Context) {
	// Listar tours programados
	tours, err := c.tourProgramadoService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByFecha lista todos los tours programados para una fecha específica
func (c *TourProgramadoController) ListByFecha(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Listar tours programados por fecha
	tours, err := c.tourProgramadoService.ListByFecha(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados por fecha", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByRangoFechas lista todos los tours programados para un rango de fechas
func (c *TourProgramadoController) ListByRangoFechas(ctx *gin.Context) {
	// Parsear fechas de los query params (formato: YYYY-MM-DD)
	fechaInicioStr := ctx.Query("fechaInicio")
	fechaFinStr := ctx.Query("fechaFin")

	fechaInicio, err := time.Parse("2006-01-02", fechaInicioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido, debe ser YYYY-MM-DD", err))
		return
	}

	fechaFin, err := time.Parse("2006-01-02", fechaFinStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Verificar que fechaInicio sea anterior o igual a fechaFin
	if fechaInicio.After(fechaFin) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de inicio debe ser anterior o igual a la fecha de fin", nil))
		return
	}

	// Listar tours programados por rango de fechas
	tours, err := c.tourProgramadoService.ListByRangoFechas(fechaInicio, fechaFin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados por rango de fechas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByEstado lista todos los tours programados por estado
func (c *TourProgramadoController) ListByEstado(ctx *gin.Context) {
	// Parsear estado de la URL
	estado := ctx.Param("estado")

	// Listar tours programados por estado
	tours, err := c.tourProgramadoService.ListByEstado(estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por estado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByEmbarcacion lista todos los tours programados por embarcación
func (c *TourProgramadoController) ListByEmbarcacion(ctx *gin.Context) {
	// Parsear ID de embarcación de la URL
	idEmbarcacion, err := strconv.Atoi(ctx.Param("idEmbarcacion"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de embarcación inválido", err))
		return
	}

	// Listar tours programados por embarcación
	tours, err := c.tourProgramadoService.ListByEmbarcacion(idEmbarcacion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por embarcación", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByChofer
// ListByChofer lista todos los tours programados asociados a un chofer
func (c *TourProgramadoController) ListByChofer(ctx *gin.Context) {
	// Parsear ID de chofer de la URL
	idChofer, err := strconv.Atoi(ctx.Param("idChofer"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Listar tours programados por chofer
	tours, err := c.tourProgramadoService.ListByChofer(idChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListToursProgramadosDisponibles lista todos los tours programados disponibles para reservación
func (c *TourProgramadoController) ListToursProgramadosDisponibles(ctx *gin.Context) {
	// Listar tours programados disponibles
	tours, err := c.tourProgramadoService.ListToursProgramadosDisponibles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados disponibles", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados disponibles listados exitosamente", tours))
}

// ListByTipoTour lista todos los tours programados por tipo de tour
func (c *TourProgramadoController) ListByTipoTour(ctx *gin.Context) {
	// Parsear ID de tipo de tour de la URL
	idTipoTour, err := strconv.Atoi(ctx.Param("idTipoTour"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Listar tours programados por tipo de tour
	tours, err := c.tourProgramadoService.ListByTipoTour(idTipoTour)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por tipo de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// GetDisponibilidadDia retorna la disponibilidad de tours para una fecha específica
func (c *TourProgramadoController) GetDisponibilidadDia(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Obtener disponibilidad para el día
	tours, err := c.tourProgramadoService.GetDisponibilidadDia(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener disponibilidad para el día", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad de tours obtenida exitosamente", tours))
}

// ListBySede lista todos los tours programados de una sede específica
func (c *TourProgramadoController) ListBySede(ctx *gin.Context) {
	// Parsear ID de sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar tours programados por sede
	tours, err := c.tourProgramadoService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ReservarCupo permite reservar cupo en un tour programado
func (c *TourProgramadoController) ReservarCupo(ctx *gin.Context) {
	// Parsear ID del tour programado de la URL
	idTour, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour inválido", err))
		return
	}

	// Parsear la cantidad de cupos a reservar del body
	var request struct {
		Cantidad int `json:"cantidad" validate:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Reservar cupo
	err = c.tourProgramadoService.ReservarCupo(idTour, request.Cantidad)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al reservar cupo", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cupo reservado exitosamente", nil))
}

// LiberarCupo permite liberar cupo en un tour programado
func (c *TourProgramadoController) LiberarCupo(ctx *gin.Context) {
	// Parsear ID del tour programado de la URL
	idTour, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour inválido", err))
		return
	}

	// Parsear la cantidad de cupos a liberar del body
	var request struct {
		Cantidad int `json:"cantidad" validate:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Liberar cupo
	err = c.tourProgramadoService.LiberarCupo(idTour, request.Cantidad)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al liberar cupo", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cupo liberado exitosamente", nil))
}
*/

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

// TourProgramadoController maneja los endpoints de tours programados
type TourProgramadoController struct {
	tourProgramadoService *servicios.TourProgramadoService
}

// NewTourProgramadoController crea una nueva instancia de TourProgramadoController
func NewTourProgramadoController(tourProgramadoService *servicios.TourProgramadoService) *TourProgramadoController {
	return &TourProgramadoController{
		tourProgramadoService: tourProgramadoService,
	}
}

// Create crea un nuevo tour programado
func (c *TourProgramadoController) Create(ctx *gin.Context) {
	var tourReq entidades.NuevoTourProgramadoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear tour programado
	id, err := c.tourProgramadoService.Create(&tourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tour programado creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un tour programado por su ID
func (c *TourProgramadoController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener tour programado
	tour, err := c.tourProgramadoService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tour programado no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado obtenido", tour))
}

// Update actualiza un tour programado
func (c *TourProgramadoController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var tourReq entidades.ActualizarTourProgramadoRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(tourReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar tour programado
	err = c.tourProgramadoService.Update(id, &tourReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado actualizado exitosamente", nil))
}

// CambiarEstado cambia el estado de un tour programado
func (c *TourProgramadoController) CambiarEstado(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var estadoReq entidades.CambiarEstadoTourRequest

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
	err = c.tourProgramadoService.CambiarEstado(id, estadoReq.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado del tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado del tour programado actualizado exitosamente", nil))
}

// Delete elimina un tour programado
func (c *TourProgramadoController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar tour programado
	err = c.tourProgramadoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar tour programado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado eliminado exitosamente", nil))
}

// List lista todos los tours programados
func (c *TourProgramadoController) List(ctx *gin.Context) {
	// Listar tours programados
	tours, err := c.tourProgramadoService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByFecha lista todos los tours programados para una fecha específica
func (c *TourProgramadoController) ListByFecha(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Listar tours programados por fecha
	tours, err := c.tourProgramadoService.ListByFecha(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados por fecha", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByRangoFechas lista todos los tours programados para un rango de fechas
func (c *TourProgramadoController) ListByRangoFechas(ctx *gin.Context) {
	// Parsear fechas de los query params (formato: YYYY-MM-DD)
	fechaInicioStr := ctx.Query("fecha_inicio")
	fechaFinStr := ctx.Query("fecha_fin")

	if fechaInicioStr == "" || fechaFinStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros fecha_inicio y fecha_fin", nil))
		return
	}

	fechaInicio, err := time.Parse("2006-01-02", fechaInicioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido, debe ser YYYY-MM-DD", err))
		return
	}

	fechaFin, err := time.Parse("2006-01-02", fechaFinStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Verificar que fechaInicio sea anterior o igual a fechaFin
	if fechaInicio.After(fechaFin) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de inicio debe ser anterior o igual a la fecha de fin", nil))
		return
	}

	// Listar tours programados por rango de fechas
	tours, err := c.tourProgramadoService.ListByRangoFechas(fechaInicio, fechaFin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados por rango de fechas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByEstado lista todos los tours programados por estado
func (c *TourProgramadoController) ListByEstado(ctx *gin.Context) {
	// Parsear estado de la URL
	estado := ctx.Param("estado")

	// Listar tours programados por estado
	tours, err := c.tourProgramadoService.ListByEstado(estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por estado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByEmbarcacion lista todos los tours programados por embarcación
func (c *TourProgramadoController) ListByEmbarcacion(ctx *gin.Context) {
	// Parsear ID de embarcación de la URL
	idEmbarcacion, err := strconv.Atoi(ctx.Param("idEmbarcacion"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de embarcación inválido", err))
		return
	}

	// Listar tours programados por embarcación
	tours, err := c.tourProgramadoService.ListByEmbarcacion(idEmbarcacion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por embarcación", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListByChofer lista todos los tours programados asociados a un chofer
func (c *TourProgramadoController) ListByChofer(ctx *gin.Context) {
	// Parsear ID de chofer de la URL
	idChofer, err := strconv.Atoi(ctx.Param("idChofer"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Listar tours programados por chofer
	tours, err := c.tourProgramadoService.ListByChofer(idChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ListToursProgramadosDisponibles lista todos los tours programados disponibles para reservación
func (c *TourProgramadoController) ListToursProgramadosDisponibles(ctx *gin.Context) {
	// Listar tours programados disponibles
	tours, err := c.tourProgramadoService.ListToursProgramadosDisponibles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar tours programados disponibles", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados disponibles listados exitosamente", tours))
}

// ListByTipoTour lista todos los tours programados por tipo de tour
func (c *TourProgramadoController) ListByTipoTour(ctx *gin.Context) {
	// Parsear ID de tipo de tour de la URL
	idTipoTour, err := strconv.Atoi(ctx.Param("idTipoTour"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Listar tours programados por tipo de tour
	tours, err := c.tourProgramadoService.ListByTipoTour(idTipoTour)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por tipo de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// GetDisponibilidadDia retorna la disponibilidad de tours para una fecha específica
func (c *TourProgramadoController) GetDisponibilidadDia(ctx *gin.Context) {
	// Parsear fecha de la URL (formato: YYYY-MM-DD)
	fechaStr := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Obtener disponibilidad para el día
	tours, err := c.tourProgramadoService.GetDisponibilidadDia(fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener disponibilidad para el día", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad de tours obtenida exitosamente", tours))
}

// ListBySede lista todos los tours programados de una sede específica
func (c *TourProgramadoController) ListBySede(ctx *gin.Context) {
	// Parsear ID de sede de la URL
	idSede, err := strconv.Atoi(ctx.Param("idSede"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Listar tours programados por sede
	tours, err := c.tourProgramadoService.ListBySede(idSede)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar tours programados por sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados listados exitosamente", tours))
}

// ReservarCupo permite reservar cupo en un tour programado
func (c *TourProgramadoController) ReservarCupo(ctx *gin.Context) {
	// Parsear ID del tour programado de la URL
	idTour, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour inválido", err))
		return
	}

	// Parsear la cantidad de cupos a reservar del body
	var request struct {
		Cantidad int `json:"cantidad" validate:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Reservar cupo
	err = c.tourProgramadoService.ReservarCupo(idTour, request.Cantidad)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al reservar cupo", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cupo reservado exitosamente", nil))
}

// LiberarCupo permite liberar cupo en un tour programado
func (c *TourProgramadoController) LiberarCupo(ctx *gin.Context) {
	// Parsear ID del tour programado de la URL
	idTour, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour inválido", err))
		return
	}

	// Parsear la cantidad de cupos a liberar del body
	var request struct {
		Cantidad int `json:"cantidad" validate:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Liberar cupo
	err = c.tourProgramadoService.LiberarCupo(idTour, request.Cantidad)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al liberar cupo", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Cupo liberado exitosamente", nil))
}
