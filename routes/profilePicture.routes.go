package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetProfilePictures obtiene todas las fotos de perfil
func GetProfilePictures(w http.ResponseWriter, r *http.Request) {
	var profilePictures []models.ProfilePicture
    if err := db.DB.Find(&profilePictures).Error; err != nil {
        log.Printf("Error fetching profile pictures: %v", err)
        http.Error(w, "Error fetching profile pictures", http.StatusInternalServerError)
        return
    }
    if err := json.NewEncoder(w).Encode(&profilePictures); err != nil {
        log.Printf("Failed to encode profile pictures: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    } 
}

// GetProfilepicture obtiene una foto de perfil segun id
func GetProfilePictureByID(w http.ResponseWriter, r *http.Request) {
    var picture models.ProfilePicture
    params := mux.Vars(r)
    if err := db.DB.First(&picture, params["id"]).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            http.Error(w, "404: Profile Picture Not Found", http.StatusNotFound)
            return
        }
        log.Printf("Failed to fetch profile picture: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if _, err := w.Write([]byte(picture.Referenced)); err != nil {
        log.Printf("Failed to write response: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}