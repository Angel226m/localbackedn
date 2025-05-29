package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoCanalVenta prueba la validación de los datos de un nuevo canal de venta
func TestValidacionNuevoCanalVenta(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		canalVenta    entidades.NuevoCanalVentaRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "CanalVenta válido",
			canalVenta: entidades.NuevoCanalVentaRequest{
				IDSede:      1,
				Nombre:      "Agencia Online",
				Descripcion: "Venta de tours a través de una plataforma digital.",
			},
			debeSerValido: true,
		},
		{
			nombre: "CanalVenta sin IDSede",
			canalVenta: entidades.NuevoCanalVentaRequest{
				Nombre:      "Agencia Online",
				Descripcion: "Venta de tours a través de una plataforma digital.",
			},
			debeSerValido: false,
			campoInvalido: "id_sede",
		},
		{
			nombre: "CanalVenta sin Nombre",
			canalVenta: entidades.NuevoCanalVentaRequest{
				IDSede:      1,
				Descripcion: "Venta de tours a través de una plataforma digital.",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.canalVenta)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
