package entidades_test

import (
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/utils"
	"testing"
	"time"
)

// TestValidacionNuevoUsuario prueba la validación de los datos de un nuevo usuario
func TestValidacionNuevoUsuario(t *testing.T) {
	// Inicializar el validador
	utils.InitValidator()

	// Tabla de casos de prueba
	tests := []struct {
		nombre        string
		usuario       entidades.NuevoUsuarioRequest
		debeSerValido bool
		campoInvalido string
	}{
		{
			nombre: "Usuario válido",
			usuario: entidades.NuevoUsuarioRequest{

				IdSede:          nil,
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				Telefono:        "123456789",
				Direccion:       "Calle Test 123",
				FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "ADMIN",
				Nacionalidad:    "Peruana",
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Contrasena:      "password123",
			},
			debeSerValido: true,
		},
		{
			nombre: "Usuario sin nombres",
			usuario: entidades.NuevoUsuarioRequest{
				Nombres:         "",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "ADMIN",
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Contrasena:      "password123",
			},
			debeSerValido: false,
			campoInvalido: "nombres", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "Email inválido",
			usuario: entidades.NuevoUsuarioRequest{
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan-test.com", // Email inválido
				FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "ADMIN",
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Contrasena:      "password123",
			},
			debeSerValido: false,
			campoInvalido: "correo", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "Rol inválido",
			usuario: entidades.NuevoUsuarioRequest{
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "INVALIDO", // Rol que no existe
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Contrasena:      "password123",
			},
			debeSerValido: false,
			campoInvalido: "rol", // Usar minúsculas para que coincida con el campo JSON
		},
		{
			nombre: "Contraseña corta",
			usuario: entidades.NuevoUsuarioRequest{
				Nombres:         "Juan",
				Apellidos:       "Pérez",
				Correo:          "juan@test.com",
				FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "ADMIN",
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				Contrasena:      "123", // Contraseña muy corta
			},
			debeSerValido: false,
			campoInvalido: "contrasena", // Usar minúsculas para que coincida con el campo JSON
		},
	}

	// Ejecutar casos de prueba
	for _, tc := range tests {
		t.Run(tc.nombre, func(t *testing.T) {
			err := utils.ValidateStruct(tc.usuario)

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

// TestUsuarioMethods prueba los métodos adicionales de la entidad Usuario
func TestUsuarioMethods(t *testing.T) {
	// Crear un usuario de prueba
	usuario := entidades.Usuario{
		ID:              1,
		IdSede:          nil,
		Nombres:         "Juan",
		Apellidos:       "Pérez",
		Correo:          "juan@test.com",
		Telefono:        "123456789",
		Direccion:       "Calle Test 123",
		FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Rol:             "ADMIN",
		Nacionalidad:    "Peruana",
		TipoDocumento:   "DNI",
		NumeroDocumento: "12345678",
		FechaRegistro:   time.Now(),
		Contrasena:      "password123",
		Eliminado:       true,
	}

	// Verifica que el usuario tenga los datos correctos
	if usuario.Nombres != "Juan" {
		t.Errorf("Se esperaba nombre 'Juan', pero se obtuvo: %s", usuario.Nombres)
	}

	if usuario.Apellidos != "Pérez" {
		t.Errorf("Se esperaba apellido 'Pérez', pero se obtuvo: %s", usuario.Apellidos)
	}

	if usuario.Rol != "ADMIN" {
		t.Errorf("Se esperaba rol 'ADMIN', pero se obtuvo: %s", usuario.Rol)
	}

	// En un caso real, aquí probarías métodos específicos de la entidad
	// Por ejemplo, si tuvieras un método GetNombreCompleto():
	// if usuario.GetNombreCompleto() != "Juan Pérez" {
	//    t.Errorf("Se esperaba nombre completo 'Juan Pérez', pero se obtuvo: %s", usuario.GetNombreCompleto())
	// }
}
