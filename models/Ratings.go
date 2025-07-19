package models

import "time"

type Rating struct {
	ID        int8
	CreatedAt time.Time
	Value int8
	UserId int8
	ProjectId int8
	UpdatedAt time.Time
}