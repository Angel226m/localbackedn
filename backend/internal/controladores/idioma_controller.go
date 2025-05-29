package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IdiomaController maneja los endpoints de idiomas
type IdiomaController struct {
	idiomaService *servicios.IdiomaService
}

// NewIdiomaController crea una nueva instancia de IdiomaController
func NewIdiomaController(idiomaService *servicios.IdiomaService) *IdiomaController {
	return &IdiomaController{
		idiomaService: idiomaService,
	}
}

// Create crea un nuevo idioma
func (c *IdiomaController) Create(ctx *gin.Context) {
	var idioma entidades.Idioma

	// Parsear request
	if err := ctx.ShouldBindJSON(&idioma); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(idioma); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear idioma
	id, err := c.idiomaService.Create(&idioma)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Idioma creado exitosamente", gin.H{"id": id}))
}

// GetByID obtiene un idioma por su ID
func (c *IdiomaController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener idioma
	idioma, err := c.idiomaService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Idioma no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma obtenido", idioma))
}

// List lista todos los idiomas
func (c *IdiomaController) List(ctx *gin.Context) {
	// Listar idiomas
	idiomas, err := c.idiomaService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar idiomas", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idiomas listados exitosamente", idiomas))
}

// Update actualiza un idioma
func (c *IdiomaController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var idioma entidades.Idioma

	// Parsear request
	if err := ctx.ShouldBindJSON(&idioma); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(idioma); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar idioma
	err = c.idiomaService.Update(id, &idioma)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma actualizado exitosamente", nil))
}

// Delete elimina un idioma
func (c *IdiomaController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar idioma
	err = c.idiomaService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al eliminar idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma eliminado exitosamente", nil))
}

// ListDeleted lista todos los idiomas eliminados
func (c *IdiomaController) ListDeleted(ctx *gin.Context) {
	// Listar idiomas eliminados
	idiomas, err := c.idiomaService.ListDeleted()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar idiomas eliminados", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idiomas eliminados listados exitosamente", idiomas))
}

// Restore restaura un idioma eliminado
func (c *IdiomaController) Restore(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Restaurar idioma
	err = c.idiomaService.Restore(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Error al restaurar idioma", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma restaurado exitosamente", nil))
}

// GetByNombre obtiene un idioma por su nombre
func (c *IdiomaController) GetByNombre(ctx *gin.Context) {
	nombre := ctx.Param("nombre")

	// Obtener idioma
	idioma, err := c.idiomaService.GetByNombre(nombre)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Idioma no encontrado", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Idioma obtenido", idioma))
}
