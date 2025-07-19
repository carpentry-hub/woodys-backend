package models

import "gorm.io/gorm"

type ProjectListItem struct {
	gorm.Model

	Project int `gorm:"not null" json:"project"`
	ProjectList int `gorm:"not null" json:"project_list"`
}