package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HorarioTourController maneja los endpoints de horarios de tour
type HorarioTourController struct {
	horarioTourService *servicios.HorarioTourService
}

// NewHorarioTourController crea una nueva instancia de HorarioTourController
func NewHorarioTourController(horarioTourService *servicios.HorarioTourService) *HorarioTourController {
	return &HorarioTourController{
		horarioTourService: horarioTourService,
	}
}

// Create crea un nuevo horario de tour
func (c *HorarioTourController) Create(ctx *gin.Context) {
	var horarioReq entidades.NuevoHorarioTourRequest

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

	// Crear horario de tour
	id, err := c.horarioTourService.Create(&horarioReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear horario de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Horario de tour creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un horario de tour por su ID
func (c *HorarioTourController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener horario de tour
	horario, err := c.horarioTourService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Horario de tour no encontrado", err))
		return
	}

	// Formatear horas para la respuesta (opcional)
	horarioResponse := map[string]interface{}{
		"id_horario":            horario.ID,
		"id_tipo_tour":          horario.IDTipoTour,
		"id_sede":               horario.IDSede,
		"hora_inicio":           horario.HoraInicio.Format("15:04"),
		"hora_fin":              horario.HoraFin.Format("15:04"),
		"disponible_lunes":      horario.DisponibleLunes,
		"disponible_martes":     horario.DisponibleMartes,
		"disponible_miercoles":  horario.DisponibleMiercoles,
		"disponible_jueves":     horario.DisponibleJueves,
		"disponible_viernes":    horario.DisponibleViernes,
		"disponible_sabado":     horario.DisponibleSabado,
		"disponible_domingo":    horario.DisponibleDomingo,
		"eliminado":             horario.Eliminado,
		"nombre_tipo_tour":      horario.NombreTipoTour,
		"descripcion_tipo_tour": horario.DescripcionTipoTour,
		"nombre_sede":           horario.NombreSede,
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de tour obtenido", horarioResponse))
}

// Update actualiza un horario de tour
func (c *HorarioTourController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var horarioReq entidades.ActualizarHorarioTourRequest

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

	// Actualizar horario de tour
	err = c.horarioTourService.Update(id, &horarioReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar horario de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de tour actualizado exitosamente", nil))
}

// Delete elimina un horario de tour (borrado lógico)
func (c *HorarioTourController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar horario de tour
	err = c.horarioTourService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar horario de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de tour eliminado exitosamente", nil))
}

// List lista todos los horarios de tour
func (c *HorarioTourController) List(ctx *gin.Context) {
	// Listar horarios de tour
	horarios, err := c.horarioTourService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar horarios de tour", err))
		return
	}

	// Formatear horas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario":            horario.ID,
			"id_tipo_tour":          horario.IDTipoTour,
			"id_sede":               horario.IDSede,
			"hora_inicio":           horario.HoraInicio.Format("15:04"),
			"hora_fin":              horario.HoraFin.Format("15:04"),
			"disponible_lunes":      horario.DisponibleLunes,
			"disponible_martes":     horario.DisponibleMartes,
			"disponible_miercoles":  horario.DisponibleMiercoles,
			"disponible_jueves":     horario.DisponibleJueves,
			"disponible_viernes":    horario.DisponibleViernes,
			"disponible_sabado":     horario.DisponibleSabado,
			"disponible_domingo":    horario.DisponibleDomingo,
			"eliminado":             horario.Eliminado,
			"nombre_tipo_tour":      horario.NombreTipoTour,
			"descripcion_tipo_tour": horario.DescripcionTipoTour,
			"nombre_sede":           horario.NombreSede,
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de tour listados exitosamente", horariosResponse))
}

// ListByTipoTour lista todos los horarios asociados a un tipo de tour específico
func (c *HorarioTourController) ListByTipoTour(ctx *gin.Context) {
	// Parsear ID del tipo de tour de la URL
	idTipoTour, err := strconv.Atoi(ctx.Param("idTipoTour"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Listar horarios por tipo de tour
	horarios, err := c.horarioTourService.ListByTipoTour(idTipoTour)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios de tour por tipo", err))
		return
	}

	// Formatear horas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario":            horario.ID,
			"id_tipo_tour":          horario.IDTipoTour,
			"id_sede":               horario.IDSede,
			"hora_inicio":           horario.HoraInicio.Format("15:04"),
			"hora_fin":              horario.HoraFin.Format("15:04"),
			"disponible_lunes":      horario.DisponibleLunes,
			"disponible_martes":     horario.DisponibleMartes,
			"disponible_miercoles":  horario.DisponibleMiercoles,
			"disponible_jueves":     horario.DisponibleJueves,
			"disponible_viernes":    horario.DisponibleViernes,
			"disponible_sabado":     horario.DisponibleSabado,
			"disponible_domingo":    horario.DisponibleDomingo,
			"eliminado":             horario.Eliminado,
			"nombre_tipo_tour":      horario.NombreTipoTour,
			"descripcion_tipo_tour": horario.DescripcionTipoTour,
			"nombre_sede":           horario.NombreSede,
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de tour por tipo listados exitosamente", horariosResponse))
}

// ListByDia lista todos los horarios disponibles para un día específico
func (c *HorarioTourController) ListByDia(ctx *gin.Context) {
	// Parsear día de la semana de la URL
	dia := ctx.Param("dia")

	// Listar horarios por día
	horarios, err := c.horarioTourService.ListByDia(dia)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al listar horarios de tour por día", err))
		return
	}

	// Formatear horas para la respuesta
	horariosResponse := make([]map[string]interface{}, len(horarios))
	for i, horario := range horarios {
		horariosResponse[i] = map[string]interface{}{
			"id_horario":            horario.ID,
			"id_tipo_tour":          horario.IDTipoTour,
			"id_sede":               horario.IDSede,
			"hora_inicio":           horario.HoraInicio.Format("15:04"),
			"hora_fin":              horario.HoraFin.Format("15:04"),
			"disponible_lunes":      horario.DisponibleLunes,
			"disponible_martes":     horario.DisponibleMartes,
			"disponible_miercoles":  horario.DisponibleMiercoles,
			"disponible_jueves":     horario.DisponibleJueves,
			"disponible_viernes":    horario.DisponibleViernes,
			"disponible_sabado":     horario.DisponibleSabado,
			"disponible_domingo":    horario.DisponibleDomingo,
			"eliminado":             horario.Eliminado,
			"nombre_tipo_tour":      horario.NombreTipoTour,
			"descripcion_tipo_tour": horario.DescripcionTipoTour,
			"nombre_sede":           horario.NombreSede,
		}
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de tour por día listados exitosamente", horariosResponse))
}
