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
			writeJSONMessage(w, http.StatusNotFound, "404: User Not Found")
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		log.Printf("Failed to encode: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
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
			writeJSONMessage(w, http.StatusNotFound, "404: User Not Found")
			return
		}
		log.Printf("Failed to fetch user by UID: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]int8{"id": user.ID}); err != nil {
		log.Printf("Failed to encode: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
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
		writeJSONMessage(w, http.StatusNotFound, "User not found")
		return
	}

	// realizacion de la query y manejo de errores
	var projects []models.Project
	if err := db.DB.Where("owner = ?", userID).Find(&projects).Error; err != nil {
		writeJSONMessage(w, http.StatusInternalServerError, "Error fetching projects")
		return
	}

	if err := json.NewEncoder(w).Encode(&projects); err != nil {
		log.Printf("Failed to encode: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// PostUser postea un usuario
func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		writeJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		log.Printf("Failed to encode: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
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
			writeJSONMessage(w, http.StatusNotFound, "User Not Found")
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// lee el usuario updated
	var updated models.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		log.Printf("Failed to decode: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// actualizar campos
	existing.Username = updated.Username
	existing.Reputation = updated.Reputation
	existing.ProfilePicture = updated.ProfilePicture

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save user: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Failed to save the user")
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// DeleteUser borra un usuario - Requiere id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	if err := db.DB.First(&user, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "User Not Found")
			return
		}
		log.Printf("Failed to fetch user: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := db.DB.Unscoped().Delete(&user).Error; err != nil {
		log.Printf("Failed to delete user: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	writeJSONMessage(w, http.StatusOK, "User deleted successfully")
}
