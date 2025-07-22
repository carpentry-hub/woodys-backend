package routes

import (
	"encoding/json"
	"net/http"

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
	
}