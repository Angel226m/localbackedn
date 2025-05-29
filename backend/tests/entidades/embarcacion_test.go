package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevaEmbarcacion prueba la validación de los datos de una nueva embarcación
func TestValidacionNuevaEmbarcacion(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		embarcacion   entidades.NuevaEmbarcacionRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Embarcación válida",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:    1,
				Nombre:    "Embarcación Azul",
				Capacidad: 20,
				Estado:    "DISPONIBLE",
			},
			debeSerValido: true,
		},
		{
			nombre: "Embarcación sin Nombre",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:    1,
				Capacidad: 20,
				Estado:    "DISPONIBLE",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.embarcacion)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
