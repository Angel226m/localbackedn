package entidades

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Usuario representa la estructura de un usuario en el sistema
type Usuario struct {
	ID              int              `json:"id_usuario" db:"id_usuario"`
	IdSede          *int             `json:"id_sede" db:"id_sede"`
	Nombres         string           `json:"nombres" db:"nombres"`
	Apellidos       string           `json:"apellidos" db:"apellidos"`
	Correo          string           `json:"correo" db:"correo"`
	Telefono        string           `json:"telefono" db:"telefono"`
	Direccion       string           `json:"direccion" db:"direccion"`
	FechaNacimiento time.Time        `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Rol             string           `json:"rol" db:"rol"` // ADMIN, VENDEDOR, CHOFER
	Nacionalidad    string           `json:"nacionalidad" db:"nacionalidad"`
	TipoDocumento   string           `json:"tipo_documento" db:"tipo_de_documento"`
	NumeroDocumento string           `json:"numero_documento" db:"numero_documento"`
	FechaRegistro   time.Time        `json:"fecha_registro" db:"fecha_registro"`
	Contrasena      string           `json:"-" db:"contrasena"`        // No se devuelve en JSON
	Eliminado       bool             `json:"eliminado" db:"eliminado"` // Campo para soft delete
	Idiomas         []*UsuarioIdioma `json:"idiomas,omitempty" db:"-"` // Nueva relación muchos a muchos
}

// NuevoUsuarioRequest representa los datos necesarios para crear un nuevo usuario
type NuevoUsuarioRequest struct {
	IdSede          *int      `json:"id_sede"`
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
	IdiomasIDs      []int     `json:"idiomas_ids,omitempty"` // Nueva propiedad para IDs de idiomas
}

// LoginRequest representa los datos necesarios para iniciar sesión
type LoginRequest struct {
	Correo     string `json:"correo" validate:"required,email"`
	Contrasena string `json:"contrasena" validate:"required"`
}

// LoginResponse es la respuesta del endpoint de login
type LoginResponse struct {
	Token        string   `json:"token,omitempty"`         // Token JWT
	RefreshToken string   `json:"refresh_token,omitempty"` // Token de actualización
	Usuario      *Usuario `json:"usuario"`                 // Datos del usuario
}

// JWTClaims contiene los claims para el token JWT
type JWTClaims struct {
	UserID int    `json:"user_id"`
	SedeID int    `json:"sede_id,omitempty"` // Solo relevante para administradores que seleccionan sede
	Role   string `json:"role,omitempty"`    // Rol del usuario (ADMIN, VENDEDOR, CHOFER, etc.)
	jwt.StandardClaims
}
