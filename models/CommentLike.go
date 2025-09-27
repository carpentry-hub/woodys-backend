// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// CommentLike representa a un like dado a un comentario/respuesta con sus respectivos datos
type CommentLike struct {
	ID        int8      `json:"id"`
	UserID    int8      `json:"user_id"`
	CommentID int8      `json:"comment_id"`
	Value     int8      `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
