package servicios

import (
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
)

// MetodoPagoService maneja la lógica de negocio para métodos de pago
type MetodoPagoService struct {
	metodoPagoRepo *repositorios.MetodoPagoRepository
}

// NewMetodoPagoService crea una nueva instancia de MetodoPagoService
func NewMetodoPagoService(metodoPagoRepo *repositorios.MetodoPagoRepository) *MetodoPagoService {
	return &MetodoPagoService{
		metodoPagoRepo: metodoPagoRepo,
	}
}

// Create crea un nuevo método de pago
func (s *MetodoPagoService) Create(metodoPago *entidades.NuevoMetodoPagoRequest) (int, error) {
	// Verificar si ya existe método de pago con el mismo nombre
	existing, err := s.metodoPagoRepo.GetByNombre(metodoPago.Nombre)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un método de pago con ese nombre")
	}

	// Crear método de pago
	return s.metodoPagoRepo.Create(metodoPago)
}

// GetByID obtiene un método de pago por su ID
func (s *MetodoPagoService) GetByID(id int) (*entidades.MetodoPago, error) {
	return s.metodoPagoRepo.GetByID(id)
}

// Update actualiza un método de pago existente
func (s *MetodoPagoService) Update(id int, metodoPago *entidades.ActualizarMetodoPagoRequest) error {
	// Verificar que el método de pago existe
	existing, err := s.metodoPagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro método de pago con el mismo nombre
	if metodoPago.Nombre != existing.Nombre {
		existingNombre, err := s.metodoPagoRepo.GetByNombre(metodoPago.Nombre)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro método de pago con ese nombre")
		}
	}

	// Actualizar método de pago
	return s.metodoPagoRepo.Update(id, metodoPago)
}

// Delete elimina un método de pago
func (s *MetodoPagoService) Delete(id int) error {
	// Verificar que el método de pago existe
	_, err := s.metodoPagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar método de pago
	return s.metodoPagoRepo.Delete(id)
}

// List lista todos los métodos de pago
func (s *MetodoPagoService) List() ([]*entidades.MetodoPago, error) {
	return s.metodoPagoRepo.List()
}
