package models

import "gorm.io/gorm"

type CommentLike struct {
	gorm.Model

	User int `gorm:"not null" json:"user"`
	Comment int `gorm:"not null" json:"comment"`
	Like bool `gorm:"not null" json:"like"`
}