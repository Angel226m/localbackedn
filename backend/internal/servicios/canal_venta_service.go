package servicios

import (
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
)

// CanalVentaService maneja la lógica de negocio para canales de venta
type CanalVentaService struct {
	canalVentaRepo *repositorios.CanalVentaRepository
}

// NewCanalVentaService crea una nueva instancia de CanalVentaService
func NewCanalVentaService(canalVentaRepo *repositorios.CanalVentaRepository) *CanalVentaService {
	return &CanalVentaService{
		canalVentaRepo: canalVentaRepo,
	}
}

// Create crea un nuevo canal de venta
func (s *CanalVentaService) Create(canal *entidades.NuevoCanalVentaRequest) (int, error) {
	// Verificar si ya existe canal con el mismo nombre
	existing, err := s.canalVentaRepo.GetByNombre(canal.Nombre)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un canal de venta con ese nombre")
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

	// Verificar si ya existe otro canal con el mismo nombre
	if canal.Nombre != existing.Nombre {
		existingNombre, err := s.canalVentaRepo.GetByNombre(canal.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro canal de venta con ese nombre")
		}
	}

	// Actualizar canal
	return s.canalVentaRepo.Update(id, canal)
}

// Delete elimina un canal de venta
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
