package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoHorarioTour prueba la validación de los datos de un nuevo horario de tour
func TestValidacionNuevoHorarioTour(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		horario       entidades.NuevoHorarioTourRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Horario válido",
			horario: entidades.NuevoHorarioTourRequest{
				IDTipoTour: 1,
				IDSede:     2,
				HoraInicio: "08:00",
				HoraFin:    "12:00",
			},
			debeSerValido: true,
		},
		{
			nombre: "Horario sin IDTipoTour",
			horario: entidades.NuevoHorarioTourRequest{
				IDSede:     2,
				HoraInicio: "08:00",
				HoraFin:    "12:00",
			},
			debeSerValido: false,
			campoInvalido: "id_tipo_tour",
		},
		{
			nombre: "Horario sin HoraInicio",
			horario: entidades.NuevoHorarioTourRequest{
				IDTipoTour: 1,
				IDSede:     2,
				HoraFin:    "12:00",
			},
			debeSerValido: false,
			campoInvalido: "hora_inicio",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.horario)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
