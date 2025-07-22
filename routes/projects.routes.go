package routes

import (
	"encoding/json"
	"net/http"

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
func GetFilterProject(w http.ResponseWriter, r *http.Request){
	
}

// postea un proyecto
func PostProject(w http.ResponseWriter, r *http.Request){
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
func PutProject(w http.ResponseWriter, r *http.Request){
	
}

// borra un proyecto - Requiere id
func DeleteProject(w http.ResponseWriter, r *http.Request){
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