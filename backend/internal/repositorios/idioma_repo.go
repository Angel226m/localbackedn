package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// IdiomaRepository maneja las operaciones de base de datos para idiomas
type IdiomaRepository struct {
	db *sql.DB
}

// NewIdiomaRepository crea una nueva instancia del repositorio
func NewIdiomaRepository(db *sql.DB) *IdiomaRepository {
	return &IdiomaRepository{
		db: db,
	}
}

// GetByID obtiene un idioma por su ID
func (r *IdiomaRepository) GetByID(id int) (*entidades.Idioma, error) {
	idioma := &entidades.Idioma{}
	query := `SELECT id_idioma, nombre, eliminado FROM idioma WHERE id_idioma = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(&idioma.ID, &idioma.Nombre, &idioma.Eliminado)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("idioma no encontrado")
		}
		return nil, err
	}

	return idioma, nil
}

// GetByNombre obtiene un idioma por su nombre
func (r *IdiomaRepository) GetByNombre(nombre string) (*entidades.Idioma, error) {
	idioma := &entidades.Idioma{}
	query := `SELECT id_idioma, nombre, eliminado FROM idioma WHERE nombre = $1 AND eliminado = false`

	err := r.db.QueryRow(query, nombre).Scan(&idioma.ID, &idioma.Nombre, &idioma.Eliminado)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("idioma no encontrado")
		}
		return nil, err
	}

	return idioma, nil
}

// Create guarda un nuevo idioma en la base de datos
func (r *IdiomaRepository) Create(idioma *entidades.Idioma) (int, error) {
	var id int
	query := `INSERT INTO idioma (nombre, eliminado) VALUES ($1, false) RETURNING id_idioma`

	err := r.db.QueryRow(query, idioma.Nombre).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un idioma
func (r *IdiomaRepository) Update(idioma *entidades.Idioma) error {
	query := `UPDATE idioma SET nombre = $1 WHERE id_idioma = $2 AND eliminado = false`
	result, err := r.db.Exec(query, idioma.Nombre, idioma.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("idioma no encontrado o ya eliminado")
	}

	return nil
}

// SoftDelete marca un idioma como eliminado (soft delete)
func (r *IdiomaRepository) SoftDelete(id int) error {
	query := `UPDATE idioma SET eliminado = true WHERE id_idioma = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("idioma no encontrado o ya eliminado")
	}

	return nil
}

// Restore restaura un idioma eliminado
func (r *IdiomaRepository) Restore(id int) error {
	query := `UPDATE idioma SET eliminado = false WHERE id_idioma = $1 AND eliminado = true`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("idioma no encontrado o no está eliminado")
	}

	return nil
}

// List lista todos los idiomas activos
func (r *IdiomaRepository) List() ([]*entidades.Idioma, error) {
	query := `SELECT id_idioma, nombre, eliminado FROM idioma WHERE eliminado = false ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idiomas := []*entidades.Idioma{}

	for rows.Next() {
		idioma := &entidades.Idioma{}
		err := rows.Scan(&idioma.ID, &idioma.Nombre, &idioma.Eliminado)
		if err != nil {
			return nil, err
		}
		idiomas = append(idiomas, idioma)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idiomas, nil
}

// ListDeleted lista todos los idiomas eliminados (soft deleted)
func (r *IdiomaRepository) ListDeleted() ([]*entidades.Idioma, error) {
	query := `SELECT id_idioma, nombre, eliminado FROM idioma WHERE eliminado = true ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idiomas := []*entidades.Idioma{}

	for rows.Next() {
		idioma := &entidades.Idioma{}
		err := rows.Scan(&idioma.ID, &idioma.Nombre, &idioma.Eliminado)
		if err != nil {
			return nil, err
		}
		idiomas = append(idiomas, idioma)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idiomas, nil
}
