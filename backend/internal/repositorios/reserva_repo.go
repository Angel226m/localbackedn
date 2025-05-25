package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
	"time"
)

// ReservaRepository maneja las operaciones de base de datos para reservas
// Implementa métodos para crear, actualizar, eliminar y consultar reservas
type ReservaRepository struct {
	db *sql.DB
}

// NewReservaRepository crea una nueva instancia del repositorio de reservas
// Recibe una conexión a la base de datos y retorna un puntero a ReservaRepository
func NewReservaRepository(db *sql.DB) *ReservaRepository {
	return &ReservaRepository{
		db: db,
	}
}

// GetByID obtiene una reserva por su ID
// Incluye información relacionada como cliente, vendedor, tour, etc.
// Retorna un error si la reserva no existe o hay problemas con la consulta
func (r *ReservaRepository) GetByID(id int) (*entidades.Reserva, error) {
	// Inicializar objeto de reserva
	reserva := &entidades.Reserva{}

	// Consulta para obtener datos de la reserva y entidades relacionadas
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_reserva = $1 AND r.eliminado = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
		&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
		&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
		&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
		&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("reserva no encontrada")
		}
		return nil, err
	}

	// Obtener las cantidades de pasajes asociados a la reserva
	queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                     FROM pasajes_cantidad pc
                     INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                     WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

	rowsPasajes, err := r.db.Query(queryPasajes, id)
	if err != nil {
		return nil, err
	}
	defer rowsPasajes.Close()

	// Inicializar el slice de cantidades de pasajes
	reserva.CantidadPasajes = []entidades.PasajeCantidad{}

	// Iterar por cada registro de pasajes
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

	// Verificar errores durante la iteración
	if err = rowsPasajes.Err(); err != nil {
		return nil, err
	}

	return reserva, nil
}

// Create guarda una nueva reserva en la base de datos
// Recibe una transacción y los datos de la nueva reserva
// Retorna el ID de la reserva creada o un error en caso de fallo
func (r *ReservaRepository) Create(tx *sql.Tx, reserva *entidades.NuevaReservaRequest) (int, error) {
	var id int
	// Consulta SQL para insertar una nueva reserva
	query := `INSERT INTO reserva (id_vendedor, id_cliente, id_tour_programado, id_canal, id_sede, total_pagar, notas, estado, eliminado)
              VALUES ($1, $2, $3, $4, $5, $6, $7, 'RESERVADO', FALSE)
              RETURNING id_reserva`

	// Ejecutar la consulta con los datos de la reserva
	err := tx.QueryRow(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDTourProgramado,
		reserva.IDCanal,
		reserva.IDSede,
		reserva.TotalPagar,
		reserva.Notas,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Insertar las cantidades de pasajes asociados a la reserva
	for _, pasaje := range reserva.CantidadPasajes {
		queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad, eliminado)
                       VALUES ($1, $2, $3, FALSE)`

		_, err = tx.Exec(queryPasaje, id, pasaje.IDTipoPasaje, pasaje.Cantidad)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// Update actualiza la información de una reserva existente
// Recibe una transacción, el ID de la reserva y los datos actualizados
// Retorna error en caso de que la actualización falle
func (r *ReservaRepository) Update(tx *sql.Tx, id int, reserva *entidades.ActualizarReservaRequest) error {
	// Actualizar la reserva con los nuevos datos
	query := `UPDATE reserva SET
              id_vendedor = $1,
              id_cliente = $2,
              id_tour_programado = $3,
              id_canal = $4,
              id_sede = $5,
              total_pagar = $6,
              notas = $7,
              estado = $8
              WHERE id_reserva = $9 AND eliminado = FALSE`

	// Ejecutar la actualización
	_, err := tx.Exec(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDTourProgramado,
		reserva.IDCanal,
		reserva.IDSede,
		reserva.TotalPagar,
		reserva.Notas,
		reserva.Estado,
		id,
	)

	if err != nil {
		return err
	}

	// Eliminar los registros de pasajes_cantidad existentes (eliminación lógica)
	queryDeletePasajes := `UPDATE pasajes_cantidad SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Insertar nuevas cantidades de pasajes
	for _, pasaje := range reserva.CantidadPasajes {
		queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad, eliminado)
                       VALUES ($1, $2, $3, FALSE)`

		_, err = tx.Exec(queryPasaje, id, pasaje.IDTipoPasaje, pasaje.Cantidad)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateEstado actualiza solo el estado de una reserva
// Recibe el ID de la reserva y el nuevo estado
// Retorna error si la actualización falla
func (r *ReservaRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE reserva SET estado = $1 WHERE id_reserva = $2 AND eliminado = FALSE`
	_, err := r.db.Exec(query, estado, id)
	return err
}

// Delete realiza una eliminación lógica de una reserva
// Verifica primero si existen pagos o comprobantes asociados
// Retorna error si no se puede eliminar o hay un problema durante el proceso
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

	// Verificar si hay pagos asociados a esta reserva (que no estén eliminados)
	var countPagos int
	queryCheckPagos := `SELECT COUNT(*) FROM pago WHERE id_reserva = $1 AND eliminado = FALSE`
	err = tx.QueryRow(queryCheckPagos, id).Scan(&countPagos)
	if err != nil {
		return err
	}

	if countPagos > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene pagos asociados")
	}

	// Verificar si hay comprobantes asociados a esta reserva (que no estén eliminados)
	var countComprobantes int
	queryCheckComprobantes := `SELECT COUNT(*) FROM comprobante_pago WHERE id_reserva = $1 AND eliminado = FALSE`
	err = tx.QueryRow(queryCheckComprobantes, id).Scan(&countComprobantes)
	if err != nil {
		return err
	}

	if countComprobantes > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene comprobantes asociados")
	}

	// Marcar los registros de pasajes_cantidad como eliminados (eliminación lógica)
	queryDeletePasajes := `UPDATE pasajes_cantidad SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Marcar la reserva como eliminada (eliminación lógica)
	queryDeleteReserva := `UPDATE reserva SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeleteReserva, id)
	if err != nil {
		return err
	}

	// Commit de la transacción
	return tx.Commit()
}

// GetCantidadPasajerosByReserva obtiene la cantidad total de pasajeros en una reserva
// Suma todas las cantidades de pasajes asociados a la reserva
// Retorna el total de pasajeros o error si hay problemas con la consulta
func (r *ReservaRepository) GetCantidadPasajerosByReserva(id int) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(cantidad), 0)
              FROM pasajes_cantidad
              WHERE id_reserva = $1 AND eliminado = FALSE`

	err := r.db.QueryRow(query, id).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// List obtiene todas las reservas activas del sistema
// Incluye información relacionada como cliente, vendedor, tour, etc.
// Retorna un slice de reservas o error si hay problemas con la consulta
func (r *ReservaRepository) List() ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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

// ListByCliente lista todas las reservas activas de un cliente específico
// Recibe el ID del cliente y retorna todas sus reservas
// Incluye información relacionada y detalles de pasajes
func (r *ReservaRepository) ListByCliente(idCliente int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_cliente = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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

// ListByTourProgramado lista todas las reservas asociadas a un tour programado específico
// Recibe el ID del tour programado y retorna todas las reservas relacionadas
// Incluye información completa de cada reserva y sus pasajes
func (r *ReservaRepository) ListByTourProgramado(idTourProgramado int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_tour_programado = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idTourProgramado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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

// ListByFecha lista todas las reservas para una fecha específica de tour
// Recibe la fecha y retorna todas las reservas para tours programados en esa fecha
// Ordenadas por hora de inicio del tour y fecha de reserva
func (r *ReservaRepository) ListByFecha(fecha time.Time) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE tp.fecha = $1 AND r.eliminado = FALSE
              ORDER BY ht.hora_inicio ASC, r.fecha_reserva DESC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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

// ListByEstado lista todas las reservas por estado específico (RESERVADO, CANCELADA)
// Recibe el estado y retorna todas las reservas que tengan ese estado
// Incluye información completa de cada reserva y sus pasajes
func (r *ReservaRepository) ListByEstado(estado string) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.estado = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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

// ListBySede lista todas las reservas de una sede específica o todas las reservas si es ADMIN
func (r *ReservaRepository) ListBySede(idSede *int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_tour_programado, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(tp.fecha, 'DD/MM/YYYY') as fecha_tour,
              to_char(ht.hora_inicio, 'HH24:MI') as hora_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN horario_tour ht ON tp.id_horario = ht.id_horario
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.eliminado = FALSE`

	// Si se proporciona un ID de sede, filtrar por ella
	if idSede != nil {
		query += " AND r.id_sede = $1"
	}

	query += " ORDER BY r.fecha_reserva DESC"

	var rows *sql.Rows
	var err error

	if idSede != nil {
		rows, err = r.db.Query(query, *idSede)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDTourProgramado,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraTour, &reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                         FROM pasajes_cantidad pc
                         INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                         WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
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
