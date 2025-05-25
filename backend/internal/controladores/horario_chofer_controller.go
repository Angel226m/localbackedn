package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HorarioChoferController maneja los endpoints de horarios de chofer
type HorarioChoferController struct {
	horarioChoferService *servicios.HorarioChoferService
}

// NewHorarioChoferController crea una nueva instancia de HorarioChoferController
func NewHorarioChoferController(horarioChoferService *servicios.HorarioChoferService) *HorarioChoferController {
	return &HorarioChoferController{
		horarioChoferService: horarioChoferService,
	}
}

// Create crea un nuevo horario de chofer
func (c *HorarioChoferController) Create(ctx *gin.Context) {
	var horarioReq entidades.NuevoHorarioChoferRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&horarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(horarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear horario de chofer
	id, err := c.horarioChoferService.Create(&horarioReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear horario de chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Horario de chofer creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un horario de chofer por su ID
func (c *HorarioChoferController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener horario de chofer
	horario, err := c.horarioChoferService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Horario de chofer no encontrado", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horarioResponse := map[string]interface{}{
		"id_horario_chofer":    horario.ID,
		"id_usuario":           horario.IDUsuario,
		"id_sede":              horario.IDSede,
		"hora_inicio":          horario.HoraInicio.Format("15:04"),
		"hora_fin":             horario.HoraFin.Format("15:04"),
		"disponible_lunes":     horario.DisponibleLunes,
		"disponible_martes":    horario.DisponibleMartes,
		"disponible_miercoles": horario.DisponibleMiercoles,
		"disponible_jueves":    horario.DisponibleJueves,
		"disponible_viernes":   horario.DisponibleViernes,
		"disponible_sabado":    horario.DisponibleSabado,
		"disponible_domingo":   horario.DisponibleDomingo,
		"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
		"eliminado":            horario.Eliminado,
		"nombre_chofer":        horario.NombreChofer,
		"apellidos_chofer":     horario.ApellidosChofer,
		"documento_chofer":     horario.DocumentoChofer,
		"telefono_chofer":      horario.TelefonoChofer,
		"nombre_sede":          horario.NombreSede,
	}

	// Formatear fecha_fin si existe
	if horario.FechaFin != nil {
		horarioResponse["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
	} else {
		horarioResponse["fecha_fin"] = nil
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de chofer obtenido", horarioResponse))
}

// Update actualiza un horario de chofer
func (c *HorarioChoferController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var horarioReq entidades.ActualizarHorarioChoferRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&horarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(horarioReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar horario de chofer
	err = c.horarioChoferService.Update(id, &horarioReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar horario de chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de chofer actualizado exitosamente", nil))
}

// Delete elimina un horario de chofer (borrado lógico)
func (c *HorarioChoferController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar horario de chofer
	err = c.horarioChoferService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar horario de chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de chofer eliminado exitosamente", nil))
}

// List lista todos los horarios de chofer
func (c *HorarioChoferController) List(ctx *gin.Context) {
	// Listar horarios de chofer
	horarios, err := c.horarioChoferService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar horarios de chofer", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario_chofer":    horario.ID,
			"id_usuario":           horario.IDUsuario,
			"id_sede":              horario.IDSede,
			"hora_inicio":          horario.HoraInicio.Format("15:04"),
			"hora_fin":             horario.HoraFin.Format("15:04"),
			"disponible_lunes":     horario.DisponibleLunes,
			"disponible_martes":    horario.DisponibleMartes,
			"disponible_miercoles": horario.DisponibleMiercoles,
			"disponible_jueves":    horario.DisponibleJueves,
			"disponible_viernes":   horario.DisponibleViernes,
			"disponible_sabado":    horario.DisponibleSabado,
			"disponible_domingo":   horario.DisponibleDomingo,
			"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
			"eliminado":            horario.Eliminado,
			"nombre_chofer":        horario.NombreChofer,
			"apellidos_chofer":     horario.ApellidosChofer,
			"documento_chofer":     horario.DocumentoChofer,
			"telefono_chofer":      horario.TelefonoChofer,
			"nombre_sede":          horario.NombreSede,
		}

		// Formatear fecha_fin si existe
		if horario.FechaFin != nil {
			horariosResponse[i]["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
		} else {
			horariosResponse[i]["fecha_fin"] = nil
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de chofer listados exitosamente", horariosResponse))
}

// ListByChofer lista todos los horarios de un chofer específico
func (c *HorarioChoferController) ListByChofer(ctx *gin.Context) {
	// Parsear ID del chofer de la URL
	idChofer, err := strconv.Atoi(ctx.Param("idChofer"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Listar horarios del chofer
	horarios, err := c.horarioChoferService.ListByChofer(idChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios del chofer", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario_chofer":    horario.ID,
			"id_usuario":           horario.IDUsuario,
			"id_sede":              horario.IDSede,
			"hora_inicio":          horario.HoraInicio.Format("15:04"),
			"hora_fin":             horario.HoraFin.Format("15:04"),
			"disponible_lunes":     horario.DisponibleLunes,
			"disponible_martes":    horario.DisponibleMartes,
			"disponible_miercoles": horario.DisponibleMiercoles,
			"disponible_jueves":    horario.DisponibleJueves,
			"disponible_viernes":   horario.DisponibleViernes,
			"disponible_sabado":    horario.DisponibleSabado,
			"disponible_domingo":   horario.DisponibleDomingo,
			"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
			"eliminado":            horario.Eliminado,
			"nombre_chofer":        horario.NombreChofer,
			"apellidos_chofer":     horario.ApellidosChofer,
			"documento_chofer":     horario.DocumentoChofer,
			"telefono_chofer":      horario.TelefonoChofer,
			"nombre_sede":          horario.NombreSede,
		}

		// Formatear fecha_fin si existe
		if horario.FechaFin != nil {
			horariosResponse[i]["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
		} else {
			horariosResponse[i]["fecha_fin"] = nil
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios del chofer listados exitosamente", horariosResponse))
}

// ListActiveByChofer lista los horarios activos de un chofer
func (c *HorarioChoferController) ListActiveByChofer(ctx *gin.Context) {
	// Parsear ID del chofer de la URL
	idChofer, err := strconv.Atoi(ctx.Param("idChofer"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de chofer inválido", err))
		return
	}

	// Listar horarios activos del chofer
	horarios, err := c.horarioChoferService.ListActiveByChofer(idChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios activos del chofer", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario_chofer":    horario.ID,
			"id_usuario":           horario.IDUsuario,
			"id_sede":              horario.IDSede,
			"hora_inicio":          horario.HoraInicio.Format("15:04"),
			"hora_fin":             horario.HoraFin.Format("15:04"),
			"disponible_lunes":     horario.DisponibleLunes,
			"disponible_martes":    horario.DisponibleMartes,
			"disponible_miercoles": horario.DisponibleMiercoles,
			"disponible_jueves":    horario.DisponibleJueves,
			"disponible_viernes":   horario.DisponibleViernes,
			"disponible_sabado":    horario.DisponibleSabado,
			"disponible_domingo":   horario.DisponibleDomingo,
			"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
			"eliminado":            horario.Eliminado,
			"nombre_chofer":        horario.NombreChofer,
			"apellidos_chofer":     horario.ApellidosChofer,
			"documento_chofer":     horario.DocumentoChofer,
			"telefono_chofer":      horario.TelefonoChofer,
			"nombre_sede":          horario.NombreSede,
		}

		// Formatear fecha_fin si existe
		if horario.FechaFin != nil {
			horariosResponse[i]["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
		} else {
			horariosResponse[i]["fecha_fin"] = nil
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios activos del chofer listados exitosamente", horariosResponse))
}

// ListByDia lista todos los horarios de choferes disponibles para un día específico
func (c *HorarioChoferController) ListByDia(ctx *gin.Context) {
	// Parsear día de la semana de la URL (1=Lunes, 7=Domingo)
	diaSemana, err := strconv.Atoi(ctx.Param("dia"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Día de la semana inválido", err))
		return
	}

	// Listar horarios de choferes por día
	horarios, err := c.horarioChoferService.ListByDia(diaSemana)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios de choferes por día", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario_chofer":    horario.ID,
			"id_usuario":           horario.IDUsuario,
			"id_sede":              horario.IDSede,
			"hora_inicio":          horario.HoraInicio.Format("15:04"),
			"hora_fin":             horario.HoraFin.Format("15:04"),
			"disponible_lunes":     horario.DisponibleLunes,
			"disponible_martes":    horario.DisponibleMartes,
			"disponible_miercoles": horario.DisponibleMiercoles,
			"disponible_jueves":    horario.DisponibleJueves,
			"disponible_viernes":   horario.DisponibleViernes,
			"disponible_sabado":    horario.DisponibleSabado,
			"disponible_domingo":   horario.DisponibleDomingo,
			"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
			"eliminado":            horario.Eliminado,
			"nombre_chofer":        horario.NombreChofer,
			"apellidos_chofer":     horario.ApellidosChofer,
			"documento_chofer":     horario.DocumentoChofer,
			"telefono_chofer":      horario.TelefonoChofer,
			"nombre_sede":          horario.NombreSede,
		}

		// Formatear fecha_fin si existe
		if horario.FechaFin != nil {
			horariosResponse[i]["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
		} else {
			horariosResponse[i]["fecha_fin"] = nil
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de choferes por día listados exitosamente", horariosResponse))
}

// GetMyActiveHorarios obtiene los horarios activos del chofer autenticado
func (c *HorarioChoferController) GetMyActiveHorarios(ctx *gin.Context) {
	// Obtener ID del usuario autenticado del contexto
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Usuario no autenticado", nil))
		return
	}

	// Listar horarios activos del chofer
	horarios, err := c.horarioChoferService.ListActiveByChofer(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios activos", err))
		return
	}

	// Formatear horas y fechas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario_chofer":    horario.ID,
			"id_usuario":           horario.IDUsuario,
			"id_sede":              horario.IDSede,
			"hora_inicio":          horario.HoraInicio.Format("15:04"),
			"hora_fin":             horario.HoraFin.Format("15:04"),
			"disponible_lunes":     horario.DisponibleLunes,
			"disponible_martes":    horario.DisponibleMartes,
			"disponible_miercoles": horario.DisponibleMiercoles,
			"disponible_jueves":    horario.DisponibleJueves,
			"disponible_viernes":   horario.DisponibleViernes,
			"disponible_sabado":    horario.DisponibleSabado,
			"disponible_domingo":   horario.DisponibleDomingo,
			"fecha_inicio":         horario.FechaInicio.Format("2006-01-02"),
			"eliminado":            horario.Eliminado,
			"nombre_chofer":        horario.NombreChofer,
			"apellidos_chofer":     horario.ApellidosChofer,
			"documento_chofer":     horario.DocumentoChofer,
			"telefono_chofer":      horario.TelefonoChofer,
			"nombre_sede":          horario.NombreSede,
		}

		// Formatear fecha_fin si existe
		if horario.FechaFin != nil {
			horariosResponse[i]["fecha_fin"] = horario.FechaFin.Format("2006-01-02")
		} else {
			horariosResponse[i]["fecha_fin"] = nil
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Mis horarios activos obtenidos exitosamente", horariosResponse))
}
