// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
    ErrCodeForeignKeyViolation = "23503" // "Not Found"
)

// GetUsersProjectLists obtiene todas las listas de un usuario - Requiere user_id
func GetUsersProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	// chequeo existencia del usuario
	userID, err := strconv.Atoi(userIDString) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "User not found"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
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
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Project list not found"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
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
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON format"})
	}

	trimmedName := strings.TrimSpace(list.Name)
    nameLength := utf8.RuneCountInString(trimmedName)

    if nameLength == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "name cannot be empty"})
        return
    }

    if nameLength > 50 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "name cannot exceed 50 characters"})
        return
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

// AddProjectToList postea un project list item (anadir un proyecto a una lista)
 func AddProjectToList(w http.ResponseWriter, r *http.Request) {
    var item models.ProjectListItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        log.Fatalf("Failed to Decode json: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON format"})
        return
    }


    createdItem := db.DB.Create(&item)
    err := createdItem.Error
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            switch pgErr.Code {
            case ErrCodeForeignKeyViolation: // Error para "Not Found"
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"message": "Project or List not found"})
                return
            case ErrCodeUniqueViolation: // Error para "Duplicate"
                w.WriteHeader(http.StatusConflict)
                json.NewEncoder(w).Encode(map[string]string{"message": "Project is already in this list"})
                return
            }
        }

        // Para otros errores
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"message": "Could not add project to list"})
        log.Printf("Failed to write response: %v", err)
        return
    }

    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(&item); err != nil {
        log.Printf("Failed to Encode json: %v", err)
    }

} 

// PutProjectLists actualiza una lista - Requiere id
func PutProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.ProjectList
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Project list not found"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// leo el updated
	var updated models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Error on json file"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// actualizar campos
	existing.Name = updated.Name
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save project list"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
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
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Project list not found"}); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&list)
	}
}

// DeleteProjectFromList borra un proyecto de una lista - Requiere id
func DeleteProjectFromList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	params := mux.Vars(r)
	listIDStr := params["list_id"]
    projectIDStr := params["project_id"]

	listID, err := strconv.Atoi(listIDStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Invalid list ID format"})
        return
    }

	projectID, err := strconv.Atoi(projectIDStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Invalid project ID format"})
        return
    }

	result := db.DB.Where("project_list_id = ? AND project_id = ?", listID, projectID).First(&item)

	// Respuestas errores 
	if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(map[string]string{"message": "Project is not in this list"}) // 404 not found
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"message": "Database error"})
            log.Printf("DB error finding item: %v", result.Error)
        }
        return
    }

	// Eliminar item
	if err := db.DB.Unscoped().Delete(&item).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete item from list"}) // Respuesta de error
        log.Printf("DB error deleting item: %v", err)
        return
    }

    // Respuesta de exito
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Project removed from list successfully"})
}
