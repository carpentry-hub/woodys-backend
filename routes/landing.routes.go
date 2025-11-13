package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
)

type LandingStats struct {
	ProjectsCount  int64   `json:"projects_count"`
	UsersCount     int64   `json:"users_count"`
	TotalRatings   int64   `json:"total_ratings"`
	AverageRating float64 `json:"average_rating"`
}

// GetStats calcula y devuelve las estadisticas clave para la landing page
func GetStats(w http.ResponseWriter, r *http.Request){
	var stats LandingStats

	// Contar proyectos
	if err := db.DB.Model(&models.Project{}).Count(&stats.ProjectsCount).Error; err != nil{
		log.Printf("Error al contar proyectos: %v", err)
		http.Error(w, "Error al contar proyectos", http.StatusInternalServerError)
		return
	}

	// Contar usuarios
	if err := db.DB.Model(&models.User{}).Count(&stats.UsersCount).Error; err != nil {
		log.Printf("Error al contar usuarios: %v", err)
		http.Error(w, "Error al contar usuarios", http.StatusInternalServerError)
		return
	}

	// Contar valoraciones
	if err := db.DB.Model(&models.Rating{}).Count(&stats.TotalRatings).Error; err != nil {
		log.Printf("Error al contar valoraciones: %v", err)
		http.Error(w, "Error al contar valoraciones", http.StatusInternalServerError)
		return
	}

	// Calcular promedio de valoraciones
	if err := db.DB.Model(&models.Rating{}).Select("COALESCE(AVG(value), 0)").Row().Scan(&stats.AverageRating); err != nil {
		log.Printf("Error al calcular el rating promedio: %v", err)
		http.Error(w, "Error al calcular el rating promedio", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Fatalf("Failed to encode stats: %v", err)
	}
}