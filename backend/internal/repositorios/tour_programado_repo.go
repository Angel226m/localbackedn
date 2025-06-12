package repositorios

import (
	"database/sql"
	"errors"
	"fmt"
	"sistema-toursseft/internal/entidades"
	"strings"
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
	tourProgramado := &entidades.TourProgramado{}
	query := `
		SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
			   tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
			   tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
			   tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
			   u.nombres || ' ' || u.apellidos as nombre_chofer,
			   ht.hora_inicio, ht.hora_fin
		FROM tour_programado tp
		INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
		INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
		INNER JOIN sede s ON tp.id_sede = s.id_sede
		INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
		LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
		WHERE tp.id_tour_programado = $1 AND tp.eliminado = false
	`

	var idChofer sql.NullInt64
	var notasExcepcion sql.NullString
	var nombreChofer sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
		&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
		&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
		&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
		&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
		&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
		&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tour programado no encontrado")
		}
		return nil, err
	}

	tourProgramado.IDChofer = idChofer
	tourProgramado.NotasExcepcion = notasExcepcion

	if nombreChofer.Valid {
		tourProgramado.NombreChofer = nombreChofer.String
	} else {
		tourProgramado.NombreChofer = "Sin asignar"
	}

	return tourProgramado, nil
}

// Create guarda un nuevo tour programado en la base de datos
func (r *TourProgramadoRepository) Create(tourProgramado *entidades.NuevoTourProgramadoRequest) (int, error) {
	// Verificar si ya existe un tour programado con la misma embarcación, fecha y horario
	var count int
	checkQuery := `
		SELECT COUNT(*) FROM tour_programado 
		WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 AND eliminado = false
	`

	fecha, err := time.Parse("2006-01-02", tourProgramado.Fecha)
	if err != nil {
		return 0, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	vigenciaDesde, err := time.Parse("2006-01-02", tourProgramado.VigenciaDesde)
	if err != nil {
		return 0, errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
	}

	vigenciaHasta, err := time.Parse("2006-01-02", tourProgramado.VigenciaHasta)
	if err != nil {
		return 0, errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
	}

	// Validar que la vigencia sea coherente
	if vigenciaHasta.Before(vigenciaDesde) {
		return 0, errors.New("la fecha de vigencia hasta debe ser posterior a la fecha de vigencia desde")
	}

	err = r.db.QueryRow(checkQuery, tourProgramado.IDEmbarcacion, fecha, tourProgramado.IDHorario).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("ya existe un tour programado con la misma embarcación, fecha y horario")
	}

	// Obtener la capacidad de la embarcación para establecer el cupo máximo si no se proporcionó
	if tourProgramado.CupoMaximo <= 0 {
		var capacidad int
		capacidadQuery := `SELECT capacidad FROM embarcacion WHERE id_embarcacion = $1 AND eliminado = false`
		err = r.db.QueryRow(capacidadQuery, tourProgramado.IDEmbarcacion).Scan(&capacidad)
		if err != nil {
			return 0, err
		}
		tourProgramado.CupoMaximo = capacidad
	}

	// Crear el tour programado
	var id int
	query := `
		INSERT INTO tour_programado (
			id_tipo_tour, id_embarcacion, id_horario, id_sede, id_chofer, 
			fecha, vigencia_desde, vigencia_hasta, cupo_maximo, cupo_disponible, 
			estado, eliminado, es_excepcion, notas_excepcion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false, $12, $13)
		RETURNING id_tour_programado
	`

	var idChofer interface{}
	if tourProgramado.IDChofer != nil {
		idChofer = *tourProgramado.IDChofer
	} else {
		idChofer = nil
	}

	var notasExcepcion interface{}
	if tourProgramado.NotasExcepcion != nil {
		notasExcepcion = *tourProgramado.NotasExcepcion
	} else {
		notasExcepcion = nil
	}

	err = r.db.QueryRow(
		query,
		tourProgramado.IDTipoTour,
		tourProgramado.IDEmbarcacion,
		tourProgramado.IDHorario,
		tourProgramado.IDSede,
		idChofer,
		fecha,
		vigenciaDesde,
		vigenciaHasta,
		tourProgramado.CupoMaximo,
		tourProgramado.CupoMaximo, // Al crear, el cupo disponible es igual al máximo
		"PROGRAMADO",              // Estado inicial
		tourProgramado.EsExcepcion,
		notasExcepcion,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un tour programado
func (r *TourProgramadoRepository) Update(id int, tourProgramado *entidades.ActualizarTourProgramadoRequest) error {
	// Verificar que el tour existe
	existeQuery := `SELECT COUNT(*) FROM tour_programado WHERE id_tour_programado = $1 AND eliminado = false`
	var count int
	err := r.db.QueryRow(existeQuery, id).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("tour programado no encontrado o fue eliminado")
	}

	// Si se cambia embarcación, fecha u horario, verificar que no colisione con otro tour
	if tourProgramado.IDEmbarcacion > 0 && tourProgramado.Fecha != "" && tourProgramado.IDHorario > 0 {
		fecha, err := time.Parse("2006-01-02", tourProgramado.Fecha)
		if err != nil {
			return errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
		}

		checkQuery := `
			SELECT COUNT(*) FROM tour_programado 
			WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 
			AND id_tour_programado != $4 AND eliminado = false
		`

		err = r.db.QueryRow(checkQuery, tourProgramado.IDEmbarcacion, fecha, tourProgramado.IDHorario, id).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("ya existe otro tour programado con la misma embarcación, fecha y horario")
		}
	}

	// Verificar y validar las fechas de vigencia si se proporcionan
	var vigenciaDesde, vigenciaHasta time.Time
	var errVigenciaDesde, errVigenciaHasta error

	if tourProgramado.VigenciaDesde != "" {
		vigenciaDesde, errVigenciaDesde = time.Parse("2006-01-02", tourProgramado.VigenciaDesde)
		if errVigenciaDesde != nil {
			return errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
		}
	}

	if tourProgramado.VigenciaHasta != "" {
		vigenciaHasta, errVigenciaHasta = time.Parse("2006-01-02", tourProgramado.VigenciaHasta)
		if errVigenciaHasta != nil {
			return errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
		}
	}

	// Si ambas fechas de vigencia se proporcionan, validar que sean coherentes
	if tourProgramado.VigenciaDesde != "" && tourProgramado.VigenciaHasta != "" {
		if vigenciaHasta.Before(vigenciaDesde) {
			return errors.New("la fecha de vigencia hasta debe ser posterior a la fecha de vigencia desde")
		}
	}

	// Construir la consulta dinámica según los campos proporcionados
	setClausulas := []string{}
	args := []interface{}{}
	argCount := 1

	if tourProgramado.IDTipoTour > 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("id_tipo_tour = $%d", argCount))
		args = append(args, tourProgramado.IDTipoTour)
		argCount++
	}

	if tourProgramado.IDEmbarcacion > 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("id_embarcacion = $%d", argCount))
		args = append(args, tourProgramado.IDEmbarcacion)
		argCount++
	}

	if tourProgramado.IDHorario > 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("id_horario = $%d", argCount))
		args = append(args, tourProgramado.IDHorario)
		argCount++
	}

	if tourProgramado.IDSede > 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("id_sede = $%d", argCount))
		args = append(args, tourProgramado.IDSede)
		argCount++
	}

	if tourProgramado.IDChofer != nil {
		setClausulas = append(setClausulas, fmt.Sprintf("id_chofer = $%d", argCount))
		args = append(args, *tourProgramado.IDChofer)
		argCount++
	} else {
		// Si explícitamente se quiere quitar el chofer
		setClausulas = append(setClausulas, "id_chofer = NULL")
	}

	if tourProgramado.Fecha != "" {
		fecha, err := time.Parse("2006-01-02", tourProgramado.Fecha)
		if err != nil {
			return errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
		}
		setClausulas = append(setClausulas, fmt.Sprintf("fecha = $%d", argCount))
		args = append(args, fecha)
		argCount++
	}

	if tourProgramado.VigenciaDesde != "" {
		setClausulas = append(setClausulas, fmt.Sprintf("vigencia_desde = $%d", argCount))
		args = append(args, vigenciaDesde)
		argCount++
	}

	if tourProgramado.VigenciaHasta != "" {
		setClausulas = append(setClausulas, fmt.Sprintf("vigencia_hasta = $%d", argCount))
		args = append(args, vigenciaHasta)
		argCount++
	}

	if tourProgramado.CupoMaximo > 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("cupo_maximo = $%d", argCount))
		args = append(args, tourProgramado.CupoMaximo)
		argCount++
	}

	if tourProgramado.CupoDisponible >= 0 {
		setClausulas = append(setClausulas, fmt.Sprintf("cupo_disponible = $%d", argCount))
		args = append(args, tourProgramado.CupoDisponible)
		argCount++
	}

	if tourProgramado.Estado != "" {
		setClausulas = append(setClausulas, fmt.Sprintf("estado = $%d", argCount))
		args = append(args, tourProgramado.Estado)
		argCount++
	}

	setClausulas = append(setClausulas, fmt.Sprintf("es_excepcion = $%d", argCount))
	args = append(args, tourProgramado.EsExcepcion)
	argCount++

	if tourProgramado.NotasExcepcion != nil {
		setClausulas = append(setClausulas, fmt.Sprintf("notas_excepcion = $%d", argCount))
		args = append(args, *tourProgramado.NotasExcepcion)
		argCount++
	} else {
		setClausulas = append(setClausulas, "notas_excepcion = NULL")
	}

	if len(setClausulas) == 0 {
		return errors.New("no se proporcionaron campos para actualizar")
	}

	query := fmt.Sprintf("UPDATE tour_programado SET %s WHERE id_tour_programado = $%d AND eliminado = false",
		strings.Join(setClausulas, ", "), argCount)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
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

// SoftDelete marca un tour programado como eliminado (borrado lógico)
func (r *TourProgramadoRepository) SoftDelete(id int) error {
	query := `UPDATE tour_programado SET eliminado = true WHERE id_tour_programado = $1 AND eliminado = false`
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

// AsignarChofer asigna un chofer a un tour programado
func (r *TourProgramadoRepository) AsignarChofer(idTour int, idChofer int) error {
	// Verificar que el chofer tenga rol de chofer
	var esChofer bool
	chkQuery := `SELECT COUNT(*) > 0 FROM usuario WHERE id_usuario = $1 AND rol = 'CHOFER' AND eliminado = false`
	err := r.db.QueryRow(chkQuery, idChofer).Scan(&esChofer)

	if err != nil {
		return err
	}

	if !esChofer {
		return errors.New("el usuario seleccionado no tiene rol de chofer")
	}

	// Asignar chofer
	query := `UPDATE tour_programado SET id_chofer = $1 WHERE id_tour_programado = $2 AND eliminado = false`
	result, err := r.db.Exec(query, idChofer, idTour)
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

// CambiarEstado cambia el estado de un tour programado
func (r *TourProgramadoRepository) CambiarEstado(id int, estado string) error {
	// Validar que el estado sea válido
	estadosValidos := []string{"PROGRAMADO", "EN_CURSO", "COMPLETADO", "CANCELADO"}
	estadoValido := false

	for _, e := range estadosValidos {
		if e == estado {
			estadoValido = true
			break
		}
	}

	if !estadoValido {
		return errors.New("estado no válido. Debe ser: PROGRAMADO, EN_CURSO, COMPLETADO o CANCELADO")
	}

	query := `UPDATE tour_programado SET estado = $1 WHERE id_tour_programado = $2 AND eliminado = false`
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

// List lista todos los tours programados con filtros opcionales
func (r *TourProgramadoRepository) List(filtros entidades.FiltrosTourProgramado) ([]*entidades.TourProgramado, error) {
	query := `
		SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
			   tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
			   tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
			   tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
			   u.nombres || ' ' || u.apellidos as nombre_chofer,
			   ht.hora_inicio, ht.hora_fin
		FROM tour_programado tp
		INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
		INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
		INNER JOIN sede s ON tp.id_sede = s.id_sede
		INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
		LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
		WHERE tp.eliminado = false
	`

	whereConditions := []string{}
	args := []interface{}{}
	argCount := 1

	if filtros.IDSede != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tp.id_sede = $%d", argCount))
		args = append(args, *filtros.IDSede)
		argCount++
	}

	if filtros.IDTipoTour != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tp.id_tipo_tour = $%d", argCount))
		args = append(args, *filtros.IDTipoTour)
		argCount++
	}

	if filtros.IDChofer != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tp.id_chofer = $%d", argCount))
		args = append(args, *filtros.IDChofer)
		argCount++
	}

	if filtros.IDEmbarcacion != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tp.id_embarcacion = $%d", argCount))
		args = append(args, *filtros.IDEmbarcacion)
		argCount++
	}

	if filtros.Estado != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tp.estado = $%d", argCount))
		args = append(args, *filtros.Estado)
		argCount++
	}

	if filtros.FechaInicio != nil {
		fechaInicio, err := time.Parse("2006-01-02", *filtros.FechaInicio)
		if err != nil {
			return nil, errors.New("formato de fecha inicio inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.fecha >= $%d", argCount))
		args = append(args, fechaInicio)
		argCount++
	}

	if filtros.FechaFin != nil {
		fechaFin, err := time.Parse("2006-01-02", *filtros.FechaFin)
		if err != nil {
			return nil, errors.New("formato de fecha fin inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.fecha <= $%d", argCount))
		args = append(args, fechaFin)
		argCount++
	}

	// Filtros para vigencia_desde
	if filtros.VigenciaDesdeIni != nil {
		vigenciaDesdeIni, err := time.Parse("2006-01-02", *filtros.VigenciaDesdeIni)
		if err != nil {
			return nil, errors.New("formato de vigencia desde inicio inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.vigencia_desde >= $%d", argCount))
		args = append(args, vigenciaDesdeIni)
		argCount++
	}

	if filtros.VigenciaDesdefin != nil {
		vigenciaDesdefin, err := time.Parse("2006-01-02", *filtros.VigenciaDesdefin)
		if err != nil {
			return nil, errors.New("formato de vigencia desde fin inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.vigencia_desde <= $%d", argCount))
		args = append(args, vigenciaDesdefin)
		argCount++
	}

	// Filtros para vigencia_hasta
	if filtros.VigenciaHastaIni != nil {
		vigenciaHastaIni, err := time.Parse("2006-01-02", *filtros.VigenciaHastaIni)
		if err != nil {
			return nil, errors.New("formato de vigencia hasta inicio inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.vigencia_hasta >= $%d", argCount))
		args = append(args, vigenciaHastaIni)
		argCount++
	}

	if filtros.VigenciaHastaFin != nil {
		vigenciaHastaFin, err := time.Parse("2006-01-02", *filtros.VigenciaHastaFin)
		if err != nil {
			return nil, errors.New("formato de vigencia hasta fin inválido, debe ser YYYY-MM-DD")
		}
		whereConditions = append(whereConditions, fmt.Sprintf("tp.vigencia_hasta <= $%d", argCount))
		args = append(args, vigenciaHastaFin)
		argCount++
	}

	if len(whereConditions) > 0 {
		query += " AND " + strings.Join(whereConditions, " AND ")
	}

	query += " ORDER BY tp.fecha, ht.hora_inicio"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toursProgramados := []*entidades.TourProgramado{}

	for rows.Next() {
		tourProgramado := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		var notasExcepcion sql.NullString
		var nombreChofer sql.NullString

		err := rows.Scan(
			&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
			&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
			&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
			&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
			&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
			&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
			&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
		)
		if err != nil {
			return nil, err
		}

		tourProgramado.IDChofer = idChofer
		tourProgramado.NotasExcepcion = notasExcepcion

		if nombreChofer.Valid {
			tourProgramado.NombreChofer = nombreChofer.String
		} else {
			tourProgramado.NombreChofer = "Sin asignar"
		}

		toursProgramados = append(toursProgramados, tourProgramado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toursProgramados, nil
}

// GetProgramacionSemanal obtiene los tours programados para una semana específica
func (r *TourProgramadoRepository) GetProgramacionSemanal(fechaInicio string, idSede int) ([]*entidades.TourProgramado, error) {
	inicio, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Calcular fecha fin (7 días después)
	fin := inicio.AddDate(0, 0, 6)

	query := `
		SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
			   tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
			   tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
			   tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
			   u.nombres || ' ' || u.apellidos as nombre_chofer,
			   ht.hora_inicio, ht.hora_fin
		FROM tour_programado tp
		INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
		INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
		INNER JOIN sede s ON tp.id_sede = s.id_sede
		INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
		LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
		WHERE tp.eliminado = false 
        AND tp.fecha BETWEEN $1 AND $2
        AND NOW() BETWEEN tp.vigencia_desde AND tp.vigencia_hasta
	`

	args := []interface{}{inicio, fin}
	argCount := 3

	if idSede > 0 {
		query += fmt.Sprintf(" AND tp.id_sede = $%d", argCount)
		args = append(args, idSede)
	}

	query += " ORDER BY tp.fecha, ht.hora_inicio"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toursProgramados := []*entidades.TourProgramado{}

	for rows.Next() {
		tourProgramado := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		var notasExcepcion sql.NullString
		var nombreChofer sql.NullString

		err := rows.Scan(
			&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
			&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
			&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
			&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
			&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
			&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
			&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
		)
		if err != nil {
			return nil, err
		}

		tourProgramado.IDChofer = idChofer
		tourProgramado.NotasExcepcion = notasExcepcion

		if nombreChofer.Valid {
			tourProgramado.NombreChofer = nombreChofer.String
		} else {
			tourProgramado.NombreChofer = "Sin asignar"
		}

		toursProgramados = append(toursProgramados, tourProgramado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toursProgramados, nil
}

// GetToursDisponiblesEnFecha obtiene tours disponibles para una fecha específica
func (r *TourProgramadoRepository) GetToursDisponiblesEnFecha(fecha string, idSede int) ([]*entidades.TourProgramado, error) {
	fechaBusqueda, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Obtener el día de la semana (0 = domingo, 1 = lunes, ..., 6 = sábado)
	diaSemana := fechaBusqueda.Weekday()

	// Construir la condición para el día de la semana específico
	var condicionDia string
	switch diaSemana {
	case 0: // Domingo
		condicionDia = "ht.disponible_domingo = true"
	case 1: // Lunes
		condicionDia = "ht.disponible_lunes = true"
	case 2: // Martes
		condicionDia = "ht.disponible_martes = true"
	case 3: // Miércoles
		condicionDia = "ht.disponible_miercoles = true"
	case 4: // Jueves
		condicionDia = "ht.disponible_jueves = true"
	case 5: // Viernes
		condicionDia = "ht.disponible_viernes = true"
	case 6: // Sábado
		condicionDia = "ht.disponible_sabado = true"
	}

	query := `
		SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
			   tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
			   tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
			   tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
			   u.nombres || ' ' || u.apellidos as nombre_chofer,
			   ht.hora_inicio, ht.hora_fin
		FROM tour_programado tp
		INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
		INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
		INNER JOIN sede s ON tp.id_sede = s.id_sede
		INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
		LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
		WHERE tp.eliminado = false
		AND $1 BETWEEN tp.vigencia_desde AND tp.vigencia_hasta
		AND tp.estado = 'PROGRAMADO'
		AND tp.cupo_disponible > 0
		AND ` + condicionDia

	args := []interface{}{fechaBusqueda}
	argCount := 2

	if idSede > 0 {
		query += fmt.Sprintf(" AND tp.id_sede = $%d", argCount)
		args = append(args, idSede)
		argCount++
	}

	query += " ORDER BY ht.hora_inicio"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toursProgramados := []*entidades.TourProgramado{}

	for rows.Next() {
		tourProgramado := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		var notasExcepcion sql.NullString
		var nombreChofer sql.NullString

		err := rows.Scan(
			&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
			&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
			&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
			&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
			&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
			&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
			&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
		)
		if err != nil {
			return nil, err
		}

		tourProgramado.IDChofer = idChofer
		tourProgramado.NotasExcepcion = notasExcepcion

		if nombreChofer.Valid {
			tourProgramado.NombreChofer = nombreChofer.String
		} else {
			tourProgramado.NombreChofer = "Sin asignar"
		}

		toursProgramados = append(toursProgramados, tourProgramado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toursProgramados, nil
}

// GetToursDisponiblesEnRangoFechas obtiene los tours disponibles para reserva en un rango de fechas
func (r *TourProgramadoRepository) GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin string, idSede int) ([]*entidades.TourProgramado, error) {
	inicio, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		return nil, errors.New("formato de fecha inicio inválido, debe ser YYYY-MM-DD")
	}

	fin, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		return nil, errors.New("formato de fecha fin inválido, debe ser YYYY-MM-DD")
	}

	if fin.Before(inicio) {
		return nil, errors.New("la fecha fin no puede ser anterior a la fecha inicio")
	}

	// Fecha actual para validar que solo se muestren tours futuros
	fechaActual := time.Now().Truncate(24 * time.Hour)
	if inicio.Before(fechaActual) {
		inicio = fechaActual
	}

	query := `
		WITH fechas AS (
			SELECT generate_series($1::date, $2::date, '1 day'::interval)::date AS fecha
		)
		SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
			tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
			tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
			tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
			u.nombres || ' ' || u.apellidos as nombre_chofer,
			ht.hora_inicio, ht.hora_fin
		FROM tour_programado tp
		INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
		INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
		INNER JOIN sede s ON tp.id_sede = s.id_sede
		INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
		LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
		INNER JOIN fechas f ON 
			f.fecha BETWEEN tp.vigencia_desde AND tp.vigencia_hasta AND
			(
				(EXTRACT(DOW FROM f.fecha) = 0 AND ht.disponible_domingo) OR 
				(EXTRACT(DOW FROM f.fecha) = 1 AND ht.disponible_lunes) OR
				(EXTRACT(DOW FROM f.fecha) = 2 AND ht.disponible_martes) OR
				(EXTRACT(DOW FROM f.fecha) = 3 AND ht.disponible_miercoles) OR
				(EXTRACT(DOW FROM f.fecha) = 4 AND ht.disponible_jueves) OR
				(EXTRACT(DOW FROM f.fecha) = 5 AND ht.disponible_viernes) OR
				(EXTRACT(DOW FROM f.fecha) = 6 AND ht.disponible_sabado)
			)
		WHERE tp.eliminado = false
		AND tp.estado = 'PROGRAMADO'
		AND tp.cupo_disponible > 0
	`

	args := []interface{}{inicio, fin}
	argCount := 3

	if idSede > 0 {
		query += fmt.Sprintf(" AND tp.id_sede = $%d", argCount)
		args = append(args, idSede)
		argCount++
	}

	query += " ORDER BY f.fecha, ht.hora_inicio"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toursProgramados := []*entidades.TourProgramado{}

	for rows.Next() {
		tourProgramado := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		var notasExcepcion sql.NullString
		var nombreChofer sql.NullString

		err := rows.Scan(
			&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
			&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
			&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
			&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
			&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
			&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
			&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
		)
		if err != nil {
			return nil, err
		}

		tourProgramado.IDChofer = idChofer
		tourProgramado.NotasExcepcion = notasExcepcion

		if nombreChofer.Valid {
			tourProgramado.NombreChofer = nombreChofer.String
		} else {
			tourProgramado.NombreChofer = "Sin asignar"
		}

		toursProgramados = append(toursProgramados, tourProgramado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toursProgramados, nil
}

// VerificarDisponibilidadHorario verifica si un horario está disponible para una fecha específica
func (r *TourProgramadoRepository) VerificarDisponibilidadHorario(idHorario int, fecha string) (bool, error) {
	// 1. Verificar si el horario existe y si está disponible para el día de la semana
	fechaObj, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return false, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Obtener el horario
	var disponibilidadDia bool
	diaSemana := int(fechaObj.Weekday()) // 0 = domingo, 1 = lunes, ..., 6 = sábado

	// Construir consulta para verificar disponibilidad del día
	var diaQuery string
	switch diaSemana {
	case 0:
		diaQuery = "disponible_domingo"
	case 1:
		diaQuery = "disponible_lunes"
	case 2:
		diaQuery = "disponible_martes"
	case 3:
		diaQuery = "disponible_miercoles"
	case 4:
		diaQuery = "disponible_jueves"
	case 5:
		diaQuery = "disponible_viernes"
	case 6:
		diaQuery = "disponible_sabado"
	}

	// Verificar si el horario existe y si está disponible para el día
	horarioQuery := fmt.Sprintf(`
		SELECT %s FROM horario_tour
		WHERE id_horario = $1 AND eliminado = false
	`, diaQuery)

	err = r.db.QueryRow(horarioQuery, idHorario).Scan(&disponibilidadDia)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("horario no encontrado")
		}
		return false, err
	}

	if !disponibilidadDia {
		// El día de la semana no está disponible en este horario
		return false, nil
	}

	// 2. Verificar si ya hay tours programados con este horario en esta fecha
	var count int
	conflictoQuery := `
		SELECT COUNT(*) FROM tour_programado
		WHERE id_horario = $1 
		AND fecha = $2
		AND estado != 'CANCELADO'
		AND eliminado = false
	`

	err = r.db.QueryRow(conflictoQuery, idHorario, fechaObj).Scan(&count)
	if err != nil {
		return false, err
	}

	// Si hay tours programados, el horario no está disponible
	return count == 0, nil
}

// ProgramarToursSemanal crea múltiples tours programados en un rango de fechas
func (r *TourProgramadoRepository) ProgramarToursSemanal(tourBase *entidades.NuevoTourProgramadoRequest, fechas []time.Time) ([]int, error) {
	// Array para almacenar los IDs de los tours creados
	tourIDs := []int{}

	// Iniciar una transacción para garantizar la atomicidad
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic después de rollback
		}
	}()

	// Query para insertar tours
	query := `
		INSERT INTO tour_programado (
			id_tipo_tour, id_embarcacion, id_horario, id_sede, id_chofer, 
			fecha, vigencia_desde, vigencia_hasta, cupo_maximo, cupo_disponible, 
			estado, eliminado, es_excepcion, notas_excepcion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false, $12, $13)
		RETURNING id_tour_programado
	`

	// Establecer valores básicos para todos los tours
	var idChofer interface{}
	if tourBase.IDChofer != nil {
		idChofer = *tourBase.IDChofer
	} else {
		idChofer = nil
	}

	var notasExcepcion interface{}
	if tourBase.NotasExcepcion != nil {
		notasExcepcion = *tourBase.NotasExcepcion
	} else {
		notasExcepcion = nil
	}

	// Parsear las fechas de vigencia una sola vez
	vigenciaDesde, err := time.Parse("2006-01-02", tourBase.VigenciaDesde)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
	}

	vigenciaHasta, err := time.Parse("2006-01-02", tourBase.VigenciaHasta)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
	}

	// Crear un tour para cada fecha
	for _, fecha := range fechas {
		// Verificar si ya existe un tour para esta embarcación, fecha y horario
		var count int
		checkQuery := `
			SELECT COUNT(*) FROM tour_programado 
			WHERE id_embarcacion = $1 AND fecha = $2 AND id_horario = $3 AND eliminado = false
		`

		err = tx.QueryRow(checkQuery, tourBase.IDEmbarcacion, fecha, tourBase.IDHorario).Scan(&count)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if count > 0 {
			// Ya existe un tour para esta fecha, saltar
			continue
		}

		// Insertar el nuevo tour
		var id int
		err = tx.QueryRow(
			query,
			tourBase.IDTipoTour,
			tourBase.IDEmbarcacion,
			tourBase.IDHorario,
			tourBase.IDSede,
			idChofer,
			fecha,
			vigenciaDesde,
			vigenciaHasta,
			tourBase.CupoMaximo,
			tourBase.CupoMaximo, // Al crear, el cupo disponible es igual al máximo
			"PROGRAMADO",        // Estado inicial
			tourBase.EsExcepcion,
			notasExcepcion,
		).Scan(&id)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tourIDs = append(tourIDs, id)
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return tourIDs, nil
}

// En repositorios/TourProgramadoRepository.go

// GetToursDisponibles obtiene tours disponibles para reserva
func (r *TourProgramadoRepository) GetToursDisponibles() ([]*entidades.TourProgramado, error) {
	// Fecha actual para validar que solo se muestren tours vigentes
	fechaActual := time.Now().Format("2006-01-02")

	query := `
        SELECT DISTINCT ON (tp.id_tour_programado) 
            tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, tp.id_sede, 
            tp.id_chofer, tp.fecha, tp.vigencia_desde, tp.vigencia_hasta, tp.cupo_maximo, tp.cupo_disponible, 
            tp.estado, tp.eliminado, tp.es_excepcion, tp.notas_excepcion,
            tt.nombre as nombre_tipo_tour, e.nombre as nombre_embarcacion, s.nombre as nombre_sede,
            u.nombres || ' ' || u.apellidos as nombre_chofer,
            ht.hora_inicio, ht.hora_fin
        FROM tour_programado tp
        INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
        INNER JOIN embarcacion e ON tp.id_embarcacion = e.id_embarcacion
        INNER JOIN sede s ON tp.id_sede = s.id_sede
        INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
        LEFT JOIN usuario u ON tp.id_chofer = u.id_usuario
        WHERE tp.eliminado = false
        AND tp.estado = 'PROGRAMADO'
        AND tp.cupo_disponible > 0
        AND $1::date BETWEEN tp.vigencia_desde AND tp.vigencia_hasta
        ORDER BY tp.id_tour_programado, ht.hora_inicio
    `

	rows, err := r.db.Query(query, fechaActual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toursProgramados := []*entidades.TourProgramado{}

	for rows.Next() {
		tourProgramado := &entidades.TourProgramado{}
		var idChofer sql.NullInt64
		var notasExcepcion sql.NullString
		var nombreChofer sql.NullString

		err := rows.Scan(
			&tourProgramado.ID, &tourProgramado.IDTipoTour, &tourProgramado.IDEmbarcacion,
			&tourProgramado.IDHorario, &tourProgramado.IDSede, &idChofer,
			&tourProgramado.Fecha, &tourProgramado.VigenciaDesde, &tourProgramado.VigenciaHasta,
			&tourProgramado.CupoMaximo, &tourProgramado.CupoDisponible,
			&tourProgramado.Estado, &tourProgramado.Eliminado, &tourProgramado.EsExcepcion, &notasExcepcion,
			&tourProgramado.NombreTipoTour, &tourProgramado.NombreEmbarcacion, &tourProgramado.NombreSede,
			&nombreChofer, &tourProgramado.HoraInicio, &tourProgramado.HoraFin,
		)
		if err != nil {
			return nil, err
		}

		tourProgramado.IDChofer = idChofer
		tourProgramado.NotasExcepcion = notasExcepcion

		if nombreChofer.Valid {
			tourProgramado.NombreChofer = nombreChofer.String
		} else {
			tourProgramado.NombreChofer = "Sin asignar"
		}

		toursProgramados = append(toursProgramados, tourProgramado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toursProgramados, nil
}
