// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// ProjectList representa una lista de proyectos con sus respectivos datos
type ProjectList struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	ProjectCount int64  `json:"project_count"`
}
