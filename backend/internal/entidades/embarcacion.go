package entidades

// Embarcacion representa la estructura de una embarcación en el sistema
type Embarcacion struct {
	ID          int    `json:"id_embarcacion" db:"id_embarcacion"`
	IDSede      int    `json:"id_sede" db:"id_sede"`
	Nombre      string `json:"nombre" db:"nombre"`
	Capacidad   int    `json:"capacidad" db:"capacidad"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Eliminado   bool   `json:"eliminado" db:"eliminado"`
	IDUsuario   int    `json:"id_usuario" db:"id_usuario"` // El chofer asignado
	Estado      string `json:"estado" db:"estado"`         // DISPONIBLE, OCUPADA, MANTENIMIENTO, FUERA_DE_SERVICIO
	// Campos adicionales para mostrar información del chofer
	NombreChofer    string `json:"nombre_chofer,omitempty" db:"-"`
	ApellidosChofer string `json:"apellidos_chofer,omitempty" db:"-"`
	DocumentoChofer string `json:"documento_chofer,omitempty" db:"-"`
	TelefonoChofer  string `json:"telefono_chofer,omitempty" db:"-"`
}

// NuevaEmbarcacionRequest representa los datos necesarios para crear una nueva embarcación
type NuevaEmbarcacionRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	IDUsuario   int    `json:"id_usuario" validate:"required"` // El chofer asignado
	Estado      string `json:"estado" validate:"required,oneof=DISPONIBLE OCUPADA MANTENIMIENTO FUERA_DE_SERVICIO"`
}

// ActualizarEmbarcacionRequest representa los datos para actualizar una embarcación
type ActualizarEmbarcacionRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	IDUsuario   int    `json:"id_usuario" validate:"required"` // El chofer asignado
	Estado      string `json:"estado" validate:"required,oneof=DISPONIBLE OCUPADA MANTENIMIENTO FUERA_DE_SERVICIO"`
	Eliminado   bool   `json:"eliminado"`
}
