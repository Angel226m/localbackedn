package entidades

// CanalVenta representa la estructura de un canal de venta en el sistema
type CanalVenta struct {
	ID          int    `json:"id_canal" db:"id_canal"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
}

// NuevoCanalVentaRequest representa los datos necesarios para crear un nuevo canal de venta
type NuevoCanalVentaRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}

// ActualizarCanalVentaRequest representa los datos para actualizar un canal de venta
type ActualizarCanalVentaRequest struct {
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion"`
}
