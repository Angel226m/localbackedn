package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
)

// UsuarioIdiomaRepository maneja las operaciones de base de datos para la relación usuario-idioma
type UsuarioIdiomaRepository struct {
	db *sql.DB
}

// NewUsuarioIdiomaRepository crea una nueva instancia del repositorio
func NewUsuarioIdiomaRepository(db *sql.DB) *UsuarioIdiomaRepository {
	return &UsuarioIdiomaRepository{
		db: db,
	}
}

// GetByUsuarioID obtiene todos los idiomas de un usuario
func (r *UsuarioIdiomaRepository) GetByUsuarioID(usuarioID int) ([]*entidades.UsuarioIdioma, error) {
	query := `
		SELECT ui.id_usuario_idioma, ui.id_usuario, ui.id_idioma, ui.nivel, ui.eliminado,
		       i.nombre, i.eliminado
		FROM usuario_idioma ui
		JOIN idioma i ON ui.id_idioma = i.id_idioma
		WHERE ui.id_usuario = $1 AND ui.eliminado = false
	`

	rows, err := r.db.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idiomasUsuario []*entidades.UsuarioIdioma

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
			return nil, err
		}

		usuarioIdioma.Idioma.ID = usuarioIdioma.IDIdioma
		usuarioIdioma.Idioma.Eliminado = idiomaEliminado
		idiomasUsuario = append(idiomasUsuario, usuarioIdioma)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idiomasUsuario, nil
}

// AsignarIdioma asigna un idioma a un usuario
func (r *UsuarioIdiomaRepository) AsignarIdioma(usuarioID, idiomaID int, nivel string) error {
	query := `
		INSERT INTO usuario_idioma (id_usuario, id_idioma, nivel, eliminado)
		VALUES ($1, $2, $3, false)
		ON CONFLICT (id_usuario, id_idioma)
		DO UPDATE SET nivel = $3, eliminado = false
	`

	_, err := r.db.Exec(query, usuarioID, idiomaID, nivel)
	if err != nil {
		return err
	}

	return nil
}

// DesasignarIdioma elimina la asignación de un idioma a un usuario
func (r *UsuarioIdiomaRepository) DesasignarIdioma(usuarioID, idiomaID int) error {
	query := `
		UPDATE usuario_idioma 
		SET eliminado = true
		WHERE id_usuario = $1 AND id_idioma = $2 AND eliminado = false
	`

	result, err := r.db.Exec(query, usuarioID, idiomaID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("relación usuario-idioma no encontrada o ya eliminada")
	}

	return nil
}

// ActualizarIdiomasUsuario actualiza todos los idiomas de un usuario
func (r *UsuarioIdiomaRepository) ActualizarIdiomasUsuario(usuarioID int, idiomasIDs []int) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Marcar todos los idiomas actuales del usuario como eliminados
	_, err = tx.Exec("UPDATE usuario_idioma SET eliminado = true WHERE id_usuario = $1", usuarioID)
	if err != nil {
		return err
	}

	// Insertar o actualizar los nuevos idiomas
	for _, idiomaID := range idiomasIDs {
		_, err = tx.Exec(`
			INSERT INTO usuario_idioma (id_usuario, id_idioma, nivel, eliminado)
			VALUES ($1, $2, 'básico', false)
			ON CONFLICT (id_usuario, id_idioma) 
			DO UPDATE SET eliminado = false
		`, usuarioID, idiomaID)

		if err != nil {
			return err
		}
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetByIdiomaID obtiene todos los usuarios con un idioma específico
func (r *UsuarioIdiomaRepository) GetByIdiomaID(idiomaID int) ([]*entidades.Usuario, error) {
	query := `
		SELECT u.id_usuario, u.id_sede, u.nombres, u.apellidos, u.correo, u.telefono, u.direccion, 
			u.fecha_nacimiento, u.rol, u.nacionalidad, u.tipo_de_documento, u.numero_documento, 
			u.fecha_registro, u.contrasena, u.eliminado,
			ui.id_usuario_idioma, ui.nivel
		FROM usuario u
		JOIN usuario_idioma ui ON u.id_usuario = ui.id_usuario
		WHERE ui.id_idioma = $1 AND ui.eliminado = false AND u.eliminado = false
		ORDER BY u.apellidos, u.nombres
	`

	rows, err := r.db.Query(query, idiomaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []*entidades.Usuario{}
	usuariosMap := make(map[int]*entidades.Usuario)

	for rows.Next() {
		usuario := &entidades.Usuario{}
		usuarioIdioma := &entidades.UsuarioIdioma{
			Idioma: &entidades.Idioma{ID: idiomaID},
		}
		var idSedeNullable sql.NullInt64

		err := rows.Scan(
			&usuario.ID, &idSedeNullable, &usuario.Nombres, &usuario.Apellidos, &usuario.Correo,
			&usuario.Telefono, &usuario.Direccion, &usuario.FechaNacimiento, &usuario.Rol,
			&usuario.Nacionalidad, &usuario.TipoDocumento, &usuario.NumeroDocumento,
			&usuario.FechaRegistro, &usuario.Contrasena, &usuario.Eliminado,
			&usuarioIdioma.ID, &usuarioIdioma.Nivel,
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

		usuarioIdioma.IDUsuario = usuario.ID
		usuarioIdioma.IDIdioma = idiomaID
		usuarioIdioma.Eliminado = false

		// Si el usuario ya existe en el mapa, solo agregamos el idioma
		if existingUsuario, found := usuariosMap[usuario.ID]; found {
			existingUsuario.Idiomas = append(existingUsuario.Idiomas, usuarioIdioma)
		} else {
			// Si es un nuevo usuario, inicializamos la lista de idiomas
			usuario.Idiomas = []*entidades.UsuarioIdioma{usuarioIdioma}
			usuariosMap[usuario.ID] = usuario
			usuarios = append(usuarios, usuario)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}
