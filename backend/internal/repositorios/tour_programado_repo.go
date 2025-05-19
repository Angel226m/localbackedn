package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
	"time"
)

// TourProgramadoRepository maneja las operaciones de base de datos para tours programados
type TourProgramadoRepository struct {
	db *sql.DB
}

// NewTourProgramadoRepository crea una nueva instancia del repositorio
func NewTourProgramadoRepository(db *sql.DB) *TourProgramadoRepository {
	return &TourProgramadoRepository{
		db: db,
	}
}

// GetByID obtiene un tour programado por su ID
func (r *TourProgramadoRepository) GetByID(id int) (*entidades.TourProgramado, error) {
	tour := &entidades.TourProgramado{}
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.id_tour_programado = $1`

	err := r.db.QueryRow(query, id).Scan(
		&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
		&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
		&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
		&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
		&tour.NombreChofer, &tour.ApellidosChofer,
		&tour.HoraInicio, &tour.HoraFin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tour programado no encontrado")
		}
		return nil, err
	}

	return tour, nil
}

// Create guarda un nuevo tour programado en la base de datos
func (r *TourProgramadoRepository) Create(tour *entidades.NuevoTourProgramadoRequest) (int, error) {
	// Verificar que la combinación embarcación-fecha-horario no exista
	var count int
	queryCheck := `SELECT COUNT(*) FROM tour_programado 
                  WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3`

	err := r.db.QueryRow(queryCheck, tour.IDEmbarcacion, tour.Fecha, tour.IDHorario).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("ya existe un tour programado para esta embarcación, fecha y horario")
	}

	// Determinar estado si no se proporcionó
	estado := tour.Estado
	if estado == "" {
		estado = "PROGRAMADO"
	}

	// Crear tour programado
	var id int
	query := `INSERT INTO tour_programado (id_tipo_tour, id_embarcacion, id_horario, 
              fecha, cupo_maximo, cupo_disponible, estado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING id_tour_programado`

	err = r.db.QueryRow(
		query,
		tour.IDTipoTour,
		tour.IDEmbarcacion,
		tour.IDHorario,
		tour.Fecha,
		tour.CupoMaximo,
		tour.CupoMaximo, // Inicialmente cupo_disponible = cupo_maximo
		estado,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un tour programado
func (r *TourProgramadoRepository) Update(id int, tour *entidades.ActualizarTourProgramadoRequest) error {
	// Verificar que la combinación embarcación-fecha-horario no exista para otros tours
	var count int
	queryCheck := `SELECT COUNT(*) FROM tour_programado 
                  WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 AND id_tour_programado != $4`

	err := r.db.QueryRow(queryCheck, tour.IDEmbarcacion, tour.Fecha, tour.IDHorario, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe otro tour programado para esta embarcación, fecha y horario")
	}

	// Actualizar tour programado
	query := `UPDATE tour_programado SET 
              id_tipo_tour = $1, 
              id_embarcacion = $2, 
              id_horario = $3, 
              fecha = $4, 
              cupo_maximo = $5,
              cupo_disponible = $6,
              estado = $7
              WHERE id_tour_programado = $8`

	_, err = r.db.Exec(
		query,
		tour.IDTipoTour,
		tour.IDEmbarcacion,
		tour.IDHorario,
		tour.Fecha,
		tour.CupoMaximo,
		tour.CupoDisponible,
		tour.Estado,
		id,
	)

	return err
}

// UpdateEstado actualiza solo el estado de un tour programado
func (r *TourProgramadoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE tour_programado SET estado = $1 WHERE id_tour_programado = $2`
	_, err := r.db.Exec(query, estado, id)
	return err
}

// UpdateCupoDisponible actualiza el cupo disponible de un tour programado
func (r *TourProgramadoRepository) UpdateCupoDisponible(id int, nuevoDisponible int) error {
	query := `UPDATE tour_programado SET cupo_disponible = $1 WHERE id_tour_programado = $2`
	_, err := r.db.Exec(query, nuevoDisponible, id)
	return err
}

// Delete elimina un tour programado
func (r *TourProgramadoRepository) Delete(id int) error {
	// Verificar si hay reservas asociadas a este tour
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_tour_programado = $1`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar este tour programado porque tiene reservas asociadas")
	}

	// Eliminar tour programado
	query := `DELETE FROM tour_programado WHERE id_tour_programado = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los tours programados
func (r *TourProgramadoRepository) List() ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByFecha lista todos los tours programados para una fecha específica
func (r *TourProgramadoRepository) ListByFecha(fecha time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.fecha = $1
              ORDER BY ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByRangoFechas lista todos los tours programados para un rango de fechas
func (r *TourProgramadoRepository) ListByRangoFechas(fechaInicio, fechaFin time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.fecha BETWEEN $1 AND $2
              ORDER BY tp.fecha ASC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fechaInicio, fechaFin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByEstado lista todos los tours programados por estado
func (r *TourProgramadoRepository) ListByEstado(estado string) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.estado = $1
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByEmbarcacion lista todos los tours programados por embarcación
func (r *TourProgramadoRepository) ListByEmbarcacion(idEmbarcacion int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.id_embarcacion = $1
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idEmbarcacion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByChofer lista todos los tours programados asociados a un chofer
func (r *TourProgramadoRepository) ListByChofer(idChofer int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE e.id_usuario = $1
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idChofer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListToursProgramadosDisponibles lista todos los tours programados disponibles para reservación (estado PROGRAMADO y con cupo)
func (r *TourProgramadoRepository) ListToursProgramadosDisponibles() ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.estado = 'PROGRAMADO' 
              AND tp.cupo_disponible > 0 
              AND tp.fecha >= CURRENT_DATE
              ORDER BY tp.fecha ASC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// ListByTipoTour lista todos los tours programados por tipo de tour
func (r *TourProgramadoRepository) ListByTipoTour(idTipoTour int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.id_tipo_tour = $1
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idTipoTour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

// GetDisponibilidadDia retorna la disponibilidad de tours para una fecha específica por tipo de tour
func (r *TourProgramadoRepository) GetDisponibilidadDia(fecha time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado,
              tt.nombre, tt.precio_base, tt.duracion_minutos,
              e.nombre, e.capacidad,
              u.nombres, u.apellidos,
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI')
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              WHERE tp.fecha = $1
              AND tp.estado = 'PROGRAMADO'
              AND tp.cupo_disponible > 0
              ORDER BY ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado,
			&tour.NombreTipoTour, &tour.PrecioBase, &tour.DuracionMinutos,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}
