package models

import "time"

type Project struct {
	ID        int8
	CreatedAt time.Time
	Owner     int
	Title     string
	UpdatedAt time.Time
	AverageRating float32
	RatingCount   int
	Materials     []string
	Tools         []string
	Description   string
	Style         []string
	Portrait      string
	Images        []string
	Tutorial      string
	TimeToBuild   int
}