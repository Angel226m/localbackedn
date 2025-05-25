package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
	"time"
)

// HorarioTourRepository maneja las operaciones de base de datos para horarios de tour
type HorarioTourRepository struct {
	db *sql.DB
}

// NewHorarioTourRepository crea una nueva instancia del repositorio
func NewHorarioTourRepository(db *sql.DB) *HorarioTourRepository {
	return &HorarioTourRepository{
		db: db,
	}
}

// parseTime convierte una cadena HH:MM a time.Time
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// GetByID obtiene un horario de tour por su ID
func (r *HorarioTourRepository) GetByID(id int) (*entidades.HorarioTour, error) {
	horario := &entidades.HorarioTour{}
	query := `SELECT h.id_horario, h.id_tipo_tour, h.id_sede, h.hora_inicio, h.hora_fin, 
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, h.eliminado, t.nombre, t.descripcion, s.nombre
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              INNER JOIN sede s ON h.id_sede = s.id_sede
              WHERE h.id_horario = $1`

	var descripcion sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&horario.ID, &horario.IDTipoTour, &horario.IDSede, &horario.HoraInicio, &horario.HoraFin,
		&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
		&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
		&horario.DisponibleDomingo, &horario.Eliminado, &horario.NombreTipoTour, &descripcion, &horario.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("horario de tour no encontrado")
		}
		return nil, err
	}

	if descripcion.Valid {
		horario.DescripcionTipoTour = descripcion.String
	}

	return horario, nil
}

// Create guarda un nuevo horario de tour en la base de datos
func (r *HorarioTourRepository) Create(horario *entidades.NuevoHorarioTourRequest) (int, error) {
	// Convertir strings HH:MM a time.Time para la base de datos
	horaInicio, err := parseTime(horario.HoraInicio)
	if err != nil {
		return 0, errors.New("formato de hora de inicio inválido, debe ser HH:MM")
	}

	horaFin, err := parseTime(horario.HoraFin)
	if err != nil {
		return 0, errors.New("formato de hora de fin inválido, debe ser HH:MM")
	}

	// Verificar que al menos un día esté disponible
	if !horario.DisponibleLunes && !horario.DisponibleMartes && !horario.DisponibleMiercoles &&
		!horario.DisponibleJueves && !horario.DisponibleViernes && !horario.DisponibleSabado &&
		!horario.DisponibleDomingo {
		return 0, errors.New("debe seleccionar al menos un día disponible")
	}

	var id int
	query := `INSERT INTO horario_tour (id_tipo_tour, id_sede, hora_inicio, hora_fin, 
              disponible_lunes, disponible_martes, disponible_miercoles, 
              disponible_jueves, disponible_viernes, disponible_sabado, disponible_domingo, eliminado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false) 
              RETURNING id_horario`

	err = r.db.QueryRow(
		query,
		horario.IDTipoTour,
		horario.IDSede,
		horaInicio,
		horaFin,
		horario.DisponibleLunes,
		horario.DisponibleMartes,
		horario.DisponibleMiercoles,
		horario.DisponibleJueves,
		horario.DisponibleViernes,
		horario.DisponibleSabado,
		horario.DisponibleDomingo,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un horario de tour
func (r *HorarioTourRepository) Update(id int, horario *entidades.ActualizarHorarioTourRequest) error {
	// Convertir strings HH:MM a time.Time para la base de datos
	horaInicio, err := parseTime(horario.HoraInicio)
	if err != nil {
		return errors.New("formato de hora de inicio inválido, debe ser HH:MM")
	}

	horaFin, err := parseTime(horario.HoraFin)
	if err != nil {
		return errors.New("formato de hora de fin inválido, debe ser HH:MM")
	}

	// Verificar que al menos un día esté disponible
	if !horario.DisponibleLunes && !horario.DisponibleMartes && !horario.DisponibleMiercoles &&
		!horario.DisponibleJueves && !horario.DisponibleViernes && !horario.DisponibleSabado &&
		!horario.DisponibleDomingo {
		return errors.New("debe seleccionar al menos un día disponible")
	}

	query := `UPDATE horario_tour SET 
              id_tipo_tour = $1, 
              id_sede = $2,
              hora_inicio = $3, 
              hora_fin = $4, 
              disponible_lunes = $5, 
              disponible_martes = $6, 
              disponible_miercoles = $7, 
              disponible_jueves = $8, 
              disponible_viernes = $9, 
              disponible_sabado = $10, 
              disponible_domingo = $11,
              eliminado = $12
              WHERE id_horario = $13`

	_, err = r.db.Exec(
		query,
		horario.IDTipoTour,
		horario.IDSede,
		horaInicio,
		horaFin,
		horario.DisponibleLunes,
		horario.DisponibleMartes,
		horario.DisponibleMiercoles,
		horario.DisponibleJueves,
		horario.DisponibleViernes,
		horario.DisponibleSabado,
		horario.DisponibleDomingo,
		horario.Eliminado,
		id,
	)

	return err
}

// Delete marca un horario de tour como eliminado (borrado lógico)
func (r *HorarioTourRepository) Delete(id int) error {
	// Comprobar si hay tours programados que dependen de este horario
	var count int
	queryCheck := `SELECT COUNT(*) FROM tour_programado WHERE id_horario = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este horario porque hay tours programados que dependen de él")
	}

	// Si no hay dependencias, procedemos a marcar como eliminado
	query := `UPDATE horario_tour SET eliminado = true WHERE id_horario = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los horarios de tour no eliminados
func (r *HorarioTourRepository) List() ([]*entidades.HorarioTour, error) {
	query := `SELECT h.id_horario, h.id_tipo_tour, h.id_sede, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, h.eliminado, t.nombre, t.descripcion, s.nombre
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              INNER JOIN sede s ON h.id_sede = s.id_sede
              WHERE h.eliminado = false
              ORDER BY t.nombre, h.hora_inicio`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		var descripcion sql.NullString

		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.IDSede, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.Eliminado, &horario.NombreTipoTour, &descripcion, &horario.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		if descripcion.Valid {
			horario.DescripcionTipoTour = descripcion.String
		}

		horarios = append(horarios, horario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return horarios, nil
}

// ListByTipoTour lista todos los horarios asociados a un tipo de tour específico
func (r *HorarioTourRepository) ListByTipoTour(idTipoTour int) ([]*entidades.HorarioTour, error) {
	query := `SELECT h.id_horario, h.id_tipo_tour, h.id_sede, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, h.eliminado, t.nombre, t.descripcion, s.nombre
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              INNER JOIN sede s ON h.id_sede = s.id_sede
              WHERE h.id_tipo_tour = $1 AND h.eliminado = false
              ORDER BY h.hora_inicio`

	rows, err := r.db.Query(query, idTipoTour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		var descripcion sql.NullString

		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.IDSede, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.Eliminado, &horario.NombreTipoTour, &descripcion, &horario.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		if descripcion.Valid {
			horario.DescripcionTipoTour = descripcion.String
		}

		horarios = append(horarios, horario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return horarios, nil
}

// ListByDia lista todos los horarios disponibles para un día específico (1=Lunes, 7=Domingo)
func (r *HorarioTourRepository) ListByDia(diaSemana int) ([]*entidades.HorarioTour, error) {
	var condition string
	switch diaSemana {
	case 1:
		condition = "h.disponible_lunes = true"
	case 2:
		condition = "h.disponible_martes = true"
	case 3:
		condition = "h.disponible_miercoles = true"
	case 4:
		condition = "h.disponible_jueves = true"
	case 5:
		condition = "h.disponible_viernes = true"
	case 6:
		condition = "h.disponible_sabado = true"
	case 7:
		condition = "h.disponible_domingo = true"
	default:
		return nil, errors.New("día de la semana inválido, debe ser un número entre 1 (Lunes) y 7 (Domingo)")
	}

	query := `SELECT h.id_horario, h.id_tipo_tour, h.id_sede, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, h.eliminado, t.nombre, t.descripcion, s.nombre
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              INNER JOIN sede s ON h.id_sede = s.id_sede
              WHERE ` + condition + ` AND h.eliminado = false
              ORDER BY t.nombre, h.hora_inicio`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		var descripcion sql.NullString

		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.IDSede, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.Eliminado, &horario.NombreTipoTour, &descripcion, &horario.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		if descripcion.Valid {
			horario.DescripcionTipoTour = descripcion.String
		}

		horarios = append(horarios, horario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return horarios, nil
}
