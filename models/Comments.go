package models

import "time"

type Comment struct {
	ID              int8
	CreatedAt       time.Time
	ProjectId       int8
	Content         string
	Rating          int
	UserId          int8
	ParentCommentId int //replies
}