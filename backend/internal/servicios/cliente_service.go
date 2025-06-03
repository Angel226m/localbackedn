package servicios

import (
	"errors"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"sistema-toursseft/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ClienteService maneja la lógica de negocio para clientes
type ClienteService struct {
	clienteRepo *repositorios.ClienteRepository
	config      *config.Config
}

// NewClienteService crea una nueva instancia de ClienteService
func NewClienteService(clienteRepo *repositorios.ClienteRepository, config *config.Config) *ClienteService {
	return &ClienteService{
		clienteRepo: clienteRepo,
		config:      config,
	}
}

// GetByID obtiene un cliente por su ID
func (s *ClienteService) GetByID(id int) (*entidades.Cliente, error) {
	return s.clienteRepo.GetByID(id)
}

// GetByDocumento obtiene un cliente por su tipo y número de documento
func (s *ClienteService) GetByDocumento(tipoDocumento, numeroDocumento string) (*entidades.Cliente, error) {
	return s.clienteRepo.GetByDocumento(tipoDocumento, numeroDocumento)
}

// Create crea un nuevo cliente
func (s *ClienteService) Create(cliente *entidades.NuevoClienteRequest) (int, error) {
	// Verificar si ya existe cliente con el mismo correo
	if cliente.Correo != "" {
		existingEmail, err := s.clienteRepo.GetByCorreo(cliente.Correo)
		if err == nil && existingEmail != nil {
			return 0, errors.New("ya existe un cliente con ese correo electrónico")
		}
	}

	// Verificar si ya existe cliente con el mismo documento
	existingDoc, err := s.clienteRepo.GetByDocumento(cliente.TipoDocumento, cliente.NumeroDocumento)
	if err == nil && existingDoc != nil {
		return 0, errors.New("ya existe un cliente con ese documento")
	}

	// Hash de la contraseña si se proporcionó una
	if cliente.Contrasena != "" {
		hashedPassword, err := utils.HashPassword(cliente.Contrasena)
		if err != nil {
			return 0, err
		}
		cliente.Contrasena = hashedPassword
	}

	// Crear cliente
	return s.clienteRepo.Create(cliente)
}

// Update actualiza un cliente existente
func (s *ClienteService) Update(id int, cliente *entidades.ActualizarClienteRequest) error {
	// Verificar que el cliente existe
	existing, err := s.clienteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro cliente con el mismo correo
	if cliente.Correo != "" && cliente.Correo != existing.Correo {
		existingEmail, err := s.clienteRepo.GetByCorreo(cliente.Correo)
		if err == nil && existingEmail != nil && existingEmail.ID != id {
			return errors.New("ya existe otro cliente con ese correo electrónico")
		}
	}

	// Verificar si ya existe otro cliente con el mismo documento
	if cliente.NumeroDocumento != existing.NumeroDocumento || cliente.TipoDocumento != existing.TipoDocumento {
		existingDoc, err := s.clienteRepo.GetByDocumento(cliente.TipoDocumento, cliente.NumeroDocumento)
		if err == nil && existingDoc != nil && existingDoc.ID != id {
			return errors.New("ya existe otro cliente con ese documento")
		}
	}

	// Actualizar cliente
	return s.clienteRepo.Update(id, cliente)
}

// Delete elimina lógicamente un cliente
func (s *ClienteService) Delete(id int) error {
	return s.clienteRepo.Delete(id)
}

// List obtiene todos los clientes no eliminados
func (s *ClienteService) List() ([]*entidades.Cliente, error) {
	return s.clienteRepo.List()
}

// SearchByName busca clientes por nombre o apellido
func (s *ClienteService) SearchByName(query string) ([]*entidades.Cliente, error) {
	return s.clienteRepo.SearchByName(query)
}

// Login autentica a un cliente y lo retorna si las credenciales son válidas
func (s *ClienteService) Login(correo, contrasena string, rememberMe bool) (*entidades.Cliente, string, string, error) {
	// Verificar si existe el cliente con ese correo
	cliente, err := s.clienteRepo.GetByCorreo(correo)
	if err != nil {
		return nil, "", "", errors.New("credenciales incorrectas")
	}

	// Obtener la contraseña hash
	hashedPassword, err := s.clienteRepo.GetPasswordByCorreo(correo)
	if err != nil {
		return nil, "", "", errors.New("credenciales incorrectas")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(contrasena, hashedPassword) {
		return nil, "", "", errors.New("credenciales incorrectas")
	}

	// Generar token JWT (15 minutos)
	token, err := s.generateAccessToken(cliente.ID, correo)
	if err != nil {
		return nil, "", "", errors.New("error al generar token")
	}

	// Generar refresh token (duración basada en rememberMe)
	refreshToken, err := s.generateRefreshToken(cliente.ID, correo, rememberMe)
	if err != nil {
		return nil, "", "", errors.New("error al generar refresh token")
	}

	return cliente, token, refreshToken, nil
}

// generateAccessToken genera un token de acceso para un cliente (15 minutos)
func (s *ClienteService) generateAccessToken(clienteID int, correo string) (string, error) {
	// Tiempo de expiración del token principal (15 minutos)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Crear los claims para el token JWT
	claims := &jwt.StandardClaims{
		// Usar el formato "cliente:{id}" para el subject para distinguir de usuarios regulares
		Subject:   "cliente:" + strconv.Itoa(clienteID),
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "sistema-tours",
	}

	// Crear token con los claims estándar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la clave secreta
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateRefreshToken genera un token de refresco para un cliente (duración variable)
func (s *ClienteService) generateRefreshToken(clienteID int, correo string, rememberMe bool) (string, error) {
	// Determinar la duración basada en rememberMe
	var expirationTime time.Time
	if rememberMe {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 días
	} else {
		expirationTime = time.Now().Add(1 * time.Hour) // 1 hora
	}

	// Crear los claims para el refresh token
	claims := &jwt.StandardClaims{
		// Usar el formato "cliente:{id}" para el subject para distinguir de usuarios regulares
		Subject:   "cliente:" + strconv.Itoa(clienteID),
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "sistema-tours",
	}

	// Crear token con los claims estándar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la clave secreta para refresh tokens
	tokenString, err := token.SignedString([]byte(s.config.JWTRefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RefreshClienteToken renueva los tokens de un cliente usando su refresh token
func (s *ClienteService) RefreshClienteToken(refreshToken string) (string, string, *entidades.Cliente, error) {
	// Validar refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTRefreshSecret), nil
	})

	if err != nil {
		return "", "", nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", "", nil, errors.New("token inválido")
	}

	// Extraer ID de cliente del subject (formato "cliente:{id}")
	subjectParts := strings.Split(claims.Subject, ":")
	if len(subjectParts) != 2 || subjectParts[0] != "cliente" {
		return "", "", nil, errors.New("token inválido: no es un token de cliente")
	}

	clienteID, err := strconv.Atoi(subjectParts[1])
	if err != nil {
		return "", "", nil, errors.New("token inválido: ID de cliente no válido")
	}

	// Obtener cliente
	cliente, err := s.clienteRepo.GetByID(clienteID)
	if err != nil {
		return "", "", nil, err
	}

	// Determinar si el token original fue creado con rememberMe
	// Si el token expira en más de 24 horas desde su emisión, consideramos que tiene rememberMe activo
	isRememberMe := false

	// Convertir Unix timestamps a time.Time para poder compararlos
	issuedAt := time.Unix(claims.IssuedAt, 0)
	expiresAt := time.Unix(claims.ExpiresAt, 0)

	// Comprobar si la duración es mayor a 24 horas
	if expiresAt.Sub(issuedAt) > 24*time.Hour {
		isRememberMe = true
	}

	// Generar nuevo token JWT
	newToken, err := s.generateAccessToken(clienteID, cliente.Correo)
	if err != nil {
		return "", "", nil, err
	}

	// Generar nuevo refresh token manteniendo la misma configuración de rememberMe
	newRefreshToken, err := s.generateRefreshToken(clienteID, cliente.Correo, isRememberMe)
	if err != nil {
		return "", "", nil, err
	}

	return newToken, newRefreshToken, cliente, nil
}

// ChangePassword cambia la contraseña de un cliente
func (s *ClienteService) ChangePassword(id int, currentPassword, newPassword string) error {
	// Obtener cliente por ID
	cliente, err := s.clienteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Obtener la contraseña hash actual
	currentHashedPassword, err := s.clienteRepo.GetPasswordByCorreo(cliente.Correo)
	if err != nil {
		return err
	}

	// Verificar contraseña actual
	if !utils.CheckPasswordHash(currentPassword, currentHashedPassword) {
		return errors.New("contraseña actual incorrecta")
	}

	// Hash de la nueva contraseña
	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	return s.clienteRepo.UpdatePassword(id, newHashedPassword)
}
