package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// MetodoPagoService maneja la lógica de negocio para métodos de pago
type MetodoPagoService struct {
	metodoPagoRepo *repositorios.MetodoPagoRepository
	sedeRepo       *repositorios.SedeRepository
}

// NewMetodoPagoService crea una nueva instancia de MetodoPagoService
func NewMetodoPagoService(
	metodoPagoRepo *repositorios.MetodoPagoRepository,
	sedeRepo *repositorios.SedeRepository,
) *MetodoPagoService {
	return &MetodoPagoService{
		metodoPagoRepo: metodoPagoRepo,
		sedeRepo:       sedeRepo,
	}
}

// Create crea un nuevo método de pago
func (s *MetodoPagoService) Create(metodoPago *entidades.NuevoMetodoPagoRequest) (int, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(metodoPago.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe método de pago con el mismo nombre en la misma sede
	existing, err := s.metodoPagoRepo.GetByNombre(metodoPago.Nombre, metodoPago.IDSede)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un método de pago con ese nombre en esta sede")
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

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(metodoPago.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar si ya existe otro método de pago con el mismo nombre en la misma sede
	if metodoPago.Nombre != existing.Nombre || metodoPago.IDSede != existing.IDSede {
		existingNombre, err := s.metodoPagoRepo.GetByNombre(metodoPago.Nombre, metodoPago.IDSede)
		if err == nil && existingNombre != nil && existingNombre.ID != id {
			return errors.New("ya existe otro método de pago con ese nombre en esta sede")
		}
	}

	// Actualizar método de pago
	return s.metodoPagoRepo.Update(id, metodoPago)
}

// Delete elimina un método de pago (borrado lógico)
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

// ListBySede lista todos los métodos de pago de una sede específica
func (s *MetodoPagoService) ListBySede(idSede int) ([]*entidades.MetodoPago, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar métodos de pago por sede
	return s.metodoPagoRepo.ListBySede(idSede)
}
