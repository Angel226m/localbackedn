package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
)

// CanalVentaRepository maneja las operaciones de base de datos para canales de venta
type CanalVentaRepository struct {
	db *sql.DB
}

// NewCanalVentaRepository crea una nueva instancia del repositorio
func NewCanalVentaRepository(db *sql.DB) *CanalVentaRepository {
	return &CanalVentaRepository{
		db: db,
	}
}

// GetByID obtiene un canal de venta por su ID
func (r *CanalVentaRepository) GetByID(id int) (*entidades.CanalVenta, error) {
	canal := &entidades.CanalVenta{}
	query := `SELECT id_canal, nombre, descripcion
              FROM canal_venta
              WHERE id_canal = $1`

	err := r.db.QueryRow(query, id).Scan(
		&canal.ID, &canal.Nombre, &canal.Descripcion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("canal de venta no encontrado")
		}
		return nil, err
	}

	return canal, nil
}

// GetByNombre obtiene un canal de venta por su nombre
func (r *CanalVentaRepository) GetByNombre(nombre string) (*entidades.CanalVenta, error) {
	canal := &entidades.CanalVenta{}
	query := `SELECT id_canal, nombre, descripcion
              FROM canal_venta
              WHERE nombre = $1`

	err := r.db.QueryRow(query, nombre).Scan(
		&canal.ID, &canal.Nombre, &canal.Descripcion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("canal de venta no encontrado")
		}
		return nil, err
	}

	return canal, nil
}

// Create guarda un nuevo canal de venta en la base de datos
func (r *CanalVentaRepository) Create(canal *entidades.NuevoCanalVentaRequest) (int, error) {
	var id int
	query := `INSERT INTO canal_venta (nombre, descripcion)
              VALUES ($1, $2)
              RETURNING id_canal`

	err := r.db.QueryRow(
		query,
		canal.Nombre,
		canal.Descripcion,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un canal de venta
func (r *CanalVentaRepository) Update(id int, canal *entidades.ActualizarCanalVentaRequest) error {
	query := `UPDATE canal_venta SET
              nombre = $1,
              descripcion = $2
              WHERE id_canal = $3`

	_, err := r.db.Exec(
		query,
		canal.Nombre,
		canal.Descripcion,
		id,
	)

	return err
}

// Delete elimina un canal de venta
func (r *CanalVentaRepository) Delete(id int) error {
	// Verificar si hay reservas que usan este canal
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_canal = $1`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar este canal de venta porque está siendo utilizado en reservas")
	}

	// Verificar si hay pagos que usan este canal
	var countPagos int
	queryCheckPagos := `SELECT COUNT(*) FROM pago WHERE id_canal = $1`
	err = r.db.QueryRow(queryCheckPagos, id).Scan(&countPagos)
	if err != nil {
		return err
	}

	if countPagos > 0 {
		return errors.New("no se puede eliminar este canal de venta porque está siendo utilizado en pagos")
	}

	// Si no hay dependencias, eliminar el canal
	query := `DELETE FROM canal_venta WHERE id_canal = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los canales de venta
func (r *CanalVentaRepository) List() ([]*entidades.CanalVenta, error) {
	query := `SELECT id_canal, nombre, descripcion
              FROM canal_venta
              ORDER BY nombre ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	canales := []*entidades.CanalVenta{}

	for rows.Next() {
		canal := &entidades.CanalVenta{}
		err := rows.Scan(
			&canal.ID, &canal.Nombre, &canal.Descripcion,
		)
		if err != nil {
			return nil, err
		}
		canales = append(canales, canal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return canales, nil
}
