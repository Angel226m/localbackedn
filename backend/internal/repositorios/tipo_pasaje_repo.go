package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
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
	query := `SELECT id_tipo_pasaje, id_sede, nombre, costo, edad, eliminado
              FROM tipo_pasaje
              WHERE id_tipo_pasaje = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&tipoPasaje.ID, &tipoPasaje.IDSede, &tipoPasaje.Nombre,
		&tipoPasaje.Costo, &tipoPasaje.Edad, &tipoPasaje.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tipo de pasaje no encontrado")
		}
		return nil, err
	}

	return tipoPasaje, nil
}

// GetByNombre obtiene un tipo de pasaje por su nombre y sede
func (r *TipoPasajeRepository) GetByNombre(nombre string, idSede int) (*entidades.TipoPasaje, error) {
	tipoPasaje := &entidades.TipoPasaje{}
	query := `SELECT id_tipo_pasaje, id_sede, nombre, costo, edad, eliminado
              FROM tipo_pasaje
              WHERE nombre = $1 AND id_sede = $2 AND eliminado = false`

	err := r.db.QueryRow(query, nombre, idSede).Scan(
		&tipoPasaje.ID, &tipoPasaje.IDSede, &tipoPasaje.Nombre,
		&tipoPasaje.Costo, &tipoPasaje.Edad, &tipoPasaje.Eliminado,
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
	query := `INSERT INTO tipo_pasaje (id_sede, nombre, costo, edad, eliminado)
              VALUES ($1, $2, $3, $4, false)
              RETURNING id_tipo_pasaje`

	err := r.db.QueryRow(
		query,
		tipoPasaje.IDSede,
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
              WHERE id_tipo_pasaje = $4 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		tipoPasaje.Nombre,
		tipoPasaje.Costo,
		tipoPasaje.Edad,
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
		return errors.New("tipo de pasaje no encontrado o ya eliminado")
	}

	return nil
}

// Delete marca un tipo de pasaje como eliminado (eliminación lógica)
func (r *TipoPasajeRepository) Delete(id int) error {
	// Verificar si hay pasajes que usan este tipo de pasaje
	var count int
	queryCheck := `SELECT COUNT(*) FROM pasajes_cantidad WHERE id_tipo_pasaje = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheck, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("no se puede eliminar este tipo de pasaje porque está siendo utilizado por reservas")
	}

	// Eliminación lógica del tipo de pasaje
	query := `UPDATE tipo_pasaje SET eliminado = true WHERE id_tipo_pasaje = $4 AND eliminado = false`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("tipo de pasaje no encontrado o ya eliminado")
	}

	return nil
}

// ListBySede lista todos los tipos de pasaje de una sede específica
func (r *TipoPasajeRepository) ListBySede(idSede int) ([]*entidades.TipoPasaje, error) {
	query := `SELECT id_tipo_pasaje, id_sede, nombre, costo, edad, eliminado
              FROM tipo_pasaje
              WHERE id_sede = $1 AND eliminado = false
              ORDER BY costo ASC`

	rows, err := r.db.Query(query, idSede)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposPasaje := []*entidades.TipoPasaje{}

	for rows.Next() {
		tipoPasaje := &entidades.TipoPasaje{}
		err := rows.Scan(
			&tipoPasaje.ID, &tipoPasaje.IDSede, &tipoPasaje.Nombre,
			&tipoPasaje.Costo, &tipoPasaje.Edad, &tipoPasaje.Eliminado,
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

// List lista todos los tipos de pasaje
func (r *TipoPasajeRepository) List() ([]*entidades.TipoPasaje, error) {
	query := `SELECT id_tipo_pasaje, id_sede, nombre, costo, edad, eliminado
              FROM tipo_pasaje
              WHERE eliminado = false
              ORDER BY id_sede, costo ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiposPasaje := []*entidades.TipoPasaje{}

	for rows.Next() {
		tipoPasaje := &entidades.TipoPasaje{}
		err := rows.Scan(
			&tipoPasaje.ID, &tipoPasaje.IDSede, &tipoPasaje.Nombre,
			&tipoPasaje.Costo, &tipoPasaje.Edad, &tipoPasaje.Eliminado,
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
