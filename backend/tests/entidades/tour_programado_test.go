package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
	"time"
)

// TestValidacionNuevoTourProgramado prueba la validación de los datos de un nuevo tour programado
func TestValidacionNuevoTourProgramado(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre         string
		tourProgramado entidades.NuevoTourProgramadoRequest
		debeSerValido  bool
		campoInvalido  string
	}{
		{
			nombre: "Tour válido",
			tourProgramado: entidades.NuevoTourProgramadoRequest{
				IDTipoTour:    1,
				IDEmbarcacion: 2,
				IDHorario:     3,
				IDSede:        4,
				Fecha:         time.Now(),
				CupoMaximo:    10,
				Estado:        "PROGRAMADO",
			},
			debeSerValido: true,
		},
		{
			nombre: "Tour sin IDTipoTour",
			tourProgramado: entidades.NuevoTourProgramadoRequest{
				IDEmbarcacion: 2,
				IDHorario:     3,
				IDSede:        4,
				Fecha:         time.Now(),
				CupoMaximo:    10,
				Estado:        "PROGRAMADO",
			},
			debeSerValido: false,
			campoInvalido: "id_tipo_tour",
		},
		{
			nombre: "Tour sin fecha",
			tourProgramado: entidades.NuevoTourProgramadoRequest{
				IDTipoTour:    1,
				IDEmbarcacion: 2,
				IDHorario:     3,
				IDSede:        4,
				CupoMaximo:    10,
				Estado:        "PROGRAMADO",
			},
			debeSerValido: false,
			campoInvalido: "fecha",
		},
		{
			nombre: "Cupo máximo inválido",
			tourProgramado: entidades.NuevoTourProgramadoRequest{
				IDTipoTour:    1,
				IDEmbarcacion: 2,
				IDHorario:     3,
				IDSede:        4,
				Fecha:         time.Now(),
				CupoMaximo:    0, // Debe ser mínimo 1
				Estado:        "PROGRAMADO",
			},
			debeSerValido: false,
			campoInvalido: "cupo_maximo",
		},
		{
			nombre: "Estado inválido",
			tourProgramado: entidades.NuevoTourProgramadoRequest{
				IDTipoTour:    1,
				IDEmbarcacion: 2,
				IDHorario:     3,
				IDSede:        4,
				Fecha:         time.Now(),
				CupoMaximo:    10,
				Estado:        "INVALIDO", // Estado debe ser PROGRAMADO, COMPLETADO o CANCELADO
			},
			debeSerValido: false,
			campoInvalido: "estado",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.tourProgramado)

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
