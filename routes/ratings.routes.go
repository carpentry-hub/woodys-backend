package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
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
	} else {
		json.NewEncoder(w).Encode(&rating)
	}
}


// actualizar el rating de un proyecto - Require id
func PutRating(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Rating
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Rating Not Found"))
		return
	}

	// leo el updated
	var updated models.Rating
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte("Error on json file"))
		return
	}

	// actualizar campos
	existing.Value = updated.Value
	existing.UpdatedAt = updated.UpdatedAt
	

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError) 
		w.Write([]byte("Failed to save the rating"))
		return
	}

	json.NewEncoder(w).Encode(&existing)
}

// Obtener lista de todos los ratings de un proyecto - Requiere project_id
func GetRating(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del proyecto
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Project not found"))
		return
	}

	// realizacion de la query y manejo de errores
	var ratings []models.Rating
	if err := db.DB.Where("project_id = ?", projectID).Find(&ratings).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching ratings"))
		return
	}

	json.NewEncoder(w).Encode(&ratings)
}