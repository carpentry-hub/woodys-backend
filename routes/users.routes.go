package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// obtener un usuario - Requiere id
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

// obtener un usuario con firebase_uid - Requiere firebase_uid
func GetUserByUID(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	uid := params["firebase_uid"]

	db.DB.Where("firebase_uid = ?", uid).First(&user)

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
	} else {
		json.NewEncoder(w).Encode(map[string]int8{"id": user.ID})
	}
}

// obtener lista de todos los proyectos de un usuario - Requiere id
func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	// chequeo existencia del usuario
	userID, err := strconv.Atoi(userIDString) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// realizacion de la query y manejo de errores
	var projects []models.Project
	if err := db.DB.Where("owner = ?", userID).Find(&projects).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching projects"))
		return
	}

	json.NewEncoder(w).Encode(&projects)
}

// postear un usuario
func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

// actualizar un usuario - Requiere id
func PutUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chqueo que el usuario exista
	var existing models.User
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("User Not Found"))
		return
	}

	// lee el usuario updated
	var updated models.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
		return
	}

	// actualizar campos
	existing.Username = updated.Username
	existing.Reputation = updated.Reputation
	existing.ProfilePicture = updated.ProfilePicture

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save the user"))
		return
	}

	json.NewEncoder(w).Encode(&existing)
}

// borrar un usuario - Requiere id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("User Not Found"))
	} else {
		db.DB.Unscoped().Delete(&user)
	}
}

