package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoCliente prueba la validación de los datos de un nuevo cliente
func TestValidacionNuevoCliente(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		cliente       entidades.NuevoClienteRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Cliente válido",
			cliente: entidades.NuevoClienteRequest{
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				Contrasena:      "password123",
			},
			debeSerValido: true,
		},
		{
			nombre: "Cliente sin TipoDocumento",
			cliente: entidades.NuevoClienteRequest{
				NumeroDocumento: "12345678",
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				Contrasena:      "password123",
			},
			debeSerValido: false,
			campoInvalido: "tipo_documento",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.cliente)

			if tc.debeSerValido && err != nil {
				t.Errorf("Esperaba que fuera válido, pero hubo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Esperaba error de validación en %s, pero no ocurrió", tc.campoInvalido)
			}
		})
	}
}
