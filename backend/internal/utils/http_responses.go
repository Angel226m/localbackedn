package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError envía una respuesta HTTP con un error
func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	response := ErrorResponse(message, err)
	RespondWithJSON(w, code, response)
}

// RespondWithJSON envía una respuesta HTTP con un objeto JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error al convertir respuesta a JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
