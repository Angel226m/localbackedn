package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
	"time"
)

// TestValidacionNuevoHorarioChofer prueba la validación de los datos de un nuevo horario de chofer
func TestValidacionNuevoHorarioChofer(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		horario       entidades.NuevoHorarioChoferRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Horario válido",
			horario: entidades.NuevoHorarioChoferRequest{
				IDUsuario:   1,
				IDSede:      2,
				HoraInicio:  "08:00",
				HoraFin:     "12:00",
				FechaInicio: time.Now(),
			},
			debeSerValido: true,
		},
		{
			nombre: "Horario sin IDUsuario",
			horario: entidades.NuevoHorarioChoferRequest{
				IDSede:      2,
				HoraInicio:  "08:00",
				HoraFin:     "12:00",
				FechaInicio: time.Now(),
			},
			debeSerValido: false,
			campoInvalido: "id_usuario",
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
