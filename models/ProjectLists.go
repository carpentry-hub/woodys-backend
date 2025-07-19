package models

import "gorm.io/gorm"

type ProjectList struct {
	gorm.Model

	User int `gorm:"not null" json:"user"`
	Name int `gorm:"not null" json:"name"`
	IsPublic bool `gorm:"not null" json:"is_public"`
}