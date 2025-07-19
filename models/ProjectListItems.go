package models

import "time"

type ProjectListItem struct {
	ID            int8
	CreatedAt     time.Time
	ProjectListId int8
	ProjectId     int8
}