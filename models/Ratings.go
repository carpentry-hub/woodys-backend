package models

import "time"

type Rating struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Value     int8      `json:"value"`
	UserId    int8		`json:"user_id"`
	ProjectId int8		`json:"project_id"`
	UpdatedAt time.Time `json:"updated_id"`
}