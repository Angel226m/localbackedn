package entidades

import (
	"database/sql"
	"time"
)

// InstanciaTour representa una instancia específica de un tour programado en una fecha determinada
type InstanciaTour struct {
	ID               int           `json:"id_instancia"`
	IDTourProgramado int           `json:"id_tour_programado"`
	FechaEspecifica  time.Time     `json:"fecha_especifica"`
	HoraInicio       time.Time     `json:"hora_inicio"`
	HoraFin          time.Time     `json:"hora_fin"`
	IDChofer         sql.NullInt64 `json:"id_chofer"`
	IDEmbarcacion    int           `json:"id_embarcacion"`
	CupoDisponible   int           `json:"cupo_disponible"`
	Estado           string        `json:"estado"`
	Eliminado        bool          `json:"eliminado"`

	// Campos adicionales para mostrar información relacionada
	NombreTipoTour     string `json:"nombre_tipo_tour,omitempty"`
	NombreEmbarcacion  string `json:"nombre_embarcacion,omitempty"`
	NombreSede         string `json:"nombre_sede,omitempty"`
	NombreChofer       string `json:"nombre_chofer,omitempty"`
	HoraInicioStr      string `json:"hora_inicio_str,omitempty"`
	HoraFinStr         string `json:"hora_fin_str,omitempty"`
	FechaEspecificaStr string `json:"fecha_especifica_str,omitempty"`
}

// NuevaInstanciaTourRequest representa los datos para crear una nueva instancia de tour
type NuevaInstanciaTourRequest struct {
	IDTourProgramado int    `json:"id_tour_programado" validate:"required"`
	FechaEspecifica  string `json:"fecha_especifica" validate:"required"`
	HoraInicio       string `json:"hora_inicio" validate:"required"`
	HoraFin          string `json:"hora_fin" validate:"required"`
	IDChofer         *int   `json:"id_chofer"`
	IDEmbarcacion    int    `json:"id_embarcacion" validate:"required"`
	CupoDisponible   int    `json:"cupo_disponible" validate:"required,min=1"`
}

// ActualizarInstanciaTourRequest representa los datos para actualizar una instancia de tour
type ActualizarInstanciaTourRequest struct {
	IDTourProgramado *int    `json:"id_tour_programado"`
	FechaEspecifica  *string `json:"fecha_especifica"`
	HoraInicio       *string `json:"hora_inicio"`
	HoraFin          *string `json:"hora_fin"`
	IDChofer         *int    `json:"id_chofer"`
	IDEmbarcacion    *int    `json:"id_embarcacion"`
	CupoDisponible   *int    `json:"cupo_disponible" validate:"omitempty,min=0"`
	Estado           *string `json:"estado" validate:"omitempty,oneof=PROGRAMADO EN_CURSO COMPLETADO CANCELADO"`
}

// AsignarChoferInstanciaRequest representa los datos para asignar un chofer a una instancia de tour
type AsignarChoferInstanciaRequest struct {
	IDChofer int `json:"id_chofer" validate:"required"`
}

// FiltrosInstanciaTour representa los filtros para buscar instancias de tour
type FiltrosInstanciaTour struct {
	IDTourProgramado *int    `json:"id_tour_programado"`
	FechaInicio      *string `json:"fecha_inicio"`
	FechaFin         *string `json:"fecha_fin"`
	Estado           *string `json:"estado"`
	IDChofer         *int    `json:"id_chofer"`
	IDEmbarcacion    *int    `json:"id_embarcacion"`
	IDSede           *int    `json:"id_sede"`
	IDTipoTour       *int    `json:"id_tipo_tour"`
}
