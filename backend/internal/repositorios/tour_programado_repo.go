package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
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
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.id_tour_programado = $1 AND tp.eliminado = FALSE`

	var idChofer sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
		&tour.IDSede, &idChofer, &tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado, &tour.Eliminado,
		&tour.NombreTipoTour, &tour.DuracionMinutos, &tour.CantidadPasajeros,
		&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
		&tour.NombreChofer, &tour.ApellidosChofer,
		&tour.HoraInicio, &tour.HoraFin,
		&tour.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tour programado no encontrado")
		}
		return nil, err
	}

	// Convertir sql.NullInt64 a *int
	if idChofer.Valid {
		choferID := int(idChofer.Int64)
		tour.IDChofer = &choferID
	}

	return tour, nil
}

// Create guarda un nuevo tour programado en la base de datos
func (r *TourProgramadoRepository) Create(tour *entidades.NuevoTourProgramadoRequest) (int, error) {
	// Verificar que la combinación embarcación-fecha-horario no exista
	var count int
	queryCheck := `SELECT COUNT(*) FROM tour_programado 
                  WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 AND eliminado = FALSE`

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
	query := `INSERT INTO tour_programado (id_tipo_tour, id_embarcacion, id_horario, id_sede, id_chofer,
              fecha, cupo_maximo, cupo_disponible, estado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
              RETURNING id_tour_programado`

	err = r.db.QueryRow(
		query,
		tour.IDTipoTour,
		tour.IDEmbarcacion,
		tour.IDHorario,
		tour.IDSede,
		tour.IDChofer,
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
                  WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 
                  AND id_tour_programado != $4 AND eliminado = FALSE`

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
              id_sede = $4,
              id_chofer = $5,
              fecha = $6, 
              cupo_maximo = $7,
              cupo_disponible = $8,
              estado = $9
              WHERE id_tour_programado = $10 AND eliminado = FALSE`

	result, err := r.db.Exec(
		query,
		tour.IDTipoTour,
		tour.IDEmbarcacion,
		tour.IDHorario,
		tour.IDSede,
		tour.IDChofer,
		tour.Fecha,
		tour.CupoMaximo,
		tour.CupoDisponible,
		tour.Estado,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("tour programado no encontrado o ya fue eliminado")
	}

	return nil
}

// UpdateEstado actualiza solo el estado de un tour programado
func (r *TourProgramadoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE tour_programado SET estado = $1 WHERE id_tour_programado = $2 AND eliminado = FALSE`
	result, err := r.db.Exec(query, estado, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("tour programado no encontrado o ya fue eliminado")
	}

	return nil
}

// UpdateCupoDisponible actualiza el cupo disponible de un tour programado
func (r *TourProgramadoRepository) UpdateCupoDisponible(id int, nuevoDisponible int) error {
	query := `UPDATE tour_programado SET cupo_disponible = $1 WHERE id_tour_programado = $2 AND eliminado = FALSE`
	result, err := r.db.Exec(query, nuevoDisponible, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("tour programado no encontrado o ya fue eliminado")
	}

	return nil
}

// SoftDelete marca un tour programado como eliminado (eliminado lógico)
func (r *TourProgramadoRepository) SoftDelete(id int) error {
	// Verificar si hay reservas asociadas a este tour
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_tour_programado = $1 AND eliminado = FALSE`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar este tour programado porque tiene reservas asociadas")
	}

	// Marcar como eliminado en lugar de eliminar físicamente
	query := `UPDATE tour_programado SET eliminado = TRUE WHERE id_tour_programado = $1 AND eliminado = FALSE`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("tour programado no encontrado o ya fue eliminado")
	}

	return nil
}

// List lista todos los tours programados no eliminados
func (r *TourProgramadoRepository) List() ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	return r.executeTourQuery(query)
}

// ListByFecha lista todos los tours programados para una fecha específica
func (r *TourProgramadoRepository) ListByFecha(fecha time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.fecha = $1 AND tp.eliminado = FALSE
              ORDER BY ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListByRangoFechas lista todos los tours programados para un rango de fechas
func (r *TourProgramadoRepository) ListByRangoFechas(fechaInicio, fechaFin time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.fecha BETWEEN $1 AND $2 AND tp.eliminado = FALSE
              ORDER BY tp.fecha ASC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fechaInicio, fechaFin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListByEstado lista todos los tours programados por estado
func (r *TourProgramadoRepository) ListByEstado(estado string) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.estado = $1 AND tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListByEmbarcacion lista todos los tours programados por embarcación
func (r *TourProgramadoRepository) ListByEmbarcacion(idEmbarcacion int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.id_embarcacion = $1 AND tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idEmbarcacion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListByChofer lista todos los tours programados asociados a un chofer
func (r *TourProgramadoRepository) ListByChofer(idChofer int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.id_chofer = $1 AND tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idChofer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListToursProgramadosDisponibles lista todos los tours programados disponibles para reservación
func (r *TourProgramadoRepository) ListToursProgramadosDisponibles() ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.estado = 'PROGRAMADO' 
              AND tp.cupo_disponible > 0 
              AND tp.fecha >= CURRENT_DATE
              AND tp.eliminado = FALSE
              ORDER BY tp.fecha ASC, ht.hora_inicio ASC`

	return r.executeTourQuery(query)
}

// ListByTipoTour lista todos los tours programados por tipo de tour
func (r *TourProgramadoRepository) ListByTipoTour(idTipoTour int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.id_tipo_tour = $1 AND tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idTipoTour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// GetDisponibilidadDia retorna la disponibilidad de tours para una fecha específica
func (r *TourProgramadoRepository) GetDisponibilidadDia(fecha time.Time) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.fecha = $1
              AND tp.estado = 'PROGRAMADO'
              AND tp.cupo_disponible > 0
              AND tp.eliminado = FALSE
              ORDER BY ht.hora_inicio ASC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// ListBySede lista todos los tours programados de una sede específica
func (r *TourProgramadoRepository) ListBySede(idSede int) ([]*entidades.TourProgramado, error) {
	query := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
              tp.id_sede, tp.id_chofer, tp.fecha, tp.cupo_maximo, tp.cupo_disponible, tp.estado, tp.eliminado,
              tt.nombre, tt.duracion_minutos, tt.cantidad_pasajeros,
              e.nombre, e.capacidad,
              COALESCE(c.nombres, ''), COALESCE(c.apellidos, ''),
              TO_CHAR(ht.hora_inicio, 'HH24:MI'), TO_CHAR(ht.hora_fin, 'HH24:MI'),
              s.nombre
              FROM tour_programado tp
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              LEFT JOIN usuario c ON tp.id_chofer = c.id_usuario
              WHERE tp.id_sede = $1 AND tp.eliminado = FALSE
              ORDER BY tp.fecha DESC, ht.hora_inicio ASC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// Método auxiliar para ejecutar consultas de tours
func (r *TourProgramadoRepository) executeTourQuery(query string, args ...interface{}) ([]*entidades.TourProgramado, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTours(rows)
}

// Método auxiliar para escanear resultados de consultas de tours
func (r *TourProgramadoRepository) scanTours(rows *sql.Rows) ([]*entidades.TourProgramado, error) {
	tours := []*entidades.TourProgramado{}

	for rows.Next() {
		tour := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		err := rows.Scan(
			&tour.ID, &tour.IDTipoTour, &tour.IDEmbarcacion, &tour.IDHorario,
			&tour.IDSede, &idChofer, &tour.Fecha, &tour.CupoMaximo, &tour.CupoDisponible, &tour.Estado, &tour.Eliminado,
			&tour.NombreTipoTour, &tour.DuracionMinutos, &tour.CantidadPasajeros,
			&tour.NombreEmbarcacion, &tour.CapacidadEmbarcacion,
			&tour.NombreChofer, &tour.ApellidosChofer,
			&tour.HoraInicio, &tour.HoraFin,
			&tour.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Convertir sql.NullInt64 a *int
		if idChofer.Valid {
			choferID := int(idChofer.Int64)
			tour.IDChofer = &choferID
		}

		tours = append(tours, tour)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}
