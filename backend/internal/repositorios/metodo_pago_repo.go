package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
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
	query := `SELECT mp.id_metodo_pago, mp.id_sede, mp.nombre, mp.descripcion, mp.eliminado, s.nombre as nombre_sede
              FROM metodo_pago mp
              INNER JOIN sede s ON mp.id_sede = s.id_sede
              WHERE mp.id_metodo_pago = $1`

	err := r.db.QueryRow(query, id).Scan(
		&metodoPago.ID, &metodoPago.IDSede, &metodoPago.Nombre, &metodoPago.Descripcion,
		&metodoPago.Eliminado, &metodoPago.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("método de pago no encontrado")
		}
		return nil, err
	}

	return metodoPago, nil
}

// GetByNombre obtiene un método de pago por su nombre en una sede específica
func (r *MetodoPagoRepository) GetByNombre(nombre string, idSede int) (*entidades.MetodoPago, error) {
	metodoPago := &entidades.MetodoPago{}
	query := `SELECT id_metodo_pago, id_sede, nombre, descripcion, eliminado
              FROM metodo_pago
              WHERE nombre = $1 AND id_sede = $2 AND eliminado = false`

	err := r.db.QueryRow(query, nombre, idSede).Scan(
		&metodoPago.ID, &metodoPago.IDSede, &metodoPago.Nombre, &metodoPago.Descripcion, &metodoPago.Eliminado,
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
	query := `INSERT INTO metodo_pago (id_sede, nombre, descripcion, eliminado)
              VALUES ($1, $2, $3, false)
              RETURNING id_metodo_pago`

	err := r.db.QueryRow(
		query,
		metodoPago.IDSede,
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
              id_sede = $1,
              nombre = $2,
              descripcion = $3,
              eliminado = $4
              WHERE id_metodo_pago = $5`

	_, err := r.db.Exec(
		query,
		metodoPago.IDSede,
		metodoPago.Nombre,
		metodoPago.Descripcion,
		metodoPago.Eliminado,
		id,
	)

	return err
}

// Delete marca un método de pago como eliminado (borrado lógico)
func (r *MetodoPagoRepository) Delete(id int) error {
	// Verificar si hay pagos que usan este método de pago
	var count int
	queryCheck := `SELECT COUNT(*) FROM pago WHERE id_metodo_pago = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este método de pago porque está siendo utilizado en pagos")
	}

	// Marcar como eliminado
	query := `UPDATE metodo_pago SET eliminado = true WHERE id_metodo_pago = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los métodos de pago no eliminados
func (r *MetodoPagoRepository) List() ([]*entidades.MetodoPago, error) {
	query := `SELECT mp.id_metodo_pago, mp.id_sede, mp.nombre, mp.descripcion, mp.eliminado, s.nombre as nombre_sede
              FROM metodo_pago mp
              INNER JOIN sede s ON mp.id_sede = s.id_sede
              WHERE mp.eliminado = false
              ORDER BY mp.nombre ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metodosPago := []*entidades.MetodoPago{}

	for rows.Next() {
		metodoPago := &entidades.MetodoPago{}
		err := rows.Scan(
			&metodoPago.ID, &metodoPago.IDSede, &metodoPago.Nombre, &metodoPago.Descripcion,
			&metodoPago.Eliminado, &metodoPago.NombreSede,
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

// ListBySede lista todos los métodos de pago de una sede específica y no eliminados
func (r *MetodoPagoRepository) ListBySede(idSede int) ([]*entidades.MetodoPago, error) {
	query := `SELECT mp.id_metodo_pago, mp.id_sede, mp.nombre, mp.descripcion, mp.eliminado, s.nombre as nombre_sede
              FROM metodo_pago mp
              INNER JOIN sede s ON mp.id_sede = s.id_sede
              WHERE mp.id_sede = $1 AND mp.eliminado = false
              ORDER BY mp.nombre ASC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metodosPago := []*entidades.MetodoPago{}

	for rows.Next() {
		metodoPago := &entidades.MetodoPago{}
		err := rows.Scan(
			&metodoPago.ID, &metodoPago.IDSede, &metodoPago.Nombre, &metodoPago.Descripcion,
			&metodoPago.Eliminado, &metodoPago.NombreSede,
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
