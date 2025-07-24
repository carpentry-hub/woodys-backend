package models

import (
	"time"
	"github.com/lib/pq"
)
type Project struct {
	ID        	  int8
	CreatedAt     time.Time
	Owner         int
	Title         string
	UpdatedAt     time.Time
	AverageRating float32
	RatingCount   int
	Materials     pq.StringArray `gorm:"type:varchar[]"`
	Tools         pq.StringArray `gorm:"type:varchar[]"`
	Description   string
	Style         pq.StringArray `gorm:"type:varchar[]"`
	Portrait      string
	Images        pq.StringArray `gorm:"type:varchar[]"`
	Tutorial      string
	TimeToBuild   int
}