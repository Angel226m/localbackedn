package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoMetodoPago prueba la validación de los datos de un nuevo método de pago
func TestValidacionNuevoMetodoPago(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		metodoPago    entidades.NuevoMetodoPagoRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "MétodoPago válido",
			metodoPago: entidades.NuevoMetodoPagoRequest{
				IDSede:      1,
				Nombre:      "Tarjeta de crédito",
				Descripcion: "Pago con tarjeta",
			},
			debeSerValido: true,
		},
		{
			nombre: "MétodoPago sin IDSede",
			metodoPago: entidades.NuevoMetodoPagoRequest{
				Nombre:      "Tarjeta de crédito",
				Descripcion: "Pago con tarjeta",
			},
			debeSerValido: false,
			campoInvalido: "id_sede",
		},
		{
			nombre: "MétodoPago sin Nombre",
			metodoPago: entidades.NuevoMetodoPagoRequest{
				IDSede:      1,
				Descripcion: "Pago con tarjeta",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.metodoPago)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
