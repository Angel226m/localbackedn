package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// ClienteRepository maneja las operaciones de base de datos para clientes
type ClienteRepository struct {
	db *sql.DB
}

// NewClienteRepository crea una nueva instancia del repositorio
func NewClienteRepository(db *sql.DB) *ClienteRepository {
	return &ClienteRepository{
		db: db,
	}
}

// GetByID obtiene un cliente por su ID
func (r *ClienteRepository) GetByID(id int) (*entidades.Cliente, error) {
	cliente := &entidades.Cliente{}
	query := `SELECT id_cliente, tipo_documento, numero_documento, nombres, apellidos, correo, eliminado
              FROM cliente
              WHERE id_cliente = $1 AND eliminado = false`

	err := r.db.QueryRow(query, id).Scan(
		&cliente.ID, &cliente.TipoDocumento, &cliente.NumeroDocumento,
		&cliente.Nombres, &cliente.Apellidos, &cliente.Correo, &cliente.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cliente no encontrado")
		}
		return nil, err
	}

	// Establecer nombre completo
	cliente.NombreCompleto = cliente.Nombres + " " + cliente.Apellidos

	return cliente, nil
}

// GetByDocumento obtiene un cliente por tipo y número de documento
func (r *ClienteRepository) GetByDocumento(tipoDocumento, numeroDocumento string) (*entidades.Cliente, error) {
	cliente := &entidades.Cliente{}
	query := `SELECT id_cliente, tipo_documento, numero_documento, nombres, apellidos, correo, eliminado
              FROM cliente
              WHERE tipo_documento = $1 AND numero_documento = $2 AND eliminado = false`

	err := r.db.QueryRow(query, tipoDocumento, numeroDocumento).Scan(
		&cliente.ID, &cliente.TipoDocumento, &cliente.NumeroDocumento,
		&cliente.Nombres, &cliente.Apellidos, &cliente.Correo, &cliente.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cliente no encontrado")
		}
		return nil, err
	}

	// Establecer nombre completo
	cliente.NombreCompleto = cliente.Nombres + " " + cliente.Apellidos

	return cliente, nil
}

// GetByCorreo obtiene un cliente por su correo electrónico
func (r *ClienteRepository) GetByCorreo(correo string) (*entidades.Cliente, error) {
	cliente := &entidades.Cliente{}
	query := `SELECT id_cliente, tipo_documento, numero_documento, nombres, apellidos, correo, eliminado
              FROM cliente
              WHERE correo = $1 AND eliminado = false`

	err := r.db.QueryRow(query, correo).Scan(
		&cliente.ID, &cliente.TipoDocumento, &cliente.NumeroDocumento,
		&cliente.Nombres, &cliente.Apellidos, &cliente.Correo, &cliente.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cliente no encontrado")
		}
		return nil, err
	}

	// Establecer nombre completo
	cliente.NombreCompleto = cliente.Nombres + " " + cliente.Apellidos

	return cliente, nil
}

// GetPasswordByCorreo obtiene la contraseña de un cliente por su correo
func (r *ClienteRepository) GetPasswordByCorreo(correo string) (string, error) {
	var contrasena string
	query := `SELECT contrasena
              FROM cliente
              WHERE correo = $1 AND eliminado = false`

	err := r.db.QueryRow(query, correo).Scan(&contrasena)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("cliente no encontrado")
		}
		return "", err
	}

	return contrasena, nil
}

// Create guarda un nuevo cliente en la base de datos
func (r *ClienteRepository) Create(cliente *entidades.NuevoClienteRequest) (int, error) {
	var id int
	query := `INSERT INTO cliente (tipo_documento, numero_documento, nombres, apellidos, correo, contrasena, eliminado)
              VALUES ($1, $2, $3, $4, $5, $6, false)
              RETURNING id_cliente`

	err := r.db.QueryRow(
		query,
		cliente.TipoDocumento,
		cliente.NumeroDocumento,
		cliente.Nombres,
		cliente.Apellidos,
		cliente.Correo,
		cliente.Contrasena,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un cliente
func (r *ClienteRepository) Update(id int, cliente *entidades.ActualizarClienteRequest) error {
	query := `UPDATE cliente SET
              tipo_documento = $1,
              numero_documento = $2,
              nombres = $3,
              apellidos = $4,
              correo = $5
              WHERE id_cliente = $6 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		cliente.TipoDocumento,
		cliente.NumeroDocumento,
		cliente.Nombres,
		cliente.Apellidos,
		cliente.Correo,
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
		return errors.New("cliente no encontrado o ya eliminado")
	}

	return nil
}

// UpdatePassword actualiza la contraseña de un cliente
func (r *ClienteRepository) UpdatePassword(id int, contrasena string) error {
	query := `UPDATE cliente SET
              contrasena = $1
              WHERE id_cliente = $2 AND eliminado = false`

	result, err := r.db.Exec(query, contrasena, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("cliente no encontrado o ya eliminado")
	}

	return nil
}

// Delete marca un cliente como eliminado (eliminación lógica)
func (r *ClienteRepository) Delete(id int) error {
	// Verificar si hay reservas asociadas a este cliente
	var countReservas int
	queryCheckReservas := `SELECT COUNT(*) FROM reserva WHERE id_cliente = $1 AND eliminado = false`
	err := r.db.QueryRow(queryCheckReservas, id).Scan(&countReservas)
	if err != nil {
		return err
	}

	if countReservas > 0 {
		return errors.New("no se puede eliminar este cliente porque tiene reservas asociadas")
	}

	// Eliminación lógica del cliente
	query := `UPDATE cliente SET eliminado = true WHERE id_cliente = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("cliente no encontrado o ya eliminado")
	}

	return nil
}

// List lista todos los clientes no eliminados
func (r *ClienteRepository) List() ([]*entidades.Cliente, error) {
	query := `SELECT id_cliente, tipo_documento, numero_documento, nombres, apellidos, correo, eliminado
              FROM cliente
              WHERE eliminado = false
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientes := []*entidades.Cliente{}

	for rows.Next() {
		cliente := &entidades.Cliente{}
		err := rows.Scan(
			&cliente.ID, &cliente.TipoDocumento, &cliente.NumeroDocumento,
			&cliente.Nombres, &cliente.Apellidos, &cliente.Correo, &cliente.Eliminado,
		)
		if err != nil {
			return nil, err
		}

		// Establecer nombre completo
		cliente.NombreCompleto = cliente.Nombres + " " + cliente.Apellidos

		clientes = append(clientes, cliente)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}

// SearchByName busca clientes por nombre o apellido
func (r *ClienteRepository) SearchByName(query string) ([]*entidades.Cliente, error) {
	sqlQuery := `SELECT id_cliente, tipo_documento, numero_documento, nombres, apellidos, correo, eliminado
              FROM cliente
              WHERE (nombres ILIKE $1 OR apellidos ILIKE $1) AND eliminado = false
              ORDER BY apellidos, nombres`

	searchPattern := "%" + query + "%"

	rows, err := r.db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientes := []*entidades.Cliente{}

	for rows.Next() {
		cliente := &entidades.Cliente{}
		err := rows.Scan(
			&cliente.ID, &cliente.TipoDocumento, &cliente.NumeroDocumento,
			&cliente.Nombres, &cliente.Apellidos, &cliente.Correo, &cliente.Eliminado,
		)
		if err != nil {
			return nil, err
		}

		// Establecer nombre completo
		cliente.NombreCompleto = cliente.Nombres + " " + cliente.Apellidos

		clientes = append(clientes, cliente)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}
