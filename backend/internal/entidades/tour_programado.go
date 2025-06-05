package entidades

import (
	"database/sql"
	"time"
)

// TourProgramado representa un tour programado en el sistema
type TourProgramado struct {
	ID             int            `json:"id_tour_programado"`
	IDTipoTour     int            `json:"id_tipo_tour"`
	IDEmbarcacion  int            `json:"id_embarcacion"`
	IDHorario      int            `json:"id_horario"`
	IDSede         int            `json:"id_sede"`
	IDChofer       sql.NullInt64  `json:"id_chofer"`
	Fecha          time.Time      `json:"fecha"`
	VigenciaDesde  time.Time      `json:"vigencia_desde"`
	VigenciaHasta  time.Time      `json:"vigencia_hasta"`
	CupoMaximo     int            `json:"cupo_maximo"`
	CupoDisponible int            `json:"cupo_disponible"`
	Estado         string         `json:"estado"`
	Eliminado      bool           `json:"eliminado"`
	EsExcepcion    bool           `json:"es_excepcion"`
	NotasExcepcion sql.NullString `json:"notas_excepcion"`

	// Campos adicionales para mostrar información relacionada
	NombreTipoTour    string `json:"nombre_tipo_tour,omitempty"`
	NombreEmbarcacion string `json:"nombre_embarcacion,omitempty"`
	NombreSede        string `json:"nombre_sede,omitempty"`
	NombreChofer      string `json:"nombre_chofer,omitempty"`
	HoraInicio        string `json:"hora_inicio,omitempty"`
	HoraFin           string `json:"hora_fin,omitempty"`
}

// NuevoTourProgramadoRequest representa los datos para crear un nuevo tour programado
type NuevoTourProgramadoRequest struct {
	IDTipoTour     int     `json:"id_tipo_tour" validate:"required"`
	IDEmbarcacion  int     `json:"id_embarcacion" validate:"required"`
	IDHorario      int     `json:"id_horario" validate:"required"`
	IDSede         int     `json:"id_sede" validate:"required"`
	IDChofer       *int    `json:"id_chofer"`
	Fecha          string  `json:"fecha" validate:"required"`
	VigenciaDesde  string  `json:"vigencia_desde" validate:"required"`
	VigenciaHasta  string  `json:"vigencia_hasta" validate:"required"`
	CupoMaximo     int     `json:"cupo_maximo" validate:"required,min=1"`
	EsExcepcion    bool    `json:"es_excepcion"`
	NotasExcepcion *string `json:"notas_excepcion"`
}

// ActualizarTourProgramadoRequest representa los datos para actualizar un tour programado
type ActualizarTourProgramadoRequest struct {
	IDTipoTour     int     `json:"id_tipo_tour"`
	IDEmbarcacion  int     `json:"id_embarcacion"`
	IDHorario      int     `json:"id_horario"`
	IDSede         int     `json:"id_sede"`
	IDChofer       *int    `json:"id_chofer"`
	Fecha          string  `json:"fecha"`
	VigenciaDesde  string  `json:"vigencia_desde"`
	VigenciaHasta  string  `json:"vigencia_hasta"`
	CupoMaximo     int     `json:"cupo_maximo" validate:"min=1"`
	CupoDisponible int     `json:"cupo_disponible" validate:"min=0"`
	Estado         string  `json:"estado" validate:"oneof=PROGRAMADO COMPLETADO CANCELADO EN_CURSO"`
	EsExcepcion    bool    `json:"es_excepcion"`
	NotasExcepcion *string `json:"notas_excepcion"`
}

// AsignarChoferRequest representa los datos para asignar un chofer a un tour programado
type AsignarChoferRequest struct {
	IDChofer int `json:"id_chofer" validate:"required"`
}

// FiltrosTourProgramado representa los filtros para buscar tours programados
type FiltrosTourProgramado struct {
	IDSede           *int    `json:"id_sede"`
	IDTipoTour       *int    `json:"id_tipo_tour"`
	FechaInicio      *string `json:"fecha_inicio"`
	FechaFin         *string `json:"fecha_fin"`
	VigenciaDesdeIni *string `json:"vigencia_desde_ini"`
	VigenciaDesdefin *string `json:"vigencia_desde_fin"`
	VigenciaHastaIni *string `json:"vigencia_hasta_ini"`
	VigenciaHastaFin *string `json:"vigencia_hasta_fin"`
	Estado           *string `json:"estado"`
	IDChofer         *int    `json:"id_chofer"`
	IDEmbarcacion    *int    `json:"id_embarcacion"`
}
