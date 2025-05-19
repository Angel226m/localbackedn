package entidades

// Embarcacion representa la estructura de una embarcación en el sistema
type Embarcacion struct {
	ID          int    `json:"id_embarcacion" db:"id_embarcacion"`
	Nombre      string `json:"nombre" db:"nombre"`
	Capacidad   int    `json:"capacidad" db:"capacidad"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Estado      bool   `json:"estado" db:"estado"`
	IDUsuario   int    `json:"id_usuario" db:"id_usuario"` // El chofer asignado
	// Campos adicionales para mostrar información del chofer
	NombreChofer    string `json:"nombre_chofer,omitempty" db:"-"`
	ApellidosChofer string `json:"apellidos_chofer,omitempty" db:"-"`
	DocumentoChofer string `json:"documento_chofer,omitempty" db:"-"`
	TelefonoChofer  string `json:"telefono_chofer,omitempty" db:"-"`
}

// NuevaEmbarcacionRequest representa los datos necesarios para crear una nueva embarcación
type NuevaEmbarcacionRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	IDUsuario   int    `json:"id_usuario" validate:"required"` // El chofer asignado
}

// ActualizarEmbarcacionRequest representa los datos para actualizar una embarcación
type ActualizarEmbarcacionRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Capacidad   int    `json:"capacidad" validate:"required,min=1"`
	Descripcion string `json:"descripcion"`
	IDUsuario   int    `json:"id_usuario" validate:"required"` // El chofer asignado
	Estado      bool   `json:"estado"`
}
