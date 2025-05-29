package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// TipoTourRepository maneja las operaciones de base de datos para tipos de tour
type TipoTourRepository struct {
	db *sql.DB
}

// NewTipoTourRepository crea una nueva instancia del repositorio
func NewTipoTourRepository(db *sql.DB) *TipoTourRepository {
	return &TipoTourRepository{
		db: db,
	}
}

// GetByID obtiene un tipo de tour por su ID
func (r *TipoTourRepository) GetByID(id int) (*entidades.TipoTour, error) {
	tipoTour := &entidades.TipoTour{}
	query := `SELECT t.id_tipo_tour, t.id_sede, t.nombre, t.descripcion, 
              t.duracion_minutos, t.url_imagen, t.eliminado, s.nombre as nombre_sede
              FROM tipo_tour t
              INNER JOIN sede s ON t.id_sede = s.id_sede
              WHERE t.id_tipo_tour = $1`

	var descripcion sql.NullString
	var urlImagen sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&tipoTour.ID, &tipoTour.IDSede, &tipoTour.Nombre, &descripcion,
		&tipoTour.DuracionMinutos, &urlImagen, &tipoTour.Eliminado, &tipoTour.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de tour no encontrado")
		}
		return nil, err
	}

	tipoTour.Descripcion = descripcion
	tipoTour.URLImagen = urlImagen

	return tipoTour, nil
}

// GetByNombre obtiene un tipo de tour por su nombre
func (r *TipoTourRepository) GetByNombre(nombre string, idSede int) (*entidades.TipoTour, error) {
	tipoTour := &entidades.TipoTour{}
	query := `SELECT id_tipo_tour, id_sede, nombre, descripcion, 
              duracion_minutos, url_imagen, eliminado
              FROM tipo_tour
              WHERE nombre = $1 AND id_sede = $2`

	var descripcion sql.NullString
	var urlImagen sql.NullString

	err := r.db.QueryRow(query, nombre, idSede).Scan(
		&tipoTour.ID, &tipoTour.IDSede, &tipoTour.Nombre, &descripcion,
		&tipoTour.DuracionMinutos, &urlImagen, &tipoTour.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de tour no encontrado")
		}
		return nil, err
	}

	tipoTour.Descripcion = descripcion
	tipoTour.URLImagen = urlImagen

	return tipoTour, nil
}

// Create guarda un nuevo tipo de tour en la base de datos
func (r *TipoTourRepository) Create(tipoTour *entidades.NuevoTipoTourRequest) (int, error) {
	var id int
	query := `INSERT INTO tipo_tour (id_sede, nombre, descripcion, duracion_minutos, 
              url_imagen, eliminado)
              VALUES ($1, $2, $3, $4, $5, false)
              RETURNING id_tipo_tour`

	err := r.db.QueryRow(
		query,
		tipoTour.IDSede,
		tipoTour.Nombre,
		tipoTour.Descripcion,
		tipoTour.DuracionMinutos,
		tipoTour.URLImagen,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un tipo de tour
func (r *TipoTourRepository) Update(id int, tipoTour *entidades.ActualizarTipoTourRequest) error {
	query := `UPDATE tipo_tour SET
              id_sede = $1,
              nombre = $2,
              descripcion = $3,
              duracion_minutos = $4,
              url_imagen = $5,
              eliminado = $6
              WHERE id_tipo_tour = $7`

	_, err := r.db.Exec(
		query,
		tipoTour.IDSede,
		tipoTour.Nombre,
		tipoTour.Descripcion,
		tipoTour.DuracionMinutos,
		tipoTour.URLImagen,
		tipoTour.Eliminado,
		id,
	)

	return err
}

// Delete marca un tipo de tour como eliminado (borrado lógico)
func (r *TipoTourRepository) Delete(id int) error {
	query := `UPDATE tipo_tour SET eliminado = true WHERE id_tipo_tour = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List lista todos los tipos de tour no eliminados
func (r *TipoTourRepository) List() ([]*entidades.TipoTour, error) {
	query := `SELECT t.id_tipo_tour, t.id_sede, t.nombre, t.descripcion, 
              t.duracion_minutos, t.url_imagen, t.eliminado, s.nombre as nombre_sede
              FROM tipo_tour t
              INNER JOIN sede s ON t.id_sede = s.id_sede
              WHERE t.eliminado = false
              ORDER BY t.nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposTour := []*entidades.TipoTour{}

	for rows.Next() {
		tipoTour := &entidades.TipoTour{}
		var descripcion sql.NullString
		var urlImagen sql.NullString

		err := rows.Scan(
			&tipoTour.ID, &tipoTour.IDSede, &tipoTour.Nombre, &descripcion,
			&tipoTour.DuracionMinutos, &urlImagen, &tipoTour.Eliminado, &tipoTour.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		tipoTour.Descripcion = descripcion
		tipoTour.URLImagen = urlImagen
		tiposTour = append(tiposTour, tipoTour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tiposTour, nil
}

// ListBySede lista todos los tipos de tour de una sede específica
func (r *TipoTourRepository) ListBySede(idSede int) ([]*entidades.TipoTour, error) {
	query := `SELECT t.id_tipo_tour, t.id_sede, t.nombre, t.descripcion, 
              t.duracion_minutos, t.url_imagen, t.eliminado, s.nombre as nombre_sede
              FROM tipo_tour t
              INNER JOIN sede s ON t.id_sede = s.id_sede
              WHERE t.id_sede = $1 AND t.eliminado = false
              ORDER BY t.nombre`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposTour := []*entidades.TipoTour{}

	for rows.Next() {
		tipoTour := &entidades.TipoTour{}
		var descripcion sql.NullString
		var urlImagen sql.NullString

		err := rows.Scan(
			&tipoTour.ID, &tipoTour.IDSede, &tipoTour.Nombre, &descripcion,
			&tipoTour.DuracionMinutos, &urlImagen, &tipoTour.Eliminado, &tipoTour.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		tipoTour.Descripcion = descripcion
		tipoTour.URLImagen = urlImagen
		tiposTour = append(tiposTour, tipoTour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tiposTour, nil
}
