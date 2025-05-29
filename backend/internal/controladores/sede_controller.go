package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SedeController maneja los endpoints de sedes
type SedeController struct {
	sedeService *servicios.SedeService
}

// NewSedeController crea una nueva instancia de SedeController
func NewSedeController(sedeService *servicios.SedeService) *SedeController {
	return &SedeController{
		sedeService: sedeService,
	}
}

// Create crea una nueva sede
func (c *SedeController) Create(ctx *gin.Context) {
	var sedeReq entidades.NuevaSedeRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&sedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(sedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear sede
	id, err := c.sedeService.Create(&sedeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Sede creada exitosamente", gin.H{"id": id}))
}

// GetByID obtiene una sede por su ID
func (c *SedeController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener sede
	sede, err := c.sedeService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Sede no encontrada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede obtenida", sede))
}

// Update actualiza una sede
func (c *SedeController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var sedeReq entidades.ActualizarSedeRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&sedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(sedeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar sede
	err = c.sedeService.Update(id, &sedeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede actualizada exitosamente", nil))
}

// Delete elimina una sede (borrado lógico)
func (c *SedeController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar sede
	err = c.sedeService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede eliminada exitosamente", nil))
}

// Restore restaura una sede eliminada lógicamente
func (c *SedeController) Restore(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Restaurar sede
	err = c.sedeService.Restore(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al restaurar sede", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sede restaurada exitosamente", nil))
}

// List lista todas las sedes
func (c *SedeController) List(ctx *gin.Context) {
	// Listar sedes
	sedes, err := c.sedeService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar sedes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sedes listadas exitosamente", sedes))
}

// GetByDistrito obtiene sedes por distrito
func (c *SedeController) GetByDistrito(ctx *gin.Context) {
	// Obtener distrito de la URL
	distrito := ctx.Param("distrito")

	// Listar sedes por distrito
	sedes, err := c.sedeService.GetByDistrito(distrito)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener sedes por distrito", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sedes obtenidas por distrito", sedes))
}

// GetByPais obtiene sedes por país
func (c *SedeController) GetByPais(ctx *gin.Context) {
	// Obtener país de la URL
	pais := ctx.Param("pais")

	// Listar sedes por país
	sedes, err := c.sedeService.GetByPais(pais)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener sedes por país", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Sedes obtenidas por país", sedes))
}
