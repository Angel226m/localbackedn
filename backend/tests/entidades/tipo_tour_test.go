package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoTipoTour prueba la validación de los datos de un nuevo tipo de tour
func TestValidacionNuevoTipoTour(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		tipoTour      entidades.NuevoTipoTourRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "TipoTour válido",
			tipoTour: entidades.NuevoTipoTourRequest{
				IDSede:          1,
				Nombre:          "Tour Aventura",
				Descripcion:     "Un recorrido emocionante.",
				DuracionMinutos: 90,
				URLImagen:       "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: true,
		},
		{
			nombre: "TipoTour sin IDSede",
			tipoTour: entidades.NuevoTipoTourRequest{
				Nombre:          "Tour Aventura",
				Descripcion:     "Un recorrido emocionante.",
				DuracionMinutos: 90,
				URLImagen:       "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "id_sede",
		},
		{
			nombre: "TipoTour sin Nombre",
			tipoTour: entidades.NuevoTipoTourRequest{
				IDSede:          1,
				Descripcion:     "Un recorrido emocionante.",
				DuracionMinutos: 90,
				URLImagen:       "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
		{
			nombre: "Duración inválida",
			tipoTour: entidades.NuevoTipoTourRequest{
				IDSede:          1,
				Nombre:          "Tour Aventura",
				Descripcion:     "Un recorrido emocionante.",
				DuracionMinutos: 0, // Debe ser mínimo 1
				URLImagen:       "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "duracion_minutos",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.tipoTour)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido {
				if err == nil {
					t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
				} else {
					t.Logf("Errores de validación detectados: %v", err)
				}
			}
		})
	}
}
