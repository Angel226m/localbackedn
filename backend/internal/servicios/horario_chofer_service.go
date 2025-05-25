package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// parseTime convierte una cadena HH:MM a time.Time
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// HorarioChoferService maneja la lógica de negocio para horarios de chofer
type HorarioChoferService struct {
	horarioChoferRepo *repositorios.HorarioChoferRepository
	usuarioRepo       *repositorios.UsuarioRepository
	sedeRepo          *repositorios.SedeRepository
}

// NewHorarioChoferService crea una nueva instancia de HorarioChoferService
func NewHorarioChoferService(
	horarioChoferRepo *repositorios.HorarioChoferRepository,
	usuarioRepo *repositorios.UsuarioRepository,
	sedeRepo *repositorios.SedeRepository,
) *HorarioChoferService {
	return &HorarioChoferService{
		horarioChoferRepo: horarioChoferRepo,
		usuarioRepo:       usuarioRepo,
		sedeRepo:          sedeRepo,
	}
}

// Create crea un nuevo horario de chofer
func (s *HorarioChoferService) Create(horario *entidades.NuevoHorarioChoferRequest) (int, error) {
	// Verificar que el usuario exista y sea un chofer
	usuario, err := s.usuarioRepo.GetByID(horario.IDUsuario)
	if err != nil {
		return 0, errors.New("el usuario especificado no existe")
	}

	if usuario.Rol != "CHOFER" {
		return 0, errors.New("el usuario especificado no es un chofer")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(horario.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Convertir strings HH:MM a time.Time para verificar solapamiento
	horaInicio, err := parseTime(horario.HoraInicio)
	if err != nil {
		return 0, errors.New("formato de hora de inicio inválido, debe ser HH:MM")
	}

	horaFin, err := parseTime(horario.HoraFin)
	if err != nil {
		return 0, errors.New("formato de hora de fin inválido, debe ser HH:MM")
	}

	// Verificar que la hora de fin sea posterior a la hora de inicio
	if horaInicio.Equal(horaFin) || horaInicio.After(horaFin) {
		return 0, errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Verificar que la fecha de inicio no sea posterior a la fecha de fin (si existe)
	if horario.FechaFin != nil && horario.FechaInicio.After(*horario.FechaFin) {
		return 0, errors.New("la fecha de inicio no puede ser posterior a la fecha de fin")
	}

	// Verificar que no haya solapamiento de horarios para el mismo chofer
	overlap, err := s.horarioChoferRepo.VerifyHorarioOverlap(
		horario.IDUsuario,
		horaInicio,
		horaFin,
		&horario.FechaInicio,
		horario.FechaFin,
		0, // No excluimos ningún ID porque estamos creando uno nuevo
	)
	if err != nil {
		return 0, err
	}

	if overlap {
		return 0, errors.New("el horario se solapa con otro horario existente del mismo chofer")
	}

	// Crear horario de chofer
	return s.horarioChoferRepo.Create(horario)
}

// GetByID obtiene un horario de chofer por su ID
func (s *HorarioChoferService) GetByID(id int) (*entidades.HorarioChofer, error) {
	return s.horarioChoferRepo.GetByID(id)
}

// Update actualiza un horario de chofer existente
func (s *HorarioChoferService) Update(id int, horario *entidades.ActualizarHorarioChoferRequest) error {
	// Verificar que el horario de chofer existe
	existingHorario, err := s.horarioChoferRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el usuario exista y sea un chofer
	usuario, err := s.usuarioRepo.GetByID(horario.IDUsuario)
	if err != nil {
		return errors.New("el usuario especificado no existe")
	}

	if usuario.Rol != "CHOFER" {
		return errors.New("el usuario especificado no es un chofer")
	}

	// Verificar que la sede exista
	_, err = s.sedeRepo.GetByID(horario.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Convertir strings HH:MM a time.Time para verificar solapamiento
	horaInicio, err := parseTime(horario.HoraInicio)
	if err != nil {
		return errors.New("formato de hora de inicio inválido, debe ser HH:MM")
	}

	horaFin, err := parseTime(horario.HoraFin)
	if err != nil {
		return errors.New("formato de hora de fin inválido, debe ser HH:MM")
	}

	// Verificar que la hora de fin sea posterior a la hora de inicio
	if horaInicio.Equal(horaFin) || horaInicio.After(horaFin) {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Verificar que la fecha de inicio no sea posterior a la fecha de fin (si existe)
	if horario.FechaFin != nil && horario.FechaInicio.After(*horario.FechaFin) {
		return errors.New("la fecha de inicio no puede ser posterior a la fecha de fin")
	}

	// Verificar que no haya solapamiento de horarios para el mismo chofer
	// Solo verificamos si cambia algún dato relevante
	if horario.IDUsuario != existingHorario.IDUsuario ||
		horaInicio != existingHorario.HoraInicio ||
		horaFin != existingHorario.HoraFin ||
		horario.FechaInicio != existingHorario.FechaInicio ||
		(horario.FechaFin == nil && existingHorario.FechaFin != nil) ||
		(horario.FechaFin != nil && (existingHorario.FechaFin == nil || *horario.FechaFin != *existingHorario.FechaFin)) {

		overlap, err := s.horarioChoferRepo.VerifyHorarioOverlap(
			horario.IDUsuario,
			horaInicio,
			horaFin,
			&horario.FechaInicio,
			horario.FechaFin,
			id, // Excluimos el ID actual para no detectarlo como solapamiento consigo mismo
		)
		if err != nil {
			return err
		}

		if overlap {
			return errors.New("el horario se solapa con otro horario existente del mismo chofer")
		}
	}

	// Actualizar horario de chofer
	return s.horarioChoferRepo.Update(id, horario)
}

// Delete elimina un horario de chofer (borrado lógico)
func (s *HorarioChoferService) Delete(id int) error {
	// Verificar que el horario de chofer existe
	_, err := s.horarioChoferRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar horario de chofer (borrado lógico)
	return s.horarioChoferRepo.Delete(id)
}

// List lista todos los horarios de chofer
func (s *HorarioChoferService) List() ([]*entidades.HorarioChofer, error) {
	return s.horarioChoferRepo.List()
}

// ListByChofer lista todos los horarios de un chofer específico
func (s *HorarioChoferService) ListByChofer(idChofer int) ([]*entidades.HorarioChofer, error) {
	// Verificar que el usuario exista y sea un chofer
	usuario, err := s.usuarioRepo.GetByID(idChofer)
	if err != nil {
		return nil, errors.New("el chofer especificado no existe")
	}

	if usuario.Rol != "CHOFER" {
		return nil, errors.New("el usuario especificado no es un chofer")
	}

	// Listar horarios del chofer
	return s.horarioChoferRepo.ListByChofer(idChofer)
}

// ListActiveByChofer lista los horarios activos de un chofer
func (s *HorarioChoferService) ListActiveByChofer(idChofer int) ([]*entidades.HorarioChofer, error) {
	// Verificar que el usuario exista y sea un chofer
	usuario, err := s.usuarioRepo.GetByID(idChofer)
	if err != nil {
		return nil, errors.New("el chofer especificado no existe")
	}

	if usuario.Rol != "CHOFER" {
		return nil, errors.New("el usuario especificado no es un chofer")
	}

	// Listar horarios activos del chofer
	return s.horarioChoferRepo.ListActiveByChofer(idChofer)
}

// ListByDia lista todos los horarios de choferes disponibles para un día específico
func (s *HorarioChoferService) ListByDia(diaSemana int) ([]*entidades.HorarioChofer, error) {
	if diaSemana < 1 || diaSemana > 7 {
		return nil, errors.New("día de la semana inválido, debe ser un número entre 1 (Lunes) y 7 (Domingo)")
	}

	// Listar horarios de choferes por día
	return s.horarioChoferRepo.ListByDia(diaSemana)
}
