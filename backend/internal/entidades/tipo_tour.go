package entidades

// TipoTour representa la estructura de un tipo de tour en el sistema
type TipoTour struct {
	ID                int     `json:"id_tipo_tour" db:"id_tipo_tour"`
	Nombre            string  `json:"nombre" db:"nombre"`
	Descripcion       string  `json:"descripcion" db:"descripcion"`
	DuracionMinutos   int     `json:"duracion_minutos" db:"duracion_minutos"`
	PrecioBase        float64 `json:"precio_base" db:"precio_base"`
	CantidadPasajeros int     `json:"cantidad_pasajeros" db:"cantidad_pasajeros"`
	URLImagen         string  `json:"url_imagen" db:"url_imagen"`
}

// NuevoTipoTourRequest representa los datos necesarios para crear un nuevo tipo de tour
type NuevoTipoTourRequest struct {
	Nombre            string  `json:"nombre" validate:"required"`
	Descripcion       string  `json:"descripcion"`
	DuracionMinutos   int     `json:"duracion_minutos" validate:"required,min=1"`
	PrecioBase        float64 `json:"precio_base" validate:"required,min=0"`
	CantidadPasajeros int     `json:"cantidad_pasajeros" validate:"required,min=1"`
	URLImagen         string  `json:"url_imagen"`
}

// ActualizarTipoTourRequest representa los datos para actualizar un tipo de tour
type ActualizarTipoTourRequest struct {
	Nombre            string  `json:"nombre" validate:"required"`
	Descripcion       string  `json:"descripcion"`
	DuracionMinutos   int     `json:"duracion_minutos" validate:"required,min=1"`
	PrecioBase        float64 `json:"precio_base" validate:"required,min=0"`
	CantidadPasajeros int     `json:"cantidad_pasajeros" validate:"required,min=1"`
	URLImagen         string  `json:"url_imagen"`
}
