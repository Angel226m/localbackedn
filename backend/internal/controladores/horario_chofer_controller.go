package controladores

import (
	"net/http"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horario de chofer obtenido", horario))
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

// Delete elimina un horario de chofer
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de chofer listados exitosamente", horarios))
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios del chofer listados exitosamente", horarios))
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios activos del chofer listados exitosamente", horarios))
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios de choferes listados exitosamente", horarios))
}

// GetMyActiveHorarios obtiene los horarios activos del chofer autenticado
func (c *HorarioChoferController) GetMyActiveHorarios(ctx *gin.Context) {
	// Obtener ID del usuario autenticado del contexto
	userID, exists := ctx.Get("user_id")
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

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Horarios activos obtenidos exitosamente", horarios))
}
