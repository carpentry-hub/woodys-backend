package models

import "time"

type User struct {
	ID             int8      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Reputation     float32   `json:"reputation"`
	ProfilePicture int8      `json:"profile_picture"`
	FirebaseUid	   string    `json:"firebase_uid"`
}