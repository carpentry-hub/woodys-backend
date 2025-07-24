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
func GetUser(w http.ResponseWriter, r *http.Request){
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])
	if user.ID == 0{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

// obtener lista de todos los proyectos de un usuario - Requiere id
func GetUserProjects(w http.ResponseWriter, r *http.Request){
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