package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// PaquetePasajesService maneja la lógica de negocio para paquetes de pasajes
type PaquetePasajesService struct {
	paquetePasajesRepo *repositorios.PaquetePasajesRepository
	sedeRepo           *repositorios.SedeRepository
	tipoTourRepo       *repositorios.TipoTourRepository
}

// NewPaquetePasajesService crea una nueva instancia de PaquetePasajesService
func NewPaquetePasajesService(
	paquetePasajesRepo *repositorios.PaquetePasajesRepository,
	sedeRepo *repositorios.SedeRepository,
	tipoTourRepo *repositorios.TipoTourRepository,
) *PaquetePasajesService {
	return &PaquetePasajesService{
		paquetePasajesRepo: paquetePasajesRepo,
		sedeRepo:           sedeRepo,
		tipoTourRepo:       tipoTourRepo,
	}
}

// Create crea un nuevo paquete de pasajes
func (s *PaquetePasajesService) Create(paquete *entidades.NuevoPaquetePasajesRequest) (int, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(paquete.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar que el tipo de tour existe
	_, err = s.tipoTourRepo.GetByID(paquete.IDTipoTour)
	if err != nil {
		return 0, errors.New("el tipo de tour especificado no existe")
	}

	// Verificar si ya existe paquete con el mismo nombre en esta sede
	existing, err := s.paquetePasajesRepo.GetByNombre(paquete.Nombre, paquete.IDSede)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un paquete de pasajes con ese nombre en esta sede")
	}

	// Crear paquete de pasajes
	return s.paquetePasajesRepo.Create(paquete)
}

// GetByID obtiene un paquete de pasajes por su ID
func (s *PaquetePasajesService) GetByID(id int) (*entidades.PaquetePasajes, error) {
	return s.paquetePasajesRepo.GetByID(id)
}

// Update actualiza un paquete de pasajes existente
func (s *PaquetePasajesService) Update(id int, paquete *entidades.ActualizarPaquetePasajesRequest) error {
	// Verificar que el paquete de pasajes existe
	existing, err := s.paquetePasajesRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el tipo de tour existe
	_, err = s.tipoTourRepo.GetByID(paquete.IDTipoTour)
	if err != nil {
		return errors.New("el tipo de tour especificado no existe")
	}

	// Verificar si ya existe otro paquete con el mismo nombre en la misma sede
	if paquete.Nombre != existing.Nombre {
		existingNombre, err := s.paquetePasajesRepo.GetByNombre(paquete.Nombre, existing.IDSede)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro paquete de pasajes con ese nombre en esta sede")
		}
	}

	// Actualizar paquete de pasajes
	return s.paquetePasajesRepo.Update(id, paquete)
}

// Delete elimina un paquete de pasajes
func (s *PaquetePasajesService) Delete(id int) error {
	// Verificar que el paquete de pasajes existe
	_, err := s.paquetePasajesRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar paquete de pasajes
	return s.paquetePasajesRepo.Delete(id)
}

// ListBySede lista todos los paquetes de pasajes de una sede específica
func (s *PaquetePasajesService) ListBySede(idSede int) ([]*entidades.PaquetePasajes, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	return s.paquetePasajesRepo.ListBySede(idSede)
}

// ListByTipoTour lista todos los paquetes de pasajes de un tipo de tour específico
func (s *PaquetePasajesService) ListByTipoTour(idTipoTour int) ([]*entidades.PaquetePasajes, error) {
	// Verificar que el tipo de tour existe
	_, err := s.tipoTourRepo.GetByID(idTipoTour)
	if err != nil {
		return nil, errors.New("el tipo de tour especificado no existe")
	}

	return s.paquetePasajesRepo.ListByTipoTour(idTipoTour)
}

// List lista todos los paquetes de pasajes
func (s *PaquetePasajesService) List() ([]*entidades.PaquetePasajes, error) {
	return s.paquetePasajesRepo.List()
}
