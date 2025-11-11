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
	MainMaterial  string         `json:"main_material"`
	Materials     pq.StringArray `json:"materials" gorm:"type:varchar[]"`
	Height        float32        `json:"height"`
	Length        float32        `json:"length"`
	Width		  float32        `json:"width"`
	Tools         pq.StringArray `json:"tools" gorm:"type:varchar[]"`
	Description   string         `json:"description"`
	Style         pq.StringArray `json:"style" gorm:"type:varchar[]"`
	Environment   string         `json:"environment"`
	Portrait      string         `json:"portrait"`
	Images        pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Tutorial      string         `json:"tutorial"`
	TimeToBuild   int            `json:"time_to_build"`
	IsPublic      bool           `json:"is_public"`
}
