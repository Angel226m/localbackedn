package servicios

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// UsuarioIdiomaService maneja la lógica de negocio para la relación usuario-idioma
type UsuarioIdiomaService struct {
	usuarioIdiomaRepo *repositorios.UsuarioIdiomaRepository
	idiomaRepo        *repositorios.IdiomaRepository
	usuarioRepo       *repositorios.UsuarioRepository
}

// NewUsuarioIdiomaService crea una nueva instancia de UsuarioIdiomaService
func NewUsuarioIdiomaService(
	usuarioIdiomaRepo *repositorios.UsuarioIdiomaRepository,
	idiomaRepo *repositorios.IdiomaRepository,
	usuarioRepo *repositorios.UsuarioRepository,
) *UsuarioIdiomaService {
	return &UsuarioIdiomaService{
		usuarioIdiomaRepo: usuarioIdiomaRepo,
		idiomaRepo:        idiomaRepo,
		usuarioRepo:       usuarioRepo,
	}
}

// GetIdiomasByUsuarioID obtiene todos los idiomas de un usuario
func (s *UsuarioIdiomaService) GetIdiomasByUsuarioID(usuarioID int) ([]*entidades.UsuarioIdioma, error) {
	// Verificar que el usuario existe
	_, err := s.usuarioRepo.GetByID(usuarioID)
	if err != nil {
		return nil, err
	}

	return s.usuarioIdiomaRepo.GetByUsuarioID(usuarioID)
}

// AsignarIdioma asigna un idioma a un usuario
func (s *UsuarioIdiomaService) AsignarIdioma(usuarioID, idiomaID int, nivel string) error {
	// Verificar que el usuario existe
	_, err := s.usuarioRepo.GetByID(usuarioID)
	if err != nil {
		return err
	}

	// Verificar que el idioma existe
	_, err = s.idiomaRepo.GetByID(idiomaID)
	if err != nil {
		return err
	}

	// Si no se especifica nivel, usar 'básico' como valor predeterminado
	if nivel == "" {
		nivel = "básico"
	}

	// Asignar idioma
	return s.usuarioIdiomaRepo.AsignarIdioma(usuarioID, idiomaID, nivel)
}

// DesasignarIdioma elimina la asignación de un idioma a un usuario
func (s *UsuarioIdiomaService) DesasignarIdioma(usuarioID, idiomaID int) error {
	// Verificar que el usuario existe
	_, err := s.usuarioRepo.GetByID(usuarioID)
	if err != nil {
		return err
	}

	// Verificar que el idioma existe
	_, err = s.idiomaRepo.GetByID(idiomaID)
	if err != nil {
		return err
	}

	// Desasignar idioma
	return s.usuarioIdiomaRepo.DesasignarIdioma(usuarioID, idiomaID)
}

// ActualizarIdiomasUsuario actualiza todos los idiomas de un usuario
func (s *UsuarioIdiomaService) ActualizarIdiomasUsuario(usuarioID int, idiomasIDs []int) error {
	// Verificar que el usuario existe
	_, err := s.usuarioRepo.GetByID(usuarioID)
	if err != nil {
		return err
	}

	// Verificar que todos los idiomas existen
	for _, idiomaID := range idiomasIDs {
		_, err := s.idiomaRepo.GetByID(idiomaID)
		if err != nil {
			return err
		}
	}

	// Actualizar idiomas
	return s.usuarioIdiomaRepo.ActualizarIdiomasUsuario(usuarioID, idiomasIDs)
}

// GetUsuariosByIdiomaID obtiene todos los usuarios con un idioma específico
func (s *UsuarioIdiomaService) GetUsuariosByIdiomaID(idiomaID int) ([]*entidades.Usuario, error) {
	// Verificar que el idioma existe
	_, err := s.idiomaRepo.GetByID(idiomaID)
	if err != nil {
		return nil, err
	}

	return s.usuarioIdiomaRepo.GetByIdiomaID(idiomaID)
}
