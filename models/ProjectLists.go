package models

import "time"

type ProjectList struct {
	ID        int8
	CreatedAt time.Time
	UserID    int
	Name      int
	IsPublic  bool
}