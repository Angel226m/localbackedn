package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoComprobantePago prueba la validación de los datos de un nuevo comprobante de pago
func TestValidacionNuevoComprobantePago(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		comprobante   entidades.NuevoComprobantePagoRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Comprobante válido",
			comprobante: entidades.NuevoComprobantePagoRequest{
				IDReserva:         1,
				IDSede:            2,
				Tipo:              "FACTURA",
				NumeroComprobante: "123456",
				Subtotal:          200.0,
				IGV:               36.0,
				Total:             236.0,
			},
			debeSerValido: true,
		},
		{
			nombre: "Comprobante sin Tipo",
			comprobante: entidades.NuevoComprobantePagoRequest{
				IDReserva:         1,
				IDSede:            2,
				NumeroComprobante: "123456",
				Subtotal:          200.0,
				IGV:               36.0,
				Total:             236.0,
			},
			debeSerValido: false,
			campoInvalido: "tipo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.comprobante)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
