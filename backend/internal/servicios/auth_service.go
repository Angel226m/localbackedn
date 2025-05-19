package servicios

import (
	"errors"
	"sistema-tours/internal/config"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
	"sistema-tours/internal/utils"
	"time"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	usuarioRepo *repositorios.UsuarioRepository
	config      *config.Config
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(usuarioRepo *repositorios.UsuarioRepository, config *config.Config) *AuthService {
	return &AuthService{
		usuarioRepo: usuarioRepo,
		config:      config,
	}
}

// Login autentica a un usuario y genera tokens JWT
func (s *AuthService) Login(loginReq *entidades.LoginRequest) (*entidades.LoginResponse, error) {
	// SOLO PARA DESARROLLO: Usuario hardcodeado para admin
	if loginReq.Correo == "admin@sistema-tours.com" && loginReq.Contrasena == "admin123" {
		// Intentar obtener el usuario de la BD para tener todos los datos
		usuario, err := s.usuarioRepo.GetByEmail(loginReq.Correo)
		if err != nil {
			// Si no podemos obtenerlo, creamos uno temporal
			usuario = &entidades.Usuario{
				ID:              1,
				Nombres:         "Admin",
				Apellidos:       "Sistema",
				Correo:          "admin@sistema-tours.com",
				Telefono:        "123456789",
				Direccion:       "Dirección Admin",
				FechaNacimiento: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Rol:             "ADMIN",
				Nacionalidad:    "Peruana",
				TipoDocumento:   "DNI",
				NumeroDocumento: "12345678",
				FechaRegistro:   time.Now(),
				Estado:          true,
			}
		}

		// Generar token JWT
		token, err := utils.GenerateJWT(usuario, s.config)
		if err != nil {
			return nil, err
		}

		// Generar refresh token
		refreshToken, err := utils.GenerateRefreshToken(usuario, s.config)
		if err != nil {
			return nil, err
		}

		// Ocultar contraseña hash
		usuario.Contrasena = ""

		// Crear respuesta
		loginResp := &entidades.LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			Usuario:      usuario,
		}

		return loginResp, nil
	}

	// Código original para otros usuarios
	// Buscar usuario por correo
	usuario, err := s.usuarioRepo.GetByEmail(loginReq.Correo)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar si el usuario está activo
	if !usuario.Estado {
		return nil, errors.New("usuario desactivado")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(loginReq.Contrasena, usuario.Contrasena) {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := utils.GenerateJWT(usuario, s.config)
	if err != nil {
		return nil, err
	}

	// Generar refresh token
	refreshToken, err := utils.GenerateRefreshToken(usuario, s.config)
	if err != nil {
		return nil, err
	}

	// Ocultar contraseña hash
	usuario.Contrasena = ""

	// Crear respuesta
	loginResp := &entidades.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Usuario:      usuario,
	}

	return loginResp, nil
}

// RefreshToken regenera el token de acceso usando un refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*entidades.LoginResponse, error) {
	// Validar refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken, s.config)
	if err != nil {
		return nil, err
	}

	// Obtener usuario
	usuario, err := s.usuarioRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Verificar si el usuario está activo
	if !usuario.Estado {
		return nil, errors.New("usuario desactivado")
	}

	// Generar nuevo token JWT
	newToken, err := utils.GenerateJWT(usuario, s.config)
	if err != nil {
		return nil, err
	}

	// Generar nuevo refresh token
	newRefreshToken, err := utils.GenerateRefreshToken(usuario, s.config)
	if err != nil {
		return nil, err
	}

	// Crear respuesta
	loginResp := &entidades.LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		Usuario:      usuario,
	}

	return loginResp, nil
}

// ChangePassword cambia la contraseña de un usuario
func (s *AuthService) ChangePassword(userID int, currentPassword, newPassword string) error {
	// Obtener usuario por ID
	user, err := s.usuarioRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Obtener contraseña actual (necesitamos el hash)
	userWithPassword, err := s.usuarioRepo.GetByEmail(user.Correo)
	if err != nil {
		return err
	}

	// Verificar contraseña actual
	if !utils.CheckPasswordHash(currentPassword, userWithPassword.Contrasena) {
		return errors.New("contraseña actual incorrecta")
	}

	// Hash de la nueva contraseña
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	return s.usuarioRepo.UpdatePassword(userID, hashedPassword)
}
