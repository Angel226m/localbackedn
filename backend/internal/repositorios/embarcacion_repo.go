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
	query := `SELECT e.id_embarcacion, e.id_sede, e.nombre, e.capacidad, e.descripcion, 
              e.eliminado, e.id_usuario, e.estado,
              u.nombres, u.apellidos, u.numero_documento, u.telefono
              FROM embarcacion e
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              WHERE e.id_embarcacion = $1`

	err := r.db.QueryRow(query, id).Scan(
		&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
		&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.IDUsuario, &embarcacion.Estado,
		&embarcacion.NombreChofer, &embarcacion.ApellidosChofer,
		&embarcacion.DocumentoChofer, &embarcacion.TelefonoChofer,
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
	query := `SELECT id_embarcacion, id_sede, nombre, capacidad, descripcion, eliminado, id_usuario, estado
              FROM embarcacion
              WHERE nombre = $1`

	err := r.db.QueryRow(query, nombre).Scan(
		&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
		&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.IDUsuario, &embarcacion.Estado,
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
	query := `INSERT INTO embarcacion (id_sede, nombre, capacidad, descripcion, id_usuario, estado, eliminado)
              VALUES ($1, $2, $3, $4, $5, $6, false)
              RETURNING id_embarcacion`

	err := r.db.QueryRow(
		query,
		embarcacion.IDSede,
		embarcacion.Nombre,
		embarcacion.Capacidad,
		embarcacion.Descripcion,
		embarcacion.IDUsuario,
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
              id_usuario = $5,
              estado = $6,
              eliminado = $7
              WHERE id_embarcacion = $8`

	_, err := r.db.Exec(
		query,
		embarcacion.IDSede,
		embarcacion.Nombre,
		embarcacion.Capacidad,
		embarcacion.Descripcion,
		embarcacion.IDUsuario,
		embarcacion.Estado,
		embarcacion.Eliminado,
		id,
	)

	return err
}

// Delete marca una embarcación como eliminada (borrado lógico)
func (r *EmbarcacionRepository) Delete(id int) error {
	query := `UPDATE embarcacion SET eliminado = true WHERE id_embarcacion = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List lista todas las embarcaciones no eliminadas con información del chofer
func (r *EmbarcacionRepository) List() ([]*entidades.Embarcacion, error) {
	query := `SELECT e.id_embarcacion, e.id_sede, e.nombre, e.capacidad, e.descripcion, 
              e.eliminado, e.id_usuario, e.estado,
              u.nombres, u.apellidos, u.numero_documento, u.telefono
              FROM embarcacion e
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              WHERE e.eliminado = false
              ORDER BY e.nombre`

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
			&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.IDUsuario, &embarcacion.Estado,
			&embarcacion.NombreChofer, &embarcacion.ApellidosChofer,
			&embarcacion.DocumentoChofer, &embarcacion.TelefonoChofer,
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

// ListByChofer lista todas las embarcaciones asignadas a un chofer específico
func (r *EmbarcacionRepository) ListByChofer(idChofer int) ([]*entidades.Embarcacion, error) {
	query := `SELECT e.id_embarcacion, e.id_sede, e.nombre, e.capacidad, e.descripcion, 
              e.eliminado, e.id_usuario, e.estado,
              u.nombres, u.apellidos, u.numero_documento, u.telefono
              FROM embarcacion e
              INNER JOIN usuario u ON e.id_usuario = u.id_usuario
              WHERE e.id_usuario = $1 AND e.eliminado = false
              ORDER BY e.nombre`

	rows, err := r.db.Query(query, idChofer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embarcaciones := []*entidades.Embarcacion{}

	for rows.Next() {
		embarcacion := &entidades.Embarcacion{}
		err := rows.Scan(
			&embarcacion.ID, &embarcacion.IDSede, &embarcacion.Nombre, &embarcacion.Capacidad,
			&embarcacion.Descripcion, &embarcacion.Eliminado, &embarcacion.IDUsuario, &embarcacion.Estado,
			&embarcacion.NombreChofer, &embarcacion.ApellidosChofer,
			&embarcacion.DocumentoChofer, &embarcacion.TelefonoChofer,
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

//cambiar si fallan
