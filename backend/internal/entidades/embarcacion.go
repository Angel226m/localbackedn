package entidades

// Embarcacion representa la estructura de una embarcación en el sistema
type Embarcacion struct {
	ID          int    `json:"id_embarcacion" db:"id_embarcacion"`
	IDSede      int    `json:"id_sede" db:"id_sede"`
	Nombre      string `json:"nombre" db:"nombre"`
	Capacidad   int    `json:"capacidad" db:"capacidad"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Eliminado   bool   `json:"eliminado" db:"eliminado"`
	Estado      string `json:"estado" db:"estado"` // DISPONIBLE, OCUPADA, MANTENIMIENTO, FUERA_DE_SERVICIO
}

// NuevaEmbarcacionRequest representa los datos necesarios para crear una nueva embarcación
type NuevaEmbarcacionRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado" validate:"required,oneof=DISPONIBLE OCUPADA MANTENIMIENTO FUERA_DE_SERVICIO"`
}

// ActualizarEmbarcacionRequest representa los datos para actualizar una embarcación
type ActualizarEmbarcacionRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado" validate:"required,oneof=DISPONIBLE OCUPADA MANTENIMIENTO FUERA_DE_SERVICIO"`
}
