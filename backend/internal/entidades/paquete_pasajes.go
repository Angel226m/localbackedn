package entidades

// PaquetePasajes representa la estructura de un paquete de pasajes en el sistema
type PaquetePasajes struct {
	ID            int     `json:"id_paquete" db:"id_paquete"`
	IDSede        int     `json:"id_sede" db:"id_sede"`
	IDTipoTour    int     `json:"id_tipo_tour" db:"id_tipo_tour"`
	Nombre        string  `json:"nombre" db:"nombre"`
	Descripcion   string  `json:"descripcion" db:"descripcion"`
	PrecioTotal   float64 `json:"precio_total" db:"precio_total"`
	CantidadTotal int     `json:"cantidad_total" db:"cantidad_total"`
	Eliminado     bool    `json:"eliminado" db:"eliminado"`
}

// NuevoPaquetePasajesRequest representa los datos necesarios para crear un nuevo paquete de pasajes
type NuevoPaquetePasajesRequest struct {
	IDSede        int     `json:"id_sede" validate:"required"`
	IDTipoTour    int     `json:"id_tipo_tour" validate:"required"`
	Nombre        string  `json:"nombre" validate:"required"`
	Descripcion   string  `json:"descripcion"`
	PrecioTotal   float64 `json:"precio_total" validate:"required,min=0"`
	CantidadTotal int     `json:"cantidad_total" validate:"required,min=1"`
}

// ActualizarPaquetePasajesRequest representa los datos para actualizar un paquete de pasajes
type ActualizarPaquetePasajesRequest struct {
	IDTipoTour    int     `json:"id_tipo_tour" validate:"required"`
	Nombre        string  `json:"nombre" validate:"required"`
	Descripcion   string  `json:"descripcion"`
	PrecioTotal   float64 `json:"precio_total" validate:"required,min=0"`
	CantidadTotal int     `json:"cantidad_total" validate:"required,min=1"`
}
