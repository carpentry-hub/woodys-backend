package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// postear un rating de un proyecto
func PostRating(w http.ResponseWriter, r *http.Request) {
	var rating models.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil{
		log.Fatalf("Failed to decode json: %v",err)
	}

	createdRating := db.DB.Create(&rating)
	err := createdRating.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		if _,err := w.Write([]byte(err.Error())); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&rating); err != nil{
			log.Fatalf("Failed to encode json: %v",err)
		}
	}
}

// actualizar el rating de un proyecto - Require id
func PutRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Rating
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _,err := w.Write([]byte("Rating Not Found")); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
		return
	}

	// leo el updated
	var updated models.Rating
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _,err := w.Write([]byte("Error on json file")); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
		return
	}

	// actualizar campos
	existing.Value = updated.Value
	existing.UpdatedAt = updated.UpdatedAt

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _,err := w.Write([]byte("Failed to save the rating")); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil{
		log.Fatalf("Failed to encode json: %v",err)
	}
}

// Obtener lista de todos los ratings de un proyecto - Requiere project_id
func GetRating(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del proyecto
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _,err := w.Write([]byte("Project not found")); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var ratings []models.Rating
	if err := db.DB.Where("project_id = ?", projectID).Find(&ratings).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _,err := w.Write([]byte("Error fetching ratings")); err != nil{
			log.Fatalf("Failed to write Response: %v",err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&ratings); err != nil{
			log.Fatalf("Failed to encode json: %v",err)
		}
}

