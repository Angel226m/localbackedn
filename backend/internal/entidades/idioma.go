package entidades

// Idioma representa la estructura de un idioma en el sistema
type Idioma struct {
	ID        int    `json:"id_idioma" db:"id_idioma"`
	Nombre    string `json:"nombre" db:"nombre" validate:"required,max=50"`
	Eliminado bool   `json:"eliminado" db:"eliminado"`
}

// NuevoIdiomaRequest representa los datos necesarios para crear un nuevo idioma
// NuevoIdiomaRequest representa los datos para crear un nuevo idioma
type NuevoIdiomaRequest struct {
	Nombre string `json:"nombre" validate:"required"`
}
