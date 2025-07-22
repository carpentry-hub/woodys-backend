package routes

import (
	"encoding/json"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
)

// postear un rating de un proyecto
func PostRating(w http.ResponseWriter, r *http.Request) {
	var rating models.Rating
	json.NewDecoder(r.Body).Decode(&rating)

	createdRating := db.DB.Create(&rating)
	err := createdRating.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&rating)
}


// actualizar el rating de un proyecto - Require id
func PutRating(w http.ResponseWriter, r *http.Request){
	
}

// Obtener lista de todos los ratings de un proyecto - Requiere project_id
func GetRating(w http.ResponseWriter, r *http.Request){
	
}