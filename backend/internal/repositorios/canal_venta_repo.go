package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
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
	query := `SELECT cv.id_canal, cv.id_sede, cv.nombre, cv.descripcion, cv.eliminado, s.nombre as nombre_sede
              FROM canal_venta cv
              INNER JOIN sede s ON cv.id_sede = s.id_sede
              WHERE cv.id_canal = $1`

	err := r.db.QueryRow(query, id).Scan(
		&canal.ID, &canal.IDSede, &canal.Nombre, &canal.Descripcion,
		&canal.Eliminado, &canal.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("canal de venta no encontrado")
		}
		return nil, err
	}

	return canal, nil
}

// GetByNombre obtiene un canal de venta por su nombre en una sede específica
func (r *CanalVentaRepository) GetByNombre(nombre string, idSede int) (*entidades.CanalVenta, error) {
	canal := &entidades.CanalVenta{}
	query := `SELECT id_canal, id_sede, nombre, descripcion, eliminado
              FROM canal_venta
              WHERE nombre = $1 AND id_sede = $2 AND eliminado = false`

	err := r.db.QueryRow(query, nombre, idSede).Scan(
		&canal.ID, &canal.IDSede, &canal.Nombre, &canal.Descripcion, &canal.Eliminado,
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
	query := `INSERT INTO canal_venta (id_sede, nombre, descripcion, eliminado)
              VALUES ($1, $2, $3, false)
              RETURNING id_canal`

	err := r.db.QueryRow(
		query,
		canal.IDSede,
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
              id_sede = $1,
              nombre = $2,
              descripcion = $3,
              eliminado = $4
              WHERE id_canal = $5`

	_, err := r.db.Exec(
		query,
		canal.IDSede,
		canal.Nombre,
		canal.Descripcion,
		canal.Eliminado,
		id,
	)

	return err
}

// Delete marca un canal de venta como eliminado (borrado lógico)
func (r *CanalVentaRepository) Delete(id int) error {
	// Verificar si hay reservas que usan este canal
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_canal = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar este canal de venta porque está siendo utilizado en reservas")
	}

	// Verificar si hay pagos que usan este canal
	var countPagos int
	queryCheckPagos := `SELECT COUNT(*) FROM pago WHERE id_canal = $1 AND eliminado = false`
	err = r.db.QueryRow(queryCheckPagos, id).Scan(&countPagos)
	if err != nil {
		return err
	}

	if countPagos > 0 {
		return errors.New("no se puede eliminar este canal de venta porque está siendo utilizado en pagos")
	}

	// Marcar como eliminado
	query := `UPDATE canal_venta SET eliminado = true WHERE id_canal = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los canales de venta no eliminados
func (r *CanalVentaRepository) List() ([]*entidades.CanalVenta, error) {
	query := `SELECT cv.id_canal, cv.id_sede, cv.nombre, cv.descripcion, cv.eliminado, s.nombre as nombre_sede
              FROM canal_venta cv
              INNER JOIN sede s ON cv.id_sede = s.id_sede
              WHERE cv.eliminado = false
              ORDER BY cv.nombre ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	canales := []*entidades.CanalVenta{}

	for rows.Next() {
		canal := &entidades.CanalVenta{}
		err := rows.Scan(
			&canal.ID, &canal.IDSede, &canal.Nombre, &canal.Descripcion,
			&canal.Eliminado, &canal.NombreSede,
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

// ListBySede lista todos los canales de venta de una sede específica y no eliminados
func (r *CanalVentaRepository) ListBySede(idSede int) ([]*entidades.CanalVenta, error) {
	query := `SELECT cv.id_canal, cv.id_sede, cv.nombre, cv.descripcion, cv.eliminado, s.nombre as nombre_sede
              FROM canal_venta cv
              INNER JOIN sede s ON cv.id_sede = s.id_sede
              WHERE cv.id_sede = $1 AND cv.eliminado = false
              ORDER BY cv.nombre ASC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	canales := []*entidades.CanalVenta{}

	for rows.Next() {
		canal := &entidades.CanalVenta{}
		err := rows.Scan(
			&canal.ID, &canal.IDSede, &canal.Nombre, &canal.Descripcion,
			&canal.Eliminado, &canal.NombreSede,
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
