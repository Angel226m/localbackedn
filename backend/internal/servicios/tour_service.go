package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// TipoTourService maneja la lógica de negocio para tipos de tour
type TipoTourService struct {
	tipoTourRepo *repositorios.TipoTourRepository
	sedeRepo     *repositorios.SedeRepository
}

// NewTipoTourService crea una nueva instancia de TipoTourService
func NewTipoTourService(
	tipoTourRepo *repositorios.TipoTourRepository,
	sedeRepo *repositorios.SedeRepository,
) *TipoTourService {
	return &TipoTourService{
		tipoTourRepo: tipoTourRepo,
		sedeRepo:     sedeRepo,
	}
}

// Create crea un nuevo tipo de tour
func (s *TipoTourService) Create(tipoTour *entidades.NuevoTipoTourRequest) (int, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(tipoTour.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe un tipo de tour con el mismo nombre en la misma sede
	existingNombre, err := s.tipoTourRepo.GetByNombre(tipoTour.Nombre, tipoTour.IDSede)
	if err == nil && existingNombre != nil {
		return 0, errors.New("ya existe un tipo de tour con ese nombre en esta sede")
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

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(tipoTour.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe otro tipo de tour con el mismo nombre en la misma sede
	if tipoTour.Nombre != existing.Nombre || tipoTour.IDSede != existing.IDSede {
		existingNombre, err := s.tipoTourRepo.GetByNombre(tipoTour.Nombre, tipoTour.IDSede)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro tipo de tour con ese nombre en esta sede")
		}
	}

	// Actualizar tipo de tour
	return s.tipoTourRepo.Update(id, tipoTour)
}

// Delete elimina un tipo de tour (borrado lógico)
func (s *TipoTourService) Delete(id int) error {
	// Verificar que el tipo de tour existe
	_, err := s.tipoTourRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar tipo de tour (borrado lógico)
	return s.tipoTourRepo.Delete(id)
}

// List lista todos los tipos de tour
func (s *TipoTourService) List() ([]*entidades.TipoTour, error) {
	return s.tipoTourRepo.List()
}

// ListBySede lista todos los tipos de tour de una sede específica
func (s *TipoTourService) ListBySede(idSede int) ([]*entidades.TipoTour, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar tipos de tour de la sede
	return s.tipoTourRepo.ListBySede(idSede)
}
