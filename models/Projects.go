package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model

	Owner int `gorm:"not null" json:"owner"`
	Title string `gorm:"not null" son:"title"`
	AvarageRating float32 `gorm:"default:0" json:"avarage_rating"`
	RatingCount int `gorm:"default:0" json:"rating_count"`
	Materials []string `gorm:"not null" json:"materials"`
	Tools []string `gorm:"not null" json:"tools"`
	Style []string `gorm:"not null" json:"style"`
	Portrait string `gorm:"not null" json:"portrait"`
	Images []string `json:"images"`
	Tutorial string `gorm:"not null" json:"tutorial"`
	Description string `gorm:"not null" json:"description"`
	TimeToBuild int `gorm:"not null" json:"time_to_build"`
}