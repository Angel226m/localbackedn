package servicios

import (
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
)

// TipoTourService maneja la lógica de negocio para tipos de tour
type TipoTourService struct {
	tipoTourRepo *repositorios.TipoTourRepository
}

// NewTipoTourService crea una nueva instancia de TipoTourService
func NewTipoTourService(tipoTourRepo *repositorios.TipoTourRepository) *TipoTourService {
	return &TipoTourService{
		tipoTourRepo: tipoTourRepo,
	}
}

// Create crea un nuevo tipo de tour
func (s *TipoTourService) Create(tipoTour *entidades.NuevoTipoTourRequest) (int, error) {
	// Verificar si ya existe un tipo de tour con el mismo nombre
	existing, err := s.tipoTourRepo.GetByNombre(tipoTour.Nombre)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un tipo de tour con ese nombre")
	}

	// Crear tipo de tour
	return s.tipoTourRepo.Create(tipoTour)
}

// GetByID obtiene un tipo de tour por su ID
func (s *TipoTourService) GetByID(id int) (*entidades.TipoTour, error) {
	return s.tipoTourRepo.GetByID(id)
}

// Update actualiza un tipo de tour existente
func (s *TipoTourService) Update(id int, tipoTour *entidades.ActualizarTipoTourRequest) error {
	// Verificar que el tipo de tour existe
	existing, err := s.tipoTourRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro tipo de tour con el mismo nombre
	if tipoTour.Nombre != existing.Nombre {
		existingNombre, err := s.tipoTourRepo.GetByNombre(tipoTour.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro tipo de tour con ese nombre")
		}
	}

	// Actualizar tipo de tour
	return s.tipoTourRepo.Update(id, tipoTour)
}

// Delete elimina un tipo de tour
func (s *TipoTourService) Delete(id int) error {
	// Verificar que el tipo de tour existe
	_, err := s.tipoTourRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar tipo de tour
	return s.tipoTourRepo.Delete(id)
}

// List lista todos los tipos de tour
func (s *TipoTourService) List() ([]*entidades.TipoTour, error) {
	return s.tipoTourRepo.List()
}
