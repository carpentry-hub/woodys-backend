// Package models proporciona todos los modelos de datos del sistema
package models

import (
	"time"

	"github.com/lib/pq"
)

// Project representa a un proyecto de carpinteria con sus respectivos datos
type Project struct {
	ID            int8           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	Owner         int            `json:"owner"`
	Title         string         `json:"title"`
	UpdatedAt     time.Time      `json:"updated_at"`
	AverageRating float32        `json:"average_rating"`
	RatingCount   int            `json:"rating_count"`
	Materials     pq.StringArray `json:"materials" gorm:"type:varchar[]"`
	Tools         pq.StringArray `json:"tools" gorm:"type:varchar[]"`
	Description   string         `json:"description"`
	Style         pq.StringArray `json:"style" gorm:"type:varchar[]"`
	Environment    pq.StringArray `json:"environment" gorm:"type:text[]"`
	Portrait      string         `json:"portrait"`
	Images        pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Tutorial      string         `json:"tutorial"`
	TimeToBuild   int            `json:"time_to_build"`
	IsPublic      bool           `json:"is_public"`
}
