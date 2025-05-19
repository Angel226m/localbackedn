package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
)

// TipoPasajeRepository maneja las operaciones de base de datos para tipos de pasaje
type TipoPasajeRepository struct {
	db *sql.DB
}

// NewTipoPasajeRepository crea una nueva instancia del repositorio
func NewTipoPasajeRepository(db *sql.DB) *TipoPasajeRepository {
	return &TipoPasajeRepository{
		db: db,
	}
}

// GetByID obtiene un tipo de pasaje por su ID
func (r *TipoPasajeRepository) GetByID(id int) (*entidades.TipoPasaje, error) {
	tipoPasaje := &entidades.TipoPasaje{}
	query := `SELECT id_tipo_pasaje, nombre, costo, edad
              FROM tipo_pasaje
              WHERE id_tipo_pasaje = $1`

	err := r.db.QueryRow(query, id).Scan(
		&tipoPasaje.ID, &tipoPasaje.Nombre, &tipoPasaje.Costo, &tipoPasaje.Edad,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de pasaje no encontrado")
		}
		return nil, err
	}

	return tipoPasaje, nil
}

// GetByNombre obtiene un tipo de pasaje por su nombre
func (r *TipoPasajeRepository) GetByNombre(nombre string) (*entidades.TipoPasaje, error) {
	tipoPasaje := &entidades.TipoPasaje{}
	query := `SELECT id_tipo_pasaje, nombre, costo, edad
              FROM tipo_pasaje
              WHERE nombre = $1`

	err := r.db.QueryRow(query, nombre).Scan(
		&tipoPasaje.ID, &tipoPasaje.Nombre, &tipoPasaje.Costo, &tipoPasaje.Edad,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de pasaje no encontrado")
		}
		return nil, err
	}

	return tipoPasaje, nil
}

// Create guarda un nuevo tipo de pasaje en la base de datos
func (r *TipoPasajeRepository) Create(tipoPasaje *entidades.NuevoTipoPasajeRequest) (int, error) {
	var id int
	query := `INSERT INTO tipo_pasaje (nombre, costo, edad)
              VALUES ($1, $2, $3)
              RETURNING id_tipo_pasaje`

	err := r.db.QueryRow(
		query,
		tipoPasaje.Nombre,
		tipoPasaje.Costo,
		tipoPasaje.Edad,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un tipo de pasaje
func (r *TipoPasajeRepository) Update(id int, tipoPasaje *entidades.ActualizarTipoPasajeRequest) error {
	query := `UPDATE tipo_pasaje SET
              nombre = $1,
              costo = $2,
              edad = $3
              WHERE id_tipo_pasaje = $4`

	_, err := r.db.Exec(
		query,
		tipoPasaje.Nombre,
		tipoPasaje.Costo,
		tipoPasaje.Edad,
		id,
	)

	return err
}

// Delete elimina un tipo de pasaje
func (r *TipoPasajeRepository) Delete(id int) error {
	// Verificar si hay pasajeros que usan este tipo de pasaje
	var count int
	queryCheck := `SELECT COUNT(*) FROM pasajero WHERE id_tipo_pasaje = $1`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este tipo de pasaje porque está siendo utilizado por pasajeros")
	}

	// Eliminar tipo de pasaje
	query := `DELETE FROM tipo_pasaje WHERE id_tipo_pasaje = $1`
	_, err = r.db.Exec(query, id)
	return err
}

// List lista todos los tipos de pasaje
func (r *TipoPasajeRepository) List() ([]*entidades.TipoPasaje, error) {
	query := `SELECT id_tipo_pasaje, nombre, costo, edad
              FROM tipo_pasaje
              ORDER BY costo ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposPasaje := []*entidades.TipoPasaje{}

	for rows.Next() {
		tipoPasaje := &entidades.TipoPasaje{}
		err := rows.Scan(
			&tipoPasaje.ID, &tipoPasaje.Nombre, &tipoPasaje.Costo, &tipoPasaje.Edad,
		)
		if err != nil {
			return nil, err
		}
		tiposPasaje = append(tiposPasaje, tipoPasaje)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tiposPasaje, nil
}
