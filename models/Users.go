package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"type:varchar(30);not null;unique" json:"username"`
	PhoneNumber int `json:"phone_number"`
	Email string `gorm:"not null;unique" json:"email"`
	Reputation float32 `gorm:"default:0" json:"reputation"`
	ProfilePicture int8 `gorm:"default:0" json:"profile_picture"`
}