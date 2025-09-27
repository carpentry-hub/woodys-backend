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

// GetUsersProjectLists obtiene todas las listas de un usuario - Requiere user_id
func GetUsersProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	// chequeo existencia del usuario
	userID, err := strconv.Atoi(userIDString) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("User not found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var lists []models.ProjectList
	if err := db.DB.Where("user_id = ?", userID).Find(&lists).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error fetching Project Lists")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&lists); err != nil {
		log.Fatalf("Failed to Encode json: %v", err)
	}
}

// GetProjectLists obtiene una lista - Requier id
func GetProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	db.DB.First(&list, params["id"])

	if list.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Project List not found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&list); err != nil {
			log.Fatalf("Failed to Encode json: %v", err)
		}
	}
}

// PostProjectLists postea una lista
func PostProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Fatalf("Failed to Decode json: %v", err)
	}

	createdList := db.DB.Create(&list)
	err := createdList.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // satatus code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&list); err != nil {
			log.Fatalf("Failed to Encode json: %v", err)
		}
	}
}

// AddProjectToList postea un project list item (aniadir un proyecto a una lista)
func AddProjectToList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Fatalf("Failed to Decode json: %v", err)
	}

	createdItem := db.DB.Create(&item)
	err := createdItem.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	}

	if err := json.NewEncoder(w).Encode(&item); err != nil {
		log.Fatalf("Failed to Encode json: %v", err)
	}
}

// PutProjectLists actualiza una lista - Requiere id
func PutProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.ProjectList
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project List Not Found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// leo el updated
	var updated models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte("Error on json file")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	// actualizar campos
	existing.Name = updated.Name
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed to save the Project List")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Fatalf("Failed to Encode json: %v", err)
	}
}

// DeleteProjectList borra una lista - Requiere id
func DeleteProjectList(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	db.DB.First(&list, params["id"])

	if list.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project List not found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&list)
	}
}

// DeleteProjectFromList borra un proyecto de una lista - Requiere id
func DeleteProjectFromList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	params := mux.Vars(r)
	db.DB.First(&item, params["id"])

	if item.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Project List Item not found")); err != nil {
			log.Fatalf("Failed to write Response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&item)
	}
}
