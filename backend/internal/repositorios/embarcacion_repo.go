package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// EmbarcacionRepository maneja las operaciones de base de datos para embarcaciones
type EmbarcacionRepository struct {
	db *sql.DB
}

// NewEmbarcacionRepository crea una nueva instancia del repositorio
func NewEmbarcacionRepository(db *sql.DB) *EmbarcacionRepository {
	return &EmbarcacionRepository{
		db: db,
	}
}

// GetByID obtiene una embarcación por su ID
func (r *EmbarcacionRepository) GetByID(id int) (*entidades.Embarcacion, error) {
	embarcacion := &entidades.Embarcacion{}
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, estado
              FROM embarcacion 
              WHERE id_embarcacion = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
		&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.Estado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("embarcación no encontrada")
		}
		return nil, err
	}

	return embarcacion, nil
}

// GetByNombre obtiene una embarcación por su nombre
func (r *EmbarcacionRepository) GetByNombre(nombre string) (*entidades.Embarcacion, error) {
	embarcacion := &entidades.Embarcacion{}
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, estado
              FROM embarcacion
              WHERE nombre = $1 AND eliminado = false`

	err := r.db.QueryRow(query, nombre).Scan(
		&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
		&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.Estado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("embarcación no encontrada")
		}
		return nil, err
	}

	return embarcacion, nil
}

// Create guarda una nueva embarcación en la base de datos
func (r *EmbarcacionRepository) Create(embarcacion *entidades.NuevaEmbarcacionRequest) (int, error) {
	var id int
	query := `INSERT INTO embarcacion (id_sede, nombre, capacidad, descripcion, estado, eliminado)
              VALUES ($1, $2, $3, $4, $5, false)
              RETURNING id_embarcacion`

	err := r.db.QueryRow(
		query,
		embarcacion.IDSede,
		embarcacion.Nombre,
		embarcacion.Capacidad,
		embarcacion.Descripcion,
		embarcacion.Estado,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de una embarcación
func (r *EmbarcacionRepository) Update(id int, embarcacion *entidades.ActualizarEmbarcacionRequest) error {
	query := `UPDATE embarcacion SET
              id_sede = $1,
              nombre = $2,
              capacidad = $3,
              descripcion = $4,
              estado = $5
              WHERE id_embarcacion = $6 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		embarcacion.IDSede,
		embarcacion.Nombre,
		embarcacion.Capacidad,
		embarcacion.Descripcion,
		embarcacion.Estado,
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
		return errors.New("embarcación no encontrada o ya fue eliminada")
	}

	return nil
}

// SoftDelete marca una embarcación como eliminada (borrado lógico)
func (r *EmbarcacionRepository) SoftDelete(id int) error {
	query := `UPDATE embarcacion SET eliminado = true WHERE id_embarcacion = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("embarcación no encontrada o ya fue eliminada")
	}

	return nil
}

// List lista todas las embarcaciones no eliminadas
func (r *EmbarcacionRepository) List() ([]*entidades.Embarcacion, error) {
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, estado
              FROM embarcacion 
              WHERE eliminado = false
              ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embarcaciones := []*entidades.Embarcacion{}

	for rows.Next() {
		embarcacion := &entidades.Embarcacion{}
		err := rows.Scan(
			&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
			&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.Estado,
		)
		if err != nil {
			return nil, err
		}
		embarcaciones = append(embarcaciones, embarcacion)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return embarcaciones, nil
}

// ListBySede lista todas las embarcaciones de una sede específica
func (r *EmbarcacionRepository) ListBySede(idSede int) ([]*entidades.Embarcacion, error) {
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, estado
              FROM embarcacion 
              WHERE id_sede = $1 AND eliminado = false
              ORDER BY nombre`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embarcaciones := []*entidades.Embarcacion{}

	for rows.Next() {
		embarcacion := &entidades.Embarcacion{}
		err := rows.Scan(
			&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
			&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.Estado,
		)
		if err != nil {
			return nil, err
		}
		embarcaciones = append(embarcaciones, embarcacion)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return embarcaciones, nil
}

// ListByEstado lista todas las embarcaciones por estado
func (r *EmbarcacionRepository) ListByEstado(estado string) ([]*entidades.Embarcacion, error) {
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, estado
              FROM embarcacion 
              WHERE estado = $1 AND eliminado = false
              ORDER BY nombre`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embarcaciones := []*entidades.Embarcacion{}

	for rows.Next() {
		embarcacion := &entidades.Embarcacion{}
		err := rows.Scan(
			&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
			&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.Estado,
		)
		if err != nil {
			return nil, err
		}
		embarcaciones = append(embarcaciones, embarcacion)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return embarcaciones, nil
}
