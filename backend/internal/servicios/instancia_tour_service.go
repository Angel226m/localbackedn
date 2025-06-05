package servicios

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
)

// InstanciaTourService maneja la lógica de negocio para instancias de tour
type InstanciaTourService struct {
	instanciaTourRepo *repositorios.InstanciaTourRepository
}

// NewInstanciaTourService crea una nueva instancia de InstanciaTourService
func NewInstanciaTourService(instanciaTourRepo *repositorios.InstanciaTourRepository) *InstanciaTourService {
	return &InstanciaTourService{
		instanciaTourRepo: instanciaTourRepo,
	}
}

// GetByID obtiene una instancia de tour por su ID
func (s *InstanciaTourService) GetByID(id int) (*entidades.InstanciaTour, error) {
	return s.instanciaTourRepo.GetByID(id)
}

// Create crea una nueva instancia de tour
func (s *InstanciaTourService) Create(instancia *entidades.NuevaInstanciaTourRequest) (int, error) {
	return s.instanciaTourRepo.Create(instancia)
}

// Update actualiza una instancia de tour existente
func (s *InstanciaTourService) Update(id int, instancia *entidades.ActualizarInstanciaTourRequest) error {
	return s.instanciaTourRepo.Update(id, instancia)
}

// Delete elimina una instancia de tour (soft delete)
func (s *InstanciaTourService) Delete(id int) error {
	return s.instanciaTourRepo.Delete(id)
}

// List lista todas las instancias de tour
func (s *InstanciaTourService) List() ([]*entidades.InstanciaTour, error) {
	return s.instanciaTourRepo.List()
}

// ListByTourProgramado lista todas las instancias de un tour programado específico
func (s *InstanciaTourService) ListByTourProgramado(idTourProgramado int) ([]*entidades.InstanciaTour, error) {
	return s.instanciaTourRepo.ListByTourProgramado(idTourProgramado)
}

// ListByFiltros lista instancias de tour según filtros específicos
func (s *InstanciaTourService) ListByFiltros(filtros entidades.FiltrosInstanciaTour) ([]*entidades.InstanciaTour, error) {
	return s.instanciaTourRepo.ListByFiltros(filtros)
}

// AsignarChofer asigna un chofer a una instancia de tour
func (s *InstanciaTourService) AsignarChofer(id int, idChofer int) error {
	return s.instanciaTourRepo.AsignarChofer(id, idChofer)
}

// GenerarInstanciasDeTourProgramado genera instancias para un tour programado
func (s *InstanciaTourService) GenerarInstanciasDeTourProgramado(idTourProgramado int) (int, error) {
	return s.instanciaTourRepo.GenerarInstanciasDeTourProgramado(idTourProgramado)
}
