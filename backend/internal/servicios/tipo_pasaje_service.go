package servicios

import (
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
)

// TipoPasajeService maneja la lógica de negocio para tipos de pasaje
type TipoPasajeService struct {
	tipoPasajeRepo *repositorios.TipoPasajeRepository
}

// NewTipoPasajeService crea una nueva instancia de TipoPasajeService
func NewTipoPasajeService(tipoPasajeRepo *repositorios.TipoPasajeRepository) *TipoPasajeService {
	return &TipoPasajeService{
		tipoPasajeRepo: tipoPasajeRepo,
	}
}

// Create crea un nuevo tipo de pasaje
func (s *TipoPasajeService) Create(tipoPasaje *entidades.NuevoTipoPasajeRequest) (int, error) {
	// Verificar si ya existe tipo de pasaje con el mismo nombre
	existing, err := s.tipoPasajeRepo.GetByNombre(tipoPasaje.Nombre)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un tipo de pasaje con ese nombre")
	}

	// Crear tipo de pasaje
	return s.tipoPasajeRepo.Create(tipoPasaje)
}

// GetByID obtiene un tipo de pasaje por su ID
func (s *TipoPasajeService) GetByID(id int) (*entidades.TipoPasaje, error) {
	return s.tipoPasajeRepo.GetByID(id)
}

// Update actualiza un tipo de pasaje existente
func (s *TipoPasajeService) Update(id int, tipoPasaje *entidades.ActualizarTipoPasajeRequest) error {
	// Verificar que el tipo de pasaje existe
	existing, err := s.tipoPasajeRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro tipo de pasaje con el mismo nombre
	if tipoPasaje.Nombre != existing.Nombre {
		existingNombre, err := s.tipoPasajeRepo.GetByNombre(tipoPasaje.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro tipo de pasaje con ese nombre")
		}
	}

	// Actualizar tipo de pasaje
	return s.tipoPasajeRepo.Update(id, tipoPasaje)
}

// Delete elimina un tipo de pasaje
func (s *TipoPasajeService) Delete(id int) error {
	// Verificar que el tipo de pasaje existe
	_, err := s.tipoPasajeRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar tipo de pasaje
	return s.tipoPasajeRepo.Delete(id)
}

// List lista todos los tipos de pasaje
func (s *TipoPasajeService) List() ([]*entidades.TipoPasaje, error) {
	return s.tipoPasajeRepo.List()
}
