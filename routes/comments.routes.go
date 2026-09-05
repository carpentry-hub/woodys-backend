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
	"gorm.io/gorm"
)

// GetProjectComments obtiene todos los comentarios de un proyecto - Requiere project_id
func GetProjectComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del usuario
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("project_id = ?", projectID).Find(&comments).Error; err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Error fetching Comments", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PostProjectComment postea un comentario a un proyecto - Requiere project_id y parent_comment_id = 0
func PostProjectComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
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
			log.Printf("Failed to write response: %v", err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DeleteComment borra un comentario de un proyecto - Requiere id
// Si el comentario tiene respuestas, en lugar de borrarlo se marca su contenido como "Comentario Eliminado"
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	params := mux.Vars(r)
	if err := db.DB.First(&comment, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment not found"}); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
			return
		}
		log.Printf("Failed to fetch comment: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// contar respuestas al comentario
	var replyCount int64
	if err := db.DB.Model(&models.Comment{}).Where("parent_comment_id = ?", comment.ID).Count(&replyCount).Error; err != nil {
		log.Printf("Failed to count comment replies: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if replyCount == 0 {
		// sin respuestas: borrado fisico
		if err := db.DB.Unscoped().Delete(&comment).Error; err != nil {
			log.Printf("Failed to delete comment: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// con respuestas: marcar el contenido como eliminado para preservar las respuestas
		if err := db.DB.Model(&comment).Update("content", "Comentario Eliminado").Error; err != nil {
			log.Printf("Failed to mark comment as deleted: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// PostCommentReply postea una respuesta a un comentario - Requiere project_id y parent_comment_id
func PostCommentReply(w http.ResponseWriter, r *http.Request) {
	var commentReply models.Comment
	if err := json.NewDecoder(r.Body).Decode(&commentReply); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdComment := db.DB.Create(&commentReply)
	err := createdComment.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(&commentReply); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetCommentReplies obtiene las respuestas a un comentario - Requiere id
func GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentIDStr := params["id"]

	// chequeo existencia del usuario
	commentID, err := strconv.Atoi(commentIDStr) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("parent_comment_id = ?", commentID).Find(&comments).Error; err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Error fetching Comments", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
