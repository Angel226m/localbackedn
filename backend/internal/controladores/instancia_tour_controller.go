package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InstanciaTourController maneja los endpoints de instancias de tour
type InstanciaTourController struct {
	instanciaTourService *servicios.InstanciaTourService
}

// NewInstanciaTourController crea una nueva instancia de InstanciaTourController
func NewInstanciaTourController(instanciaTourService *servicios.InstanciaTourService) *InstanciaTourController {
	return &InstanciaTourController{
		instanciaTourService: instanciaTourService,
	}
}

// Create crea una nueva instancia de tour
func (c *InstanciaTourController) Create(ctx *gin.Context) {
	var request entidades.NuevaInstanciaTourRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear instancia de tour
	id, err := c.instanciaTourService.Create(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear instancia de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Instancia de tour creada exitosamente", gin.H{"id": id}))
}

// GetByID obtiene una instancia de tour por su ID
func (c *InstanciaTourController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener instancia de tour
	instancia, err := c.instanciaTourService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Instancia de tour no encontrada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancia de tour obtenida", instancia))
}

// List lista todas las instancias de tour
func (c *InstanciaTourController) List(ctx *gin.Context) {
	// Listar instancias de tour
	instancias, err := c.instanciaTourService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar instancias de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias de tour listadas exitosamente", instancias))
}

// Update actualiza una instancia de tour
func (c *InstanciaTourController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.ActualizarInstanciaTourRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar instancia de tour
	err = c.instanciaTourService.Update(id, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar instancia de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancia de tour actualizada exitosamente", nil))
}

// Delete elimina una instancia de tour
func (c *InstanciaTourController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar instancia de tour
	err = c.instanciaTourService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar instancia de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancia de tour eliminada exitosamente", nil))
}

// ListByTourProgramado lista todas las instancias de un tour programado específico
func (c *InstanciaTourController) ListByTourProgramado(ctx *gin.Context) {
	// Parsear ID de la URL
	idTourProgramado, err := strconv.Atoi(ctx.Param("id_tour_programado"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour programado inválido", err))
		return
	}

	// Listar instancias de tour
	instancias, err := c.instanciaTourService.ListByTourProgramado(idTourProgramado)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar instancias de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias de tour listadas exitosamente", instancias))
}

// ListByFiltros lista instancias de tour según filtros específicos
/*
func (c *InstanciaTourController) ListByFiltros(ctx *gin.Context) {
	var filtros entidades.FiltrosInstanciaTour

	// Parsear request
	if err := ctx.ShouldBindJSON(&filtros); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Listar instancias de tour por filtros
	instancias, err := c.instanciaTourService.ListByFiltros(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al filtrar instancias de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias de tour filtradas exitosamente", instancias))
}
*/

// ListByFiltros lista instancias de tour según filtros específicos
func (c *InstanciaTourController) ListByFiltros(ctx *gin.Context) {
	// Obtener filtros del contexto
	filtrosInterface, exists := ctx.Get("filtros")
	if !exists {
		// Si no hay filtros en el contexto, intentar leerlos del JSON del cuerpo
		var filtros entidades.FiltrosInstanciaTour
		if err := ctx.ShouldBindJSON(&filtros); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
			return
		}

		// Listar instancias de tour por filtros
		instancias, err := c.instanciaTourService.ListByFiltros(filtros)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al filtrar instancias de tour", err))
			return
		}

		// Respuesta exitosa
		ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias de tour filtradas exitosamente", instancias))
		return
	}

	// Convertir interface{} a FiltrosInstanciaTour
	filtros, ok := filtrosInterface.(entidades.FiltrosInstanciaTour)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tipo de filtros inválido", nil))
		return
	}

	// Listar instancias de tour por filtros
	instancias, err := c.instanciaTourService.ListByFiltros(filtros)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al filtrar instancias de tour", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias de tour filtradas exitosamente", instancias))
}

// AsignarChofer asigna un chofer a una instancia de tour
// AsignarChofer asigna un chofer a una instancia de tour
func (c *InstanciaTourController) AsignarChofer(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var request entidades.AsignarChoferInstanciaRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Asignar chofer
	err = c.instanciaTourService.AsignarChofer(id, request.IDChofer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al asignar chofer", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Chofer asignado exitosamente", nil))
}

// GenerarInstanciasDeTourProgramado genera instancias para un tour programado
func (c *InstanciaTourController) GenerarInstanciasDeTourProgramado(ctx *gin.Context) {
	// Parsear ID de la URL
	idTourProgramado, err := strconv.Atoi(ctx.Param("id_tour_programado"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tour programado inválido", err))
		return
	}

	// Generar instancias
	cantidad, err := c.instanciaTourService.GenerarInstanciasDeTourProgramado(idTourProgramado)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al generar instancias", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Instancias generadas exitosamente", gin.H{"cantidad": cantidad}))
}
