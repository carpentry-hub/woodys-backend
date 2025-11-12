// Package models proporciona todos los modelos de datos del sistema
package models

import "time"

// ProfilePicture representa el modelo de una foto de perfil por defecto del sistema.
type ProfilePicture struct {
	ID        int8  `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	Referenced string `json:"referenced"`
}