package entidades

import "time"

// Reserva representa la estructura de una reserva en el sistema
type Reserva struct {
	ID               int       `json:"id_reserva" db:"id_reserva"`
	IDVendedor       *int      `json:"id_vendedor,omitempty" db:"id_vendedor"`
	IDCliente        int       `json:"id_cliente" db:"id_cliente"`
	IDTourProgramado int       `json:"id_tour_programado" db:"id_tour_programado"`
	IDCanal          int       `json:"id_canal" db:"id_canal"`
	IDSede           int       `json:"id_sede" db:"id_sede"`
	FechaReserva     time.Time `json:"fecha_reserva" db:"fecha_reserva"`
	TotalPagar       float64   `json:"total_pagar" db:"total_pagar"`
	Notas            string    `json:"notas" db:"notas"`
	Estado           string    `json:"estado" db:"estado"` // RESERVADO, CANCELADA, etc.
	Eliminado        bool      `json:"eliminado" db:"eliminado"`

	// Campos adicionales para mostrar información relacionada
	NombreCliente   string           `json:"nombre_cliente,omitempty" db:"-"`
	NombreVendedor  string           `json:"nombre_vendedor,omitempty" db:"-"`
	NombreTour      string           `json:"nombre_tour,omitempty" db:"-"`
	FechaTour       string           `json:"fecha_tour,omitempty" db:"-"`
	HoraTour        string           `json:"hora_tour,omitempty" db:"-"`
	NombreCanal     string           `json:"nombre_canal,omitempty" db:"-"`
	NombreSede      string           `json:"nombre_sede,omitempty" db:"-"`
	CantidadPasajes []PasajeCantidad `json:"cantidad_pasajes,omitempty" db:"-"`
}

// PasajeCantidad representa la cantidad de pasajes de un tipo en la reserva
type PasajeCantidad struct {
	IDTipoPasaje int    `json:"id_tipo_pasaje" db:"id_tipo_pasaje"`
	NombreTipo   string `json:"nombre_tipo" db:"nombre"`
	Cantidad     int    `json:"cantidad" db:"cantidad"`
}

// NuevaReservaRequest representa los datos necesarios para crear una nueva reserva
type NuevaReservaRequest struct {
	IDCliente        int                     `json:"id_cliente" validate:"required"`
	IDTourProgramado int                     `json:"id_tour_programado" validate:"required"`
	IDCanal          int                     `json:"id_canal" validate:"required"`
	IDSede           int                     `json:"id_sede" validate:"required"`
	IDVendedor       *int                    `json:"id_vendedor,omitempty"` // Opcional, solo si es reserva en LOCAL
	TotalPagar       float64                 `json:"total_pagar" validate:"required,min=0"`
	Notas            string                  `json:"notas"`
	CantidadPasajes  []PasajeCantidadRequest `json:"cantidad_pasajes" validate:"required,min=1,dive"`
}

// PasajeCantidadRequest representa la cantidad de pasajes de un tipo en la solicitud
type PasajeCantidadRequest struct {
	IDTipoPasaje int `json:"id_tipo_pasaje" validate:"required"`
	Cantidad     int `json:"cantidad" validate:"required,min=1"`
}

// ActualizarReservaRequest representa los datos para actualizar una reserva
type ActualizarReservaRequest struct {
	IDCliente        int                     `json:"id_cliente" validate:"required"`
	IDTourProgramado int                     `json:"id_tour_programado" validate:"required"`
	IDCanal          int                     `json:"id_canal" validate:"required"`
	IDSede           int                     `json:"id_sede" validate:"required"`
	IDVendedor       *int                    `json:"id_vendedor,omitempty"` // Opcional, solo si es reserva en LOCAL
	TotalPagar       float64                 `json:"total_pagar" validate:"required,min=0"`
	Notas            string                  `json:"notas"`
	Estado           string                  `json:"estado" validate:"required,oneof=RESERVADO CANCELADA"`
	CantidadPasajes  []PasajeCantidadRequest `json:"cantidad_pasajes" validate:"required,min=1,dive"`
}

// CambiarEstadoReservaRequest representa los datos para cambiar el estado de una reserva
type CambiarEstadoReservaRequest struct {
	Estado string `json:"estado" validate:"required,oneof=RESERVADO CANCELADA"`
}
