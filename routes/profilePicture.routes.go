package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
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

// GetProfilepicture obtiene una foto de perfil segun id
func GetProfilePictureByID(w http.ResponseWriter, r *http.Request) {
    var picture models.ProfilePicture
    params := mux.Vars(r)
    db.DB.First(&picture, params["id"])
    if picture.ID == 0 {
        w.WriteHeader(http.StatusNotFound)
        if _, err := w.Write([]byte("404: Profile Picture Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
    } else {

		if _, err := w.Write([]byte(picture.Referenced)); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
    }
}