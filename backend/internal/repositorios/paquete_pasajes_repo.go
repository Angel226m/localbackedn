package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// PaquetePasajesRepository maneja las operaciones de base de datos para paquetes de pasajes
type PaquetePasajesRepository struct {
	db *sql.DB
}

// NewPaquetePasajesRepository crea una nueva instancia del repositorio
func NewPaquetePasajesRepository(db *sql.DB) *PaquetePasajesRepository {
	return &PaquetePasajesRepository{
		db: db,
	}
}

// GetByID obtiene un paquete de pasajes por su ID
func (r *PaquetePasajesRepository) GetByID(id int) (*entidades.PaquetePasajes, error) {
	paquete := &entidades.PaquetePasajes{}
	query := `SELECT id_paquete, id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado
              FROM paquete_pasajes
              WHERE id_paquete = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&paquete.ID, &paquete.IDSede, &paquete.IDTipoTour, &paquete.Nombre,
		&paquete.Descripcion, &paquete.PrecioTotal, &paquete.CantidadTotal, &paquete.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("paquete de pasajes no encontrado")
		}
		return nil, err
	}

	return paquete, nil
}

// GetByNombre obtiene un paquete de pasajes por su nombre y sede
func (r *PaquetePasajesRepository) GetByNombre(nombre string, idSede int) (*entidades.PaquetePasajes, error) {
	paquete := &entidades.PaquetePasajes{}
	query := `SELECT id_paquete, id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado
              FROM paquete_pasajes
              WHERE nombre = $1 AND id_sede = $2 AND eliminado = false`

	err := r.db.QueryRow(query, nombre, idSede).Scan(
		&paquete.ID, &paquete.IDSede, &paquete.IDTipoTour, &paquete.Nombre,
		&paquete.Descripcion, &paquete.PrecioTotal, &paquete.CantidadTotal, &paquete.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("paquete de pasajes no encontrado")
		}
		return nil, err
	}

	return paquete, nil
}

// Create guarda un nuevo paquete de pasajes en la base de datos
func (r *PaquetePasajesRepository) Create(paquete *entidades.NuevoPaquetePasajesRequest) (int, error) {
	var id int
	query := `INSERT INTO paquete_pasajes (id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado)
              VALUES ($1, $2, $3, $4, $5, $6, false)
              RETURNING id_paquete`

	err := r.db.QueryRow(
		query,
		paquete.IDSede,
		paquete.IDTipoTour,
		paquete.Nombre,
		paquete.Descripcion,
		paquete.PrecioTotal,
		paquete.CantidadTotal,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un paquete de pasajes
func (r *PaquetePasajesRepository) Update(id int, paquete *entidades.ActualizarPaquetePasajesRequest) error {
	query := `UPDATE paquete_pasajes SET
              id_tipo_tour = $1,
              nombre = $2,
              descripcion = $3,
              precio_total = $4,
              cantidad_total = $5
              WHERE id_paquete = $6 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		paquete.IDTipoTour,
		paquete.Nombre,
		paquete.Descripcion,
		paquete.PrecioTotal,
		paquete.CantidadTotal,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("paquete de pasajes no encontrado o ya eliminado")
	}

	return nil
}

// Delete marca un paquete de pasajes como eliminado (eliminación lógica)
func (r *PaquetePasajesRepository) Delete(id int) error {
	// Verificar si hay reservas que usan este paquete de pasajes
	var count int
	queryCheck := `SELECT COUNT(*) FROM reservas WHERE id_paquete = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este paquete de pasajes porque está siendo utilizado por reservas")
	}

	// Eliminación lógica del paquete de pasajes
	query := `UPDATE paquete_pasajes SET eliminado = true WHERE id_paquete = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("paquete de pasajes no encontrado o ya eliminado")
	}

	return nil
}

// ListBySede lista todos los paquetes de pasajes de una sede específica
func (r *PaquetePasajesRepository) ListBySede(idSede int) ([]*entidades.PaquetePasajes, error) {
	query := `SELECT id_paquete, id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado
              FROM paquete_pasajes
              WHERE id_sede = $1 AND eliminado = false
              ORDER BY precio_total ASC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paquetes := []*entidades.PaquetePasajes{}

	for rows.Next() {
		paquete := &entidades.PaquetePasajes{}
		err := rows.Scan(
			&paquete.ID, &paquete.IDSede, &paquete.IDTipoTour, &paquete.Nombre,
			&paquete.Descripcion, &paquete.PrecioTotal, &paquete.CantidadTotal, &paquete.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		paquetes = append(paquetes, paquete)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return paquetes, nil
}

// ListByTipoTour lista todos los paquetes de pasajes asociados a un tipo de tour específico
func (r *PaquetePasajesRepository) ListByTipoTour(idTipoTour int) ([]*entidades.PaquetePasajes, error) {
	query := `SELECT id_paquete, id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado
              FROM paquete_pasajes
              WHERE id_tipo_tour = $1 AND eliminado = false
              ORDER BY precio_total ASC`

	rows, err := r.db.Query(query, idTipoTour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paquetes := []*entidades.PaquetePasajes{}

	for rows.Next() {
		paquete := &entidades.PaquetePasajes{}
		err := rows.Scan(
			&paquete.ID, &paquete.IDSede, &paquete.IDTipoTour, &paquete.Nombre,
			&paquete.Descripcion, &paquete.PrecioTotal, &paquete.CantidadTotal, &paquete.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		paquetes = append(paquetes, paquete)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return paquetes, nil
}

// List lista todos los paquetes de pasajes
func (r *PaquetePasajesRepository) List() ([]*entidades.PaquetePasajes, error) {
	query := `SELECT id_paquete, id_sede, id_tipo_tour, nombre, descripcion, precio_total, cantidad_total, eliminado
              FROM paquete_pasajes
              WHERE eliminado = false
              ORDER BY id_sede, precio_total ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paquetes := []*entidades.PaquetePasajes{}

	for rows.Next() {
		paquete := &entidades.PaquetePasajes{}
		err := rows.Scan(
			&paquete.ID, &paquete.IDSede, &paquete.IDTipoTour, &paquete.Nombre,
			&paquete.Descripcion, &paquete.PrecioTotal, &paquete.CantidadTotal, &paquete.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		paquetes = append(paquetes, paquete)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return paquetes, nil
}
