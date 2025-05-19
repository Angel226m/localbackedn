package entidades

// TipoPasaje representa la estructura de un tipo de pasaje en el sistema
type TipoPasaje struct {
	ID     int     `json:"id_tipo_pasaje" db:"id_tipo_pasaje"`
	Nombre string  `json:"nombre" db:"nombre"`
	Costo  float64 `json:"costo" db:"costo"`
	Edad   string  `json:"edad" db:"edad"`
}

// NuevoTipoPasajeRequest representa los datos necesarios para crear un nuevo tipo de pasaje
type NuevoTipoPasajeRequest struct {
	Nombre string  `json:"nombre" validate:"required"`
	Costo  float64 `json:"costo" validate:"required,min=0"`
	Edad   string  `json:"edad" validate:"required"`
}

// ActualizarTipoPasajeRequest representa los datos para actualizar un tipo de pasaje
type ActualizarTipoPasajeRequest struct {
	Nombre string  `json:"nombre" validate:"required"`
	Costo  float64 `json:"costo" validate:"required,min=0"`
	Edad   string  `json:"edad" validate:"required"`
}
