package models

import "time"

type ProjectList struct {
	ID        int8
	CreatedAt time.Time
	UserID    int
	Name      string
	IsPublic  bool
}