package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
	"time"
)

// ComprobantePagoRepository maneja las operaciones de base de datos para comprobantes de pago
type ComprobantePagoRepository struct {
	db *sql.DB
}

// NewComprobantePagoRepository crea una nueva instancia del repositorio
func NewComprobantePagoRepository(db *sql.DB) *ComprobantePagoRepository {
	return &ComprobantePagoRepository{
		db: db,
	}
}

// GetByID obtiene un comprobante de pago por su ID
func (r *ComprobantePagoRepository) GetByID(id int) (*entidades.ComprobantePago, error) {
	comprobante := &entidades.ComprobantePago{}
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE cp.id_comprobante = $1`

	err := r.db.QueryRow(query, id).Scan(
		&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
		&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
		&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
		&comprobante.TourNombre, &comprobante.TourFecha,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comprobante de pago no encontrado")
		}
		return nil, err
	}

	return comprobante, nil
}

// GetByTipoAndNumero obtiene un comprobante de pago por su tipo y número
func (r *ComprobantePagoRepository) GetByTipoAndNumero(tipo, numero string) (*entidades.ComprobantePago, error) {
	comprobante := &entidades.ComprobantePago{}
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE cp.tipo = $1 AND cp.numero_comprobante = $2`

	err := r.db.QueryRow(query, tipo, numero).Scan(
		&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
		&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
		&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
		&comprobante.TourNombre, &comprobante.TourFecha,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comprobante de pago no encontrado")
		}
		return nil, err
	}

	return comprobante, nil
}

// Create guarda un nuevo comprobante de pago en la base de datos
func (r *ComprobantePagoRepository) Create(comprobante *entidades.NuevoComprobantePagoRequest) (int, error) {
	var id int
	query := `INSERT INTO comprobante_pago (id_reserva, tipo, numero_comprobante, subtotal, igv, total)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id_comprobante`

	err := r.db.QueryRow(
		query,
		comprobante.IDReserva,
		comprobante.Tipo,
		comprobante.NumeroComprobante,
		comprobante.Subtotal,
		comprobante.IGV,
		comprobante.Total,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un comprobante de pago
func (r *ComprobantePagoRepository) Update(id int, comprobante *entidades.ActualizarComprobantePagoRequest) error {
	query := `UPDATE comprobante_pago SET
              tipo = $1,
              numero_comprobante = $2,
              subtotal = $3,
              igv = $4,
              total = $5,
              estado = $6
              WHERE id_comprobante = $7`

	_, err := r.db.Exec(
		query,
		comprobante.Tipo,
		comprobante.NumeroComprobante,
		comprobante.Subtotal,
		comprobante.IGV,
		comprobante.Total,
		comprobante.Estado,
		id,
	)

	return err
}

// UpdateEstado actualiza solo el estado de un comprobante de pago
func (r *ComprobantePagoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE comprobante_pago SET estado = $1 WHERE id_comprobante = $2`
	_, err := r.db.Exec(query, estado, id)
	return err
}

// Delete elimina un comprobante de pago
func (r *ComprobantePagoRepository) Delete(id int) error {
	query := `DELETE FROM comprobante_pago WHERE id_comprobante = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List lista todos los comprobantes de pago
func (r *ComprobantePagoRepository) List() ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}

// ListByReserva lista todos los comprobantes de pago de una reserva específica
func (r *ComprobantePagoRepository) ListByReserva(idReserva int) ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE cp.id_reserva = $1
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query, idReserva)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}

// ListByFecha lista todos los comprobantes de pago de una fecha específica
func (r *ComprobantePagoRepository) ListByFecha(fecha time.Time) ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE DATE(cp.fecha_emision) = $1
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}

// ListByTipo lista todos los comprobantes de pago de un tipo específico
func (r *ComprobantePagoRepository) ListByTipo(tipo string) ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE cp.tipo = $1
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query, tipo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}

// ListByEstado lista todos los comprobantes de pago con un estado específico
func (r *ComprobantePagoRepository) ListByEstado(estado string) ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE cp.estado = $1
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}

// ListByCliente lista todos los comprobantes de un cliente específico
func (r *ComprobantePagoRepository) ListByCliente(idCliente int) ([]*entidades.ComprobantePago, error) {
	query := `SELECT cp.id_comprobante, cp.id_reserva, cp.tipo, cp.numero_comprobante, 
              cp.fecha_emision, cp.subtotal, cp.igv, cp.total, cp.estado,
              c.nombres, c.apellidos, c.numero_documento,
              tt.nombre, tp.fecha
              FROM comprobante_pago cp
              INNER JOIN reserva r ON cp.id_reserva = r.id_reserva
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              INNER JOIN tour_programado tp ON r.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              WHERE r.id_cliente = $1
              ORDER BY cp.fecha_emision DESC`

	rows, err := r.db.Query(query, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comprobantes := []*entidades.ComprobantePago{}

	for rows.Next() {
		comprobante := &entidades.ComprobantePago{}
		err := rows.Scan(
			&comprobante.ID, &comprobante.IDReserva, &comprobante.Tipo, &comprobante.NumeroComprobante,
			&comprobante.FechaEmision, &comprobante.Subtotal, &comprobante.IGV, &comprobante.Total, &comprobante.Estado,
			&comprobante.NombreCliente, &comprobante.ApellidosCliente, &comprobante.DocumentoCliente,
			&comprobante.TourNombre, &comprobante.TourFecha,
		)
		if err != nil {
			return nil, err
		}
		comprobantes = append(comprobantes, comprobante)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comprobantes, nil
}
