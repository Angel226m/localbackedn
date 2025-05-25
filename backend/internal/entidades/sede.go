package entidades

import (
	"time"
)

// Sede representa la estructura de una sede en el sistema
type Sede struct {
	ID        int       `json:"id_sede" db:"id_sede"`
	Nombre    string    `json:"nombre" db:"nombre"`
	Direccion string    `json:"direccion" db:"direccion"`
	Telefono  string    `json:"telefono" db:"telefono"`
	Correo    string    `json:"correo" db:"correo"`
	Ciudad    string    `json:"ciudad" db:"ciudad"`
	Provincia string    `json:"provincia" db:"provincia"`
	Pais      string    `json:"pais" db:"pais"`
	Eliminado bool      `json:"eliminado" db:"eliminado"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// NuevaSedeRequest representa los datos necesarios para crear una nueva sede
type NuevaSedeRequest struct {
	Nombre    string `json:"nombre" validate:"required"`
	Direccion string `json:"direccion" validate:"required"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo" validate:"omitempty,email"`
	Ciudad    string `json:"ciudad" validate:"required"`
	Provincia string `json:"provincia"`
	Pais      string `json:"pais" validate:"required"`
}

// ActualizarSedeRequest representa los datos necesarios para actualizar una sede
type ActualizarSedeRequest struct {
	Nombre    string `json:"nombre" validate:"required"`
	Direccion string `json:"direccion" validate:"required"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo" validate:"omitempty,email"`
	Ciudad    string `json:"ciudad" validate:"required"`
	Provincia string `json:"provincia"`
	Pais      string `json:"pais" validate:"required"`
}
