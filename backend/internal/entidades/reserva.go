package entidades

import "time"

// Reserva representa la estructura de una reserva en el sistema
type Reserva struct {
	ID           int       `json:"id_reserva" db:"id_reserva"`
	IDVendedor   *int      `json:"id_vendedor,omitempty" db:"id_vendedor"`
	IDCliente    int       `json:"id_cliente" db:"id_cliente"`
	IDInstancia  int       `json:"id_instancia" db:"id_instancia"`
	IDCanal      int       `json:"id_canal" db:"id_canal"`
	IDSede       int       `json:"id_sede" db:"id_sede"`
	FechaReserva time.Time `json:"fecha_reserva" db:"fecha_reserva"`
	TotalPagar   float64   `json:"total_pagar" db:"total_pagar"`
	Notas        string    `json:"notas" db:"notas"`
	Estado       string    `json:"estado" db:"estado"` // RESERVADO, CANCELADA, CONFIRMADA, etc.
	Eliminado    bool      `json:"eliminado" db:"eliminado"`

	// Campos adicionales para mostrar información relacionada
	NombreCliente   string                 `json:"nombre_cliente,omitempty" db:"-"`
	NombreVendedor  string                 `json:"nombre_vendedor,omitempty" db:"-"`
	NombreTour      string                 `json:"nombre_tour,omitempty" db:"-"`
	FechaTour       string                 `json:"fecha_tour,omitempty" db:"-"`
	HoraInicioTour  string                 `json:"hora_inicio_tour,omitempty" db:"-"`
	HoraFinTour     string                 `json:"hora_fin_tour,omitempty" db:"-"`
	NombreCanal     string                 `json:"nombre_canal,omitempty" db:"-"`
	NombreSede      string                 `json:"nombre_sede,omitempty" db:"-"`
	CantidadPasajes []PasajeCantidad       `json:"cantidad_pasajes,omitempty" db:"-"`
	Paquetes        []PaquetePasajeDetalle `json:"paquetes,omitempty" db:"-"`
}

// PasajeCantidad representa la cantidad de pasajes de un tipo en la reserva
type PasajeCantidad struct {
	IDTipoPasaje int    `json:"id_tipo_pasaje" db:"id_tipo_pasaje"`
	NombreTipo   string `json:"nombre_tipo" db:"nombre"`
	Cantidad     int    `json:"cantidad" db:"cantidad"`
}

// PaquetePasajeDetalle representa un paquete de pasajes incluido en la reserva
type PaquetePasajeDetalle struct {
	IDPaquete      int     `json:"id_paquete" db:"id_paquete"`
	NombrePaquete  string  `json:"nombre_paquete" db:"nombre_paquete"`
	Cantidad       int     `json:"cantidad" db:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" db:"precio_unitario"`
	Subtotal       float64 `json:"subtotal" db:"subtotal"`
	CantidadTotal  int     `json:"cantidad_total" db:"cantidad_total"` // Total de pasajeros en el paquete
}

// NuevaReservaRequest representa los datos necesarios para crear una nueva reserva
type NuevaReservaRequest struct {
	IDCliente       int                     `json:"id_cliente" validate:"required"`
	IDInstancia     int                     `json:"id_instancia" validate:"required"`
	IDCanal         int                     `json:"id_canal" validate:"required"`
	IDSede          int                     `json:"id_sede" validate:"required"`
	IDVendedor      *int                    `json:"id_vendedor,omitempty"` // Opcional, solo si es reserva en LOCAL
	TotalPagar      float64                 `json:"total_pagar" validate:"required,min=0"`
	Notas           string                  `json:"notas"`
	CantidadPasajes []PasajeCantidadRequest `json:"cantidad_pasajes" validate:"dive"`
	Paquetes        []PaqueteRequest        `json:"paquetes" validate:"dive"`
}

// PasajeCantidadRequest representa la cantidad de pasajes de un tipo en la solicitud
type PasajeCantidadRequest struct {
	IDTipoPasaje int `json:"id_tipo_pasaje" validate:"required"`
	Cantidad     int `json:"cantidad" validate:"min=0"`
}

// PaqueteRequest representa un paquete de pasajes en la solicitud
type PaqueteRequest struct {
	IDPaquete int `json:"id_paquete" validate:"required"`
	Cantidad  int `json:"cantidad" validate:"min=0"`
}

// ActualizarReservaRequest representa los datos para actualizar una reserva
type ActualizarReservaRequest struct {
	IDCliente       int                     `json:"id_cliente" validate:"required"`
	IDInstancia     int                     `json:"id_instancia" validate:"required"`
	IDCanal         int                     `json:"id_canal" validate:"required"`
	IDSede          int                     `json:"id_sede" validate:"required"`
	IDVendedor      *int                    `json:"id_vendedor,omitempty"` // Opcional, solo si es reserva en LOCAL
	TotalPagar      float64                 `json:"total_pagar" validate:"required,min=0"`
	Notas           string                  `json:"notas"`
	Estado          string                  `json:"estado" validate:"required,oneof=RESERVADO CANCELADA CONFIRMADA"`
	CantidadPasajes []PasajeCantidadRequest `json:"cantidad_pasajes" validate:"dive"`
	Paquetes        []PaqueteRequest        `json:"paquetes" validate:"dive"`
}

// CambiarEstadoReservaRequest representa los datos para cambiar el estado de una reserva
type CambiarEstadoReservaRequest struct {
	Estado string `json:"estado" validate:"required,oneof=RESERVADO CANCELADA CONFIRMADA"`
}

// ReservaMercadoPagoRequest representa los datos para crear una reserva desde Mercado Pago
type ReservaMercadoPagoRequest struct {
	IDCliente       int                     `json:"id_cliente" validate:"required"`
	IDInstancia     int                     `json:"id_instancia" validate:"required"`
	TotalPagar      float64                 `json:"total_pagar" validate:"required,min=0"`
	CantidadPasajes []PasajeCantidadRequest `json:"cantidad_pasajes" validate:"dive"`
	Paquetes        []PaqueteRequest        `json:"paquetes" validate:"dive"`
	Email           string                  `json:"email" validate:"required,email"`
	Telefono        string                  `json:"telefono"`
	Documento       string                  `json:"documento"`
}

// ReservaMercadoPagoResponse representa la respuesta a una solicitud de reserva por Mercado Pago
type ReservaMercadoPagoResponse struct {
	IDReserva        int    `json:"id_reserva"`
	NombreTour       string `json:"nombre_tour"`
	PreferenceID     string `json:"preference_id"`
	InitPoint        string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point"`
}
