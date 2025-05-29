package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoPago prueba la validación de los datos de un nuevo pago
func TestValidacionNuevoPago(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		pago          entidades.NuevoPagoRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Pago válido",
			pago: entidades.NuevoPagoRequest{
				IDReserva:    1,
				IDMetodoPago: 2,
				IDCanal:      3,
				IDSede:       4,
				Monto:        100.0,
				Comprobante:  "ABC123",
			},
			debeSerValido: true,
		},
		{
			nombre: "Pago sin IDReserva",
			pago: entidades.NuevoPagoRequest{
				IDMetodoPago: 2,
				IDCanal:      3,
				IDSede:       4,
				Monto:        100.0,
				Comprobante:  "ABC123",
			},
			debeSerValido: false,
			campoInvalido: "id_reserva",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.pago)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
