package entidades

// TipoPasaje representa la estructura de un tipo de pasaje en el sistema
type TipoPasaje struct {
	ID         int     `json:"id_tipo_pasaje" db:"id_tipo_pasaje"`
	IDSede     int     `json:"id_sede" db:"id_sede"`
	IDTipoTour int     `json:"id_tipo_tour" db:"id_tipo_tour"`
	Nombre     string  `json:"nombre" db:"nombre"`
	Costo      float64 `json:"costo" db:"costo"`
	Edad       string  `json:"edad" db:"edad"`
	Eliminado  bool    `json:"eliminado" db:"eliminado"`
}

// NuevoTipoPasajeRequest representa los datos necesarios para crear un nuevo tipo de pasaje
type NuevoTipoPasajeRequest struct {
	IDSede     int     `json:"id_sede" validate:"required"`
	IDTipoTour int     `json:"id_tipo_tour" validate:"required"`
	Nombre     string  `json:"nombre" validate:"required"`
	Costo      float64 `json:"costo" validate:"required,min=0"`
	Edad       string  `json:"edad" validate:"required"`
}

// ActualizarTipoPasajeRequest representa los datos para actualizar un tipo de pasaje
type ActualizarTipoPasajeRequest struct {
	IDTipoTour int     `json:"id_tipo_tour" validate:"required"`
	Nombre     string  `json:"nombre" validate:"required"`
	Costo      float64 `json:"costo" validate:"required,min=0"`
	Edad       string  `json:"edad" validate:"required"`
}
