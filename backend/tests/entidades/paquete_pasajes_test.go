package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoPaquetePasajes prueba la validación de los datos de un nuevo paquete de pasajes
func TestValidacionNuevoPaquetePasajes(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		paquete       entidades.NuevoPaquetePasajesRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Paquete válido",
			paquete: entidades.NuevoPaquetePasajesRequest{
				IDSede:        1,
				IDTipoTour:    2,
				Nombre:        "Paquete Familiar",
				Descripcion:   "Incluye varias actividades.",
				PrecioTotal:   200.0,
				CantidadTotal: 5,
			},
			debeSerValido: true,
		},
		{
			nombre: "Paquete sin IDSede",
			paquete: entidades.NuevoPaquetePasajesRequest{
				IDTipoTour:    2,
				Nombre:        "Paquete Familiar",
				Descripcion:   "Incluye varias actividades.",
				PrecioTotal:   200.0,
				CantidadTotal: 5,
			},
			debeSerValido: false,
			campoInvalido: "id_sede",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.paquete)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
