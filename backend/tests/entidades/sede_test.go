package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevaSede prueba la validación de los datos de una nueva sede
func TestValidacionNuevaSede(t *testing.T) {
	utils.InitValidator()

	tests := []struct {
		nombre        string
		sede          entidades.NuevaSedeRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Sede válida",
			sede: entidades.NuevaSedeRequest{
				Nombre:    "Sede Principal",
				Direccion: "Av. Ejemplo 123",
				Telefono:  "987654321",
				Correo:    "contacto@ejemplo.com",
				Distrito:  "Centro",
				Provincia: "Ejemplo",
				Pais:      "Perú",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: true,
		},
		{
			nombre: "Sede sin Nombre",
			sede: entidades.NuevaSedeRequest{
				Direccion: "Av. Ejemplo 123",
				Telefono:  "987654321",
				Correo:    "contacto@ejemplo.com",
				Distrito:  "Centro",
				Provincia: "Ejemplo",
				Pais:      "Perú",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "nombre",
		},
		{
			nombre: "Sede sin Dirección",
			sede: entidades.NuevaSedeRequest{
				Nombre:    "Sede Principal",
				Telefono:  "987654321",
				Correo:    "contacto@ejemplo.com",
				Distrito:  "Centro",
				Provincia: "Ejemplo",
				Pais:      "Perú",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "direccion",
		},
		{
			nombre: "Correo inválido",
			sede: entidades.NuevaSedeRequest{
				Nombre:    "Sede Principal",
				Direccion: "Av. Ejemplo 123",
				Telefono:  "987654321",
				Correo:    "correo-invalido", // No es un email válido
				Distrito:  "Centro",
				Provincia: "Ejemplo",
				Pais:      "Perú",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "correo",
		},
		{
			nombre: "Sede sin Distrito",
			sede: entidades.NuevaSedeRequest{
				Nombre:    "Sede Principal",
				Direccion: "Av. Ejemplo 123",
				Telefono:  "987654321",
				Correo:    "contacto@ejemplo.com",
				Provincia: "Ejemplo",
				Pais:      "Perú",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "distrito",
		},
		{
			nombre: "Sede sin País",
			sede: entidades.NuevaSedeRequest{
				Nombre:    "Sede Principal",
				Direccion: "Av. Ejemplo 123",
				Telefono:  "987654321",
				Correo:    "contacto@ejemplo.com",
				Distrito:  "Centro",
				Provincia: "Ejemplo",
				ImageURL:  "https://ejemplo.com/imagen.jpg",
			},
			debeSerValido: false,
			campoInvalido: "pais",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.sede)

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
