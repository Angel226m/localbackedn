package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoIdioma prueba la validación de los datos de un nuevo idioma
func TestValidacionNuevoIdioma(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		idioma        entidades.NuevoIdiomaRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Idioma válido",
			idioma: entidades.NuevoIdiomaRequest{
				Nombre: "Español",
			},
			debeSerValido: true,
		},
		{
			nombre:        "Idioma sin Nombre",
			idioma:        entidades.NuevoIdiomaRequest{},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.idioma)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
