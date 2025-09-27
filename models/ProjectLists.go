package models

import "time"

type ProjectList struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
}
