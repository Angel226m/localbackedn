package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// CanalVentaService maneja la lógica de negocio para canales de venta
type CanalVentaService struct {
	canalVentaRepo *repositorios.CanalVentaRepository
	sedeRepo       *repositorios.SedeRepository
}

// NewCanalVentaService crea una nueva instancia de CanalVentaService
func NewCanalVentaService(
	canalVentaRepo *repositorios.CanalVentaRepository,
	sedeRepo *repositorios.SedeRepository,
) *CanalVentaService {
	return &CanalVentaService{
		canalVentaRepo: canalVentaRepo,
		sedeRepo:       sedeRepo,
	}
}

// Create crea un nuevo canal de venta
func (s *CanalVentaService) Create(canal *entidades.NuevoCanalVentaRequest) (int, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(canal.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe canal con el mismo nombre en la misma sede
	existing, err := s.canalVentaRepo.GetByNombre(canal.Nombre, canal.IDSede)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un canal de venta con ese nombre en esta sede")
	}

	// Crear canal
	return s.canalVentaRepo.Create(canal)
}

// GetByID obtiene un canal de venta por su ID
func (s *CanalVentaService) GetByID(id int) (*entidades.CanalVenta, error) {
	return s.canalVentaRepo.GetByID(id)
}

// Update actualiza un canal de venta existente
func (s *CanalVentaService) Update(id int, canal *entidades.ActualizarCanalVentaRequest) error {
	// Verificar que el canal existe
	existing, err := s.canalVentaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(canal.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe otro canal con el mismo nombre en la misma sede
	if canal.Nombre != existing.Nombre || canal.IDSede != existing.IDSede {
		existingNombre, err := s.canalVentaRepo.GetByNombre(canal.Nombre, canal.IDSede)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro canal de venta con ese nombre en esta sede")
		}
	}

	// Actualizar canal
	return s.canalVentaRepo.Update(id, canal)
}

// Delete elimina un canal de venta (borrado lógico)
func (s *CanalVentaService) Delete(id int) error {
	// Verificar que el canal existe
	_, err := s.canalVentaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar canal
	return s.canalVentaRepo.Delete(id)
}

// List lista todos los canales de venta
func (s *CanalVentaService) List() ([]*entidades.CanalVenta, error) {
	return s.canalVentaRepo.List()
}

// ListBySede lista todos los canales de venta de una sede específica
func (s *CanalVentaService) ListBySede(idSede int) ([]*entidades.CanalVenta, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar canales de venta por sede
	return s.canalVentaRepo.ListBySede(idSede)
}
