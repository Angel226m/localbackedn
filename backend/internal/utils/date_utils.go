package utils

import (
	"errors"
	"time"
)

// ParseDateString convierte una cadena de fecha a un objeto Time
func ParseDateString(dateString string) (time.Time, error) {
	// Intentar varios formatos comunes
	formats := []string{
		"2006-01-02",          // ISO
		"02/01/2006",          // DD/MM/YYYY
		"01/02/2006",          // MM/DD/YYYY
		"2006-01-02 15:04:05", // Con hora
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, errors.New("formato de fecha no reconocido")
}

// FormatDate formatea una fecha en un formato específico
func FormatDate(t time.Time, format string) string {
	switch format {
	case "iso":
		return t.Format("2006-01-02")
	case "human":
		return t.Format("02/01/2006")
	case "datetime":
		return t.Format("2006-01-02 15:04:05")
	default:
		return t.Format("2006-01-02")
	}
}

// IsWeekend verifica si una fecha cae en fin de semana
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// AddBusinessDays añade días hábiles (excluyendo fines de semana)
func AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	for i := 0; i < days; {
		result = result.AddDate(0, 0, 1)
		if !IsWeekend(result) {
			i++
		}
	}
	return result
}

// GetDayName obtiene el nombre del día de la semana en español
func GetDayName(t time.Time) string {
	days := []string{
		"Domingo",
		"Lunes",
		"Martes",
		"Miércoles",
		"Jueves",
		"Viernes",
		"Sábado",
	}
	return days[t.Weekday()]
}

// GetMonthName obtiene el nombre del mes en español
func GetMonthName(t time.Time) string {
	months := []string{
		"Enero",
		"Febrero",
		"Marzo",
		"Abril",
		"Mayo",
		"Junio",
		"Julio",
		"Agosto",
		"Septiembre",
		"Octubre",
		"Noviembre",
		"Diciembre",
	}
	return months[t.Month()-1]
}

// FechaEnRango verifica si una fecha está en un rango específico
func FechaEnRango(fecha, inicio, fin time.Time) bool {
	fecha = time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	inicio = time.Date(inicio.Year(), inicio.Month(), inicio.Day(), 0, 0, 0, 0, inicio.Location())

	// Si fin es cero, solo comprobar que la fecha sea >= inicio
	if fin.IsZero() {
		return !fecha.Before(inicio)
	}

	fin = time.Date(fin.Year(), fin.Month(), fin.Day(), 0, 0, 0, 0, fin.Location())
	return !fecha.Before(inicio) && !fecha.After(fin)
}

// ParseTime convierte una cadena HH:MM a time.Time
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}
