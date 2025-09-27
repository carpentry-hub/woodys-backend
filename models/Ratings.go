// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// Rating representa a una valoracion creada por un usuario con sus respectivos datos
type Rating struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Value     int8      `json:"value"`
	UserID    int8      `json:"user_id"`
	ProjectID int8      `json:"project_id"`
	UpdatedAt time.Time `json:"updated_id"`
}
