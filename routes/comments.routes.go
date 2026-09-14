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

// GetComment obtiene un comentario - Requiere id
func GetComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	params := mux.Vars(r)
	if err := db.DB.First(&comment, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONMessage(w, http.StatusNotFound, "Comment Not Found")
			return
		}
		log.Printf("Failed to fetch comment: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// GetProjectComments obtiene todos los comentarios de un proyecto - Requiere project_id
func GetProjectComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["id"]

	// chequeo existencia del usuario
	projectID, err := strconv.Atoi(projectIDStr) // cambio de str a int para evitar errores
	if err != nil {
		writeJSONMessage(w, http.StatusNotFound, "Project not found")
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("project_id = ?", projectID).Find(&comments).Error; err != nil {
		log.Printf("Error fetching comments: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Error fetching Comments")
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// PostProjectComment postea un comentario a un proyecto - Requiere project_id y parent_comment_id = 0
func PostProjectComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// Quitar espacios en blanco para contar caracteres
	trimmedContent := strings.TrimSpace(comment.Content)
	contentLength := utf8.RuneCountInString(trimmedContent)

	// Verificar que el comentario no este vacio
	if contentLength == 0 {
		writeJSONMessage(w, http.StatusBadRequest, "content cannot be empty")
		return
	}

	// Verificar que el comentario no supere los 200 caracteres
	if contentLength > 200 {
		writeJSONMessage(w, http.StatusBadRequest, "content cannot exceed 200 characters")
		return
	}

	createdComment := db.DB.Create(&comment)
	err := createdComment.Error

	if err != nil {
		writeJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(&comment); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
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
			writeJSONMessage(w, http.StatusNotFound, "Comment not found")
			return
		}
		log.Printf("Failed to fetch comment: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// contar respuestas al comentario
	var replyCount int64
	if err := db.DB.Model(&models.Comment{}).Where("parent_comment_id = ?", comment.ID).Count(&replyCount).Error; err != nil {
		log.Printf("Failed to count comment replies: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if replyCount == 0 {
		// sin respuestas: borrado fisico
		if err := db.DB.Unscoped().Delete(&comment).Error; err != nil {
			log.Printf("Failed to delete comment: %v", err)
			writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else {
		// con respuestas: marcar el contenido como eliminado para preservar las respuestas
		if err := db.DB.Model(&comment).Update("content", "Comentario Eliminado").Error; err != nil {
			log.Printf("Failed to mark comment as deleted: %v", err)
			writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	writeJSONMessage(w, http.StatusOK, "Comment deleted successfully")
}

// PostCommentReply postea una respuesta a un comentario - Requiere project_id y parent_comment_id
func PostCommentReply(w http.ResponseWriter, r *http.Request) {
	var commentReply models.Comment
	if err := json.NewDecoder(r.Body).Decode(&commentReply); err != nil {
		log.Printf("Failed to decode json: %v", err)
		writeJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	createdComment := db.DB.Create(&commentReply)
	err := createdComment.Error

	if err != nil {
		writeJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(&commentReply); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
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
		writeJSONMessage(w, http.StatusNotFound, "Comment not found")
		return
	}

	// realizacion de la query y manejo de errores
	var comments []models.Comment
	if err := db.DB.Where("parent_comment_id = ?", commentID).Find(&comments).Error; err != nil {
		log.Printf("Error fetching comments: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Error fetching Comments")
		return
	}

	if err := json.NewEncoder(w).Encode(&comments); err != nil {
		log.Printf("Failed to encode json: %v", err)
		writeJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
