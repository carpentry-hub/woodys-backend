// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// GetProject obtiene un proyecto - Requiere id
func GetProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	db.DB.First(&project, params["id"])

	if project.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&project); err != nil {
			log.Fatalf("Failed to Encode json: %v", err)
		}
	}
}

// SearchProjects obtene  lista proyectos segun una busqueda - Requiere
func SearchProjects(w http.ResponseWriter, r *http.Request) {
	query := db.DB.Model(&models.Project{})

	style := r.URL.Query().Get("style")
	if style != "" {
		query = query.Where("? = ANY(style)", style).Where("is_public = TRUE")
	}

	env := r.URL.Query().Get("environment")
	if env != "" {
		query = query.Where("? = ANY(environment)", env).Where("is_public = TRUE")
	}

	maxTimeStr := r.URL.Query().Get("max_time_to_build")
	if maxTimeStr != "" {
		maxTime, err := strconv.Atoi(maxTimeStr)
		if err != nil {
			http.Error(w, "max_time_to_build debe ser un número", http.StatusBadRequest)
			return
		}
		query = query.Where("time_to_build <= ?", maxTime).Where("is_public = TRUE")
	}

	var results []models.Project
	if err := query.Find(&results).Error; err != nil {
		http.Error(w, "Error al buscar proyectos", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Fatalf("Failed to Encode json: %v", err)
	}
}

// PostProject postea un proyecto
func PostProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.Fatalf("Failed to Decode json: %v", err)
	}

	createdProject := db.DB.Create(&project)
	err := createdProject.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&project); err != nil {
			log.Fatalf("Failed to Encode json: %v", err)
		}
	}
}

// PutProject actualiza un proyecto
func PutProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Project
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project Not Found")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// lee el proyecto updated
	var updated models.Project
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte("Error on json file")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// actualizar campos
	existing.Title = updated.Title
	existing.Description = updated.Description
	existing.Images = updated.Images
	existing.Materials = updated.Materials
	existing.TimeToBuild = updated.TimeToBuild
	existing.Portrait = updated.Portrait
	existing.Style = updated.Style
	existing.Environment = updated.Environment
	existing.Tools = updated.Tools
	existing.Tutorial = updated.Tutorial
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed to save the project")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Fatalf("Failed to encode json: %v", err)
	}
}

// DeleteProject borra un proyecto - Requiere id
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	db.DB.First(&project, params["id"])

	if project.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project Not Found")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&project)
	}
}
