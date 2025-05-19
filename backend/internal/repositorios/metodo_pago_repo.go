package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
)

// MetodoPagoRepository maneja las operaciones de base de datos para métodos de pago
type MetodoPagoRepository struct {
	db *sql.DB
}

// NewMetodoPagoRepository crea una nueva instancia del repositorio
func NewMetodoPagoRepository(db *sql.DB) *MetodoPagoRepository {
	return &MetodoPagoRepository{
		db: db,
	}
}

// GetByID obtiene un método de pago por su ID
func (r *MetodoPagoRepository) GetByID(id int) (*entidades.MetodoPago, error) {
	metodoPago := &entidades.MetodoPago{}
	query := `SELECT id_metodo_pago, nombre, descripcion
              FROM metodo_pago
              WHERE id_metodo_pago = $1`

	err := r.db.QueryRow(query, id).Scan(
		&metodoPago.ID, &metodoPago.Nombre, &metodoPago.Descripcion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("método de pago no encontrado")
		}
		return nil, err
	}

	return metodoPago, nil
}

// GetByNombre obtiene un método de pago por su nombre
func (r *MetodoPagoRepository) GetByNombre(nombre string) (*entidades.MetodoPago, error) {
	metodoPago := &entidades.MetodoPago{}
	query := `SELECT id_metodo_pago, nombre, descripcion
              FROM metodo_pago
              WHERE nombre = $1`

	err := r.db.QueryRow(query, nombre).Scan(
		&metodoPago.ID, &metodoPago.Nombre, &metodoPago.Descripcion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("método de pago no encontrado")
		}
		return nil, err
	}

	return metodoPago, nil
}

// Create guarda un nuevo método de pago en la base de datos
func (r *MetodoPagoRepository) Create(metodoPago *entidades.NuevoMetodoPagoRequest) (int, error) {
	var id int
	query := `INSERT INTO metodo_pago (nombre, descripcion)
              VALUES ($1, $2)
              RETURNING id_metodo_pago`

	err := r.db.QueryRow(
		query,
		metodoPago.Nombre,
		metodoPago.Descripcion,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un método de pago
func (r *MetodoPagoRepository) Update(id int, metodoPago *entidades.ActualizarMetodoPagoRequest) error {
	query := `UPDATE metodo_pago SET
              nombre = $1,
              descripcion = $2
              WHERE id_metodo_pago = $3`

	_, err := r.db.Exec(
		query,
		metodoPago.Nombre,
		metodoPago.Descripcion,
		id,
	)

	return err
}

// Delete elimina un método de pago
func (r *MetodoPagoRepository) Delete(id int) error {
	// Verificar si hay pagos que usan este método de pago
	var count int
	queryCheck := `SELECT COUNT(*) FROM pago WHERE id_metodo_pago = $1`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este método de pago porque está siendo utilizado en pagos")
	}

	// Eliminar método de pago
	query := `DELETE FROM metodo_pago WHERE id_metodo_pago = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los métodos de pago
func (r *MetodoPagoRepository) List() ([]*entidades.MetodoPago, error) {
	query := `SELECT id_metodo_pago, nombre, descripcion
              FROM metodo_pago
              ORDER BY nombre ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metodosPago := []*entidades.MetodoPago{}

	for rows.Next() {
		metodoPago := &entidades.MetodoPago{}
		err := rows.Scan(
			&metodoPago.ID, &metodoPago.Nombre, &metodoPago.Descripcion,
		)
		if err != nil {
			return nil, err
		}
		metodosPago = append(metodosPago, metodoPago)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metodosPago, nil
}
