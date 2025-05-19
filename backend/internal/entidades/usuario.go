package entidades

import (
	"time"
)

// Usuario representa la estructura de un usuario en el sistema
type Usuario struct {
	ID              int       `json:"id_usuario" db:"id_usuario"`
	Nombres         string    `json:"nombres" db:"nombres"`
	Apellidos       string    `json:"apellidos" db:"apellidos"`
	Correo          string    `json:"correo" db:"correo"`
	Telefono        string    `json:"telefono" db:"telefono"`
	Direccion       string    `json:"direccion" db:"direccion"`
	FechaNacimiento time.Time `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Rol             string    `json:"rol" db:"rol"` // ADMIN, VENDEDOR, CHOFER, CLIENTE
	Nacionalidad    string    `json:"nacionalidad" db:"nacionalidad"`
	TipoDocumento   string    `json:"tipo_documento" db:"tipo_de_documento"`
	NumeroDocumento string    `json:"numero_documento" db:"numero_documento"`
	FechaRegistro   time.Time `json:"fecha_registro" db:"fecha_registro"`
	Contrasena      string    `json:"-" db:"contrasena"` // No se devuelve en JSON
	Estado          bool      `json:"estado" db:"estado"`
}

// NuevoUsuarioRequest representa los datos necesarios para crear un nuevo usuario
type NuevoUsuarioRequest struct {
	Nombres         string    `json:"nombres" validate:"required"`
	Apellidos       string    `json:"apellidos" validate:"required"`
	Correo          string    `json:"correo" validate:"required,email"`
	Telefono        string    `json:"telefono"`
	Direccion       string    `json:"direccion"`
	FechaNacimiento time.Time `json:"fecha_nacimiento" validate:"required"`
	Rol             string    `json:"rol" validate:"required,oneof=ADMIN VENDEDOR CHOFER CLIENTE"`
	Nacionalidad    string    `json:"nacionalidad"`
	TipoDocumento   string    `json:"tipo_documento" validate:"required"`
	NumeroDocumento string    `json:"numero_documento" validate:"required"`
	Contrasena      string    `json:"contrasena" validate:"required,min=8"`
}

// LoginRequest representa los datos necesarios para iniciar sesión
type LoginRequest struct {
	Correo     string `json:"correo" validate:"required,email"`
	Contrasena string `json:"contrasena" validate:"required"`
}

// LoginResponse representa la respuesta al iniciar sesión exitosamente
type LoginResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	Usuario      *Usuario `json:"usuario"`
}
