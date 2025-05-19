package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
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
	query := `SELECT id_tipo_tour, nombre, descripcion, duracion_minutos, precio_base, 
              cantidad_pasajeros, url_imagen 
              FROM tipo_tour 
              WHERE id_tipo_tour = $1`

	err := r.db.QueryRow(query, id).Scan(
		&tipoTour.ID, &tipoTour.Nombre, &tipoTour.Descripcion,
		&tipoTour.DuracionMinutos, &tipoTour.PrecioBase,
		&tipoTour.CantidadPasajeros, &tipoTour.URLImagen,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de tour no encontrado")
		}
		return nil, err
	}

	return tipoTour, nil
}

// GetByNombre obtiene un tipo de tour por su nombre
func (r *TipoTourRepository) GetByNombre(nombre string) (*entidades.TipoTour, error) {
	tipoTour := &entidades.TipoTour{}
	query := `SELECT id_tipo_tour, nombre, descripcion, duracion_minutos, precio_base, 
              cantidad_pasajeros, url_imagen 
              FROM tipo_tour 
              WHERE nombre = $1`

	err := r.db.QueryRow(query, nombre).Scan(
		&tipoTour.ID, &tipoTour.Nombre, &tipoTour.Descripcion,
		&tipoTour.DuracionMinutos, &tipoTour.PrecioBase,
		&tipoTour.CantidadPasajeros, &tipoTour.URLImagen,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de tour no encontrado")
		}
		return nil, err
	}

	return tipoTour, nil
}

// Create guarda un nuevo tipo de tour en la base de datos
func (r *TipoTourRepository) Create(tipoTour *entidades.NuevoTipoTourRequest) (int, error) {
	var id int
	query := `INSERT INTO tipo_tour (nombre, descripcion, duracion_minutos, precio_base, 
              cantidad_pasajeros, url_imagen) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id_tipo_tour`

	err := r.db.QueryRow(
		query,
		tipoTour.Nombre,
		tipoTour.Descripcion,
		tipoTour.DuracionMinutos,
		tipoTour.PrecioBase,
		tipoTour.CantidadPasajeros,
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
              nombre = $1, 
              descripcion = $2, 
              duracion_minutos = $3, 
              precio_base = $4, 
              cantidad_pasajeros = $5, 
              url_imagen = $6 
              WHERE id_tipo_tour = $7`

	_, err := r.db.Exec(
		query,
		tipoTour.Nombre,
		tipoTour.Descripcion,
		tipoTour.DuracionMinutos,
		tipoTour.PrecioBase,
		tipoTour.CantidadPasajeros,
		tipoTour.URLImagen,
		id,
	)

	return err
}

// Delete elimina un tipo de tour
func (r *TipoTourRepository) Delete(id int) error {
	// Primero verificamos si hay horarios de tour que dependen de este tipo_tour
	var count int
	queryCheck := `SELECT COUNT(*) FROM horario_tour WHERE id_tipo_tour = $1`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este tipo de tour porque hay horarios que dependen de él")
	}

	// Ahora verificamos si hay tours programados que dependen de este tipo_tour
	queryCheckTours := `SELECT COUNT(*) FROM tour_programado WHERE id_tipo_tour = $1`
	err = r.db.QueryRow(queryCheckTours, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este tipo de tour porque hay tours programados que dependen de él")
	}

	// Si no hay dependencias, procedemos a eliminar
	query := `DELETE FROM tipo_tour WHERE id_tipo_tour = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los tipos de tour
func (r *TipoTourRepository) List() ([]*entidades.TipoTour, error) {
	query := `SELECT id_tipo_tour, nombre, descripcion, duracion_minutos, precio_base, 
              cantidad_pasajeros, url_imagen 
              FROM tipo_tour 
              ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposTour := []*entidades.TipoTour{}

	for rows.Next() {
		tipoTour := &entidades.TipoTour{}
		err := rows.Scan(
			&tipoTour.ID, &tipoTour.Nombre, &tipoTour.Descripcion,
			&tipoTour.DuracionMinutos, &tipoTour.PrecioBase,
			&tipoTour.CantidadPasajeros, &tipoTour.URLImagen,
		)
		if err != nil {
			return nil, err
		}
		tiposTour = append(tiposTour, tipoTour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tiposTour, nil
}
