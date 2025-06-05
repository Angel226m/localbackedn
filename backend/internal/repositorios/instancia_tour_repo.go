package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
	"strings"
	"time"
)

// InstanciaTourRepository maneja las operaciones de base de datos para instancias de tour
type InstanciaTourRepository struct {
	db *sql.DB
}

// NewInstanciaTourRepository crea una nueva instancia del repositorio
func NewInstanciaTourRepository(db *sql.DB) *InstanciaTourRepository {
	return &InstanciaTourRepository{
		db: db,
	}
}

// GetByID obtiene una instancia de tour por su ID
func (r *InstanciaTourRepository) GetByID(id int) (*entidades.InstanciaTour, error) {
	instancia := &entidades.InstanciaTour{}
	query := `SELECT i.id_instancia, i.id_tour_programado, i.fecha_especifica, i.hora_inicio, i.hora_fin,
              i.id_chofer, i.id_embarcacion, i.cupo_disponible, i.estado, i.eliminado,
              t.nombre, e.nombre, s.nombre, 
              COALESCE(u.nombres || ' ' || u.apellidos, 'Sin asignar')
              FROM instancia_tour i
              INNER JOIN tour_programado tp ON i.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour t ON tp.id_tipo_tour = t.id_tipo_tour
              INNER JOIN embarcacion e ON i.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              LEFT JOIN usuario u ON i.id_chofer = u.id_usuario
              WHERE i.id_instancia = $1 AND i.eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&instancia.ID, &instancia.IDTourProgramado, &instancia.FechaEspecifica, &instancia.HoraInicio, &instancia.HoraFin,
		&instancia.IDChofer, &instancia.IDEmbarcacion, &instancia.CupoDisponible, &instancia.Estado, &instancia.Eliminado,
		&instancia.NombreTipoTour, &instancia.NombreEmbarcacion, &instancia.NombreSede, &instancia.NombreChofer,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("instancia de tour no encontrada")
		}
		return nil, err
	}

	// Formatear las fechas y horas para presentación
	instancia.HoraInicioStr = instancia.HoraInicio.Format("15:04")
	instancia.HoraFinStr = instancia.HoraFin.Format("15:04")
	instancia.FechaEspecificaStr = instancia.FechaEspecifica.Format("2006-01-02")

	return instancia, nil
}

// Create guarda una nueva instancia de tour en la base de datos
func (r *InstanciaTourRepository) Create(instancia *entidades.NuevaInstanciaTourRequest) (int, error) {
	// Verificar que el tour programado existe
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tour_programado WHERE id_tour_programado = $1 AND eliminado = false)",
		instancia.IDTourProgramado).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("el tour programado especificado no existe")
	}

	// Verificar que la embarcación existe
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM embarcacion WHERE id_embarcacion = $1 AND eliminado = false AND estado = 'DISPONIBLE')",
		instancia.IDEmbarcacion).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("la embarcación especificada no existe o no está disponible")
	}

	// Verificar que el chofer existe si se proporciona
	if instancia.IDChofer != nil {
		err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuario WHERE id_usuario = $1 AND rol = 'CHOFER' AND eliminado = false)",
			*instancia.IDChofer).Scan(&exists)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, errors.New("el chofer especificado no existe o no tiene rol de chofer")
		}
	}

	// Parsear fecha y horas
	fechaEspecifica, err := time.Parse("2006-01-02", instancia.FechaEspecifica)
	if err != nil {
		return 0, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	horaInicio, err := time.Parse("15:04", instancia.HoraInicio)
	if err != nil {
		return 0, errors.New("formato de hora de inicio inválido, debe ser HH:MM")
	}

	horaFin, err := time.Parse("15:04", instancia.HoraFin)
	if err != nil {
		return 0, errors.New("formato de hora de fin inválido, debe ser HH:MM")
	}

	// Verificar que la hora de fin es posterior a la de inicio
	if !horaFin.After(horaInicio) {
		return 0, errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Si se proporciona un chofer, verificar disponibilidad
	if instancia.IDChofer != nil {
		// Verificar si el chofer está disponible en esa fecha y horario según horario_chofer
		var disponible bool
		diaSemana := int(fechaEspecifica.Weekday())
		if diaSemana == 0 { // En Go, Sunday es 0, pero en nuestra DB es 7
			diaSemana = 7 // Domingo
		}

		var condicionDia string
		switch diaSemana {
		case 1:
			condicionDia = "disponible_lunes"
		case 2:
			condicionDia = "disponible_martes"
		case 3:
			condicionDia = "disponible_miercoles"
		case 4:
			condicionDia = "disponible_jueves"
		case 5:
			condicionDia = "disponible_viernes"
		case 6:
			condicionDia = "disponible_sabado"
		case 7:
			condicionDia = "disponible_domingo"
		}

		queryDisponibilidad := `
			SELECT EXISTS (
				SELECT 1 FROM horario_chofer 
				WHERE id_usuario = $1 
				AND ` + condicionDia + ` = true 
				AND hora_inicio <= $2::time 
				AND hora_fin >= $3::time
				AND fecha_inicio <= $4 
				AND (fecha_fin IS NULL OR fecha_fin >= $4)
				AND eliminado = false
			)`

		err = tx.QueryRow(queryDisponibilidad, *instancia.IDChofer, instancia.HoraInicio, instancia.HoraFin,
			instancia.FechaEspecifica).Scan(&disponible)
		if err != nil {
			return 0, err
		}
		if !disponible {
			return 0, errors.New("el chofer no está disponible en la fecha y horario especificados")
		}

		// Verificar que el chofer no esté asignado a otro tour en el mismo horario
		queryOcupado := `
			SELECT EXISTS (
				SELECT 1 FROM instancia_tour 
				WHERE id_chofer = $1 
				AND fecha_especifica = $2 
				AND (
					(hora_inicio <= $3::time AND hora_fin > $3::time) OR 
					(hora_inicio < $4::time AND hora_fin >= $4::time) OR 
					(hora_inicio >= $3::time AND hora_fin <= $4::time)
				)
				AND estado IN ('PROGRAMADO', 'EN_CURSO')
				AND eliminado = false
			)`

		var ocupado bool
		err = tx.QueryRow(queryOcupado, *instancia.IDChofer, instancia.FechaEspecifica,
			instancia.HoraInicio, instancia.HoraFin).Scan(&ocupado)
		if err != nil {
			return 0, err
		}
		if ocupado {
			return 0, errors.New("el chofer ya está asignado a otro tour en el mismo horario")
		}
	}

	// Verificar que la embarcación no esté asignada a otro tour en el mismo horario
	queryEmbarcacionOcupada := `
		SELECT EXISTS (
			SELECT 1 FROM instancia_tour 
			WHERE id_embarcacion = $1 
			AND fecha_especifica = $2 
			AND (
				(hora_inicio <= $3::time AND hora_fin > $3::time) OR 
				(hora_inicio < $4::time AND hora_fin >= $4::time) OR 
				(hora_inicio >= $3::time AND hora_fin <= $4::time)
			)
			AND estado IN ('PROGRAMADO', 'EN_CURSO')
			AND eliminado = false
		)`

	var embarcacionOcupada bool
	err = tx.QueryRow(queryEmbarcacionOcupada, instancia.IDEmbarcacion, instancia.FechaEspecifica,
		instancia.HoraInicio, instancia.HoraFin).Scan(&embarcacionOcupada)
	if err != nil {
		return 0, err
	}
	if embarcacionOcupada {
		return 0, errors.New("la embarcación ya está asignada a otro tour en el mismo horario")
	}

	// Insertar la instancia de tour
	var id int
	query := `INSERT INTO instancia_tour (id_tour_programado, fecha_especifica, hora_inicio, hora_fin, 
              id_chofer, id_embarcacion, cupo_disponible, estado, eliminado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, 'PROGRAMADO', false) 
              RETURNING id_instancia`

	var idChoferParam interface{}
	if instancia.IDChofer != nil {
		idChoferParam = *instancia.IDChofer
	} else {
		idChoferParam = nil
	}

	err = tx.QueryRow(
		query,
		instancia.IDTourProgramado,
		fechaEspecifica,
		horaInicio,
		horaFin,
		idChoferParam,
		instancia.IDEmbarcacion,
		instancia.CupoDisponible,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Actualizar estado de la embarcación si es necesario
	_, err = tx.Exec("UPDATE embarcacion SET estado = 'OCUPADA' WHERE id_embarcacion = $1", instancia.IDEmbarcacion)
	if err != nil {
		return 0, err
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de una instancia de tour
func (r *InstanciaTourRepository) Update(id int, instancia *entidades.ActualizarInstanciaTourRequest) error {
	// Verificar que la instancia existe
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM instancia_tour WHERE id_instancia = $1 AND eliminado = false)",
		id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("la instancia de tour especificada no existe")
	}

	// Obtener la instancia actual para comparaciones
	var instanciaActual entidades.InstanciaTour
	err = r.db.QueryRow(`
		SELECT id_tour_programado, fecha_especifica, hora_inicio, hora_fin, id_chofer, id_embarcacion, cupo_disponible, estado
		FROM instancia_tour
		WHERE id_instancia = $1 AND eliminado = false`, id).Scan(
		&instanciaActual.IDTourProgramado, &instanciaActual.FechaEspecifica, &instanciaActual.HoraInicio,
		&instanciaActual.HoraFin, &instanciaActual.IDChofer, &instanciaActual.IDEmbarcacion,
		&instanciaActual.CupoDisponible, &instanciaActual.Estado)
	if err != nil {
		return err
	}

	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Construir la consulta de actualización
	queryParts := []string{}
	queryParams := []interface{}{}
	paramCount := 1

	// Función auxiliar para agregar parámetros a la consulta
	addParam := func(column string, value interface{}) {
		queryParts = append(queryParts, column+" = $"+string(48+paramCount)) // 48 es el código ASCII de '0'
		queryParams = append(queryParams, value)
		paramCount++
	}

	// Actualizar tour programado si se proporciona
	if instancia.IDTourProgramado != nil {
		// Verificar que el tour programado existe
		err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tour_programado WHERE id_tour_programado = $1 AND eliminado = false)",
			*instancia.IDTourProgramado).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("el tour programado especificado no existe")
		}
		addParam("id_tour_programado", *instancia.IDTourProgramado)
	}

	// Actualizar fecha específica si se proporciona
	var fechaEspecifica time.Time
	if instancia.FechaEspecifica != nil {
		fechaEspecifica, err = time.Parse("2006-01-02", *instancia.FechaEspecifica)
		if err != nil {
			return errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
		}
		addParam("fecha_especifica", fechaEspecifica)
	} else {
		fechaEspecifica = instanciaActual.FechaEspecifica
	}

	// Actualizar hora de inicio si se proporciona
	var horaInicio time.Time
	if instancia.HoraInicio != nil {
		horaInicio, err = time.Parse("15:04", *instancia.HoraInicio)
		if err != nil {
			return errors.New("formato de hora de inicio inválido, debe ser HH:MM")
		}
		addParam("hora_inicio", horaInicio)
	} else {
		horaInicio = instanciaActual.HoraInicio
	}

	// Actualizar hora de fin si se proporciona
	var horaFin time.Time
	if instancia.HoraFin != nil {
		horaFin, err = time.Parse("15:04", *instancia.HoraFin)
		if err != nil {
			return errors.New("formato de hora de fin inválido, debe ser HH:MM")
		}
		addParam("hora_fin", horaFin)
	} else {
		horaFin = instanciaActual.HoraFin
	}

	// Verificar que la hora de fin es posterior a la de inicio
	if !horaFin.After(horaInicio) {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// Actualizar chofer si se proporciona
	if instancia.IDChofer != nil {
		// Si es un chofer diferente, verificar disponibilidad
		if instanciaActual.IDChofer.Valid && int(instanciaActual.IDChofer.Int64) != *instancia.IDChofer {
			// Verificar que el chofer existe
			err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuario WHERE id_usuario = $1 AND rol = 'CHOFER' AND eliminado = false)",
				*instancia.IDChofer).Scan(&exists)
			if err != nil {
				return err
			}
			if !exists {
				return errors.New("el chofer especificado no existe o no tiene rol de chofer")
			}

			// Verificar disponibilidad del chofer
			diaSemana := int(fechaEspecifica.Weekday())
			if diaSemana == 0 { // En Go, Sunday es 0, pero en nuestra DB es 7
				diaSemana = 7 // Domingo
			}

			var condicionDia string
			switch diaSemana {
			case 1:
				condicionDia = "disponible_lunes"
			case 2:
				condicionDia = "disponible_martes"
			case 3:
				condicionDia = "disponible_miercoles"
			case 4:
				condicionDia = "disponible_jueves"
			case 5:
				condicionDia = "disponible_viernes"
			case 6:
				condicionDia = "disponible_sabado"
			case 7:
				condicionDia = "disponible_domingo"
			}

			queryDisponibilidad := `
				SELECT EXISTS (
					SELECT 1 FROM horario_chofer 
					WHERE id_usuario = $1 
					AND ` + condicionDia + ` = true 
					AND hora_inicio <= $2::time 
					AND hora_fin >= $3::time
					AND fecha_inicio <= $4 
					AND (fecha_fin IS NULL OR fecha_fin >= $4)
					AND eliminado = false
				)`

			var disponible bool
			err = tx.QueryRow(queryDisponibilidad, *instancia.IDChofer, horaInicio, horaFin,
				fechaEspecifica).Scan(&disponible)
			if err != nil {
				return err
			}
			if !disponible {
				return errors.New("el chofer no está disponible en la fecha y horario especificados")
			}

			// Verificar que el chofer no esté asignado a otro tour en el mismo horario
			queryOcupado := `
				SELECT EXISTS (
					SELECT 1 FROM instancia_tour 
					WHERE id_chofer = $1 
					AND fecha_especifica = $2 
					AND (
						(hora_inicio <= $3::time AND hora_fin > $3::time) OR 
						(hora_inicio < $4::time AND hora_fin >= $4::time) OR 
						(hora_inicio >= $3::time AND hora_fin <= $4::time)
					)
					AND id_instancia != $5
					AND estado IN ('PROGRAMADO', 'EN_CURSO')
					AND eliminado = false
				)`

			var ocupado bool
			err = tx.QueryRow(queryOcupado, *instancia.IDChofer, fechaEspecifica,
				horaInicio, horaFin, id).Scan(&ocupado)
			if err != nil {
				return err
			}
			if ocupado {
				return errors.New("el chofer ya está asignado a otro tour en el mismo horario")
			}
		}
		addParam("id_chofer", *instancia.IDChofer)
	}

	// Actualizar embarcación si se proporciona
	if instancia.IDEmbarcacion != nil && *instancia.IDEmbarcacion != instanciaActual.IDEmbarcacion {
		// Verificar que la embarcación existe
		err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM embarcacion WHERE id_embarcacion = $1 AND eliminado = false AND estado IN ('DISPONIBLE', 'OCUPADA'))",
			*instancia.IDEmbarcacion).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("la embarcación especificada no existe o no está disponible")
		}

		// Verificar que la embarcación no esté asignada a otro tour en el mismo horario
		queryEmbarcacionOcupada := `
			SELECT EXISTS (
				SELECT 1 FROM instancia_tour 
				WHERE id_embarcacion = $1 
				AND fecha_especifica = $2 
				AND (
					(hora_inicio <= $3::time AND hora_fin > $3::time) OR 
					(hora_inicio < $4::time AND hora_fin >= $4::time) OR 
					(hora_inicio >= $3::time AND hora_fin <= $4::time)
				)
				AND id_instancia != $5
				AND estado IN ('PROGRAMADO', 'EN_CURSO')
				AND eliminado = false
			)`

		var embarcacionOcupada bool
		err = tx.QueryRow(queryEmbarcacionOcupada, *instancia.IDEmbarcacion, fechaEspecifica,
			horaInicio, horaFin, id).Scan(&embarcacionOcupada)
		if err != nil {
			return err
		}
		if embarcacionOcupada {
			return errors.New("la embarcación ya está asignada a otro tour en el mismo horario")
		}

		// Liberar la embarcación anterior
		_, err = tx.Exec("UPDATE embarcacion SET estado = 'DISPONIBLE' WHERE id_embarcacion = $1",
			instanciaActual.IDEmbarcacion)
		if err != nil {
			return err
		}

		// Ocupar la nueva embarcación
		_, err = tx.Exec("UPDATE embarcacion SET estado = 'OCUPADA' WHERE id_embarcacion = $1",
			*instancia.IDEmbarcacion)
		if err != nil {
			return err
		}

		addParam("id_embarcacion", *instancia.IDEmbarcacion)
	}

	// Actualizar cupo disponible si se proporciona
	if instancia.CupoDisponible != nil {
		addParam("cupo_disponible", *instancia.CupoDisponible)
	}

	// Actualizar estado si se proporciona
	if instancia.Estado != nil {
		if *instancia.Estado == "COMPLETADO" || *instancia.Estado == "CANCELADO" {
			// Liberar la embarcación si el tour se completa o cancela
			_, err = tx.Exec("UPDATE embarcacion SET estado = 'DISPONIBLE' WHERE id_embarcacion = $1",
				instanciaActual.IDEmbarcacion)
			if err != nil {
				return err
			}
		}
		addParam("estado", *instancia.Estado)
	}

	// Si no hay nada que actualizar, retornar
	if len(queryParts) == 0 {
		return nil
	}

	// Construir y ejecutar la consulta de actualización
	query := "UPDATE instancia_tour SET " + strings.Join(queryParts, ", ") + " WHERE id_instancia = $" + string(48+paramCount)
	queryParams = append(queryParams, id)

	_, err = tx.Exec(query, queryParams...)
	if err != nil {
		return err
	}

	// Confirmar transacción
	return tx.Commit()
}

// Delete marca una instancia de tour como eliminada (borrado lógico)
func (r *InstanciaTourRepository) Delete(id int) error {
	// Verificar si hay reservas asociadas a esta instancia
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_instancia = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar esta instancia de tour porque tiene reservas asociadas")
	}

	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Obtener ID de la embarcación para liberarla
	var idEmbarcacion int
	queryEmbarcacion := `SELECT id_embarcacion FROM instancia_tour WHERE id_instancia = $1 AND eliminado = false`
	err = tx.QueryRow(queryEmbarcacion, id).Scan(&idEmbarcacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("instancia de tour no encontrada")
		}
		return err
	}

	// Marcar como eliminada la instancia
	query := `UPDATE instancia_tour SET eliminado = true WHERE id_instancia = $1`
	result, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("instancia de tour no encontrada o ya eliminada")
	}

	// Liberar la embarcación
	_, err = tx.Exec("UPDATE embarcacion SET estado = 'DISPONIBLE' WHERE id_embarcacion = $1", idEmbarcacion)
	if err != nil {
		return err
	}

	// Confirmar transacción
	return tx.Commit()
}

// List lista todas las instancias de tour no eliminadas
func (r *InstanciaTourRepository) List() ([]*entidades.InstanciaTour, error) {
	query := `SELECT i.id_instancia, i.id_tour_programado, i.fecha_especifica, i.hora_inicio, i.hora_fin,
              i.id_chofer, i.id_embarcacion, i.cupo_disponible, i.estado, i.eliminado,
              t.nombre, e.nombre, s.nombre, 
              COALESCE(u.nombres || ' ' || u.apellidos, 'Sin asignar')
              FROM instancia_tour i
              INNER JOIN tour_programado tp ON i.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour t ON tp.id_tipo_tour = t.id_tipo_tour
              INNER JOIN embarcacion e ON i.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              LEFT JOIN usuario u ON i.id_chofer = u.id_usuario
              WHERE i.eliminado = false
              ORDER BY i.fecha_especifica, i.hora_inicio`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instancias := []*entidades.InstanciaTour{}

	for rows.Next() {
		instancia := &entidades.InstanciaTour{}
		err := rows.Scan(
			&instancia.ID, &instancia.IDTourProgramado, &instancia.FechaEspecifica, &instancia.HoraInicio, &instancia.HoraFin,
			&instancia.IDChofer, &instancia.IDEmbarcacion, &instancia.CupoDisponible, &instancia.Estado, &instancia.Eliminado,
			&instancia.NombreTipoTour, &instancia.NombreEmbarcacion, &instancia.NombreSede, &instancia.NombreChofer,
		)
		if err != nil {
			return nil, err
		}

		// Formatear las fechas y horas para presentación
		instancia.HoraInicioStr = instancia.HoraInicio.Format("15:04")
		instancia.HoraFinStr = instancia.HoraFin.Format("15:04")
		instancia.FechaEspecificaStr = instancia.FechaEspecifica.Format("2006-01-02")

		instancias = append(instancias, instancia)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instancias, nil
}

// ListByTourProgramado lista todas las instancias de un tour programado específico
func (r *InstanciaTourRepository) ListByTourProgramado(idTourProgramado int) ([]*entidades.InstanciaTour, error) {
	query := `SELECT i.id_instancia, i.id_tour_programado, i.fecha_especifica, i.hora_inicio, i.hora_fin,
              i.id_chofer, i.id_embarcacion, i.cupo_disponible, i.estado, i.eliminado,
              t.nombre, e.nombre, s.nombre, 
              COALESCE(u.nombres || ' ' || u.apellidos, 'Sin asignar')
              FROM instancia_tour i
              INNER JOIN tour_programado tp ON i.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour t ON tp.id_tipo_tour = t.id_tipo_tour
              INNER JOIN embarcacion e ON i.id_embarcacion = e.id_embarcacion
              INNER JOIN sede s ON tp.id_sede = s.id_sede
              LEFT JOIN usuario u ON i.id_chofer = u.id_usuario
              WHERE i.id_tour_programado = $1 AND i.eliminado = false
              ORDER BY i.fecha_especifica, i.hora_inicio`

	rows, err := r.db.Query(query, idTourProgramado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instancias := []*entidades.InstanciaTour{}

	for rows.Next() {
		instancia := &entidades.InstanciaTour{}
		err := rows.Scan(
			&instancia.ID, &instancia.IDTourProgramado, &instancia.FechaEspecifica, &instancia.HoraInicio, &instancia.HoraFin,
			&instancia.IDChofer, &instancia.IDEmbarcacion, &instancia.CupoDisponible, &instancia.Estado, &instancia.Eliminado,
			&instancia.NombreTipoTour, &instancia.NombreEmbarcacion, &instancia.NombreSede, &instancia.NombreChofer,
		)
		if err != nil {
			return nil, err
		}

		// Formatear las fechas y horas para presentación
		instancia.HoraInicioStr = instancia.HoraInicio.Format("15:04")
		instancia.HoraFinStr = instancia.HoraFin.Format("15:04")
		instancia.FechaEspecificaStr = instancia.FechaEspecifica.Format("2006-01-02")

		instancias = append(instancias, instancia)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instancias, nil
}

// ListByFiltros lista instancias de tour según filtros específicos
func (r *InstanciaTourRepository) ListByFiltros(filtros entidades.FiltrosInstanciaTour) ([]*entidades.InstanciaTour, error) {
	// Construir la consulta base
	queryBase := `SELECT i.id_instancia, i.id_tour_programado, i.fecha_especifica, i.hora_inicio, i.hora_fin,
                  i.id_chofer, i.id_embarcacion, i.cupo_disponible, i.estado, i.eliminado,
                  t.nombre, e.nombre, s.nombre, 
                  COALESCE(u.nombres || ' ' || u.apellidos, 'Sin asignar')
                  FROM instancia_tour i
                  INNER JOIN tour_programado tp ON i.id_tour_programado = tp.id_tour_programado
                  INNER JOIN tipo_tour t ON tp.id_tipo_tour = t.id_tipo_tour
                  INNER JOIN embarcacion e ON i.id_embarcacion = e.id_embarcacion
                  INNER JOIN sede s ON tp.id_sede = s.id_sede
                  LEFT JOIN usuario u ON i.id_chofer = u.id_usuario
                  WHERE i.eliminado = false`

	// Agregar condiciones según los filtros
	conditions := []string{}
	args := []interface{}{}
	argCount := 1

	if filtros.IDTourProgramado != nil {
		conditions = append(conditions, "i.id_tour_programado = $"+string(48+argCount))
		args = append(args, *filtros.IDTourProgramado)
		argCount++
	}

	if filtros.FechaInicio != nil {
		conditions = append(conditions, "i.fecha_especifica >= $"+string(48+argCount))
		args = append(args, *filtros.FechaInicio)
		argCount++
	}

	if filtros.FechaFin != nil {
		conditions = append(conditions, "i.fecha_especifica <= $"+string(48+argCount))
		args = append(args, *filtros.FechaFin)
		argCount++
	}

	if filtros.Estado != nil {
		conditions = append(conditions, "i.estado = $"+string(48+argCount))
		args = append(args, *filtros.Estado)
		argCount++
	}

	if filtros.IDChofer != nil {
		conditions = append(conditions, "i.id_chofer = $"+string(48+argCount))
		args = append(args, *filtros.IDChofer)
		argCount++
	}

	if filtros.IDEmbarcacion != nil {
		conditions = append(conditions, "i.id_embarcacion = $"+string(48+argCount))
		args = append(args, *filtros.IDEmbarcacion)
		argCount++
	}

	if filtros.IDSede != nil {
		conditions = append(conditions, "tp.id_sede = $"+string(48+argCount))
		args = append(args, *filtros.IDSede)
		argCount++
	}

	if filtros.IDTipoTour != nil {
		conditions = append(conditions, "tp.id_tipo_tour = $"+string(48+argCount))
		args = append(args, *filtros.IDTipoTour)
		argCount++
	}

	// Agregar condiciones a la consulta base
	query := queryBase
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY i.fecha_especifica, i.hora_inicio"

	// Ejecutar la consulta
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instancias := []*entidades.InstanciaTour{}

	for rows.Next() {
		instancia := &entidades.InstanciaTour{}
		err := rows.Scan(
			&instancia.ID, &instancia.IDTourProgramado, &instancia.FechaEspecifica, &instancia.HoraInicio, &instancia.HoraFin,
			&instancia.IDChofer, &instancia.IDEmbarcacion, &instancia.CupoDisponible, &instancia.Estado, &instancia.Eliminado,
			&instancia.NombreTipoTour, &instancia.NombreEmbarcacion, &instancia.NombreSede, &instancia.NombreChofer,
		)
		if err != nil {
			return nil, err
		}

		// Formatear las fechas y horas para presentación
		instancia.HoraInicioStr = instancia.HoraInicio.Format("15:04")
		instancia.HoraFinStr = instancia.HoraFin.Format("15:04")
		instancia.FechaEspecificaStr = instancia.FechaEspecifica.Format("2006-01-02")

		instancias = append(instancias, instancia)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instancias, nil
}

// AsignarChofer asigna un chofer a una instancia de tour
func (r *InstanciaTourRepository) AsignarChofer(id int, idChofer int) error {
	// Verificar que la instancia existe
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM instancia_tour WHERE id_instancia = $1 AND eliminado = false)",
		id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("la instancia de tour especificada no existe")
	}

	// Verificar que el chofer existe
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuario WHERE id_usuario = $1 AND rol = 'CHOFER' AND eliminado = false)",
		idChofer).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("el chofer especificado no existe o no tiene rol de chofer")
	}

	// Obtener información de la instancia
	var fechaEspecifica time.Time
	var horaInicio time.Time
	var horaFin time.Time
	var estado string

	err = r.db.QueryRow(`
		SELECT fecha_especifica, hora_inicio, hora_fin, estado
		FROM instancia_tour
		WHERE id_instancia = $1`, id).Scan(&fechaEspecifica, &horaInicio, &horaFin, &estado)
	if err != nil {
		return err
	}

	if estado != "PROGRAMADO" {
		return errors.New("solo se puede asignar un chofer a instancias en estado PROGRAMADO")
	}

	// Verificar disponibilidad del chofer
	diaSemana := int(fechaEspecifica.Weekday())
	if diaSemana == 0 { // En Go, Sunday es 0, pero en nuestra DB es 7
		diaSemana = 7 // Domingo
	}

	var condicionDia string
	switch diaSemana {
	case 1:
		condicionDia = "disponible_lunes"
	case 2:
		condicionDia = "disponible_martes"
	case 3:
		condicionDia = "disponible_miercoles"
	case 4:
		condicionDia = "disponible_jueves"
	case 5:
		condicionDia = "disponible_viernes"
	case 6:
		condicionDia = "disponible_sabado"
	case 7:
		condicionDia = "disponible_domingo"
	}

	queryDisponibilidad := `
		SELECT EXISTS (
			SELECT 1 FROM horario_chofer 
			WHERE id_usuario = $1 
			AND ` + condicionDia + ` = true 
			AND hora_inicio <= $2::time 
			AND hora_fin >= $3::time
			AND fecha_inicio <= $4 
			AND (fecha_fin IS NULL OR fecha_fin >= $4)
			AND eliminado = false
		)`

	var disponible bool
	err = r.db.QueryRow(queryDisponibilidad, idChofer, horaInicio, horaFin,
		fechaEspecifica).Scan(&disponible)
	if err != nil {
		return err
	}
	if !disponible {
		return errors.New("el chofer no está disponible en la fecha y horario especificados")
	}

	// Verificar que el chofer no esté asignado a otro tour en el mismo horario
	queryOcupado := `
		SELECT EXISTS (
			SELECT 1 FROM instancia_tour 
			WHERE id_chofer = $1 
			AND fecha_especifica = $2 
			AND (
				(hora_inicio <= $3::time AND hora_fin > $3::time) OR 
				(hora_inicio < $4::time AND hora_fin >= $4::time) OR 
				(hora_inicio >= $3::time AND hora_fin <= $4::time)
			)
			AND id_instancia != $5
			AND estado IN ('PROGRAMADO', 'EN_CURSO')
			AND eliminado = false
		)`

	var ocupado bool
	err = r.db.QueryRow(queryOcupado, idChofer, fechaEspecifica,
		horaInicio, horaFin, id).Scan(&ocupado)
	if err != nil {
		return err
	}
	if ocupado {
		return errors.New("el chofer ya está asignado a otro tour en el mismo horario")
	}

	// Asignar el chofer
	query := `UPDATE instancia_tour SET id_chofer = $1 WHERE id_instancia = $2`
	_, err = r.db.Exec(query, idChofer, id)
	return err
}

// GenerarInstanciasDeTourProgramado genera instancias para un tour programado
func (r *InstanciaTourRepository) GenerarInstanciasDeTourProgramado(idTourProgramado int) (int, error) {
	// Obtener información del tour programado
	var tp entidades.TourProgramado
	var horarioTour entidades.HorarioTour

	// Consultar tour programado
	queryTP := `SELECT tp.id_tour_programado, tp.id_tipo_tour, tp.id_embarcacion, tp.id_horario, 
				tp.id_sede, tp.id_chofer, tp.vigencia_desde, tp.vigencia_hasta, 
				tp.cupo_maximo, tp.cupo_disponible
				FROM tour_programado tp
				WHERE tp.id_tour_programado = $1 AND tp.eliminado = false`

	err := r.db.QueryRow(queryTP, idTourProgramado).Scan(
		&tp.ID, &tp.IDTipoTour, &tp.IDEmbarcacion, &tp.IDHorario,
		&tp.IDSede, &tp.IDChofer, &tp.VigenciaDesde, &tp.VigenciaHasta,
		&tp.CupoMaximo, &tp.CupoDisponible)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("tour programado no encontrado")
		}
		return 0, err
	}

	// Consultar horario del tour
	queryHorario := `SELECT h.hora_inicio, h.hora_fin, 
					h.disponible_lunes, h.disponible_martes, h.disponible_miercoles, 
					h.disponible_jueves, h.disponible_viernes, h.disponible_sabado, h.disponible_domingo
					FROM horario_tour h
					WHERE h.id_horario = $1 AND h.eliminado = false`

	err = r.db.QueryRow(queryHorario, tp.IDHorario).Scan(
		&horarioTour.HoraInicio, &horarioTour.HoraFin,
		&horarioTour.DisponibleLunes, &horarioTour.DisponibleMartes, &horarioTour.DisponibleMiercoles,
		&horarioTour.DisponibleJueves, &horarioTour.DisponibleViernes, &horarioTour.DisponibleSabado,
		&horarioTour.DisponibleDomingo)

	if err != nil {
		return 0, err
	}

	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generar instancias para todas las fechas dentro del rango de vigencia
	startDate := tp.VigenciaDesde
	endDate := tp.VigenciaHasta
	currentDate := startDate

	// Array para almacenar los días disponibles (1=Lunes, 7=Domingo)
	diasDisponibles := []bool{
		false, // Índice 0 no se usa
		horarioTour.DisponibleLunes,
		horarioTour.DisponibleMartes,
		horarioTour.DisponibleMiercoles,
		horarioTour.DisponibleJueves,
		horarioTour.DisponibleViernes,
		horarioTour.DisponibleSabado,
		horarioTour.DisponibleDomingo,
	}

	// Contador de instancias creadas
	instanciasCreadas := 0

	// Iterar sobre cada día en el rango de fechas
	for !currentDate.After(endDate) {
		// Obtener el día de la semana (1=Lunes, 7=Domingo)
		diaSemana := int(currentDate.Weekday())
		if diaSemana == 0 { // En Go, Sunday es 0
			diaSemana = 7 // Convertir a 7 para Domingo
		}

		// Verificar si este día de la semana está disponible
		if diasDisponibles[diaSemana] {
			// Crear instancia para este día
			query := `INSERT INTO instancia_tour (id_tour_programado, fecha_especifica, hora_inicio, hora_fin, 
					id_chofer, id_embarcacion, cupo_disponible, estado, eliminado) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, 'PROGRAMADO', false)`

			_, err = tx.Exec(
				query,
				tp.ID,
				currentDate,
				horarioTour.HoraInicio,
				horarioTour.HoraFin,
				tp.IDChofer.Int64, // Usar el chofer del tour programado
				tp.IDEmbarcacion,
				tp.CupoMaximo,
			)
			if err != nil {
				return 0, err
			}

			instanciasCreadas++
		}

		// Avanzar al siguiente día
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Si no se creó ninguna instancia, devolver error
	if instanciasCreadas == 0 {
		return 0, errors.New("no se pudo crear ninguna instancia: no hay días disponibles en el rango de fechas")
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return instanciasCreadas, nil
}
