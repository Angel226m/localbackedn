package entidades

import (
	"time"
)

// HorarioTour representa la estructura de un horario de tour en el sistema
type HorarioTour struct {
	ID                  int       `json:"id_horario" db:"id_horario"`
	IDTipoTour          int       `json:"id_tipo_tour" db:"id_tipo_tour"`
	IDSede              int       `json:"id_sede" db:"id_sede"`
	HoraInicio          time.Time `json:"hora_inicio" db:"hora_inicio"`
	HoraFin             time.Time `json:"hora_fin" db:"hora_fin"`
	DisponibleLunes     bool      `json:"disponible_lunes" db:"disponible_lunes"`
	DisponibleMartes    bool      `json:"disponible_martes" db:"disponible_martes"`
	DisponibleMiercoles bool      `json:"disponible_miercoles" db:"disponible_miercoles"`
	DisponibleJueves    bool      `json:"disponible_jueves" db:"disponible_jueves"`
	DisponibleViernes   bool      `json:"disponible_viernes" db:"disponible_viernes"`
	DisponibleSabado    bool      `json:"disponible_sabado" db:"disponible_sabado"`
	DisponibleDomingo   bool      `json:"disponible_domingo" db:"disponible_domingo"`
	Eliminado           bool      `json:"eliminado" db:"eliminado"`

	// Campos adicionales para mostrar información relacionada
	NombreTipoTour      string `json:"nombre_tipo_tour,omitempty" db:"-"`
	DescripcionTipoTour string `json:"descripcion_tipo_tour,omitempty" db:"-"`
	NombreSede          string `json:"nombre_sede,omitempty" db:"-"`
}

// NuevoHorarioTourRequest representa los datos necesarios para crear un nuevo horario de tour
type NuevoHorarioTourRequest struct {
	IDTipoTour          int    `json:"id_tipo_tour" validate:"required"`
	IDSede              int    `json:"id_sede" validate:"required"`
	HoraInicio          string `json:"hora_inicio" validate:"required"` // formato HH:MM
	HoraFin             string `json:"hora_fin" validate:"required"`    // formato HH:MM
	DisponibleLunes     bool   `json:"disponible_lunes"`
	DisponibleMartes    bool   `json:"disponible_martes"`
	DisponibleMiercoles bool   `json:"disponible_miercoles"`
	DisponibleJueves    bool   `json:"disponible_jueves"`
	DisponibleViernes   bool   `json:"disponible_viernes"`
	DisponibleSabado    bool   `json:"disponible_sabado"`
	DisponibleDomingo   bool   `json:"disponible_domingo"`
}

// ActualizarHorarioTourRequest representa los datos para actualizar un horario de tour
type ActualizarHorarioTourRequest struct {
	IDTipoTour          int    `json:"id_tipo_tour" validate:"required"`
	IDSede              int    `json:"id_sede" validate:"required"`
	HoraInicio          string `json:"hora_inicio" validate:"required"` // formato HH:MM
	HoraFin             string `json:"hora_fin" validate:"required"`    // formato HH:MM
	DisponibleLunes     bool   `json:"disponible_lunes"`
	DisponibleMartes    bool   `json:"disponible_martes"`
	DisponibleMiercoles bool   `json:"disponible_miercoles"`
	DisponibleJueves    bool   `json:"disponible_jueves"`
	DisponibleViernes   bool   `json:"disponible_viernes"`
	DisponibleSabado    bool   `json:"disponible_sabado"`
	DisponibleDomingo   bool   `json:"disponible_domingo"`
	Eliminado           bool   `json:"eliminado"`
}
