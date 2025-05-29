package entidades

// UsuarioIdioma representa la relación muchos a muchos entre Usuario e Idioma
type UsuarioIdioma struct {
	ID        int     `json:"id_usuario_idioma" db:"id_usuario_idioma"`
	IDUsuario int     `json:"id_usuario" db:"id_usuario"`
	IDIdioma  int     `json:"id_idioma" db:"id_idioma"`
	Nivel     string  `json:"nivel" db:"nivel"`
	Eliminado bool    `json:"eliminado" db:"eliminado"`
	Idioma    *Idioma `json:"idioma,omitempty" db:"-"` // Para incluir datos del idioma en consultas
}
