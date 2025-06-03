package servicios

import (
	"errors"
	"fmt"

	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

type GaleriaTourService struct {
	repo         *repositorios.GaleriaTourRepo
	tipoTourRepo *repositorios.TipoTourRepository
}

func NewGaleriaTourService(repo *repositorios.GaleriaTourRepo, tipoTourRepo *repositorios.TipoTourRepository) *GaleriaTourService {
	return &GaleriaTourService{repo: repo, tipoTourRepo: tipoTourRepo}
}

func (s *GaleriaTourService) CrearImagen(req *entidades.GaleriaTourRequest) (int, error) {
	// Verificar que existe el tipo de tour
	// Cambiado de ObtenerPorID a GetByID
	_, err := s.tipoTourRepo.GetByID(req.IDTipoTour)
	if err != nil {
		return 0, fmt.Errorf("tipo de tour no encontrado: %v", err)
	}

	galeria := &entidades.GaleriaTour{
		IDTipoTour:  req.IDTipoTour,
		URLImagen:   req.URLImagen,
		Descripcion: req.Descripcion,
		Orden:       req.Orden,
	}

	return s.repo.Crear(galeria)
}

func (s *GaleriaTourService) ObtenerPorID(id int) (*entidades.GaleriaTour, error) {
	return s.repo.ObtenerPorID(id)
}

func (s *GaleriaTourService) ListarPorTipoTour(idTipoTour int) ([]*entidades.GaleriaTour, error) {
	// Verificar que existe el tipo de tour
	// Cambiado de ObtenerPorID a GetByID
	_, err := s.tipoTourRepo.GetByID(idTipoTour)
	if err != nil {
		return nil, fmt.Errorf("tipo de tour no encontrado: %v", err)
	}

	return s.repo.ListarPorTipoTour(idTipoTour)
}

func (s *GaleriaTourService) ActualizarImagen(id int, req *entidades.GaleriaTourUpdateRequest) error {
	galeria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return fmt.Errorf("imagen no encontrada: %v", err)
	}

	galeria.URLImagen = req.URLImagen
	galeria.Descripcion = req.Descripcion
	galeria.Orden = req.Orden

	return s.repo.Actualizar(galeria)
}

func (s *GaleriaTourService) EliminarImagen(id int) error {
	galeria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return errors.New("imagen no encontrada")
	}

	return s.repo.Eliminar(galeria.ID)
}

func (s *GaleriaTourService) EliminarImagenesPorTipoTour(idTipoTour int) error {
	return s.repo.EliminarPorTipoTour(idTipoTour)
}
