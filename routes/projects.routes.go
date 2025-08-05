package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// obtener un proyecto - Requiere id
func GetProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	db.DB.First(&project, params["id"])

	if project.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project Not Found"))
	} else {
		json.NewEncoder(w).Encode(&project)
	}
}

// obtener  lista proyectos segun una busqueda - Requiere
func SearchProjects(w http.ResponseWriter, r *http.Request) {
	query := db.DB.Model(&models.Project{})

	style := r.URL.Query().Get("style")
	if style != "" {
		query = query.Where("? = ANY(style)", style) // ILIKE para búsqueda flexible
	}

	env := r.URL.Query().Get("environment")
	if env != "" {
		query = query.Where("? = ANY(enviroment)", env)
	}

	maxTimeStr := r.URL.Query().Get("max_time_to_build")
	if maxTimeStr != "" {
		maxTime, err := strconv.Atoi(maxTimeStr)
		if err != nil {
			http.Error(w, "max_time_to_build debe ser un número", http.StatusBadRequest)
			return
		}
		query = query.Where("time_to_build <= ?", maxTime)
	}

	var results []models.Project
	if err := query.Find(&results).Error; err != nil {
		http.Error(w, "Error al buscar proyectos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}

// postea un proyecto
func PostProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	json.NewDecoder(r.Body).Decode(&project)

	createdProject := db.DB.Create(&project)
	err := createdProject.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&project)
	}
}

// actualiza un proyecto
func PutProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.Project
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project Not Found"))
		return
	}

	// lee el proyecto updated
	var updated models.Project
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte("Error on json file"))
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
	existing.Enviroment = updated.Enviroment
	existing.Tools = updated.Tools
	existing.Tutorial = updated.Tutorial

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save the project"))
		return
	}

	json.NewEncoder(w).Encode(&existing)
}

// borra un proyecto - Requiere id
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	params := mux.Vars(r)
	db.DB.First(&project, params["id"])

	if project.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project Not Found"))
	} else {
		db.DB.Unscoped().Delete(&project)
	}
}

