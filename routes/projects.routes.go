// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetProject obtiene un proyecto - Requiere id
func GetProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	if err := db.DB.First(&project, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project Not Found")
			return
		}
		log.Printf("Failed to fetch project: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.NewEncoder(w).Encode(&project); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// SearchProjects obtene  lista proyectos segun una busqueda - Requiere
func SearchProjects(w http.ResponseWriter, r *http.Request) {
	query := db.DB.Model(&models.Project{}).Where("is_public = TRUE")

	style := r.URL.Query().Get("style")
	if style != "" {
		query = query.Where("? = ANY(style)", style)
	}

	env := r.URL.Query().Get("environment")
	if env != "" {
		query = query.Where("environment = ?", env)
	}

	mat := r.URL.Query().Get("materials")
	if mat != "" {
		query = query.Where("? = ANY(materials)", mat)
	}

	title := r.URL.Query().Get("title")
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	maxTimeStr := r.URL.Query().Get("max_time_to_build")
	if maxTimeStr != "" {
		maxTime, err := strconv.Atoi(maxTimeStr)
		if err != nil {
			writeJSONMessage(w, http.StatusBadRequest, "max_time_to_build debe ser un número")
			return
		}
		query = query.Where("time_to_build <= ?", maxTime)
	}

	var results []models.Project
	if err := query.Find(&results).Error; err != nil {
		writeJSONMessage(w, http.StatusInternalServerError, "Error al buscar proyectos")
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// PostProject postea un proyecto
func PostProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	createdProject := db.DB.Create(&project)
	err := createdProject.Error
	if err != nil {
		writeJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(&project); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// PutProject actualiza un proyecto
func PutProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Project
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project Not Found")
			return
		}
		log.Printf("Failed to fetch project: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// lee el proyecto updated
	var updated models.Project
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// actualizar campos
	existing.Title = updated.Title
	existing.Description = updated.Description
	existing.Images = updated.Images
	existing.MainMaterial = updated.MainMaterial
	existing.Materials = updated.Materials
	existing.Height = updated.Height
	existing.Width = updated.Width
	existing.Length = updated.Length
	existing.TimeToBuild = updated.TimeToBuild
	existing.Portrait = updated.Portrait
	existing.Style = updated.Style
	existing.Environment = updated.Environment
	existing.Tools = updated.Tools
	existing.Tutorial = updated.Tutorial
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save project: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Failed to save the project")
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// DeleteProject borra un proyecto - Requiere id
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	if err := db.DB.First(&project, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project Not Found")
			return
		}
		log.Printf("Failed to fetch project: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := db.DB.Unscoped().Delete(&project).Error; err != nil {
		log.Printf("Failed to delete project: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	writeJSONMessage(w, http.StatusOK, "Project deleted successfully")
}
