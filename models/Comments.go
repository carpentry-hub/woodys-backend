package models

import "time"

type Comment struct {
	ID              int8        `json:"id"`
	CreatedAt       time.Time   `json:"created_at"`
	ProjectId       int8		`json:"project_id"`
	Content         string		`json:"content"`
	Rating          int			`json:"rating"`
	UserId          int8		`json:"user_id"`
	ParentCommentId int			`json:"parent_comment_id"` // replies
}

