package entidades

// CanalVenta representa la estructura de un canal de venta en el sistema
type CanalVenta struct {
	ID          int    `json:"id_canal" db:"id_canal"`
	IDSede      int    `json:"id_sede" db:"id_sede"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Eliminado   bool   `json:"eliminado" db:"eliminado"`
	// Campos adicionales para mostrar información relacionada
	NombreSede string `json:"nombre_sede,omitempty" db:"-"`
}

// NuevoCanalVentaRequest representa los datos necesarios para crear un nuevo canal de venta
type NuevoCanalVentaRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}

// ActualizarCanalVentaRequest representa los datos para actualizar un canal de venta
type ActualizarCanalVentaRequest struct {
	IDSede      int    `json:"id_sede" validate:"required"`
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
	Eliminado   bool   `json:"eliminado"`
}
