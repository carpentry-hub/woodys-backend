package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
)

// GetProfilePictures obtiene todas las fotos de perfil
func GetProfilePictures(w http.ResponseWriter, r *http.Request) {
	var profilePictures []models.ProfilePicture
    if err := db.DB.Find(&profilePictures).Error; err != nil {
        log.Printf("Error fetching profile pictures: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching profile pictures"})
        return
    }
    if err := json.NewEncoder(w).Encode(&profilePictures); err != nil {
        log.Fatalf("Failed to encode profile pictures: %v", err)
    } 
}