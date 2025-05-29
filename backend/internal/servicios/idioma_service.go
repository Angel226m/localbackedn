package servicios

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// IdiomaService maneja la lógica de negocio para idiomas
type IdiomaService struct {
	idiomaRepo *repositorios.IdiomaRepository
}

// NewIdiomaService crea una nueva instancia de IdiomaService
func NewIdiomaService(idiomaRepo *repositorios.IdiomaRepository) *IdiomaService {
	return &IdiomaService{
		idiomaRepo: idiomaRepo,
	}
}

// GetByID obtiene un idioma por su ID
func (s *IdiomaService) GetByID(id int) (*entidades.Idioma, error) {
	return s.idiomaRepo.GetByID(id)
}

// GetByNombre obtiene un idioma por su nombre
func (s *IdiomaService) GetByNombre(nombre string) (*entidades.Idioma, error) {
	return s.idiomaRepo.GetByNombre(nombre)
}

// Create crea un nuevo idioma
func (s *IdiomaService) Create(idioma *entidades.Idioma) (int, error) {
	// Verificar si ya existe un idioma con el mismo nombre
	existingNombre, err := s.idiomaRepo.GetByNombre(idioma.Nombre)
	if err == nil && existingNombre != nil {
		return existingNombre.ID, nil // Devolvemos el ID del idioma existente
	}

	// Crear idioma
	return s.idiomaRepo.Create(idioma)
}

// Update actualiza un idioma existente
func (s *IdiomaService) Update(id int, idioma *entidades.Idioma) error {
	// Verificar que el idioma existe
	existing, err := s.idiomaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro idioma con el mismo nombre
	if idioma.Nombre != existing.Nombre {
		existingNombre, err := s.idiomaRepo.GetByNombre(idioma.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return err
		}
	}

	// Actualizar ID para asegurar que sea el correcto
	idioma.ID = id

	// Actualizar idioma
	return s.idiomaRepo.Update(idioma)
}

// Delete elimina lógicamente un idioma (soft delete)
func (s *IdiomaService) Delete(id int) error {
	// Verificar que el idioma existe
	_, err := s.idiomaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar idioma (soft delete)
	return s.idiomaRepo.SoftDelete(id)
}

// Restore restaura un idioma eliminado
func (s *IdiomaService) Restore(id int) error {
	return s.idiomaRepo.Restore(id)
}

// List lista todos los idiomas activos
func (s *IdiomaService) List() ([]*entidades.Idioma, error) {
	return s.idiomaRepo.List()
}

// ListDeleted lista todos los idiomas eliminados
func (s *IdiomaService) ListDeleted() ([]*entidades.Idioma, error) {
	return s.idiomaRepo.ListDeleted()
}
