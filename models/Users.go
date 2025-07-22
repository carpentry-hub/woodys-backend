package models

import "time"

type User struct {
	ID             int8 
	CreatedAt      time.Time 
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Reputation     float32 `son:"reputation"`
	ProfilePicture int8    `json:"profile_picture"`
}