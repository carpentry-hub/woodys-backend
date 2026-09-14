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
    ErrCodeUniqueViolation     = "23505" // "Duplicate"
)

// GetUsersProjectLists obtiene todas las listas de un usuario - Requiere user_id
func GetUsersProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		writeJSONMessage(w, http.StatusNotFound, "User not found")
		return
	}

	var lists []models.ProjectList

	err = db.DB.Model(&models.ProjectList{}).
		Select("project_lists.*, COUNT(project_list_items.project_id) as project_count").
		Joins("LEFT JOIN project_list_items ON project_list_items.project_list_id = project_lists.id").
		Where("project_lists.user_id = ?", userID).
		Group("project_lists.id").
		Order("project_lists.created_at DESC").
		Find(&lists).Error

	if err != nil {
		writeJSONMessage(w, http.StatusInternalServerError, "Error fetching Project Lists")
		return
	}

	if err := json.NewEncoder(w).Encode(&lists); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// GetProjectLists obtiene una lista - Requier id
func GetProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	if err := db.DB.First(&list, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project list not found")
			return
		}
		log.Printf("Failed to fetch project list: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.NewEncoder(w).Encode(&list); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// GetProjectsInList obtiene todos los proyectos dentro de una lista especifica
func GetProjectsInList(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    listIDStr := params["id"]

    listID, err := strconv.ParseInt(listIDStr, 10, 64)
    if err != nil {
        writeJSONMessage(w, http.StatusBadRequest, "Invalid list ID format")
        return
    }

    var items []models.ProjectListItem
    // Encontrar todos los items que pertenecen a esta lista
    if err := db.DB.Where("project_list_id = ?", listID).Find(&items).Error; err != nil {
        log.Printf("Error fetching list items: %v", err)
        writeJSONMessage(w, http.StatusInternalServerError, "Could not fetch list items")
        return
    }

    if len(items) == 0 {
        if err := json.NewEncoder(w).Encode([]models.Project{}); err != nil {
            log.Printf("Failed to encode json: %v", err)
        }
        return
    }

    // Extraer todos los project_id de esos items
    var projectIDs []int8
    for _, item := range items {
        projectIDs = append(projectIDs, item.ProjectID)
    }

    // Buscar todos los proyectos que coincidan con esos id
    var projects []models.Project
    if err := db.DB.Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
        log.Printf("Error fetching projects: %v", err)
        writeJSONMessage(w, http.StatusInternalServerError, "Could not fetch projects")
        return
    }

    if err := json.NewEncoder(w).Encode(&projects); err != nil {
        log.Printf("Failed to encode json: %v", err)
        writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }
}

// PostProjectLists postea una lista
func PostProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	trimmedName := strings.TrimSpace(list.Name)
    nameLength := utf8.RuneCountInString(trimmedName)

    if nameLength == 0 {
        writeJSONMessage(w, http.StatusBadRequest, "name cannot be empty")
        return
    }

    if nameLength > 50 {
        writeJSONMessage(w, http.StatusBadRequest, "name cannot exceed 50 characters")
        return
    }

	createdList := db.DB.Create(&list)
	err := createdList.Error

	if err != nil {
		log.Printf("Failed to create project list: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(&list); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// AddProjectToList postea un project list item (anadir un proyecto a una lista)
 func AddProjectToList(w http.ResponseWriter, r *http.Request) {
    var item models.ProjectListItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        log.Printf("Failed to decode json: %v", err)
        writeJSONMessage(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    createdItem := db.DB.Create(&item)
    err := createdItem.Error
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            switch pgErr.Code {
            case ErrCodeForeignKeyViolation: // Error para "Not Found"
                writeJSONMessage(w, http.StatusNotFound, "Project or List not found")
                return
            case ErrCodeUniqueViolation: // Error para "Duplicate"
                writeJSONMessage(w, http.StatusConflict, "Project is already in this list")
                return
            }
        }

        // Para otros errores
        log.Printf("Failed to add project to list: %v", err)
        writeJSONMessage(w, http.StatusInternalServerError, "Could not add project to list")
        return
    }

    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(&item); err != nil {
        log.Printf("Failed to encode json: %v", err)
    }

} 

// PutProjectLists actualiza una lista - Requiere id
func PutProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.ProjectList
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project list not found")
			return
		}
		log.Printf("Failed to fetch project list: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// leo el updated
	var updated models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// actualizar campos
	existing.Name = updated.Name
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save project list: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Failed to save project list")
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// DeleteProjectList borra una lista - Requiere id
func DeleteProjectList(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	if err := db.DB.First(&list, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Project list not found")
			return
		}
		log.Printf("Failed to fetch project list: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := db.DB.Unscoped().Delete(&list).Error; err != nil {
		log.Printf("Failed to delete project list: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	writeJSONMessage(w, http.StatusOK, "Project list deleted successfully")
}

// DeleteProjectFromList borra un proyecto de una lista - Requiere id
func DeleteProjectFromList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	params := mux.Vars(r)
	listIDStr := params["list_id"]
    projectIDStr := params["project_id"]

	listID, err := strconv.Atoi(listIDStr)
    if err != nil {
        writeJSONMessage(w, http.StatusBadRequest, "Invalid list ID format")
        return
    }

	projectID, err := strconv.Atoi(projectIDStr)
    if err != nil {
        writeJSONMessage(w, http.StatusBadRequest, "Invalid project ID format")
        return
    }

	result := db.DB.Where("project_list_id = ? AND project_id = ?", listID, projectID).First(&item)

	// Respuestas errores 
	if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            writeJSONMessage(w, http.StatusNotFound, "Project is not in this list")
        } else {
            log.Printf("DB error finding item: %v", result.Error)
            writeJSONMessage(w, http.StatusInternalServerError, "Database error")
        }
        return
    }

	// Eliminar item
	if err := db.DB.Unscoped().Delete(&item).Error; err != nil {
        log.Printf("DB error deleting item: %v", err)
        writeJSONMessage(w, http.StatusInternalServerError, "Failed to delete item from list")
        return
    }

    // Respuesta de exito
    writeJSONMessage(w, http.StatusOK, "Project removed from list successfully")
}
