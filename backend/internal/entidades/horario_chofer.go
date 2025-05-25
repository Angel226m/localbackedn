package entidades

import (
	"time"
)

// HorarioChofer representa la estructura de un horario de chofer en el sistema
type HorarioChofer struct {
	ID                  int        `json:"id_horario_chofer" db:"id_horario_chofer"`
	IDUsuario           int        `json:"id_usuario" db:"id_usuario"`
	IDSede              int        `json:"id_sede" db:"id_sede"`
	HoraInicio          time.Time  `json:"hora_inicio" db:"hora_inicio"`
	HoraFin             time.Time  `json:"hora_fin" db:"hora_fin"`
	DisponibleLunes     bool       `json:"disponible_lunes" db:"disponible_lunes"`
	DisponibleMartes    bool       `json:"disponible_martes" db:"disponible_martes"`
	DisponibleMiercoles bool       `json:"disponible_miercoles" db:"disponible_miercoles"`
	DisponibleJueves    bool       `json:"disponible_jueves" db:"disponible_jueves"`
	DisponibleViernes   bool       `json:"disponible_viernes" db:"disponible_viernes"`
	DisponibleSabado    bool       `json:"disponible_sabado" db:"disponible_sabado"`
	DisponibleDomingo   bool       `json:"disponible_domingo" db:"disponible_domingo"`
	FechaInicio         time.Time  `json:"fecha_inicio" db:"fecha_inicio"`
	FechaFin            *time.Time `json:"fecha_fin,omitempty" db:"fecha_fin"`
	Eliminado           bool       `json:"eliminado" db:"eliminado"`

	// Campos adicionales para mostrar información relacionada
	NombreChofer    string `json:"nombre_chofer,omitempty" db:"-"`
	ApellidosChofer string `json:"apellidos_chofer,omitempty" db:"-"`
	DocumentoChofer string `json:"documento_chofer,omitempty" db:"-"`
	TelefonoChofer  string `json:"telefono_chofer,omitempty" db:"-"`
	NombreSede      string `json:"nombre_sede,omitempty" db:"-"`
}

// NuevoHorarioChoferRequest representa los datos necesarios para crear un nuevo horario de chofer
type NuevoHorarioChoferRequest struct {
	IDUsuario           int        `json:"id_usuario" validate:"required"`
	IDSede              int        `json:"id_sede" validate:"required"`
	HoraInicio          string     `json:"hora_inicio" validate:"required"` // formato HH:MM
	HoraFin             string     `json:"hora_fin" validate:"required"`    // formato HH:MM
	DisponibleLunes     bool       `json:"disponible_lunes"`
	DisponibleMartes    bool       `json:"disponible_martes"`
	DisponibleMiercoles bool       `json:"disponible_miercoles"`
	DisponibleJueves    bool       `json:"disponible_jueves"`
	DisponibleViernes   bool       `json:"disponible_viernes"`
	DisponibleSabado    bool       `json:"disponible_sabado"`
	DisponibleDomingo   bool       `json:"disponible_domingo"`
	FechaInicio         time.Time  `json:"fecha_inicio" validate:"required"`
	FechaFin            *time.Time `json:"fecha_fin,omitempty"`
}

// ActualizarHorarioChoferRequest representa los datos para actualizar un horario de chofer
type ActualizarHorarioChoferRequest struct {
	IDUsuario           int        `json:"id_usuario" validate:"required"`
	IDSede              int        `json:"id_sede" validate:"required"`
	HoraInicio          string     `json:"hora_inicio" validate:"required"` // formato HH:MM
	HoraFin             string     `json:"hora_fin" validate:"required"`    // formato HH:MM
	DisponibleLunes     bool       `json:"disponible_lunes"`
	DisponibleMartes    bool       `json:"disponible_martes"`
	DisponibleMiercoles bool       `json:"disponible_miercoles"`
	DisponibleJueves    bool       `json:"disponible_jueves"`
	DisponibleViernes   bool       `json:"disponible_viernes"`
	DisponibleSabado    bool       `json:"disponible_sabado"`
	DisponibleDomingo   bool       `json:"disponible_domingo"`
	FechaInicio         time.Time  `json:"fecha_inicio" validate:"required"`
	FechaFin            *time.Time `json:"fecha_fin,omitempty"`
	Eliminado           bool       `json:"eliminado"`
}
