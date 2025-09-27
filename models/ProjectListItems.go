package models

import "time"

type ProjectListItem struct {
	ID            int8      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ProjectListId int8      `json:"project_list_id"`
	ProjectId     int8      `json:"project_id"`
}
