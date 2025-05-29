package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevaReserva prueba la validación de los datos de una nueva reserva
func TestValidacionNuevaReserva(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		reserva       entidades.NuevaReservaRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Reserva válida",
			reserva: entidades.NuevaReservaRequest{
				IDCliente:        1,
				IDTourProgramado: 2,
				IDCanal:          3,
				IDSede:           4,
				TotalPagar:       100.0,
				CantidadPasajes:  []entidades.PasajeCantidadRequest{{IDTipoPasaje: 5, Cantidad: 2}},
			},
			debeSerValido: true,
		},
		{
			nombre: "Reserva sin IDCliente",
			reserva: entidades.NuevaReservaRequest{
				IDTourProgramado: 2,
				IDCanal:          3,
				IDSede:           4,
				TotalPagar:       100.0,
				CantidadPasajes:  []entidades.PasajeCantidadRequest{{IDTipoPasaje: 5, Cantidad: 2}},
			},
			debeSerValido: false,
			campoInvalido: "id_cliente",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.reserva)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
