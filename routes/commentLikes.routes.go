// Package routes proporciona los servicios de la api
package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// commentLikeRequest representa el body recibido al crear una reaccion a un comentario
type commentLikeRequest struct {
	UserID int8   `json:"user_id"`
	Value  string `json:"value"` // "like" o "dislike"
}

// parseLikeValue mapea el string recibido al booleano correspondiente
func parseLikeValue(value string) (bool, bool) {
	switch strings.ToLower(value) {
	case "like":
		return true, true
	case "dislike":
		return false, true
	default:
		return false, false
	}
}

// PostCommentLike crea un like/dislike a un comentario - Requiere comment_id en el path y user_id + value en el body
func PostCommentLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	var request commentLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// mapeo del string recibido a booleano
	value, ok := parseLikeValue(request.Value)
	if !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		json.NewEncoder(w).Encode(map[string]string{"message": "value must be 'like' or 'dislike'"})
		return
	}

	// chequeo de que el usuario no haya reaccionado previamente a este comentario
	var existing models.CommentLike
	err = db.DB.Where("user_id = ? AND comment_id = ?", request.UserID, commentID).First(&existing).Error
	if err == nil {
		w.WriteHeader(http.StatusConflict) // 409
		json.NewEncoder(w).Encode(map[string]string{"message": "User has already reacted to this comment"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to check existing comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commentLike := models.CommentLike{
		UserID:    request.UserID,
		CommentID: int8(commentID),
		Value:     value,
	}
	if err := db.DB.Create(&commentLike).Error; err != nil {
		log.Printf("Failed to create comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&commentLike); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetCommentLikes obtiene los conteos de likes y dislikes de un comentario - Requiere id
func GetCommentLikes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	var likes int64
	if err := db.DB.Model(&models.CommentLike{}).Where("comment_id = ? AND value = ?", commentID, true).Count(&likes).Error; err != nil {
		log.Printf("Error counting comment likes: %v", err)
		http.Error(w, "Error fetching comment likes", http.StatusInternalServerError)
		return
	}

	var dislikes int64
	if err := db.DB.Model(&models.CommentLike{}).Where("comment_id = ? AND value = ?", commentID, false).Count(&dislikes).Error; err != nil {
		log.Printf("Error counting comment dislikes: %v", err)
		http.Error(w, "Error fetching comment dislikes", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]int64{"likes": likes, "dislikes": dislikes}); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetUserCommentLike obtiene la reaccion de un usuario a un comentario - Requiere id y user_id
func GetUserCommentLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentID, err := strconv.Atoi(params["id"]) // cambio de str a int para evitar errores
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var commentLike models.CommentLike
	if err := db.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&commentLike).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment like not found"}); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
			return
		}
		log.Printf("Failed to fetch comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&commentLike); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PutCommentLike invierte el valor de una reaccion (like <-> dislike) - Requiere id del registro
func PutCommentLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que la reaccion exista
	var existing models.CommentLike
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Comment like not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// inversion del valor
	existing.Value = !existing.Value

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Printf("Failed to save comment like: %v", err)
		http.Error(w, "Failed to save the comment like", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&existing); err != nil {
		log.Printf("Failed to encode json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DeleteCommentLike borra la reaccion de un usuario a un comentario - Requiere id del registro
func DeleteCommentLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var commentLike models.CommentLike
	if err := db.DB.First(&commentLike, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment like not found"}); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
			return
		}
		log.Printf("Failed to fetch comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Unscoped().Delete(&commentLike).Error; err != nil {
		log.Printf("Failed to delete comment like: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment like deleted successfully"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
