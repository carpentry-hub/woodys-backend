// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetUser obtiene un usuario - Requiere id
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	if err := db.DB.First(&user, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "404: User Not Found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		log.Printf("Failed to encode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetUserByUID obtiene un usuario con firebase_uid - Requiere firebase_uid
func GetUserByUID(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	uid := params["firebase_uid"]

	if err := db.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "404: User Not Found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch user by UID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]int8{"id": user.ID}); err != nil {
		log.Printf("Failed to encode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetUserProjects obtiene lista de todos los proyectos de un usuario - Requiere id
func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	// chequeo existencia del usuario
	userID, err := strconv.Atoi(userIDString) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("User not found")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var projects []models.Project
	if err := db.DB.Where("owner = ?", userID).Find(&projects).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error fetching projects")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&projects); err != nil {
		log.Printf("Failed to encode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PostUser postea un usuario
func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("Failed to write Response: %v", err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		log.Printf("Failed to encode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PutUser actualiza un usuario - Requiere id
func PutUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chqueo que el usuario exista
	var existing models.User
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// lee el usuario updated
	var updated models.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		log.Printf("Failed to decode: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// actualizar campos
	existing.Username = updated.Username
	existing.Reputation = updated.Reputation
	existing.ProfilePicture = updated.ProfilePicture

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save user: %v", err)
		http.Error(w, "Failed to save the user", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DeleteUser borra un usuario - Requiere id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	if err := db.DB.First(&user, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Unscoped().Delete(&user).Error; err != nil {
		log.Printf("Failed to delete user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
