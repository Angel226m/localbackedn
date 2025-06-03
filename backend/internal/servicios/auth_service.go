package servicios

import (
	"errors"
	"fmt"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"sistema-toursseft/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	usuarioRepo *repositorios.UsuarioRepository
	sedeRepo    *repositorios.SedeRepository
	config      *config.Config
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(usuarioRepo *repositorios.UsuarioRepository, sedeRepo *repositorios.SedeRepository, config *config.Config) *AuthService {
	return &AuthService{
		usuarioRepo: usuarioRepo,
		sedeRepo:    sedeRepo,
		config:      config,
	}
}

// Login autentica a un usuario y genera tokens JWT
func (s *AuthService) Login(loginReq *entidades.LoginRequest, rememberMe bool) (*entidades.LoginResponse, error) {
	// SOLO PARA DESARROLLO: Usuario hardcodeado para admin
	if loginReq.Correo == "admin@sistema-tours.com" && loginReq.Contrasena == "admin123" {
		// Intentar obtener el usuario de la BD para tener todos los datos
		usuario, err := s.usuarioRepo.GetByEmail(loginReq.Correo)
		if err != nil {
			// Si no podemos obtenerlo, creamos uno temporal
			usuario = &entidades.Usuario{
				ID:              1,
				IdSede:          nil, // Cambiar a nil para ADMIN
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
				Eliminado:       false,
			}
		}

		// Generar token JWT (15 minutos de duración)
		token, err := s.generateAccessToken(usuario)
		if err != nil {
			return nil, err
		}

		// Generar refresh token (duración basada en rememberMe)
		refreshToken, err := s.generateRefreshToken(usuario, rememberMe)
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

	// Verificar si el usuario está activo (no eliminado)
	if usuario.Eliminado {
		return nil, errors.New("usuario desactivado")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(loginReq.Contrasena, usuario.Contrasena) {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT (15 minutos de duración)
	token, err := s.generateAccessToken(usuario)
	if err != nil {
		return nil, err
	}

	// Generar refresh token (duración basada en rememberMe)
	refreshToken, err := s.generateRefreshToken(usuario, rememberMe)
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

// generateAccessToken genera un token de acceso (15 minutos)
func (s *AuthService) generateAccessToken(usuario *entidades.Usuario) (string, error) {
	// Tiempo de expiración del token principal (15 minutos)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Crear los claims para el token JWT
	claims := &entidades.JWTClaims{
		UserID: usuario.ID,
		Role:   usuario.Rol,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", usuario.ID),
		},
	}

	// Incluir sede si no es admin
	if usuario.Rol != "ADMIN" && usuario.IdSede != nil {
		claims.SedeID = *usuario.IdSede
	}

	// Crear el token con el algoritmo de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateRefreshToken genera un token de refresco (duración variable)
func (s *AuthService) generateRefreshToken(usuario *entidades.Usuario, rememberMe bool) (string, error) {
	// Determinar la duración basada en rememberMe
	var expirationTime time.Time
	if rememberMe {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 días
	} else {
		expirationTime = time.Now().Add(1 * time.Hour) // 1 hora
	}

	// Crear los claims para el refresh token
	claims := &entidades.JWTClaims{
		UserID: usuario.ID,
		Role:   usuario.Rol,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", usuario.ID),
		},
	}

	// Incluir sede si no es admin
	if usuario.Rol != "ADMIN" && usuario.IdSede != nil {
		claims.SedeID = *usuario.IdSede
	}

	// Crear el token con el algoritmo de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta para refresh tokens
	tokenString, err := token.SignedString([]byte(s.config.JWTRefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

	// Verificar si el usuario está activo (no eliminado)
	if usuario.Eliminado {
		return nil, errors.New("usuario desactivado")
	}

	// Determinar si el token original fue creado con rememberMe
	// Si el token expira en más de 24 horas desde su emisión, consideramos que tiene rememberMe activo

	//si cambiamos el tiempo de expiracion a 24 horas, entonces el refresh token se vuelve de 24 horas
	isRememberMe := false
	if claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) > 24*time.Hour {
		isRememberMe = true
	}

	// Generar nuevo token JWT (15 minutos)
	newToken, err := s.generateAccessToken(usuario)
	if err != nil {
		return nil, err
	}

	// Generar nuevo refresh token manteniendo la misma configuración de rememberMe
	newRefreshToken, err := s.generateRefreshToken(usuario, isRememberMe)
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

// GetUserByID obtiene un usuario por su ID
func (s *AuthService) GetUserByID(userID int) (*entidades.Usuario, error) {
	usuario, err := s.usuarioRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Asegurarse de no devolver la contraseña hash
	usuario.Contrasena = ""

	return usuario, nil
}

// GetSedeByID obtiene una sede por su ID
func (s *AuthService) GetSedeByID(sedeID int) (*entidades.Sede, error) {
	return s.sedeRepo.GetByID(sedeID)
}

// GetAllSedes obtiene todas las sedes disponibles (no eliminadas)
func (s *AuthService) GetAllSedes() ([]*entidades.Sede, error) {
	return s.sedeRepo.GetAll()
}

// GenerateTokensWithSede genera tokens incluyendo la sede seleccionada para administradores
func (s *AuthService) GenerateTokensWithSede(userID int, sedeID int, rememberMe bool) (string, string, error) {
	// Obtener usuario para verificar rol
	usuario, err := s.GetUserByID(userID)
	if err != nil {
		return "", "", err
	}

	// Verificar que sea administrador
	if usuario.Rol != "ADMIN" {
		return "", "", fmt.Errorf("solo los administradores pueden seleccionar sede temporalmente")
	}

	// Verificar que la sede exista
	sede, err := s.GetSedeByID(sedeID)
	if err != nil {
		return "", "", fmt.Errorf("sede no encontrada: %w", err)
	}

	if sede.Eliminado {
		return "", "", fmt.Errorf("la sede seleccionada no está disponible")
	}

	// Tiempo de expiración para el token de acceso (15 minutos)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Crear los claims para el token JWT
	claims := &entidades.JWTClaims{
		UserID: userID,
		SedeID: sedeID, // Incluir la sede seleccionada
		Role:   usuario.Rol,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Crear el token con el algoritmo de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Determinar la duración para el refresh token basado en rememberMe
	var refreshExpirationTime time.Time
	if rememberMe {
		refreshExpirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 días
	} else {
		refreshExpirationTime = time.Now().Add(1 * time.Hour) // 1 hora
	}

	// Crear los claims para el refresh token
	refreshClaims := &entidades.JWTClaims{
		UserID: userID,
		SedeID: sedeID, // Incluir la sede seleccionada
		Role:   usuario.Rol,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// GenerateTokensForAdminWithSede genera tokens con sede seleccionada sin buscar el usuario
func (s *AuthService) GenerateTokensForAdminWithSede(userID int, sedeID int, rememberMe bool) (string, string, error) {
	// Verificar que la sede exista
	sede, err := s.sedeRepo.GetByID(sedeID)
	if err != nil {
		return "", "", fmt.Errorf("sede no encontrada: %w", err)
	}

	if sede.Eliminado {
		return "", "", fmt.Errorf("la sede seleccionada no está disponible")
	}

	// Tiempo de expiración para el token de acceso (15 minutos)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Crear los claims para el token JWT
	claims := &utils.TokenClaims{
		UserID: userID,
		SedeID: sedeID,  // Incluir la sede seleccionada
		Role:   "ADMIN", // Ya sabemos que es admin por el middleware
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Crear el token con el algoritmo de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Determinar la duración para el refresh token basado en rememberMe
	var refreshExpirationTime time.Time
	if rememberMe {
		refreshExpirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 días
	} else {
		refreshExpirationTime = time.Now().Add(1 * time.Hour) // 1 hora
	}

	// Crear los claims para el refresh token
	refreshClaims := &utils.TokenClaims{
		UserID: userID,
		SedeID: sedeID,  // Incluir la sede seleccionada
		Role:   "ADMIN", // Ya sabemos que es admin
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// GenerateTokensWithoutDb genera tokens sin buscar el usuario en la base de datos
func (s *AuthService) GenerateTokensWithoutDb(userID int, userRole string, rememberMe bool) (string, string, error) {
	// Tiempo de expiración para el token de acceso (15 minutos)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Crear los claims para el token JWT
	claims := &utils.TokenClaims{
		UserID: userID,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Crear el token con el algoritmo de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Determinar la duración para el refresh token basado en rememberMe
	var refreshExpirationTime time.Time
	if rememberMe {
		refreshExpirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 días
	} else {
		refreshExpirationTime = time.Now().Add(1 * time.Hour) // 1 hora
	}

	// Crear los claims para el refresh token
	refreshClaims := &utils.TokenClaims{
		UserID: userID,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}
