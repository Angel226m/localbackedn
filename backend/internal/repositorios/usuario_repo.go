package repositorios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
)

// UsuarioRepository maneja las operaciones de base de datos para usuarios
type UsuarioRepository struct {
	db *sql.DB
}

// NewUsuarioRepository crea una nueva instancia del repositorio
func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{
		db: db,
	}
}

// GetByID obtiene un usuario por su ID
func (r *UsuarioRepository) GetByID(id int) (*entidades.Usuario, error) {
	usuario := &entidades.Usuario{}
	query := `SELECT id_usuario, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, estado 
              FROM usuario 
              WHERE id_usuario = $1 AND estado = true`

	err := r.db.QueryRow(query, id).Scan(
		&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
		&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
		&usuario.FechaRegistro, &usuario.Estado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return usuario, nil
}

// GetByEmail obtiene un usuario por su correo electrónico
func (r *UsuarioRepository) GetByEmail(correo string) (*entidades.Usuario, error) {
	usuario := &entidades.Usuario{}
	query := `SELECT id_usuario, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, estado 
              FROM usuario 
              WHERE correo = $1`

	err := r.db.QueryRow(query, correo).Scan(
		&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
		&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
		&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Estado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return usuario, nil
}

// GetByDocumento obtiene un usuario por su número de documento
func (r *UsuarioRepository) GetByDocumento(tipo, numero string) (*entidades.Usuario, error) {
	usuario := &entidades.Usuario{}
	query := `SELECT id_usuario, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, estado 
              FROM usuario 
              WHERE tipo_de_documento = $1 AND numero_documento = $2`

	err := r.db.QueryRow(query, tipo, numero).Scan(
		&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
		&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
		&usuario.FechaRegistro, &usuario.Estado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return usuario, nil
}

// Create guarda un nuevo usuario en la base de datos
func (r *UsuarioRepository) Create(usuario *entidades.NuevoUsuarioRequest, hashedPassword string) (int, error) {
	var id int
	query := `INSERT INTO usuario (nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, contrasena, estado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
              RETURNING id_usuario`

	err := r.db.QueryRow(
		query,
		usuario.Nombres,
		usuario.Apellidos,
		usuario.Correo,
		usuario.Telefono,
		usuario.Direccion,
		usuario.FechaNacimiento,
		usuario.Rol,
		usuario.Nacionalidad,
		usuario.TipoDocumento,
		usuario.NumeroDocumento,
		hashedPassword,
		true,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un usuario
func (r *UsuarioRepository) Update(usuario *entidades.Usuario) error {
	query := `UPDATE usuario SET 
              nombres = $1, 
              apellidos = $2, 
              correo = $3, 
              telefono = $4, 
              direccion = $5, 
              fecha_nacimiento = $6, 
              rol = $7, 
              nacionalidad = $8, 
              tipo_de_documento = $9, 
              numero_documento = $10, 
              estado = $11 
              WHERE id_usuario = $12`

	_, err := r.db.Exec(
		query,
		usuario.Nombres,
		usuario.Apellidos,
		usuario.Correo,
		usuario.Telefono,
		usuario.Direccion,
		usuario.FechaNacimiento,
		usuario.Rol,
		usuario.Nacionalidad,
		usuario.TipoDocumento,
		usuario.NumeroDocumento,
		usuario.Estado,
		usuario.ID,
	)

	return err
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *UsuarioRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `UPDATE usuario SET contrasena = $1 WHERE id_usuario = $2`
	_, err := r.db.Exec(query, hashedPassword, id)
	return err
}

// Delete marca un usuario como inactivo (borrado lógico)
func (r *UsuarioRepository) Delete(id int) error {
	query := `UPDATE usuario SET estado = false WHERE id_usuario = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListByRol lista todos los usuarios con un rol específico
func (r *UsuarioRepository) ListByRol(rol string) ([]*entidades.Usuario, error) {
	query := `SELECT id_usuario, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, estado 
              FROM usuario 
              WHERE rol = $1 AND estado = true 
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query, rol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}

	for rows.Next() {
		usuario := &entidades.Usuario{}
		err := rows.Scan(
			&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Estado,
		)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

// List lista todos los usuarios activos
func (r *UsuarioRepository) List() ([]*entidades.Usuario, error) {
	query := `SELECT id_usuario, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, estado 
              FROM usuario 
              WHERE estado = true 
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}

	for rows.Next() {
		usuario := &entidades.Usuario{}
		err := rows.Scan(
			&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Estado,
		)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}
