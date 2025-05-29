package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// SedeRepository maneja las operaciones de base de datos para sedes
type SedeRepository struct {
	db *sql.DB
}

// NewSedeRepository crea una nueva instancia del repositorio
func NewSedeRepository(db *sql.DB) *SedeRepository {
	return &SedeRepository{
		db: db,
	}
}

// GetByID obtiene una sede por su ID
func (r *SedeRepository) GetByID(id int) (*entidades.Sede, error) {
	sede := &entidades.Sede{}
	query := `SELECT id_sede, nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado 
              FROM sede 
              WHERE id_sede = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&sede.ID, &sede.Nombre, &sede.Direccion, &sede.Telefono,
		&sede.Correo, &sede.Distrito, &sede.Provincia, &sede.Pais, &sede.ImageURL, &sede.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sede no encontrada")
		}
		return nil, err
	}

	return sede, nil
}

// Create guarda una nueva sede en la base de datos
func (r *SedeRepository) Create(sede *entidades.NuevaSedeRequest) (int, error) {
	var id int
	query := `INSERT INTO sede (nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
              RETURNING id_sede`

	err := r.db.QueryRow(
		query,
		sede.Nombre,
		sede.Direccion,
		sede.Telefono,
		sede.Correo,
		sede.Distrito,
		sede.Provincia,
		sede.Pais,
		sede.ImageURL,
		false, // No eliminado por defecto
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de una sede
func (r *SedeRepository) Update(id int, sede *entidades.ActualizarSedeRequest) error {
	query := `UPDATE sede SET 
              nombre = $1, 
              direccion = $2, 
              telefono = $3, 
              correo = $4, 
              distrito = $5, 
              provincia = $6, 
              pais = $7,
              image_url = $8
              WHERE id_sede = $9 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		sede.Nombre,
		sede.Direccion,
		sede.Telefono,
		sede.Correo,
		sede.Distrito,
		sede.Provincia,
		sede.Pais,
		sede.ImageURL,
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
		return errors.New("sede no encontrada o ya fue eliminada")
	}

	return nil
}

// SoftDelete marca una sede como eliminada (borrado lógico)
func (r *SedeRepository) SoftDelete(id int) error {
	query := `UPDATE sede SET eliminado = true WHERE id_sede = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("sede no encontrada o ya fue eliminada")
	}

	return nil
}

// Restore restaura una sede eliminada lógicamente
func (r *SedeRepository) Restore(id int) error {
	query := `UPDATE sede SET eliminado = false WHERE id_sede = $1 AND eliminado = true`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("sede no encontrada o no está eliminada")
	}

	return nil
}

// List lista todas las sedes activas
func (r *SedeRepository) List() ([]*entidades.Sede, error) {
	query := `SELECT id_sede, nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado 
              FROM sede 
              WHERE eliminado = false 
              ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sedes := []*entidades.Sede{}

	for rows.Next() {
		sede := &entidades.Sede{}
		err := rows.Scan(
			&sede.ID, &sede.Nombre, &sede.Direccion, &sede.Telefono,
			&sede.Correo, &sede.Distrito, &sede.Provincia, &sede.Pais, &sede.ImageURL, &sede.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		sedes = append(sedes, sede)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sedes, nil
}

// GetByDistrito obtiene sedes por distrito
func (r *SedeRepository) GetByDistrito(distrito string) ([]*entidades.Sede, error) {
	query := `SELECT id_sede, nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado 
              FROM sede 
              WHERE distrito = $1 AND eliminado = false 
              ORDER BY nombre`

	rows, err := r.db.Query(query, distrito)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sedes := []*entidades.Sede{}

	for rows.Next() {
		sede := &entidades.Sede{}
		err := rows.Scan(
			&sede.ID, &sede.Nombre, &sede.Direccion, &sede.Telefono,
			&sede.Correo, &sede.Distrito, &sede.Provincia, &sede.Pais, &sede.ImageURL, &sede.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		sedes = append(sedes, sede)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sedes, nil
}

// GetByPais obtiene sedes por país
func (r *SedeRepository) GetByPais(pais string) ([]*entidades.Sede, error) {
	query := `SELECT id_sede, nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado 
              FROM sede 
              WHERE pais = $1 AND eliminado = false 
              ORDER BY distrito, nombre`

	rows, err := r.db.Query(query, pais)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sedes := []*entidades.Sede{}

	for rows.Next() {
		sede := &entidades.Sede{}
		err := rows.Scan(
			&sede.ID, &sede.Nombre, &sede.Direccion, &sede.Telefono,
			&sede.Correo, &sede.Distrito, &sede.Provincia, &sede.Pais, &sede.ImageURL, &sede.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		sedes = append(sedes, sede)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sedes, nil
}

// GetAll obtiene todas las sedes no eliminadas
func (r *SedeRepository) GetAll() ([]*entidades.Sede, error) {
	query := `
		SELECT id_sede, nombre, direccion, telefono, correo, distrito, provincia, pais, image_url, eliminado
		FROM sede
		WHERE eliminado = false
		ORDER BY nombre ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sedes []*entidades.Sede
	for rows.Next() {
		var sede entidades.Sede
		err := rows.Scan(
			&sede.ID,
			&sede.Nombre,
			&sede.Direccion,
			&sede.Telefono,
			&sede.Correo,
			&sede.Distrito,
			&sede.Provincia,
			&sede.Pais,
			&sede.ImageURL,
			&sede.Eliminado,
		)
		if err != nil {
			return nil, err
		}
		sedes = append(sedes, &sede)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sedes, nil
}
