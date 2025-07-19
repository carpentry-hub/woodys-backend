package models

import "time"

type CommentLike struct {
	ID        int8
	UserId    int8
	CommentId int8
	Value     int8
	CreatedAt time.Time
}