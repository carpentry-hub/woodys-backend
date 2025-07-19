package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model

	User int `gorm:"not null" json:"user"`
	Project int `gorm:"not null" json:"project"`
	Rating int `gorm:"not null" json:"rating"` 
	Content string `gorm:"not null" json:"content"`
	ParentComment int `json:"parent_comment"` //replies
}