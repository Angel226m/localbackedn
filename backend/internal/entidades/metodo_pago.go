package entidades

// MetodoPago representa la estructura de un método de pago en el sistema
type MetodoPago struct {
	ID          int    `json:"id_metodo_pago" db:"id_metodo_pago"`
	IDSede      int    `json:"id_sede" db:"id_sede"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Eliminado   bool   `json:"eliminado" db:"eliminado"`
	// Campos adicionales para mostrar información relacionada
	NombreSede string `json:"nombre_sede,omitempty" db:"-"`
}

// NuevoMetodoPagoRequest representa los datos necesarios para crear un nuevo método de pago
type NuevoMetodoPagoRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}

// ActualizarMetodoPagoRequest representa los datos para actualizar un método de pago
type ActualizarMetodoPagoRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
	Eliminado   bool   `json:"eliminado"`
}
