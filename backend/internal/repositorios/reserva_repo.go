package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
	"time"
)

// ReservaRepository maneja las operaciones de base de datos para reservas
type ReservaRepository struct {
	db *sql.DB
}

// NewReservaRepository crea una nueva instancia del repositorio
func NewReservaRepository(db *sql.DB) *ReservaRepository {
	return &ReservaRepository{
		db: db,
	}
}

// GetByID obtiene una reserva por su ID
func (r *ReservaRepository) GetByID(id int) (*entidades.Reserva, error) {
	// Inicializar objeto de reserva
	reserva := &entidades.Reserva{}

	// Consulta para obtener datos de la reserva y entidades relacionadas
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              WHERE r.id_reserva = $1`

	err := r.db.QueryRow(query, id).Scan(
		&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
		&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
		&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
		&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("reserva no encontrada")
		}
		return nil, err
	}

	// Obtener las cantidades de pasajes
	queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                     FROM pasajes_cantidad pc
                     INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                     WHERE pc.id_reserva = $1`

	rowsPasajes, err := r.db.Query(queryPasajes, id)
	if err != nil {
		return nil, err
	}
	defer rowsPasajes.Close()

	reserva.CantidadPasajes = []entidades.PasajeCantidad{}

	for rowsPasajes.Next() {
		var pasajeCantidad entidades.PasajeCantidad
		err := rowsPasajes.Scan(
			&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
		)
		if err != nil {
			return nil, err
		}
		reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
	}

	if err = rowsPasajes.Err(); err != nil {
		return nil, err
	}

	return reserva, nil
}

// Create guarda una nueva reserva en la base de datos
func (r *ReservaRepository) Create(tx *sql.Tx, reserva *entidades.NuevaReservaRequest) (int, error) {
	var id int
	query := `INSERT INTO reserva (id_vendedor, id_cliente, id_tour_programado, id_canal, total_pagar, notas)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id_reserva`

	err := tx.QueryRow(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDTourProgramado,
		reserva.IDCanal,
		reserva.TotalPagar,
		reserva.Notas,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Insertar las cantidades de pasajes
	for _, pasaje := range reserva.CantidadPasajes {
		queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad)
                       VALUES ($1, $2, $3)`

		_, err = tx.Exec(queryPasaje, id, pasaje.IDTipoPasaje, pasaje.Cantidad)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// Update actualiza la información de una reserva
func (r *ReservaRepository) Update(tx *sql.Tx, id int, reserva *entidades.ActualizarReservaRequest) error {
	// Actualizar la reserva
	query := `UPDATE reserva SET
              id_vendedor = $1,
              id_cliente = $2,
              id_tour_programado = $3,
              id_canal = $4,
              total_pagar = $5,
              notas = $6,
              estado = $7
              WHERE id_reserva = $8`

	_, err := tx.Exec(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDTourProgramado,
		reserva.IDCanal,
		reserva.TotalPagar,
		reserva.Notas,
		reserva.Estado,
		id,
	)

	if err != nil {
		return err
	}

	// Eliminar pasajes_cantidad existentes
	queryDeletePasajes := `DELETE FROM pasajes_cantidad WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Insertar nuevas cantidades de pasajes
	for _, pasaje := range reserva.CantidadPasajes {
		queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad)
                       VALUES ($1, $2, $3)`

		_, err = tx.Exec(queryPasaje, id, pasaje.IDTipoPasaje, pasaje.Cantidad)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateEstado actualiza solo el estado de una reserva
func (r *ReservaRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE reserva SET estado = $1 WHERE id_reserva = $2`
	_, err := r.db.Exec(query, estado, id)
	return err
}

// Delete elimina una reserva
func (r *ReservaRepository) Delete(id int) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar si hay pagos asociados a esta reserva
	var countPagos int
	queryCheckPagos := `SELECT COUNT(*) FROM pago WHERE id_reserva = $1`
	err = tx.QueryRow(queryCheckPagos, id).Scan(&countPagos)
	if err != nil {
		return err
	}

	if countPagos > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene pagos asociados")
	}

	// Verificar si hay comprobantes asociados a esta reserva
	var countComprobantes int
	queryCheckComprobantes := `SELECT COUNT(*) FROM comprobante_pago WHERE id_reserva = $1`
	err = tx.QueryRow(queryCheckComprobantes, id).Scan(&countComprobantes)
	if err != nil {
		return err
	}

	if countComprobantes > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene comprobantes asociados")
	}

	// Eliminar los registros de pasajes_cantidad
	queryDeletePasajes := `DELETE FROM pasajes_cantidad WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Eliminar la reserva
	queryDeleteReserva := `DELETE FROM reserva WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeleteReserva, id)
	if err != nil {
		return err
	}

	// Commit de la transacción
	return tx.Commit()
}

// GetCantidadPasajerosByReserva obtiene la cantidad total de pasajeros en una reserva
func (r *ReservaRepository) GetCantidadPasajerosByReserva(id int) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(cantidad), 0)
              FROM pasajes_cantidad
              WHERE id_reserva = $1`

	err := r.db.QueryRow(query, id).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// List lista todas las reservas
func (r *ReservaRepository) List() ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByCliente lista todas las reservas de un cliente
func (r *ReservaRepository) ListByCliente(idCliente int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              WHERE r.id_cliente = $1
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByTourProgramado lista todas las reservas para un tour programado
func (r *ReservaRepository) ListByTourProgramado(idTourProgramado int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
			                INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              WHERE r.id_tour_programado = $1
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idTourProgramado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByFecha lista todas las reservas para una fecha específica
func (r *ReservaRepository) ListByFecha(fecha time.Time) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              WHERE tp.fecha = $1
              ORDER BY ht.hora_inicio ASC, r.fecha_reserva DESC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByEstado lista todas las reservas por estado
func (r *ReservaRepository) ListByEstado(estado string) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.fecha_reserva, r.total_pagar, r.notas, r.estado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              WHERE r.estado = $1
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.FechaReserva, &reserva.TotalPagar, &reserva.Notas, &reserva.Estado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}
