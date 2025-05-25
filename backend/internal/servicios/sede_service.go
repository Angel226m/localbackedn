package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// SedeService maneja la lógica de negocio para sedes
type SedeService struct {
	sedeRepo *repositorios.SedeRepository
}

// NewSedeService crea una nueva instancia de SedeService
func NewSedeService(sedeRepo *repositorios.SedeRepository) *SedeService {
	return &SedeService{
		sedeRepo: sedeRepo,
	}
}

// Create crea una nueva sede
func (s *SedeService) Create(sede *entidades.NuevaSedeRequest) (int, error) {
	// Crear sede
	return s.sedeRepo.Create(sede)
}

// GetByID obtiene una sede por su ID
func (s *SedeService) GetByID(id int) (*entidades.Sede, error) {
	return s.sedeRepo.GetByID(id)
}

// Update actualiza una sede existente
func (s *SedeService) Update(id int, sede *entidades.ActualizarSedeRequest) error {
	// Actualizar sede
	return s.sedeRepo.Update(id, sede)
}

// Delete elimina una sede (borrado lógico)
func (s *SedeService) Delete(id int) error {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar sede (soft delete)
	return s.sedeRepo.SoftDelete(id)
}

// Restore restaura una sede eliminada lógicamente
func (s *SedeService) Restore(id int) error {
	return s.sedeRepo.Restore(id)
}

// List lista todas las sedes
func (s *SedeService) List() ([]*entidades.Sede, error) {
	return s.sedeRepo.List()
}

// GetByCiudad obtiene sedes por ciudad
func (s *SedeService) GetByCiudad(ciudad string) ([]*entidades.Sede, error) {
	if ciudad == "" {
		return nil, errors.New("la ciudad no puede estar vacía")
	}
	return s.sedeRepo.GetByCiudad(ciudad)
}

// GetByPais obtiene sedes por país
func (s *SedeService) GetByPais(pais string) ([]*entidades.Sede, error) {
	if pais == "" {
		return nil, errors.New("el país no puede estar vacío")
	}
	return s.sedeRepo.GetByPais(pais)
}
