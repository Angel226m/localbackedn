package entidades

// Cliente representa la estructura de un cliente en el sistema
type Cliente struct {
	ID              int    `json:"id_cliente" db:"id_cliente"`
	TipoDocumento   string `json:"tipo_documento" db:"tipo_documento"`
	NumeroDocumento string `json:"numero_documento" db:"numero_documento"`
	Nombres         string `json:"nombres" db:"nombres"`
	Apellidos       string `json:"apellidos" db:"apellidos"`
	Correo          string `json:"correo" db:"correo"`
	Contrasena      string `json:"-" db:"contrasena,omitempty"`
	NombreCompleto  string `json:"nombre_completo,omitempty" db:"-"` // Campo calculado
	Eliminado       bool   `json:"eliminado" db:"eliminado"`
}

// NuevoClienteRequest representa los datos necesarios para crear un nuevo cliente
type NuevoClienteRequest struct {
	TipoDocumento   string `json:"tipo_documento" validate:"required"`
	NumeroDocumento string `json:"numero_documento" validate:"required"`
	Nombres         string `json:"nombres" validate:"required"`
	Apellidos       string `json:"apellidos" validate:"required"`
	Correo          string `json:"correo" validate:"omitempty,email"`
	Contrasena      string `json:"contrasena" validate:"omitempty,min=6"`
}

// ActualizarClienteRequest representa los datos para actualizar un cliente
type ActualizarClienteRequest struct {
	TipoDocumento   string `json:"tipo_documento" validate:"required"`
	NumeroDocumento string `json:"numero_documento" validate:"required"`
	Nombres         string `json:"nombres" validate:"required"`
	Apellidos       string `json:"apellidos" validate:"required"`
	Correo          string `json:"correo" validate:"omitempty,email"`
}

// LoginClienteRequest representa los datos para el login de un cliente
type LoginClienteRequest struct {
	Correo     string `json:"correo" validate:"required,email"`
	Contrasena string `json:"contrasena" validate:"required"`
}
