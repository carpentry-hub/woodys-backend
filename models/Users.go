// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// User representa a un usuario con sus respectivos datos
type User struct {
	ID             int8      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Reputation     float32   `json:"reputation"`
	ProfilePicture int8      `json:"profile_picture"`
	FirebaseUID    string    `json:"firebase_uid"`
}
