// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/middlewares"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ratingRequest representa el body recibido al crear o actualizar un rating
type ratingRequest struct {
	UserID int8 `json:"user_id"`
	Value  int8 `json:"value"`
}

// PostRating postea un rating de un proyecto - Requiere project_id en el path y user_id + value en el body
func PostRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var request ratingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed to decode json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON format"})
		return
	}

	// chequeo de que el usuario no haya valorado previamente este proyecto
	var existing models.Rating
	err = db.DB.Where("user_id = ? AND project_id = ?", request.UserID, projectID).First(&existing).Error
	if err == nil {
		w.WriteHeader(http.StatusConflict) // 409
		json.NewEncoder(w).Encode(map[string]string{"message": "You have already rated this project"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to check existing rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rating := models.Rating{
		UserID:    request.UserID,
		ProjectID: int8(projectID),
		Value:     request.Value,
	}
	if err := db.DB.Create(&rating).Error; err != nil {
		log.Printf("Failed to create rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ante nuevo rating actualizo average_rating y rating_count en proyecto
	middlewares.UpdateAverageRating(rating.ProjectID)
	middlewares.UpdateRatingCount(rating.ProjectID)

	if err := json.NewEncoder(w).Encode(&rating); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PutRating actualiza el rating de un usuario a un proyecto - Requiere project_id en el path y user_id + value en el body
func PutRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var request ratingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// chequeo que el rating del usuario a este proyecto ya exista
	var existing models.Rating
	if err := db.DB.Where("project_id = ? AND user_id = ?", projectID, request.UserID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// actualizar campos
	existing.Value = request.Value
	existing.UpdatedAt = time.Now()

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save rating: %v", err)
		http.Error(w, "Failed to save the rating", http.StatusInternalServerError)
		return
	}

	// Ante actualizacion de rating actualizo average_rating en proyecto
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

// GetUserProjectRating obtiene el rating de un usuario a un proyecto - Requiere project_id y user_id
func GetUserProjectRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var rating models.Rating
	if err := db.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "Rating not found"}); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
			return
		}
		log.Printf("Failed to fetch rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&rating); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DeleteRating borra el rating de un usuario a un proyecto - Requiere project_id y user_id
func DeleteRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var rating models.Rating
	if err := db.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "Rating not found"}); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
			return
		}
		log.Printf("Failed to fetch rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Unscoped().Delete(&rating).Error; err != nil {
		log.Printf("Failed to delete rating: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ante borrado de rating actualizo average_rating y rating_count en proyecto
	middlewares.UpdateAverageRating(rating.ProjectID)
	middlewares.UpdateRatingCount(rating.ProjectID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Rating deleted successfully"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
