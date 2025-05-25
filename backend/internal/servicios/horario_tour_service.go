package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"strconv"
)

// HorarioTourService maneja la lógica de negocio para horarios de tour
type HorarioTourService struct {
	horarioTourRepo *repositorios.HorarioTourRepository
	tipoTourRepo    *repositorios.TipoTourRepository
	sedeRepo        *repositorios.SedeRepository
}

// NewHorarioTourService crea una nueva instancia de HorarioTourService
func NewHorarioTourService(
	horarioTourRepo *repositorios.HorarioTourRepository,
	tipoTourRepo *repositorios.TipoTourRepository,
	sedeRepo *repositorios.SedeRepository,
) *HorarioTourService {
	return &HorarioTourService{
		horarioTourRepo: horarioTourRepo,
		tipoTourRepo:    tipoTourRepo,
		sedeRepo:        sedeRepo,
	}
}

// Create crea un nuevo horario de tour
func (s *HorarioTourService) Create(horario *entidades.NuevoHorarioTourRequest) (int, error) {
	// Verificar que el tipo de tour exista
	_, err := s.tipoTourRepo.GetByID(horario.IDTipoTour)
	if err != nil {
		return 0, errors.New("el tipo de tour especificado no existe")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(horario.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Validar que la hora de fin sea posterior a la hora de inicio
	if horario.HoraInicio >= horario.HoraFin {
		return 0, errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Crear horario de tour
	return s.horarioTourRepo.Create(horario)
}

// GetByID obtiene un horario de tour por su ID
func (s *HorarioTourService) GetByID(id int) (*entidades.HorarioTour, error) {
	return s.horarioTourRepo.GetByID(id)
}

// Update actualiza un horario de tour existente
func (s *HorarioTourService) Update(id int, horario *entidades.ActualizarHorarioTourRequest) error {
	// Verificar que el horario de tour existe
	_, err := s.horarioTourRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el tipo de tour exista
	_, err = s.tipoTourRepo.GetByID(horario.IDTipoTour)
	if err != nil {
		return errors.New("el tipo de tour especificado no existe")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(horario.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Validar que la hora de fin sea posterior a la hora de inicio
	if horario.HoraInicio >= horario.HoraFin {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Actualizar horario de tour
	return s.horarioTourRepo.Update(id, horario)
}

// Delete elimina un horario de tour (borrado lógico)
func (s *HorarioTourService) Delete(id int) error {
	// Verificar que el horario de tour existe
	_, err := s.horarioTourRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar horario de tour
	return s.horarioTourRepo.Delete(id)
}

// List lista todos los horarios de tour
func (s *HorarioTourService) List() ([]*entidades.HorarioTour, error) {
	return s.horarioTourRepo.List()
}

// ListByTipoTour lista todos los horarios asociados a un tipo de tour específico
func (s *HorarioTourService) ListByTipoTour(idTipoTour int) ([]*entidades.HorarioTour, error) {
	// Verificar que el tipo de tour exista
	_, err := s.tipoTourRepo.GetByID(idTipoTour)
	if err != nil {
		return nil, errors.New("el tipo de tour especificado no existe")
	}

	// Listar horarios por tipo de tour
	return s.horarioTourRepo.ListByTipoTour(idTipoTour)
}

// ListByDia lista todos los horarios disponibles para un día específico
func (s *HorarioTourService) ListByDia(dia string) ([]*entidades.HorarioTour, error) {
	// Convertir el día de string a int
	diaSemana, err := strconv.Atoi(dia)
	if err != nil {
		return nil, errors.New("formato de día inválido, debe ser un número entre 1 (Lunes) y 7 (Domingo)")
	}

	// Validar el rango del día de la semana
	if diaSemana < 1 || diaSemana > 7 {
		return nil, errors.New("día inválido, debe estar entre 1 (Lunes) y 7 (Domingo)")
	}

	// Listar horarios por día
	return s.horarioTourRepo.ListByDia(diaSemana)
}
