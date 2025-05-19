package utils

import (
	"errors"
	"fmt"
	"sistema-tours/internal/config"
	"sistema-tours/internal/entidades"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims define los claims del token JWT
type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT genera un nuevo token JWT para el usuario
func GenerateJWT(usuario *entidades.Usuario, config *config.Config) (string, error) {
	// Configurar claims estándar
	claims := TokenClaims{
		UserID: usuario.ID,
		Email:  usuario.Correo,
		Role:   usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", usuario.ID),
		},
	}

	// Crear token con los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la llave secreta
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken genera un token de actualización
func GenerateRefreshToken(usuario *entidades.Usuario, config *config.Config) (string, error) {
	// Configurar claims estándar para refresh token (mayor duración)
	claims := TokenClaims{
		UserID: usuario.ID,
		Email:  usuario.Correo,
		Role:   usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 días
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sistema-tours",
			Subject:   fmt.Sprintf("%d", usuario.ID),
		},
	}

	// Crear token con los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la llave secreta
	tokenString, err := token.SignedString([]byte(config.JWTRefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT
func ValidateToken(tokenString string, config *config.Config) (*TokenClaims, error) {
	// Parse del token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}

		// Retornar llave de firma
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar si el token es válido
	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Obtener claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("no se pudieron extraer los claims del token")
	}

	return claims, nil
}

// ValidateRefreshToken valida un token de actualización
func ValidateRefreshToken(tokenString string, config *config.Config) (*TokenClaims, error) {
	// Parse del token de actualización
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}

		// Retornar llave de firma para refresh token
		return []byte(config.JWTRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar si el token es válido
	if !token.Valid {
		return nil, errors.New("token de actualización inválido")
	}

	// Obtener claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("no se pudieron extraer los claims del token")
	}

	return claims, nil
}
