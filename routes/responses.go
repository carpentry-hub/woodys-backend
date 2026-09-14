// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSONMessage envia una respuesta JSON con el formato {"message": "..."} y el status indicado
func writeJSONMessage(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": message}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
