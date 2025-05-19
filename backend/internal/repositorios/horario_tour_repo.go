package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
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
	query := `SELECT h.id_horario, h.id_tipo_tour, h.hora_inicio, h.hora_fin, 
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, t.nombre, t.descripcion
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              WHERE h.id_horario = $1`

	err := r.db.QueryRow(query, id).Scan(
		&horario.ID, &horario.IDTipoTour, &horario.HoraInicio, &horario.HoraFin,
		&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
		&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
		&horario.DisponibleDomingo, &horario.NombreTipoTour, &horario.DescripcionTipoTour,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("horario de tour no encontrado")
		}
		return nil, err
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
	query := `INSERT INTO horario_tour (id_tipo_tour, hora_inicio, hora_fin, 
              disponible_lunes, disponible_martes, disponible_miercoles, 
              disponible_jueves, disponible_viernes, disponible_sabado, disponible_domingo) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
              RETURNING id_horario`

	err = r.db.QueryRow(
		query,
		horario.IDTipoTour,
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
              hora_inicio = $2, 
              hora_fin = $3, 
              disponible_lunes = $4, 
              disponible_martes = $5, 
              disponible_miercoles = $6, 
              disponible_jueves = $7, 
              disponible_viernes = $8, 
              disponible_sabado = $9, 
              disponible_domingo = $10 
              WHERE id_horario = $11`

	_, err = r.db.Exec(
		query,
		horario.IDTipoTour,
		horaInicio,
		horaFin,
		horario.DisponibleLunes,
		horario.DisponibleMartes,
		horario.DisponibleMiercoles,
		horario.DisponibleJueves,
		horario.DisponibleViernes,
		horario.DisponibleSabado,
		horario.DisponibleDomingo,
		id,
	)

	return err
}

// Delete elimina un horario de tour
func (r *HorarioTourRepository) Delete(id int) error {
	// Comprobar si hay tours programados que dependen de este horario
	var count int
	queryCheck := `SELECT COUNT(*) FROM tour_programado WHERE id_horario = $1`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este horario porque hay tours programados que dependen de él")
	}

	// Si no hay dependencias, procedemos a eliminar
	query := `DELETE FROM horario_tour WHERE id_horario = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los horarios de tour
func (r *HorarioTourRepository) List() ([]*entidades.HorarioTour, error) {
	query := `SELECT h.id_horario, h.id_tipo_tour, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, t.nombre, t.descripcion
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              ORDER BY t.nombre, h.hora_inicio`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.NombreTipoTour, &horario.DescripcionTipoTour,
		)
		if err != nil {
			return nil, err
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
	query := `SELECT h.id_horario, h.id_tipo_tour, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, t.nombre, t.descripcion
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              WHERE h.id_tipo_tour = $1
              ORDER BY h.hora_inicio`

	rows, err := r.db.Query(query, idTipoTour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.NombreTipoTour, &horario.DescripcionTipoTour,
		)
		if err != nil {
			return nil, err
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

	query := `SELECT h.id_horario, h.id_tipo_tour, h.hora_inicio, h.hora_fin,
              h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
              h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, 
              h.disponible_domingo, t.nombre, t.descripcion
              FROM horario_tour h
              INNER JOIN tipo_tour t ON h.id_tipo_tour = t.id_tipo_tour
              WHERE ` + condition + `
              ORDER BY t.nombre, h.hora_inicio`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	horarios := []*entidades.HorarioTour{}

	for rows.Next() {
		horario := &entidades.HorarioTour{}
		err := rows.Scan(
			&horario.ID, &horario.IDTipoTour, &horario.HoraInicio, &horario.HoraFin,
			&horario.DisponibleLunes, &horario.DisponibleMartes, &horario.DisponibleMiercoles,
			&horario.DisponibleJueves, &horario.DisponibleViernes, &horario.DisponibleSabado,
			&horario.DisponibleDomingo, &horario.NombreTipoTour, &horario.DescripcionTipoTour,
		)
		if err != nil {
			return nil, err
		}
		horarios = append(horarios, horario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return horarios, nil
}
