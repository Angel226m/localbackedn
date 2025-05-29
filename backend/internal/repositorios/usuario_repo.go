package repositorios

import (
	"database/sql"
	"errors"
	"fmt"
	"sistema-toursseft/internal/entidades"
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
	fmt.Printf("UsuarioRepository: Buscando usuario con ID %d\n", id)

	var usuario entidades.Usuario
	var idSedeNullable sql.NullInt64 // Para manejar NULL en id_sede

	query := `SELECT id_usuario, nombres, apellidos, correo, contrasena, rol, id_sede, 
              telefono, direccion, fecha_nacimiento, nacionalidad, tipo_de_documento, 
              numero_documento, fecha_registro, eliminado 
              FROM usuario 
              WHERE id_usuario = $1 AND eliminado = false`

	fmt.Printf("UsuarioRepository: Ejecutando consulta: %s\n", query)

	err := r.db.QueryRow(query, id).Scan(
		&usuario.ID, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Contrasena, &usuario.Rol, &idSedeNullable, &usuario.Telefono,
		&usuario.Direccion, &usuario.FechaNacimiento, &usuario.Nacionalidad,
		&usuario.TipoDocumento, &usuario.NumeroDocumento, &usuario.FechaRegistro,
		&usuario.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("UsuarioRepository: Usuario con ID %d no encontrado\n", id)
			return nil, errors.New("usuario no encontrado")
		}
		fmt.Printf("UsuarioRepository: Error en la consulta SQL: %v\n", err)
		return nil, err
	}

	// Convertir sql.NullInt64 a *int
	if idSedeNullable.Valid {
		idSedeInt := int(idSedeNullable.Int64)
		usuario.IdSede = &idSedeInt
	} else {
		usuario.IdSede = nil // Es NULL en la base de datos
	}

	// Obtener idiomas del usuario
	idiomasQuery := `
		SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
		       i.nombre, i.eliminado
		FROM usuario_idioma ui
		JOIN idioma i ON ui.id_idioma = i.id_idioma
		WHERE ui.id_usuario = $1 AND ui.eliminado = false
	`

	rows, err := r.db.Query(idiomasQuery, usuario.ID)
	if err == nil {
		defer rows.Close()
		usuario.Idiomas = []*entidades.UsuarioIdioma{}

		for rows.Next() {
			usuarioIdioma := &entidades.UsuarioIdioma{
				Idioma: &entidades.Idioma{},
			}
			var idiomaEliminado bool

			err := rows.Scan(
				&usuarioIdioma.ID,
				&usuarioIdioma.IDUsuario,
				&usuarioIdioma.IDIdioma,
				&usuarioIdioma.Nivel,
				&usuarioIdioma.Eliminado,
				&usuarioIdioma.Idioma.Nombre,
				&idiomaEliminado,
			)

			if err != nil {
				fmt.Printf("Error al escanear idiomas del usuario: %v\n", err)
				continue
			}

			usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
			usuarioIdioma.Idioma.Eliminado = idiomaEliminado
			usuario.Idiomas = append(usuario.Idiomas, usuarioIdioma)
		}
	}

	fmt.Printf("UsuarioRepository: Usuario encontrado con éxito - ID: %d, Rol: %s, IdSede: %v\n",
		usuario.ID, usuario.Rol, usuario.IdSede)

	return &usuario, nil
}

// GetByEmail obtiene un usuario por su correo electrónico
func (r *UsuarioRepository) GetByEmail(correo string) (*entidades.Usuario, error) {
	usuario := &entidades.Usuario{}
	var idSedeNullable sql.NullInt64

	query := `SELECT id_usuario, id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, eliminado 
              FROM usuario 
              WHERE correo = $1 AND eliminado = false`

	err := r.db.QueryRow(query, correo).Scan(
		&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
		&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
		&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Convertir sql.NullInt64 a *int
	if idSedeNullable.Valid {
		idSedeInt := int(idSedeNullable.Int64)
		usuario.IdSede = &idSedeInt
	} else {
		usuario.IdSede = nil
	}

	// Cargar los idiomas del usuario
	idiomasQuery := `
		SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
		       i.nombre, i.eliminado
		FROM usuario_idioma ui
		JOIN idioma i ON ui.id_idioma = i.id_idioma
		WHERE ui.id_usuario = $1 AND ui.eliminado = false
	`

	rows, err := r.db.Query(idiomasQuery, usuario.ID)
	if err == nil {
		defer rows.Close()
		usuario.Idiomas = []*entidades.UsuarioIdioma{}

		for rows.Next() {
			usuarioIdioma := &entidades.UsuarioIdioma{
				Idioma: &entidades.Idioma{},
			}
			var idiomaEliminado bool

			err := rows.Scan(
				&usuarioIdioma.ID,
				&usuarioIdioma.IDUsuario,
				&usuarioIdioma.IDIdioma,
				&usuarioIdioma.Nivel,
				&usuarioIdioma.Eliminado,
				&usuarioIdioma.Idioma.Nombre,
				&idiomaEliminado,
			)

			if err != nil {
				continue
			}

			usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
			usuarioIdioma.Idioma.Eliminado = idiomaEliminado
			usuario.Idiomas = append(usuario.Idiomas, usuarioIdioma)
		}
	}

	return usuario, nil
}

// GetByDocumento obtiene un usuario por su número de documento
func (r *UsuarioRepository) GetByDocumento(tipo, numero string) (*entidades.Usuario, error) {
	usuario := &entidades.Usuario{}
	var idSedeNullable sql.NullInt64

	query := `SELECT id_usuario, id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, eliminado 
              FROM usuario 
              WHERE tipo_de_documento = $1 AND numero_documento = $2 AND eliminado = false`

	err := r.db.QueryRow(query, tipo, numero).Scan(
		&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
		&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
		&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
		&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Convertir sql.NullInt64 a *int
	if idSedeNullable.Valid {
		idSedeInt := int(idSedeNullable.Int64)
		usuario.IdSede = &idSedeInt
	} else {
		usuario.IdSede = nil
	}

	// Cargar los idiomas del usuario
	idiomasQuery := `
		SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
		       i.nombre, i.eliminado
		FROM usuario_idioma ui
		JOIN idioma i ON ui.id_idioma = i.id_idioma
		WHERE ui.id_usuario = $1 AND ui.eliminado = false
	`

	rows, err := r.db.Query(idiomasQuery, usuario.ID)
	if err == nil {
		defer rows.Close()
		usuario.Idiomas = []*entidades.UsuarioIdioma{}

		for rows.Next() {
			usuarioIdioma := &entidades.UsuarioIdioma{
				Idioma: &entidades.Idioma{},
			}
			var idiomaEliminado bool

			err := rows.Scan(
				&usuarioIdioma.ID,
				&usuarioIdioma.IDUsuario,
				&usuarioIdioma.IDIdioma,
				&usuarioIdioma.Nivel,
				&usuarioIdioma.Eliminado,
				&usuarioIdioma.Idioma.Nombre,
				&idiomaEliminado,
			)

			if err != nil {
				continue
			}

			usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
			usuarioIdioma.Idioma.Eliminado = idiomaEliminado
			usuario.Idiomas = append(usuario.Idiomas, usuarioIdioma)
		}
	}

	return usuario, nil
}

// Create guarda un nuevo usuario en la base de datos
func (r *UsuarioRepository) Create(usuario *entidades.NuevoUsuarioRequest, hashedPassword string) (int, error) {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int
	query := `INSERT INTO usuario (id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, contrasena, eliminado) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false) 
              RETURNING id_usuario`

	err = tx.QueryRow(
		query,
		usuario.IdSede,
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
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Asignar idiomas si existen
	if len(usuario.IdiomasIDs) > 0 {
		for _, idiomaID := range usuario.IdiomasIDs {
			_, err = tx.Exec(`
				INSERT INTO usuario_idioma (id_usuario, id_idioma, nivel, eliminado)
				VALUES ($1, $2, 'básico', false)
			`, id, idiomaID)

			if err != nil {
				return 0, err
			}
		}
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update actualiza la información de un usuario
func (r *UsuarioRepository) Update(usuario *entidades.Usuario) error {
	query := `UPDATE usuario SET 
              id_sede = $1,
              nombres = $2, 
              apellidos = $3, 
              correo = $4, 
              telefono = $5, 
              direccion = $6, 
              fecha_nacimiento = $7, 
              rol = $8, 
              nacionalidad = $9, 
              tipo_de_documento = $10, 
              numero_documento = $11
              WHERE id_usuario = $12 AND eliminado = false`

	result, err := r.db.Exec(
		query,
		usuario.IdSede,
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
		usuario.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado o ya eliminado")
	}

	return nil
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *UsuarioRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `UPDATE usuario SET contrasena = $1 WHERE id_usuario = $2 AND eliminado = false`
	result, err := r.db.Exec(query, hashedPassword, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado o ya eliminado")
	}

	return nil
}

// SoftDelete marca un usuario como eliminado (soft delete)
func (r *UsuarioRepository) SoftDelete(id int) error {
	query := `UPDATE usuario SET eliminado = true WHERE id_usuario = $1 AND eliminado = false`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado o ya eliminado")
	}

	return nil
}

// Restore restaura un usuario eliminado
func (r *UsuarioRepository) Restore(id int) error {
	query := `UPDATE usuario SET eliminado = false WHERE id_usuario = $1 AND eliminado = true`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado o no está eliminado")
	}

	return nil
}

// ListByRol lista todos los usuarios con un rol específico
func (r *UsuarioRepository) ListByRol(rol string) ([]*entidades.Usuario, error) {
	query := `SELECT id_usuario, id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, eliminado 
              FROM usuario 
              WHERE rol = $1 AND eliminado = false 
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query, rol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}

	for rows.Next() {
		usuario := &entidades.Usuario{}
		var idSedeNullable sql.NullInt64

		err := rows.Scan(
			&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
		)
		if err != nil {
			return nil, err
		}

		// Convertir sql.NullInt64 a *int
		if idSedeNullable.Valid {
			idSedeInt := int(idSedeNullable.Int64)
			usuario.IdSede = &idSedeInt
		} else {
			usuario.IdSede = nil
		}

		// Cargar los idiomas de cada usuario (podría hacerse en batch para optimizar)
		idiomasQuery := `
			SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
				i.nombre, i.eliminado
			FROM usuario_idioma ui
			JOIN idioma i ON ui.id_idioma = i.id_idioma
			WHERE ui.id_usuario = $1 AND ui.eliminado = false
		`

		idiomasRows, err := r.db.Query(idiomasQuery, usuario.ID)
		if err == nil {
			defer idiomasRows.Close()
			usuario.Idiomas = []*entidades.UsuarioIdioma{}

			for idiomasRows.Next() {
				usuarioIdioma := &entidades.UsuarioIdioma{
					Idioma: &entidades.Idioma{},
				}
				var idiomaEliminado bool

				err := idiomasRows.Scan(
					&usuarioIdioma.ID,
					&usuarioIdioma.IDUsuario,
					&usuarioIdioma.IDIdioma,
					&usuarioIdioma.Nivel,
					&usuarioIdioma.Eliminado,
					&usuarioIdioma.Idioma.Nombre,
					&idiomaEliminado,
				)

				if err != nil {
					continue
				}

				usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
				usuarioIdioma.Idioma.Eliminado = idiomaEliminado
				usuario.Idiomas = append(usuario.Idiomas, usuarioIdioma)
			}
			idiomasRows.Close()
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
	query := `SELECT id_usuario, id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, eliminado 
              FROM usuario 
              WHERE eliminado = false 
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}

	for rows.Next() {
		usuario := &entidades.Usuario{}
		var idSedeNullable sql.NullInt64

		err := rows.Scan(
			&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
		)
		if err != nil {
			return nil, err
		}

		// Convertir sql.NullInt64 a *int
		if idSedeNullable.Valid {
			idSedeInt := int(idSedeNullable.Int64)
			usuario.IdSede = &idSedeInt
		} else {
			usuario.IdSede = nil
		}

		// Cargar los idiomas de cada usuario (podría hacerse en batch para optimizar)
		idiomasQuery := `
			SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
				i.nombre, i.eliminado
			FROM usuario_idioma ui
			JOIN idioma i ON ui.id_idioma = i.id_idioma
			WHERE ui.id_usuario = $1 AND ui.eliminado = false
		`

		idiomasRows, err := r.db.Query(idiomasQuery, usuario.ID)
		if err == nil {
			defer idiomasRows.Close()
			usuario.Idiomas = []*entidades.UsuarioIdioma{}

			for idiomasRows.Next() {
				usuarioIdioma := &entidades.UsuarioIdioma{
					Idioma: &entidades.Idioma{},
				}
				var idiomaEliminado bool

				err := idiomasRows.Scan(
					&usuarioIdioma.ID,
					&usuarioIdioma.IDUsuario,
					&usuarioIdioma.IDIdioma,
					&usuarioIdioma.Nivel,
					&usuarioIdioma.Eliminado,
					&usuarioIdioma.Idioma.Nombre,
					&idiomaEliminado,
				)

				if err != nil {
					continue
				}

				usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
				usuarioIdioma.Idioma.Eliminado = idiomaEliminado
				usuario.Idiomas = append(usuario.Idiomas, usuarioIdioma)
			}
			idiomasRows.Close()
		}

		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

// ListDeleted lista todos los usuarios eliminados (soft deleted)
func (r *UsuarioRepository) ListDeleted() ([]*entidades.Usuario, error) {
	query := `SELECT id_usuario, id_sede, nombres, apellidos, correo, telefono, direccion, 
              fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, 
              fecha_registro, contrasena, eliminado 
              FROM usuario 
              WHERE eliminado = true 
              ORDER BY apellidos, nombres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}

	for rows.Next() {
		usuario := &entidades.Usuario{}
		var idSedeNullable sql.NullInt64

		err := rows.Scan(
			&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
		)
		if err != nil {
			return nil, err
		}

		// Convertir sql.NullInt64 a *int
		if idSedeNullable.Valid {
			idSedeInt := int(idSedeNullable.Int64)
			usuario.IdSede = &idSedeInt
		} else {
			usuario.IdSede = nil
		}

		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}
