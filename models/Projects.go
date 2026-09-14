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
	Owner         int            `json:"owner" gorm:"not null" validate:"required"`
	Title         string         `json:"title" gorm:"not null" validate:"required"`
	UpdatedAt     time.Time      `json:"updated_at"`
	AverageRating float32        `json:"average_rating"`
	RatingCount   int            `json:"rating_count"`
	MainMaterial  string         `json:"main_material" gorm:"not null" validate:"required"`
	Materials     pq.StringArray `json:"materials" gorm:"not null;type:varchar[]" validate:"required"`
	Height        float32        `json:"height" gorm:"not null" validate:"required"`
	Length        float32        `json:"length" gorm:"not null" validate:"required"`
	Width         float32        `json:"width" gorm:"not null" validate:"required"`
	Tools         pq.StringArray `json:"tools" gorm:"not null;type:varchar[]" validate:"required"`
	Description   string         `json:"description" gorm:"not null" validate:"required"`
	Style         pq.StringArray `json:"style" gorm:"not null;type:varchar[]" validate:"required"`
	Environment   string         `json:"environment" gorm:"not null" validate:"required"`
	Portrait      string         `json:"portrait" gorm:"not null" validate:"required"`
	Images        pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Tutorial      string         `json:"tutorial" gorm:"not null" validate:"required"`
	TimeToBuild   int            `json:"time_to_build" gorm:"not null" validate:"required"`
	IsPublic      bool           `json:"is_public" gorm:"not null;default:true"`
}
