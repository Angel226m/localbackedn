package controladores

import (
	"net/http"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GaleriaTourController maneja los endpoints para la galería de imágenes de tours
type GaleriaTourController struct {
	galeriaTourService *servicios.GaleriaTourService
}

// NewGaleriaTourController crea una nueva instancia de GaleriaTourController
func NewGaleriaTourController(galeriaTourService *servicios.GaleriaTourService) *GaleriaTourController {
	return &GaleriaTourController{
		galeriaTourService: galeriaTourService,
	}
}

// Create crea una nueva imagen en la galería
func (c *GaleriaTourController) Create(ctx *gin.Context) {
	var galeriaReq entidades.GaleriaTourRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&galeriaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(galeriaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Crear imagen
	id, err := c.galeriaTourService.CrearImagen(&galeriaReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al crear imagen de galería", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusCreated, utils.SuccessResponse("Imagen de galería creada exitosamente", gin.H{"id": id}))
}

// GetByID obtiene una imagen por su ID
func (c *GaleriaTourController) GetByID(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Obtener imagen
	imagen, err := c.galeriaTourService.ObtenerPorID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Imagen no encontrada", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Imagen obtenida exitosamente", imagen))
}

// Update actualiza una imagen
func (c *GaleriaTourController) Update(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	var galeriaReq entidades.GaleriaTourUpdateRequest

	// Parsear request
	if err := ctx.ShouldBindJSON(&galeriaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
		return
	}

	// Validar datos
	if err := utils.ValidateStruct(galeriaReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error de validación", err))
		return
	}

	// Actualizar imagen
	err = c.galeriaTourService.ActualizarImagen(id, &galeriaReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al actualizar imagen", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Imagen actualizada exitosamente", nil))
}

// Delete elimina una imagen
func (c *GaleriaTourController) Delete(ctx *gin.Context) {
	// Parsear ID de la URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID inválido", err))
		return
	}

	// Eliminar imagen
	err = c.galeriaTourService.EliminarImagen(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Error al eliminar imagen", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Imagen eliminada exitosamente", nil))
}

// ListByTipoTour lista todas las imágenes de un tipo de tour específico
func (c *GaleriaTourController) ListByTipoTour(ctx *gin.Context) {
	// Parsear ID del tipo de tour de la URL
	idTipoTour, err := strconv.Atoi(ctx.Param("id_tipo_tour"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("ID de tipo de tour inválido", err))
		return
	}

	// Listar imágenes por tipo de tour
	imagenes, err := c.galeriaTourService.ListarPorTipoTour(idTipoTour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al listar imágenes", err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Imágenes listadas exitosamente", imagenes))
}
