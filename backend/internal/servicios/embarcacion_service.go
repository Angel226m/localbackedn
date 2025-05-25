package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// EmbarcacionService maneja la lógica de negocio para embarcaciones
type EmbarcacionService struct {
	embarcacionRepo *repositorios.EmbarcacionRepository
	usuarioRepo     *repositorios.UsuarioRepository
}

// NewEmbarcacionService crea una nueva instancia de EmbarcacionService
func NewEmbarcacionService(
	embarcacionRepo *repositorios.EmbarcacionRepository,
	usuarioRepo *repositorios.UsuarioRepository,
) *EmbarcacionService {
	return &EmbarcacionService{
		embarcacionRepo: embarcacionRepo,
		usuarioRepo:     usuarioRepo,
	}
}

// Create crea una nueva embarcación
func (s *EmbarcacionService) Create(embarcacion *entidades.NuevaEmbarcacionRequest) (int, error) {
	// Verificar que la sede exista
	// Esta verificación debe implementarse si tienes un repositorio de sedes

	// Verificar si ya existe embarcación con el mismo nombre
	existingNombre, err := s.embarcacionRepo.GetByNombre(embarcacion.Nombre)
	if err == nil && existingNombre != nil {
		return 0, errors.New("ya existe una embarcación con ese nombre")
	}

	// Verificar que el chofer exista y tenga rol CHOFER
	chofer, err := s.usuarioRepo.GetByID(embarcacion.IDUsuario)
	if err != nil {
		return 0, errors.New("el chofer especificado no existe")
	}

	if chofer.Rol != "CHOFER" {
		return 0, errors.New("el usuario especificado no es un chofer")
	}

	// Validar que el estado sea uno de los permitidos
	if embarcacion.Estado != "DISPONIBLE" &&
		embarcacion.Estado != "OCUPADA" &&
		embarcacion.Estado != "MANTENIMIENTO" &&
		embarcacion.Estado != "FUERA_DE_SERVICIO" {
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

	// Verificar si ya existe otra embarcación con el mismo nombre
	if embarcacion.Nombre != existing.Nombre {
		existingNombre, err := s.embarcacionRepo.GetByNombre(embarcacion.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otra embarcación con ese nombre")
		}
	}

	// Verificar que el chofer exista y tenga rol CHOFER
	chofer, err := s.usuarioRepo.GetByID(embarcacion.IDUsuario)
	if err != nil {
		return errors.New("el chofer especificado no existe")
	}

	if chofer.Rol != "CHOFER" {
		return errors.New("el usuario especificado no es un chofer")
	}

	// Validar que el estado sea uno de los permitidos
	if embarcacion.Estado != "DISPONIBLE" &&
		embarcacion.Estado != "OCUPADA" &&
		embarcacion.Estado != "MANTENIMIENTO" &&
		embarcacion.Estado != "FUERA_DE_SERVICIO" {
		return errors.New("estado de embarcación no válido")
	}

	// Actualizar embarcación
	return s.embarcacionRepo.Update(id, embarcacion)
}

// Delete elimina una embarcación (borrado lógico)
func (s *EmbarcacionService) Delete(id int) error {
	// Verificar que la embarcación existe
	_, err := s.embarcacionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar embarcación (borrado lógico)
	return s.embarcacionRepo.Delete(id)
}

// List lista todas las embarcaciones
func (s *EmbarcacionService) List() ([]*entidades.Embarcacion, error) {
	return s.embarcacionRepo.List()
}

// ListByChofer lista todas las embarcaciones de un chofer específico
func (s *EmbarcacionService) ListByChofer(idChofer int) ([]*entidades.Embarcacion, error) {
	// Verificar que el chofer exista y tenga rol CHOFER
	chofer, err := s.usuarioRepo.GetByID(idChofer)
	if err != nil {
		return nil, errors.New("el chofer especificado no existe")
	}

	if chofer.Rol != "CHOFER" {
		return nil, errors.New("el usuario especificado no es un chofer")
	}

	// Listar embarcaciones del chofer
	return s.embarcacionRepo.ListByChofer(idChofer)
}
