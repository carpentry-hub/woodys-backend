package models

import "time"

type CommentLike struct {
	ID        int8      `json:"id"`
	UserId    int8      `json:"user_id"`
	CommentId int8      `json:"comment_id"`
	Value     int8      `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
