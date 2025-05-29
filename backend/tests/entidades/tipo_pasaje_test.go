package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoTipoPasaje prueba la validación de los datos de un nuevo tipo de pasaje
func TestValidacionNuevoTipoPasaje(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		tipoPasaje    entidades.NuevoTipoPasajeRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "TipoPasaje válido",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDSede:     1,
				IDTipoTour: 2,
				Nombre:     "Pasaje Adulto",
				Costo:      50.0,
				Edad:       "Adulto",
			},
			debeSerValido: true,
		},
		{
			nombre: "TipoPasaje sin IDSede",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDTipoTour: 2,
				Nombre:     "Pasaje Adulto",
				Costo:      50.0,
				Edad:       "Adulto",
			},
			debeSerValido: false,
			campoInvalido: "id_sede",
		},
		{
			nombre: "TipoPasaje sin IDTipoTour",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDSede: 1,
				Nombre: "Pasaje Adulto",
				Costo:  50.0,
				Edad:   "Adulto",
			},
			debeSerValido: false,
			campoInvalido: "id_tipo_tour",
		},
		{
			nombre: "TipoPasaje sin Nombre",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDSede:     1,
				IDTipoTour: 2,
				Costo:      50.0,
				Edad:       "Adulto",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
		{
			nombre: "Costo inválido",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDSede:     1,
				IDTipoTour: 2,
				Nombre:     "Pasaje Adulto",
				Costo:      -5.0, // El costo no puede ser negativo
				Edad:       "Adulto",
			},
			debeSerValido: false,
			campoInvalido: "costo",
		},
		{
			nombre: "TipoPasaje sin Edad",
			tipoPasaje: entidades.NuevoTipoPasajeRequest{
				IDSede:     1,
				IDTipoTour: 2,
				Nombre:     "Pasaje Adulto",
				Costo:      50.0,
			},
			debeSerValido: false,
			campoInvalido: "edad",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.tipoPasaje)

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
