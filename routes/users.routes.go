// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// GetUser obtiene un usuario - Requiere id
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404: User Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&user); err != nil {
			log.Fatalf("Failed to encode: %v", err)
		}
	}
}

// GetUserByUID obtiene un usuario con firebase_uid - Requiere firebase_uid
func GetUserByUID(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	uid := params["firebase_uid"]

	db.DB.Where("firebase_uid = ?", uid).First(&user)

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404: User Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(map[string]int8{"id": user.ID}); err != nil {
			log.Fatalf("Failed to encode: %v", err)
		}
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
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var projects []models.Project
	if err := db.DB.Where("owner = ?", userID).Find(&projects).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error fetching projects")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&projects); err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
}

// PostUser postea un usuario
func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&user); err != nil {
			log.Fatalf("Failed to encode: %v", err)
		}
	}
}

// PutUser actualiza un usuario - Requiere id
func PutUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chqueo que el usuario exista
	var existing models.User
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("User Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// lee el usuario updated
	var updated models.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// actualizar campos
	existing.Username = updated.Username
	existing.Reputation = updated.Reputation
	existing.ProfilePicture = updated.ProfilePicture

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed to save the user")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
}

// DeleteUser borra un usuario - Requiere id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("User Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&user)
	}
}
