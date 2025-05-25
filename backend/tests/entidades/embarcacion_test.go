package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
)

// TestValidacionNuevoUsuario prueba la validación de los datos de un nuevo usuario
func TestValidacionNuevoembarcacion(t *testing.T) {
	// Inicializar el validador
	utils.InitValidator()

	// Tabla de casos de prueba
	tests := []struct {
		nombre        string
		embarcacion   entidades.NuevaEmbarcacionRequest
		debeSerValido bool
		campoInvalido string
		Nombre        string
		Capacidad     int
	}{
		{
			nombre: "Embarcación válida",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:      1,
				Capacidad:   5,
				Nombre:      "as ssa",
				Descripcion: "Descripción de la embarcación",
				IDUsuario:   1,
				Estado:      "DISPONIBLE",
			},
			debeSerValido: true,
		},
		{
			nombre: "Embarcación sin ID de sede",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:      0,
				Capacidad:   5,
				Nombre:      "ssas asa",
				Descripcion: "Descripción de la embarcación",
				IDUsuario:   1,
				Estado:      "DISPONIBLE",
			},
			debeSerValido: false,
			campoInvalido: "id_sede incompleto", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "capacidad  inválido",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:    1,
				Capacidad: 0,
				Nombre:    "ssas asa",

				IDUsuario:   1,
				Descripcion: "Descripción de la embarcación",
				Estado:      "DISPONIBLE",
			},
			debeSerValido: false,
			campoInvalido: "capacidad", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "descripción inválido",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:      1,
				Capacidad:   5,
				Descripcion: "",
				Nombre:      "ssas asa",
				IDUsuario:   1,
				Estado:      "DISPONIBLE",
			},
			debeSerValido: false,
			campoInvalido: "descripcion", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "estado",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:      1,
				Capacidad:   5,
				IDUsuario:   1,
				Nombre:      "ssas asa",
				Descripcion: "Descripción de la embarcación",
				Estado:      "", // Estado inválido
			},
			debeSerValido: false,
			campoInvalido: "estado", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "Nombre",
			embarcacion: entidades.NuevaEmbarcacionRequest{
				IDSede:      1,
				Capacidad:   5,
				Nombre:      "",
				IDUsuario:   1,
				Descripcion: "Descripción de la embarcación",
				Estado:      "DISPONIBLE", // Estado inválido
			},
			debeSerValido: false,
			campoInvalido: "nombre", // Usar minúsculas para que coincida con el campo JSON
		},
	}

	// Ejecutar casos de prueba
	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.embarcacion)

			if tc.debeSerValido && err != nil {
				t.Errorf("Se esperaba que fuera válido, pero se obtuvo error: %v", err)
			}

			if !tc.debeSerValido && err == nil {
				t.Errorf("Se esperaba error de validación para el campo %s, pero no se obtuvo ninguno", tc.campoInvalido)
			}

			// Verificar que el error mencionado contiene el nombre del campo inválido
			if !tc.debeSerValido && err != nil {
				validationErrors, ok := err.(utils.ValidationErrors)
				if !ok {
					t.Errorf("Se esperaba un error de tipo ValidationErrors, pero se obtuvo: %T", err)
					return
				}

				encontrado := false
				for _, fieldErr := range validationErrors {
					// Imprimir cada error para depuración
					t.Logf("Error de validación: Campo '%s': %s", fieldErr.Field, fieldErr.Message)

					// Comparar exactamente los nombres de campo
					if fieldErr.Field == tc.campoInvalido {
						encontrado = true
						break
					}
				}

				if !encontrado {
					t.Errorf("Error de validación no contiene el campo esperado %s. Errores: %v", tc.campoInvalido, validationErrors)
				}
			}
		})
	}
}
