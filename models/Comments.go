// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// Comment representa a un comentario o a una respuesta con sus respectivos datos
type Comment struct {
	ID              int8      `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	ProjectID       int8      `json:"project_id"`
	Content         string    `json:"content"`
	Rating          int       `json:"rating"`
	UserID          int8      `json:"user_id"`
	ParentCommentID int       `json:"parent_comment_id"` // replies
}
