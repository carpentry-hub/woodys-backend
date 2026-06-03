// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/middlewares"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Constante para chequear no haya ratings con el par user_id y project_id duplicados
const (
    ErrCodeUniqueViolation = "23505"
)

// PostRating postea un rating de un proyecto
func PostRating(w http.ResponseWriter, r *http.Request) {
	var rating models.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
        log.Printf("Failed to decode json: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON format"})
        return
    }

	createdRating := db.DB.Create(&rating)
	err := createdRating.Error
	if err != nil {
		var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            if pgErr.Code == ErrCodeUniqueViolation {
                w.WriteHeader(http.StatusConflict)
                json.NewEncoder(w).Encode(map[string]string{"message": "You have already rated this project"})
                return
            }
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"message": "Could not create rating"})
        log.Printf("Failed to write response: %v", err)
	} else {
		// Ante nuevo rating actualizo average_rating y rating_count en proyecto
		middlewares.UpdateAverageRating(rating.ProjectID)
		middlewares.UpdateRatingCount(rating.ProjectID)
		if err := json.NewEncoder(w).Encode(&rating); err != nil {
			log.Printf("Failed to encode json: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// PutRating actualiza el rating de un proyecto - Require id
func PutRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Rating
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// leo el updated
	var updated models.Rating
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// actualizar campos
	existing.Value = updated.Value
	existing.UpdatedAt = updated.UpdatedAt

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save rating: %v", err)
		http.Error(w, "Failed to save the rating", http.StatusInternalServerError)
		return
	}

	// Ante actulizacion de rating actualizo average_rating en proyecto
	middlewares.UpdateAverageRating(existing.ProjectID)
	
	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetRating obtiene lista de todos los ratings de un proyecto - Requiere project_id
func GetRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del proyecto
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// realizacion de la query y manejo de errores
	var ratings []models.Rating
	if err := db.DB.Where("project_id = ?", projectID).Find(&ratings).Error; err != nil {
		log.Printf("Error fetching ratings: %v", err)
		http.Error(w, "Error fetching ratings", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&ratings); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
