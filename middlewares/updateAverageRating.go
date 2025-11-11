package middlewares

import (
	"log"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
)

// UpdateAverageRating actualizara la valoracion promedio de un proyecto ante un post o un put de un nuevo rating
func UpdateAverageRating(projectID int8) {
	var result float32

	err := db.DB.Model(&models.Rating{}).Where("project_id = ?", projectID).Select("COALESCE(AVG(value), 0)").Scan(&result).Error

	if err != nil {
		log.Printf("Error al calcular el promedio para el proyecto %d: %v", projectID, err)
		return
	}

	err = db.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("average_rating", result).Error

	if err != nil {
		log.Printf("Error al actualizar el average_rating para el proyecto %d: %v", projectID, err)
	} else {
		log.Printf("average_rating actualizado para el proyecto %d: %f", projectID, result)
	}
}