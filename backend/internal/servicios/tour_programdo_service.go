package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// TourProgramadoService maneja la lógica de negocio para tours programados
type TourProgramadoService struct {
	tourProgramadoRepo *repositorios.TourProgramadoRepository
	tipoTourRepo       *repositorios.TipoTourRepository
	embarcacionRepo    *repositorios.EmbarcacionRepository
	horarioTourRepo    *repositorios.HorarioTourRepository
	sedeRepo           *repositorios.SedeRepository
	usuarioRepo        *repositorios.UsuarioRepository // Repositorio para validar choferes
}

// NewTourProgramadoService crea una nueva instancia de TourProgramadoService
func NewTourProgramadoService(
	tourProgramadoRepo *repositorios.TourProgramadoRepository,
	tipoTourRepo *repositorios.TipoTourRepository,
	embarcacionRepo *repositorios.EmbarcacionRepository,
	horarioTourRepo *repositorios.HorarioTourRepository,
	sedeRepo *repositorios.SedeRepository,
	usuarioRepo *repositorios.UsuarioRepository,
) *TourProgramadoService {
	return &TourProgramadoService{
		tourProgramadoRepo: tourProgramadoRepo,
		tipoTourRepo:       tipoTourRepo,
		embarcacionRepo:    embarcacionRepo,
		horarioTourRepo:    horarioTourRepo,
		sedeRepo:           sedeRepo,
		usuarioRepo:        usuarioRepo,
	}
}

// Create crea un nuevo tour programado
func (s *TourProgramadoService) Create(tour *entidades.NuevoTourProgramadoRequest) (int, error) {
	// Verificar que el tipo de tour exista
	_, err := s.tipoTourRepo.GetByID(tour.IDTipoTour)
	if err != nil {
		return 0, errors.New("el tipo de tour especificado no existe")
	}

	// Verificar que la embarcación exista
	_, err = s.embarcacionRepo.GetByID(tour.IDEmbarcacion)
	if err != nil {
		return 0, errors.New("la embarcación especificada no existe")
	}

	// Verificar que el horario de tour exista
	horario, err := s.horarioTourRepo.GetByID(tour.IDHorario)
	if err != nil {
		return 0, errors.New("el horario de tour especificado no existe")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(tour.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar que el chofer exista y tenga rol de CHOFER si se proporciona
	if tour.IDChofer != nil {
		usuario, err := s.usuarioRepo.GetByID(*tour.IDChofer)
		if err != nil {
			return 0, errors.New("el chofer especificado no existe")
		}
		if usuario.Rol != "CHOFER" {
			return 0, errors.New("el usuario especificado no tiene rol de CHOFER")
		}
	}

	// Verificar que el horario corresponde al tipo de tour
	if horario.IDTipoTour != tour.IDTipoTour {
		return 0, errors.New("el horario especificado no corresponde al tipo de tour")
	}

	// Verificar que la fecha no sea anterior a la fecha actual
	if tour.Fecha.Before(time.Now().Truncate(24 * time.Hour)) {
		return 0, errors.New("no se puede programar un tour para una fecha pasada")
	}

	// Obtener el día de la semana de la fecha (0 = domingo, 1 = lunes, ..., 6 = sábado)
	diaSemana := int(tour.Fecha.Weekday())
	if diaSemana == 0 {
		diaSemana = 7 // Ajustar para que domingo sea 7 en lugar de 0
	}

	// Verificar disponibilidad del horario para ese día de la semana
	disponible := false
	switch diaSemana {
	case 1:
		disponible = horario.DisponibleLunes
	case 2:
		disponible = horario.DisponibleMartes
	case 3:
		disponible = horario.DisponibleMiercoles
	case 4:
		disponible = horario.DisponibleJueves
	case 5:
		disponible = horario.DisponibleViernes
	case 6:
		disponible = horario.DisponibleSabado
	case 7:
		disponible = horario.DisponibleDomingo
	}

	if !disponible {
		return 0, errors.New("el horario no está disponible para el día de la semana seleccionado")
	}

	// Crear tour programado
	return s.tourProgramadoRepo.Create(tour)
}

// GetByID obtiene un tour programado por su ID
func (s *TourProgramadoService) GetByID(id int) (*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.GetByID(id)
}

// Update actualiza un tour programado existente
func (s *TourProgramadoService) Update(id int, tour *entidades.ActualizarTourProgramadoRequest) error {
	// Verificar que el tour programado existe
	existingTour, err := s.tourProgramadoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que no tenga reservas antes de permitir cambios
	if existingTour.CupoMaximo != existingTour.CupoDisponible && existingTour.CupoDisponible != tour.CupoDisponible {
		return errors.New("no se puede modificar un tour que ya tiene reservas, solo se puede cancelar")
	}

	// Verificar que el tipo de tour exista
	_, err = s.tipoTourRepo.GetByID(tour.IDTipoTour)
	if err != nil {
		return errors.New("el tipo de tour especificado no existe")
	}

	// Verificar que la embarcación exista
	_, err = s.embarcacionRepo.GetByID(tour.IDEmbarcacion)
	if err != nil {
		return errors.New("la embarcación especificada no existe")
	}

	// Verificar que el horario de tour exista
	horario, err := s.horarioTourRepo.GetByID(tour.IDHorario)
	if err != nil {
		return errors.New("el horario de tour especificado no existe")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(tour.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar que el chofer exista y tenga rol de CHOFER si se proporciona
	if tour.IDChofer != nil {
		usuario, err := s.usuarioRepo.GetByID(*tour.IDChofer)
		if err != nil {
			return errors.New("el chofer especificado no existe")
		}
		if usuario.Rol != "CHOFER" {
			return errors.New("el usuario especificado no tiene rol de CHOFER")
		}
	}

	// Verificar que el horario corresponde al tipo de tour
	if horario.IDTipoTour != tour.IDTipoTour {
		return errors.New("el horario especificado no corresponde al tipo de tour")
	}

	// Verificar que la fecha no sea anterior a la fecha actual
	if tour.Fecha.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("no se puede programar un tour para una fecha pasada")
	}

	// Obtener el día de la semana de la fecha (0 = domingo, 1 = lunes, ..., 6 = sábado)
	diaSemana := int(tour.Fecha.Weekday())
	if diaSemana == 0 {
		diaSemana = 7 // Ajustar para que domingo sea 7 en lugar de 0
	}

	// Verificar disponibilidad del horario para ese día de la semana
	disponible := false
	switch diaSemana {
	case 1:
		disponible = horario.DisponibleLunes
	case 2:
		disponible = horario.DisponibleMartes
	case 3:
		disponible = horario.DisponibleMiercoles
	case 4:
		disponible = horario.DisponibleJueves
	case 5:
		disponible = horario.DisponibleViernes
	case 6:
		disponible = horario.DisponibleSabado
	case 7:
		disponible = horario.DisponibleDomingo
	}

	if !disponible {
		return errors.New("el horario no está disponible para el día de la semana seleccionado")
	}

	// Verificar cupo máximo y disponible
	if tour.CupoDisponible > tour.CupoMaximo {
		return errors.New("el cupo disponible no puede ser mayor que el cupo máximo")
	}

	// Actualizar tour programado
	return s.tourProgramadoRepo.Update(id, tour)
}

// CambiarEstado cambia el estado de un tour programado
func (s *TourProgramadoService) CambiarEstado(id int, estado string) error {
	// Verificar estado válido
	estadosValidos := map[string]bool{
		"PROGRAMADO": true,
		"COMPLETADO": true,
		"CANCELADO":  true,
	}

	if !estadosValidos[estado] {
		return errors.New("estado inválido, debe ser PROGRAMADO, COMPLETADO o CANCELADO")
	}

	// Verificar que el tour programado existe
	existingTour, err := s.tourProgramadoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Si ya tiene ese estado, no hacer nada
	if existingTour.Estado == estado {
		return nil
	}

	// Cambiar estado
	return s.tourProgramadoRepo.UpdateEstado(id, estado)
}

// ReservarCupo disminuye el cupo disponible de un tour programado
func (s *TourProgramadoService) ReservarCupo(id int, cantidad int) error {
	// Verificar cantidad válida
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a 0")
	}

	// Verificar que el tour programado existe
	tour, err := s.tourProgramadoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el tour está programado
	if tour.Estado != "PROGRAMADO" {
		return errors.New("no se puede reservar un tour que no está programado")
	}

	// Verificar disponibilidad
	if tour.CupoDisponible < cantidad {
		return errors.New("no hay suficiente cupo disponible")
	}

	// Actualizar cupo disponible
	nuevoCupo := tour.CupoDisponible - cantidad
	return s.tourProgramadoRepo.UpdateCupoDisponible(id, nuevoCupo)
}

// LiberarCupo aumenta el cupo disponible de un tour programado
func (s *TourProgramadoService) LiberarCupo(id int, cantidad int) error {
	// Verificar cantidad válida
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a 0")
	}

	// Verificar que el tour programado existe
	tour, err := s.tourProgramadoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el tour está programado
	if tour.Estado != "PROGRAMADO" {
		return errors.New("no se puede liberar cupo de un tour que no está programado")
	}

	// Verificar que no exceda el cupo máximo
	nuevoCupo := tour.CupoDisponible + cantidad
	if nuevoCupo > tour.CupoMaximo {
		return errors.New("el cupo liberado excede el cupo máximo")
	}

	// Actualizar cupo disponible
	return s.tourProgramadoRepo.UpdateCupoDisponible(id, nuevoCupo)
}

// Delete elimina un tour programado (CORREGIDO - usa SoftDelete)
func (s *TourProgramadoService) Delete(id int) error {
	// Verificar que el tour programado existe
	existingTour, err := s.tourProgramadoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que no tenga reservas
	if existingTour.CupoMaximo != existingTour.CupoDisponible {
		return errors.New("no se puede eliminar un tour que ya tiene reservas, solo se puede cancelar")
	}

	// Eliminar tour programado (usando SoftDelete en lugar de Delete)
	return s.tourProgramadoRepo.SoftDelete(id)
}

// List lista todos los tours programados
func (s *TourProgramadoService) List() ([]*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.List()
}

// ListByFecha lista todos los tours programados para una fecha específica
func (s *TourProgramadoService) ListByFecha(fecha time.Time) ([]*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.ListByFecha(fecha)
}

// ListByRangoFechas lista todos los tours programados para un rango de fechas
func (s *TourProgramadoService) ListByRangoFechas(fechaInicio, fechaFin time.Time) ([]*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.ListByRangoFechas(fechaInicio, fechaFin)
}

// ListByEstado lista todos los tours programados por estado
func (s *TourProgramadoService) ListByEstado(estado string) ([]*entidades.TourProgramado, error) {
	// Verificar estado válido
	estadosValidos := map[string]bool{
		"PROGRAMADO": true,
		"COMPLETADO": true,
		"CANCELADO":  true,
	}

	if !estadosValidos[estado] {
		return nil, errors.New("estado inválido, debe ser PROGRAMADO, COMPLETADO o CANCELADO")
	}

	return s.tourProgramadoRepo.ListByEstado(estado)
}

// ListByEmbarcacion lista todos los tours programados por embarcación
func (s *TourProgramadoService) ListByEmbarcacion(idEmbarcacion int) ([]*entidades.TourProgramado, error) {
	// Verificar que la embarcación exista
	_, err := s.embarcacionRepo.GetByID(idEmbarcacion)
	if err != nil {
		return nil, errors.New("la embarcación especificada no existe")
	}

	return s.tourProgramadoRepo.ListByEmbarcacion(idEmbarcacion)
}

// ListByChofer lista todos los tours programados asociados a un chofer
func (s *TourProgramadoService) ListByChofer(idChofer int) ([]*entidades.TourProgramado, error) {
	// Verificar que el chofer exista y tenga rol CHOFER
	usuario, err := s.usuarioRepo.GetByID(idChofer)
	if err != nil {
		return nil, errors.New("el chofer especificado no existe")
	}
	if usuario.Rol != "CHOFER" {
		return nil, errors.New("el usuario especificado no tiene rol de CHOFER")
	}

	return s.tourProgramadoRepo.ListByChofer(idChofer)
}

// ListToursProgramadosDisponibles lista todos los tours programados disponibles para reservación
func (s *TourProgramadoService) ListToursProgramadosDisponibles() ([]*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.ListToursProgramadosDisponibles()
}

// ListByTipoTour lista todos los tours programados por tipo de tour
func (s *TourProgramadoService) ListByTipoTour(idTipoTour int) ([]*entidades.TourProgramado, error) {
	// Verificar que el tipo de tour exista
	_, err := s.tipoTourRepo.GetByID(idTipoTour)
	if err != nil {
		return nil, errors.New("el tipo de tour especificado no existe")
	}

	return s.tourProgramadoRepo.ListByTipoTour(idTipoTour)
}

// GetDisponibilidadDia retorna la disponibilidad de tours para una fecha específica
func (s *TourProgramadoService) GetDisponibilidadDia(fecha time.Time) ([]*entidades.TourProgramado, error) {
	return s.tourProgramadoRepo.GetDisponibilidadDia(fecha)
}

// ListBySede lista todos los tours programados de una sede específica
func (s *TourProgramadoService) ListBySede(idSede int) ([]*entidades.TourProgramado, error) {
	// Verificar que la sede exista
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	return s.tourProgramadoRepo.ListBySede(idSede)
}
