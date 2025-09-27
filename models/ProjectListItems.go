// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// ProjectListItem representa al item de una lista de proyectos con sus respectivos datos
type ProjectListItem struct {
	ID            int8      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ProjectListID int8      `json:"project_list_id"`
	ProjectID     int8      `json:"project_id"`
}
