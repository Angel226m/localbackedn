package entidades

// MetodoPago representa la estructura de un método de pago en el sistema
type MetodoPago struct {
	ID          int    `json:"id_metodo_pago" db:"id_metodo_pago"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
}

// NuevoMetodoPagoRequest representa los datos necesarios para crear un nuevo método de pago
type NuevoMetodoPagoRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}

// ActualizarMetodoPagoRequest representa los datos para actualizar un método de pago
type ActualizarMetodoPagoRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}
