// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// GetProjectComments obtiene todos los comentarios de un proyecto - Requiere project_id
func GetProjectComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del usuario
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Project not found")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("project_id = ?", projectID).Find(&comments).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error fetching Comments")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Fatalf("Failed to encode json: %v", err)
	}
}

// PostProjectComment postea un comentario a un proyecto - Requiere project_id y parent_comment_id = 0
func PostProjectComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Fatalf("Failed to decode json: %v", err)
	}

	// Quitar espacios en blanco para contar caracteres
	trimmedContent := strings.TrimSpace(comment.Content)
	contentLength := utf8.RuneCountInString(trimmedContent)

	// Verificar que el comentario no este vacio
	if contentLength == 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		json.NewEncoder(w).Encode(map[string]string{"message": "content cannot be empty"})
		return
	}

	// Verificar que el comentario no supere los 200 caracteres
	if contentLength > 200 {
		w.WriteHeader(http.StatusBadRequest) // 400
		json.NewEncoder(w).Encode(map[string]string{"message": "content cannot exceed 200 characters"})
		return
	}

	createdComment := db.DB.Create(&comment)
	err := createdComment.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&comment); err != nil {
			log.Fatalf("Failed to encode json: %v", err)
		}
	}
}

// DeleteComment borra un comentario de un proyecto - Requiere id
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	params := mux.Vars(r)
	db.DB.First(&comment, params["id"])

	if comment.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		if _, err := w.Write([]byte("Comment not found")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		db.DB.Unscoped().Delete(&comment)
	}
}

// PostCommentReply postea una respuesta a un comentario - Requiere project_id y parent_comment_id
func PostCommentReply(w http.ResponseWriter, r *http.Request) {
	var commentReply models.Comment
	if err := json.NewDecoder(r.Body).Decode(&commentReply); err != nil {
		log.Fatalf("Failed to decode json: %v", err)
	}

	createdComment := db.DB.Create(&commentReply)
	err := createdComment.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(&commentReply); err != nil {
			log.Fatalf("Failed to encode json: %v", err)
		}
	}
}

// GetCommentReplies obtiene las respuestas a un comentario - Requiere id
func GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentIDStr := params["id"]

	// chequeo existencia del usuario
	commentID, err := strconv.Atoi(commentIDStr) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Comment not found")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("parent_comment_id = ?", commentID).Find(&comments).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error fetching Comments")); err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Fatalf("Failed to encode json: %v", err)
	}
}
