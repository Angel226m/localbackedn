package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// TipoPasajeService maneja la lógica de negocio para tipos de pasaje
type TipoPasajeService struct {
	tipoPasajeRepo *repositorios.TipoPasajeRepository
	sedeRepo       *repositorios.SedeRepository // Añadimos referencia al repositorio de sedes
}

// NewTipoPasajeService crea una nueva instancia de TipoPasajeService
func NewTipoPasajeService(tipoPasajeRepo *repositorios.TipoPasajeRepository, sedeRepo *repositorios.SedeRepository) *TipoPasajeService {
	return &TipoPasajeService{
		tipoPasajeRepo: tipoPasajeRepo,
		sedeRepo:       sedeRepo,
	}
}

// Create crea un nuevo tipo de pasaje
func (s *TipoPasajeService) Create(tipoPasaje *entidades.NuevoTipoPasajeRequest) (int, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(tipoPasaje.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe tipo de pasaje con el mismo nombre en esta sede
	existing, err := s.tipoPasajeRepo.GetByNombre(tipoPasaje.Nombre, tipoPasaje.IDSede)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un tipo de pasaje con ese nombre en esta sede")
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

	// Verificar si ya existe otro tipo de pasaje con el mismo nombre en la misma sede
	if tipoPasaje.Nombre != existing.Nombre {
		existingNombre, err := s.tipoPasajeRepo.GetByNombre(tipoPasaje.Nombre, existing.IDSede)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro tipo de pasaje con ese nombre en esta sede")
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

// ListBySede lista todos los tipos de pasaje de una sede específica
func (s *TipoPasajeService) ListBySede(idSede int) ([]*entidades.TipoPasaje, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	return s.tipoPasajeRepo.ListBySede(idSede)
}

// List lista todos los tipos de pasaje
func (s *TipoPasajeService) List() ([]*entidades.TipoPasaje, error) {
	return s.tipoPasajeRepo.List()
}
