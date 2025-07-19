package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model

	User int `gorm:"not null" json:"user"`
	Project int `gorm:"not null" json:"project"`
	Value int `gorm:"not null" json:"value"`
}