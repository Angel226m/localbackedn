package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
	"time"
)

// PagoRepository maneja las operaciones de base de datos para pagos
type PagoRepository struct {
	db *sql.DB
}

// NewPagoRepository crea una nueva instancia del repositorio
func NewPagoRepository(db *sql.DB) *PagoRepository {
	return &PagoRepository{
		db: db,
	}
}

// GetByID obtiene un pago por su ID
func (r *PagoRepository) GetByID(id int) (*entidades.Pago, error) {
	pago := &entidades.Pago{}
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE p.id_pago = $1 AND p.eliminado = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
		&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
		&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
		&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
		&pago.TourNombre, &pago.TourFecha,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pago no encontrado")
		}
		return nil, err
	}

	return pago, nil
}

// Create guarda un nuevo pago en la base de datos
func (r *PagoRepository) Create(pago *entidades.NuevoPagoRequest) (int, error) {
	var id int
	query := `INSERT INTO pago (id_reserva, id_metodo_pago, id_canal, id_sede, monto, comprobante, eliminado)
              VALUES ($1, $2, $3, $4, $5, $6, FALSE)
              RETURNING id_pago`

	err := r.db.QueryRow(
		query,
		pago.IDReserva,
		pago.IDMetodoPago,
		pago.IDCanal,
		pago.IDSede,
		pago.Monto,
		pago.Comprobante,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un pago
func (r *PagoRepository) Update(id int, pago *entidades.ActualizarPagoRequest) error {
	query := `UPDATE pago SET
              id_metodo_pago = $1,
              id_canal = $2,
              id_sede = $3,
              monto = $4,
              comprobante = $5,
              estado = $6
              WHERE id_pago = $7 AND eliminado = FALSE`

	_, err := r.db.Exec(
		query,
		pago.IDMetodoPago,
		pago.IDCanal,
		pago.IDSede,
		pago.Monto,
		pago.Comprobante,
		pago.Estado,
		id,
	)

	return err
}

// UpdateEstado actualiza solo el estado de un pago
func (r *PagoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE pago SET estado = $1 WHERE id_pago = $2 AND eliminado = FALSE`
	_, err := r.db.Exec(query, estado, id)
	return err
}

// Delete elimina un pago
func (r *PagoRepository) Delete(id int) error {
	// Verificar si hay comprobantes asociados a este pago a través de la reserva
	var idReserva int
	queryGetReserva := `SELECT id_reserva FROM pago WHERE id_pago = $1 AND eliminado = FALSE`
	err := r.db.QueryRow(queryGetReserva, id).Scan(&idReserva)
	if err != nil {
		return err
	}

	var countComprobantes int
	queryCheckComprobantes := `SELECT COUNT(*) FROM comprobante_pago WHERE id_reserva = $1 AND eliminado = FALSE`
	err = r.db.QueryRow(queryCheckComprobantes, idReserva).Scan(&countComprobantes)
	if err != nil {
		return err
	}

	if countComprobantes > 0 {
		return errors.New("no se puede eliminar este pago porque la reserva tiene comprobantes asociados")
	}

	// Eliminación lógica
	query := `UPDATE pago SET eliminado = TRUE WHERE id_pago = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los pagos activos
func (r *PagoRepository) List() ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}

// ListByReserva lista todos los pagos de una reserva específica
func (r *PagoRepository) ListByReserva(idReserva int) ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE p.id_reserva = $1 AND p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query, idReserva)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}

// ListByFecha lista todos los pagos de una fecha específica
func (r *PagoRepository) ListByFecha(fecha time.Time) ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE DATE(p.fecha_pago) = $1 AND p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}

// GetTotalPagadoByReserva obtiene el total pagado de una reserva específica
func (r *PagoRepository) GetTotalPagadoByReserva(idReserva int) (float64, error) {
	var totalPagado float64
	query := `SELECT COALESCE(SUM(monto), 0) FROM pago WHERE id_reserva = $1 AND estado = 'PROCESADO' AND eliminado = FALSE`

	err := r.db.QueryRow(query, idReserva).Scan(&totalPagado)
	if err != nil {
		return 0, err
	}

	return totalPagado, nil
}

// ListByEstado lista todos los pagos con un estado específico
func (r *PagoRepository) ListByEstado(estado string) ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE p.estado = $1 AND p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}

// ListByCliente lista todos los pagos de un cliente específico
func (r *PagoRepository) ListByCliente(idCliente int) ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE r.id_cliente = $1 AND p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}

// ListBySede lista todos los pagos de una sede específica
func (r *PagoRepository) ListBySede(idSede int) ([]*entidades.Pago, error) {
	query := `SELECT p.id_pago, p.id_reserva, p.id_metodo_pago, p.id_canal, p.id_sede,
              p.monto, p.fecha_pago, p.comprobante, p.estado, p.eliminado,
              c.nombres, c.apellidos, c.numero_documento,
              mp.nombre, cv.nombre, s.nombre,
              tt.nombre, tp.fecha
              FROM pago p
              INNER JOIN reserva r ON p.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN metodo_pago mp ON p.id_metodo_pago = mp.id_metodo_pago
              INNER JOIN canal_venta cv ON p.id_canal = cv.id_canal
              INNER JOIN sede s ON p.id_sede = s.id_sede
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE p.id_sede = $1 AND p.eliminado = FALSE
              ORDER BY p.fecha_pago DESC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagos := []*entidades.Pago{}

	for rows.Next() {
		pago := &entidades.Pago{}
		err := rows.Scan(
			&pago.ID, &pago.IDReserva, &pago.IDMetodoPago, &pago.IDCanal, &pago.IDSede,
			&pago.Monto, &pago.FechaPago, &pago.Comprobante, &pago.Estado, &pago.Eliminado,
			&pago.NombreCliente, &pago.ApellidosCliente, &pago.DocumentoCliente,
			&pago.NombreMetodoPago, &pago.NombreCanalVenta, &pago.NombreSede,
			&pago.TourNombre, &pago.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, pago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pagos, nil
}
