package utils

import (
	"fmt"
)

// Response estructura general de respuesta
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse crea una respuesta exitosa
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse crea una respuesta de error
func ErrorResponse(message string, err error) Response {
	var errorMsg interface{}
	if err != nil {
		errorMsg = err.Error()
	} else {
		errorMsg = nil
	}

	return Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	}
}

// ValidationErrorResponse crea una respuesta con errores de validación
func ValidationErrorResponse(message string, validationErrors []ValidationError) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   validationErrors,
	}
}

// ErrorWithFields crea un error con campos
func ErrorWithFields(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
