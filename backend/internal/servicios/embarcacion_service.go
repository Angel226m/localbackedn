package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// EmbarcacionService maneja la lógica de negocio para embarcaciones
type EmbarcacionService struct {
	embarcacionRepo *repositorios.EmbarcacionRepository
	sedeRepo        *repositorios.SedeRepository
}

// NewEmbarcacionService crea una nueva instancia de EmbarcacionService
func NewEmbarcacionService(
	embarcacionRepo *repositorios.EmbarcacionRepository,
	sedeRepo *repositorios.SedeRepository,
) *EmbarcacionService {
	return &EmbarcacionService{
		embarcacionRepo: embarcacionRepo,
		sedeRepo:        sedeRepo,
	}
}

// Create crea una nueva embarcación
func (s *EmbarcacionService) Create(embarcacion *entidades.NuevaEmbarcacionRequest) (int, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(embarcacion.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe embarcación con el mismo nombre
	existingNombre, err := s.embarcacionRepo.GetByNombre(embarcacion.Nombre)
	if err == nil && existingNombre != nil {
		return 0, errors.New("ya existe una embarcación con ese nombre")
	}

	// Validar que el estado sea uno de los permitidos
	if !isValidEstado(embarcacion.Estado) {
		return 0, errors.New("estado de embarcación no válido")
	}

	// Crear embarcación
	return s.embarcacionRepo.Create(embarcacion)
}

// GetByID obtiene una embarcación por su ID
func (s *EmbarcacionService) GetByID(id int) (*entidades.Embarcacion, error) {
	return s.embarcacionRepo.GetByID(id)
}

// Update actualiza una embarcación existente
func (s *EmbarcacionService) Update(id int, embarcacion *entidades.ActualizarEmbarcacionRequest) error {
	// Verificar que la embarcación existe
	existing, err := s.embarcacionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(embarcacion.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe otra embarcación con el mismo nombre
	if embarcacion.Nombre != existing.Nombre {
		existingNombre, err := s.embarcacionRepo.GetByNombre(embarcacion.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otra embarcación con ese nombre")
		}
	}

	// Validar que el estado sea uno de los permitidos
	if !isValidEstado(embarcacion.Estado) {
		return errors.New("estado de embarcación no válido")
	}

	// Actualizar embarcación
	return s.embarcacionRepo.Update(id, embarcacion)
}

// Delete elimina una embarcación (borrado lógico)
func (s *EmbarcacionService) Delete(id int) error {
	// Verificar que la embarcación exists
	_, err := s.embarcacionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar embarcación (borrado lógico)
	return s.embarcacionRepo.SoftDelete(id)
}

// List lista todas las embarcaciones
func (s *EmbarcacionService) List() ([]*entidades.Embarcacion, error) {
	return s.embarcacionRepo.List()
}

// ListBySede lista todas las embarcaciones de una sede específica
func (s *EmbarcacionService) ListBySede(idSede int) ([]*entidades.Embarcacion, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar embarcaciones de la sede
	return s.embarcacionRepo.ListBySede(idSede)
}

// ListByEstado lista todas las embarcaciones por estado
func (s *EmbarcacionService) ListByEstado(estado string) ([]*entidades.Embarcacion, error) {
	// Validar que el estado sea válido
	if !isValidEstado(estado) {
		return nil, errors.New("estado no válido")
	}

	return s.embarcacionRepo.ListByEstado(estado)
}

// isValidEstado valida que el estado sea uno de los permitidos
func isValidEstado(estado string) bool {
	validEstados := []string{"DISPONIBLE", "OCUPADA", "MANTENIMIENTO", "FUERA_DE_SERVICIO"}
	for _, validEstado := range validEstados {
		if estado == validEstado {
			return true
		}
	}
	return false
}
