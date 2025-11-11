package middlewares

import (
	"log"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
)

func UpdateRatingCount(projectID int8) {
	var count int64

	err := db.DB.Model(&models.Rating{}).Where("project_id = ?", projectID).Count(&count).Error

	if err != nil {
		log.Printf("Error contando ratings para el proyecto %d: %v", projectID, err)
        return
	}

	err = db.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("rating_count", count).Error

	if err != nil {
		log.Printf("Error actualizando rating_count para el proyecto %d: %v", projectID, err)
	}
}