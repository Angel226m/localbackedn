package controladores

import (
	"net/http"
	"strconv"
	"time"

	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"

	"github.com/gin-gonic/gin"
)

// TourProgramadoController maneja las rutas relacionadas con tours programados
type TourProgramadoController struct {
	service *servicios.TourProgramadoService
}

// NewTourProgramadoController crea una nueva instancia del controlador
func NewTourProgramadoController(service *servicios.TourProgramadoService) *TourProgramadoController {
	return &TourProgramadoController{
		service: service,
	}
}

// GetByID obtiene un tour programado por su ID
func (c *TourProgramadoController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	tourProgramado, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tour programado no encontrado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado obtenido con éxito", tourProgramado))
}

// Create crea un nuevo tour programado
func (c *TourProgramadoController) Create(ctx *gin.Context) {
	var request entidades.NuevoTourProgramadoRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Validar las fechas de vigencia
	if request.VigenciaDesde == "" || request.VigenciaHasta == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Las fechas de vigencia son obligatorias", nil))
		return
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", request.Fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	vigenciaDesde, err := time.Parse("2006-01-02", request.VigenciaDesde)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia desde inválido, debe ser YYYY-MM-DD", err))
		return
	}

	vigenciaHasta, err := time.Parse("2006-01-02", request.VigenciaHasta)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia hasta inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Validar que la fecha de vigencia desde no sea posterior a la fecha de vigencia hasta
	if vigenciaDesde.After(vigenciaHasta) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de vigencia desde no puede ser posterior a la fecha de vigencia hasta", nil))
		return
	}

	id, err := c.service.Create(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tour programado", err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tour programado creado con éxito", gin.H{"id": id}))
}

// Update actualiza un tour programado
func (c *TourProgramadoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.ActualizarTourProgramadoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Validar la coherencia de las fechas de vigencia si se proporcionan
	if request.VigenciaDesde != "" && request.VigenciaHasta != "" {
		vigenciaDesde, err := time.Parse("2006-01-02", request.VigenciaDesde)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia desde inválido, debe ser YYYY-MM-DD", err))
			return
		}

		vigenciaHasta, err := time.Parse("2006-01-02", request.VigenciaHasta)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia hasta inválido, debe ser YYYY-MM-DD", err))
			return
		}

		if vigenciaDesde.After(vigenciaHasta) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de vigencia desde no puede ser posterior a la fecha de vigencia hasta", nil))
			return
		}

		// Validar formato de fecha si se proporciona
		if request.Fecha != "" {
			_, err := time.Parse("2006-01-02", request.Fecha)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
				return
			}
		}
	} else if (request.VigenciaDesde != "" && request.VigenciaHasta == "") || (request.VigenciaDesde == "" && request.VigenciaHasta != "") {
		// Si solo se proporciona una de las fechas de vigencia
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Debe proporcionar ambas fechas de vigencia o ninguna", nil))
		return
	}

	err = c.service.Update(id, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tour programado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado actualizado con éxito", nil))
}

// Delete elimina lógicamente un tour programado
func (c *TourProgramadoController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar tour programado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado eliminado con éxito", nil))
}

// List obtiene una lista de tours programados con filtros opcionales
func (c *TourProgramadoController) List(ctx *gin.Context) {
	var filtros entidades.FiltrosTourProgramado

	// Extraer parámetros de consulta
	if idSede := ctx.Query("id_sede"); idSede != "" {
		if id, err := strconv.Atoi(idSede); err == nil {
			filtros.IDSede = &id
		}
	}

	if idTipoTour := ctx.Query("id_tipo_tour"); idTipoTour != "" {
		if id, err := strconv.Atoi(idTipoTour); err == nil {
			filtros.IDTipoTour = &id
		}
	}

	if idChofer := ctx.Query("id_chofer"); idChofer != "" {
		if id, err := strconv.Atoi(idChofer); err == nil {
			filtros.IDChofer = &id
		}
	}

	if idEmbarcacion := ctx.Query("id_embarcacion"); idEmbarcacion != "" {
		if id, err := strconv.Atoi(idEmbarcacion); err == nil {
			filtros.IDEmbarcacion = &id
		}
	}

	if estado := ctx.Query("estado"); estado != "" {
		filtros.Estado = &estado
	}

	if fechaInicio := ctx.Query("fecha_inicio"); fechaInicio != "" {
		if _, err := time.Parse("2006-01-02", fechaInicio); err == nil {
			filtros.FechaInicio = &fechaInicio
		}
	}

	if fechaFin := ctx.Query("fecha_fin"); fechaFin != "" {
		if _, err := time.Parse("2006-01-02", fechaFin); err == nil {
			filtros.FechaFin = &fechaFin
		}
	}

	// Añadir filtros de vigencia
	if vigenciaDesdeIni := ctx.Query("vigencia_desde_ini"); vigenciaDesdeIni != "" {
		if _, err := time.Parse("2006-01-02", vigenciaDesdeIni); err == nil {
			filtros.VigenciaDesdeIni = &vigenciaDesdeIni
		}
	}

	if vigenciaDesdefin := ctx.Query("vigencia_desde_fin"); vigenciaDesdefin != "" {
		if _, err := time.Parse("2006-01-02", vigenciaDesdefin); err == nil {
			filtros.VigenciaDesdefin = &vigenciaDesdefin
		}
	}

	if vigenciaHastaIni := ctx.Query("vigencia_hasta_ini"); vigenciaHastaIni != "" {
		if _, err := time.Parse("2006-01-02", vigenciaHastaIni); err == nil {
			filtros.VigenciaHastaIni = &vigenciaHastaIni
		}
	}

	if vigenciaHastaFin := ctx.Query("vigencia_hasta_fin"); vigenciaHastaFin != "" {
		if _, err := time.Parse("2006-01-02", vigenciaHastaFin); err == nil {
			filtros.VigenciaHastaFin = &vigenciaHastaFin
		}
	}

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// AsignarChofer asigna un chofer a un tour programado
func (c *TourProgramadoController) AsignarChofer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.AsignarChoferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	err = c.service.AsignarChofer(id, request.IDChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al asignar chofer", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Chofer asignado con éxito", nil))
}

// CambiarEstado cambia el estado de un tour programado
func (c *TourProgramadoController) CambiarEstado(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request struct {
		Estado string `json:"estado" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		// Intentar obtener el estado como parámetro de consulta si no está en el cuerpo
		estado := ctx.Query("estado")
		if estado == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el parámetro 'estado'", err))
			return
		}
		request.Estado = estado
	}

	err = c.service.CambiarEstado(id, request.Estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado actualizado con éxito", nil))
}

// GetProgramacionSemanal obtiene los tours programados para una semana específica
func (c *TourProgramadoController) GetProgramacionSemanal(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	if fechaInicio == "" {
		// Si no se proporciona fecha, usar la fecha actual
		fechaInicio = time.Now().Format("2006-01-02")
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetProgramacionSemanal(fechaInicio, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener programación semanal", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Programación semanal obtenida con éxito", tours))
}

// GetToursDisponiblesEnFecha obtiene tours disponibles para una fecha específica
func (c *TourProgramadoController) GetToursDisponiblesEnFecha(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		// Si no se proporciona fecha, usar la fecha actual
		fecha = time.Now().Format("2006-01-02")
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetToursDisponiblesEnFecha(fecha, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// GetToursDisponiblesEnRangoFechas obtiene tours disponibles para un rango de fechas
func (c *TourProgramadoController) GetToursDisponiblesEnRangoFechas(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	fechaFin := ctx.Query("fecha_fin")

	if fechaInicio == "" || fechaFin == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'fecha_inicio' y 'fecha_fin' en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido", err))
		return
	}

	_, err = time.Parse("2006-01-02", fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// VerificarDisponibilidadHorario verifica si un horario está disponible para una fecha específica
func (c *TourProgramadoController) VerificarDisponibilidadHorario(ctx *gin.Context) {
	idHorarioStr := ctx.Query("id_horario")
	fecha := ctx.Query("fecha")

	if idHorarioStr == "" || fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'id_horario' y 'fecha'", nil))
		return
	}

	idHorario, err := strconv.Atoi(idHorarioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de horario inválido", err))
		return
	}

	// Verificar primero que la fecha sea válida
	_, err = time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Usar el método del servicio en lugar del método auxiliar
	disponible, err := c.service.VerificarDisponibilidadHorario(idHorario, fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al verificar disponibilidad", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad verificada con éxito", gin.H{
		"disponible": disponible,
	}))
}

// ListByFecha obtiene los tours programados para una fecha específica
func (c *TourProgramadoController) ListByFecha(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere la fecha en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.FechaInicio = &fecha
	filtros.FechaFin = &fecha

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByRangoFechas obtiene los tours programados para un rango de fechas
func (c *TourProgramadoController) ListByRangoFechas(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	fechaFin := ctx.Query("fecha_fin")

	if fechaInicio == "" || fechaFin == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'fecha_inicio' y 'fecha_fin' en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido", err))
		return
	}

	_, err = time.Parse("2006-01-02", fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.FechaInicio = &fechaInicio
	filtros.FechaFin = &fechaFin

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByEstado obtiene los tours programados por estado
func (c *TourProgramadoController) ListByEstado(ctx *gin.Context) {
	estado := ctx.Param("estado")
	if estado == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el parámetro 'estado'", nil))
		return
	}

	// Validar que el estado sea válido
	estadosValidos := []string{"PROGRAMADO", "EN_CURSO", "COMPLETADO", "CANCELADO"}
	estadoValido := false
	for _, e := range estadosValidos {
		if e == estado {
			estadoValido = true
			break
		}
	}

	if !estadoValido {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Estado no válido. Debe ser: PROGRAMADO, EN_CURSO, COMPLETADO o CANCELADO", nil))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.Estado = &estado

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListToursProgramadosDisponibles lista los tours programados disponibles para reserva
func (c *TourProgramadoController) ListToursProgramadosDisponibles(ctx *gin.Context) {
	// Por defecto, buscar desde la fecha actual en adelante
	fechaActual := time.Now().Format("2006-01-02")

	fechaInicio := ctx.Query("fecha_inicio")
	if fechaInicio == "" {
		fechaInicio = fechaActual
	}

	fechaFin := ctx.Query("fecha_fin")
	// Si no se proporciona fecha fin, usar 30 días después de la fecha inicio
	if fechaFin == "" {
		fechaInicioTime, _ := time.Parse("2006-01-02", fechaInicio)
		fechaFin = fechaInicioTime.AddDate(0, 0, 30).Format("2006-01-02")
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	// Usar el método específico para obtener tours disponibles en un rango de fechas
	tours, err := c.service.GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// GetDisponibilidadDia obtiene la disponibilidad de tours para un día específico
func (c *TourProgramadoController) GetDisponibilidadDia(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere la fecha en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	// Usar el método GetToursDisponiblesEnFecha
	tours, err := c.service.GetToursDisponiblesEnFecha(fecha, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener disponibilidad", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad obtenida con éxito", tours))
}

// GetToursVigentes obtiene los tours que están vigentes en la fecha actual
func (c *TourProgramadoController) GetToursVigentes(ctx *gin.Context) {
	// Fecha actual
	fechaActual := time.Now().Format("2006-01-02")

	// Usar filtros para obtener tours cuyo período de vigencia incluya la fecha actual
	var filtros entidades.FiltrosTourProgramado
	filtros.VigenciaDesdeIni = &fechaActual // vigencia_desde <= fechaActual
	filtros.VigenciaHastaFin = &fechaActual // vigencia_hasta >= fechaActual

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours vigentes", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours vigentes obtenidos con éxito", tours))
}

// ListByEmbarcacion obtiene los tours programados por embarcación
func (c *TourProgramadoController) ListByEmbarcacion(ctx *gin.Context) {
	idEmbarcacionStr := ctx.Param("idEmbarcacion")
	if idEmbarcacionStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de embarcación", nil))
		return
	}

	idEmbarcacion, err := strconv.Atoi(idEmbarcacionStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de embarcación inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDEmbarcacion = &idEmbarcacion

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByTipoTour obtiene los tours programados por tipo de tour
func (c *TourProgramadoController) ListByTipoTour(ctx *gin.Context) {
	idTipoTourStr := ctx.Param("idTipoTour")
	if idTipoTourStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de tipo de tour", nil))
		return
	}

	idTipoTour, err := strconv.Atoi(idTipoTourStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDTipoTour = &idTipoTour

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListBySede obtiene los tours programados por sede
func (c *TourProgramadoController) ListBySede(ctx *gin.Context) {
	idSedeStr := ctx.Param("idSede")
	if idSedeStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de sede", nil))
		return
	}

	idSede, err := strconv.Atoi(idSedeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDSede = &idSede

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByChofer obtiene los tours programados por chofer
func (c *TourProgramadoController) ListByChofer(ctx *gin.Context) {
	idChoferStr := ctx.Param("idChofer")
	if idChoferStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de chofer", nil))
		return
	}

	idChofer, err := strconv.Atoi(idChoferStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDChofer = &idChofer

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// En controladores/TourProgramadoController.go

// GetToursDisponibles obtiene tours disponibles para reserva
func (c *TourProgramadoController) GetToursDisponibles(ctx *gin.Context) {
	tours, err := c.service.GetToursDisponibles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

/*package controladores

import (
	"net/http"
	"strconv"
	"time"

	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"

	"github.com/gin-gonic/gin"
)

// TourProgramadoController maneja las rutas relacionadas con tours programados
type TourProgramadoController struct {
	service *servicios.TourProgramadoService
}

// NewTourProgramadoController crea una nueva instancia del controlador
func NewTourProgramadoController(service *servicios.TourProgramadoService) *TourProgramadoController {
	return &TourProgramadoController{
		service: service,
	}
}

// GetByID obtiene un tour programado por su ID
func (c *TourProgramadoController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	tourProgramado, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Tour programado no encontrado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado obtenido con éxito", tourProgramado))
}

// Create crea un nuevo tour programado
func (c *TourProgramadoController) Create(ctx *gin.Context) {
	var request entidades.NuevoTourProgramadoRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Validar las fechas de vigencia
	if request.VigenciaDesde == "" || request.VigenciaHasta == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Las fechas de vigencia son obligatorias", nil))
		return
	}

	// Validar formato y lógica de fechas
	fechaTour, err := time.Parse("2006-01-02", request.Fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	vigenciaDesde, err := time.Parse("2006-01-02", request.VigenciaDesde)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia desde inválido, debe ser YYYY-MM-DD", err))
		return
	}

	vigenciaHasta, err := time.Parse("2006-01-02", request.VigenciaHasta)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia hasta inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Validaciones lógicas de fechas
	if vigenciaDesde.After(vigenciaHasta) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de vigencia desde no puede ser posterior a la fecha de vigencia hasta", nil))
		return
	}

	if fechaTour.Before(vigenciaDesde) || fechaTour.After(vigenciaHasta) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha del tour debe estar dentro del rango de vigencia", nil))
		return
	}

	id, err := c.service.Create(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear tour programado", err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Tour programado creado con éxito", gin.H{"id": id}))
}

// Update actualiza un tour programado
func (c *TourProgramadoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.ActualizarTourProgramadoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Validar la coherencia de las fechas de vigencia si se proporcionan
	if request.VigenciaDesde != "" && request.VigenciaHasta != "" {
		vigenciaDesde, err := time.Parse("2006-01-02", request.VigenciaDesde)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia desde inválido, debe ser YYYY-MM-DD", err))
			return
		}

		vigenciaHasta, err := time.Parse("2006-01-02", request.VigenciaHasta)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha vigencia hasta inválido, debe ser YYYY-MM-DD", err))
			return
		}

		if vigenciaDesde.After(vigenciaHasta) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha de vigencia desde no puede ser posterior a la fecha de vigencia hasta", nil))
			return
		}

		// Si también se proporciona una fecha de tour, validar que esté dentro del rango de vigencia
		if request.Fecha != "" {
			fechaTour, err := time.Parse("2006-01-02", request.Fecha)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
				return
			}

			if fechaTour.Before(vigenciaDesde) || fechaTour.After(vigenciaHasta) {
				ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("La fecha del tour debe estar dentro del rango de vigencia", nil))
				return
			}
		}
	} else if (request.VigenciaDesde != "" && request.VigenciaHasta == "") || (request.VigenciaDesde == "" && request.VigenciaHasta != "") {
		// Si solo se proporciona una de las fechas de vigencia
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Debe proporcionar ambas fechas de vigencia o ninguna", nil))
		return
	}

	err = c.service.Update(id, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar tour programado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado actualizado con éxito", nil))
}

// Delete elimina lógicamente un tour programado
func (c *TourProgramadoController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar tour programado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tour programado eliminado con éxito", nil))
}

// List obtiene una lista de tours programados con filtros opcionales
func (c *TourProgramadoController) List(ctx *gin.Context) {
	var filtros entidades.FiltrosTourProgramado

	// Extraer parámetros de consulta
	if idSede := ctx.Query("id_sede"); idSede != "" {
		if id, err := strconv.Atoi(idSede); err == nil {
			filtros.IDSede = &id
		}
	}

	if idTipoTour := ctx.Query("id_tipo_tour"); idTipoTour != "" {
		if id, err := strconv.Atoi(idTipoTour); err == nil {
			filtros.IDTipoTour = &id
		}
	}

	if idChofer := ctx.Query("id_chofer"); idChofer != "" {
		if id, err := strconv.Atoi(idChofer); err == nil {
			filtros.IDChofer = &id
		}
	}

	if idEmbarcacion := ctx.Query("id_embarcacion"); idEmbarcacion != "" {
		if id, err := strconv.Atoi(idEmbarcacion); err == nil {
			filtros.IDEmbarcacion = &id
		}
	}

	if estado := ctx.Query("estado"); estado != "" {
		filtros.Estado = &estado
	}

	if fechaInicio := ctx.Query("fecha_inicio"); fechaInicio != "" {
		if _, err := time.Parse("2006-01-02", fechaInicio); err == nil {
			filtros.FechaInicio = &fechaInicio
		}
	}

	if fechaFin := ctx.Query("fecha_fin"); fechaFin != "" {
		if _, err := time.Parse("2006-01-02", fechaFin); err == nil {
			filtros.FechaFin = &fechaFin
		}
	}

	// Añadir filtros de vigencia
	if vigenciaDesdeIni := ctx.Query("vigencia_desde_ini"); vigenciaDesdeIni != "" {
		if _, err := time.Parse("2006-01-02", vigenciaDesdeIni); err == nil {
			filtros.VigenciaDesdeIni = &vigenciaDesdeIni
		}
	}

	if vigenciaDesdefin := ctx.Query("vigencia_desde_fin"); vigenciaDesdefin != "" {
		if _, err := time.Parse("2006-01-02", vigenciaDesdefin); err == nil {
			filtros.VigenciaDesdefin = &vigenciaDesdefin
		}
	}

	if vigenciaHastaIni := ctx.Query("vigencia_hasta_ini"); vigenciaHastaIni != "" {
		if _, err := time.Parse("2006-01-02", vigenciaHastaIni); err == nil {
			filtros.VigenciaHastaIni = &vigenciaHastaIni
		}
	}

	if vigenciaHastaFin := ctx.Query("vigencia_hasta_fin"); vigenciaHastaFin != "" {
		if _, err := time.Parse("2006-01-02", vigenciaHastaFin); err == nil {
			filtros.VigenciaHastaFin = &vigenciaHastaFin
		}
	}

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// AsignarChofer asigna un chofer a un tour programado
func (c *TourProgramadoController) AsignarChofer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.AsignarChoferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar campos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	err = c.service.AsignarChofer(id, request.IDChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al asignar chofer", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Chofer asignado con éxito", nil))
}

// CambiarEstado cambia el estado de un tour programado
func (c *TourProgramadoController) CambiarEstado(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	estado := ctx.Query("estado")
	if estado == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el parámetro 'estado'", nil))
		return
	}

	err = c.service.CambiarEstado(id, estado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al cambiar estado", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Estado actualizado con éxito", nil))
}

// GetProgramacionSemanal obtiene los tours programados para una semana específica
func (c *TourProgramadoController) GetProgramacionSemanal(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	if fechaInicio == "" {
		// Si no se proporciona fecha, usar la fecha actual
		fechaInicio = time.Now().Format("2006-01-02")
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetProgramacionSemanal(fechaInicio, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener programación semanal", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Programación semanal obtenida con éxito", tours))
}

// GetToursDisponiblesEnFecha obtiene tours disponibles para una fecha específica
func (c *TourProgramadoController) GetToursDisponiblesEnFecha(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		// Si no se proporciona fecha, usar la fecha actual
		fecha = time.Now().Format("2006-01-02")
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetToursDisponiblesEnFecha(fecha, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// GetToursDisponiblesEnRangoFechas obtiene tours disponibles para un rango de fechas
func (c *TourProgramadoController) GetToursDisponiblesEnRangoFechas(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	fechaFin := ctx.Query("fecha_fin")

	if fechaInicio == "" || fechaFin == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'fecha_inicio' y 'fecha_fin' en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido", err))
		return
	}

	_, err = time.Parse("2006-01-02", fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	tours, err := c.service.GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// VerificarDisponibilidadHorario verifica si un horario está disponible para una fecha específica
func (c *TourProgramadoController) VerificarDisponibilidadHorario(ctx *gin.Context) {
	idHorarioStr := ctx.Query("id_horario")
	fecha := ctx.Query("fecha")

	if idHorarioStr == "" || fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'id_horario' y 'fecha'", nil))
		return
	}

	idHorario, err := strconv.Atoi(idHorarioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de horario inválido", err))
		return
	}

	// Verificar primero que la fecha sea válida
	_, err = time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido, debe ser YYYY-MM-DD", err))
		return
	}

	// Usar el método del servicio en lugar del método auxiliar
	disponible, err := c.service.VerificarDisponibilidadHorario(idHorario, fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al verificar disponibilidad", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad verificada con éxito", gin.H{
		"disponible": disponible,
	}))
}

// Método auxiliar para verificar disponibilidad de horario
// (Esta función debería implementarse en el servicio)
func (c *TourProgramadoController) verificarDisponibilidadHorario(idHorario int, fecha string) (bool, error) {
	// Esta es una implementación temporal. En una implementación real,
	// deberías implementar esta funcionalidad en el servicio.
	return true, nil
}

// ListByFecha obtiene los tours programados para una fecha específica
func (c *TourProgramadoController) ListByFecha(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere la fecha en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.FechaInicio = &fecha
	filtros.FechaFin = &fecha

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByRangoFechas obtiene los tours programados para un rango de fechas
func (c *TourProgramadoController) ListByRangoFechas(ctx *gin.Context) {
	fechaInicio := ctx.Query("fecha_inicio")
	fechaFin := ctx.Query("fecha_fin")

	if fechaInicio == "" || fechaFin == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requieren los parámetros 'fecha_inicio' y 'fecha_fin' en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inicio inválido", err))
		return
	}

	_, err = time.Parse("2006-01-02", fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha fin inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.FechaInicio = &fechaInicio
	filtros.FechaFin = &fechaFin

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByEstado obtiene los tours programados por estado
func (c *TourProgramadoController) ListByEstado(ctx *gin.Context) {
	estado := ctx.Param("estado")
	if estado == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el parámetro 'estado'", nil))
		return
	}

	// Validar que el estado sea válido
	estadosValidos := []string{"PROGRAMADO", "EN_CURSO", "COMPLETADO", "CANCELADO"}
	estadoValido := false
	for _, e := range estadosValidos {
		if e == estado {
			estadoValido = true
			break
		}
	}

	if !estadoValido {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Estado no válido. Debe ser: PROGRAMADO, EN_CURSO, COMPLETADO o CANCELADO", nil))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.Estado = &estado

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListToursProgramadosDisponibles lista los tours programados disponibles para reserva
func (c *TourProgramadoController) ListToursProgramadosDisponibles(ctx *gin.Context) {
	// Por defecto, buscar desde la fecha actual en adelante
	fechaActual := time.Now().Format("2006-01-02")

	fechaInicio := ctx.Query("fecha_inicio")
	if fechaInicio == "" {
		fechaInicio = fechaActual
	}

	fechaFin := ctx.Query("fecha_fin")
	// Si no se proporciona fecha fin, usar 30 días después de la fecha inicio
	if fechaFin == "" {
		fechaInicioTime, _ := time.Parse("2006-01-02", fechaInicio)
		fechaFin = fechaInicioTime.AddDate(0, 0, 30).Format("2006-01-02")
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	// Usar el método específico para obtener tours disponibles en un rango de fechas
	tours, err := c.service.GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours disponibles", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours disponibles obtenidos con éxito", tours))
}

// GetDisponibilidadDia obtiene la disponibilidad de tours para un día específico
func (c *TourProgramadoController) GetDisponibilidadDia(ctx *gin.Context) {
	fecha := ctx.Param("fecha")
	if fecha == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere la fecha en formato YYYY-MM-DD", nil))
		return
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Formato de fecha inválido", err))
		return
	}

	idSedeStr := ctx.Query("id_sede")
	idSede := 0
	if idSedeStr != "" {
		var err error
		idSede, err = strconv.Atoi(idSedeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
			return
		}
	}

	// Usar el método GetToursDisponiblesEnFecha
	tours, err := c.service.GetToursDisponiblesEnFecha(fecha, idSede)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener disponibilidad", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Disponibilidad obtenida con éxito", tours))
}

// GetToursVigentes obtiene los tours que están vigentes en la fecha actual
func (c *TourProgramadoController) GetToursVigentes(ctx *gin.Context) {
	// Fecha actual
	fechaActual := time.Now().Format("2006-01-02")

	// Usar filtros para obtener tours cuyo período de vigencia incluya la fecha actual
	var filtros entidades.FiltrosTourProgramado
	filtros.VigenciaDesdeIni = &fechaActual // vigencia_desde <= fechaActual
	filtros.VigenciaHastaFin = &fechaActual // vigencia_hasta >= fechaActual

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours vigentes", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours vigentes obtenidos con éxito", tours))
}

// ListByEmbarcacion obtiene los tours programados por embarcación
func (c *TourProgramadoController) ListByEmbarcacion(ctx *gin.Context) {
	idEmbarcacionStr := ctx.Param("idEmbarcacion")
	if idEmbarcacionStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de embarcación", nil))
		return
	}

	idEmbarcacion, err := strconv.Atoi(idEmbarcacionStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de embarcación inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDEmbarcacion = &idEmbarcacion

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByTipoTour obtiene los tours programados por tipo de tour
func (c *TourProgramadoController) ListByTipoTour(ctx *gin.Context) {
	idTipoTourStr := ctx.Param("idTipoTour")
	if idTipoTourStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de tipo de tour", nil))
		return
	}

	idTipoTour, err := strconv.Atoi(idTipoTourStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDTipoTour = &idTipoTour

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListBySede obtiene los tours programados por sede
func (c *TourProgramadoController) ListBySede(ctx *gin.Context) {
	idSedeStr := ctx.Param("idSede")
	if idSedeStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de sede", nil))
		return
	}

	idSede, err := strconv.Atoi(idSedeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de sede inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDSede = &idSede

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}

// ListByChofer obtiene los tours programados por chofer
func (c *TourProgramadoController) ListByChofer(ctx *gin.Context) {
	idChoferStr := ctx.Param("idChofer")
	if idChoferStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Se requiere el ID de chofer", nil))
		return
	}

	idChofer, err := strconv.Atoi(idChoferStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Usar el mismo filtro que en List
	var filtros entidades.FiltrosTourProgramado
	filtros.IDChofer = &idChofer

	tours, err := c.service.List(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener tours programados", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Tours programados obtenidos con éxito", tours))
}
*/
