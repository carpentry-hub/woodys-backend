package domain

import (
	"errors"
	"strings"
	"time"
)

// Comment represents a comment on a project
type Comment struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	ProjectID       int64     `json:"project_id" gorm:"not null;index"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	Rating          int       `json:"rating" gorm:"default:0;check:rating >= 1 AND rating <= 5"`
	UserID          int64     `json:"user_id" gorm:"not null;index"`
	ParentCommentID *int64    `json:"parent_comment_id,omitempty" gorm:"index"` // For replies
	IsDeleted       bool      `json:"is_deleted" gorm:"default:false"`

	// Relationships
	User    User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Project Project   `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentCommentID"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentCommentID"`
}

// CommentCreateRequest represents the data needed to create a new comment
type CommentCreateRequest struct {
	ProjectID       int64  `json:"project_id" validate:"required"`
	Content         string `json:"content" validate:"required,min=1,max=1000"`
	Rating          int    `json:"rating" validate:"min=1,max=5"`
	UserID          int64  `json:"user_id" validate:"required"`
	ParentCommentID *int64 `json:"parent_comment_id,omitempty"`
}

// CommentUpdateRequest represents the data that can be updated for a comment
type CommentUpdateRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1,max=1000"`
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
}

// CommentResponse represents the comment data returned to clients
type CommentResponse struct {
	ID              int64             `json:"id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ProjectID       int64             `json:"project_id"`
	Content         string            `json:"content"`
	Rating          int               `json:"rating"`
	UserID          int64             `json:"user_id"`
	ParentCommentID *int64            `json:"parent_comment_id,omitempty"`
	Username        string            `json:"username"`
	UserReputation  float32           `json:"user_reputation"`
	ReplyCount      int               `json:"reply_count"`
	Replies         []CommentResponse `json:"replies,omitempty"`
}

// Validate validates the comment data
func (c *Comment) Validate() error {
	if strings.TrimSpace(c.Content) == "" {
		return errors.New("content is required")
	}
	if len(c.Content) > 1000 {
		return errors.New("content cannot exceed 1000 characters")
	}
	if c.ProjectID <= 0 {
		return errors.New("valid project_id is required")
	}
	if c.UserID <= 0 {
		return errors.New("valid user_id is required")
	}
	if c.Rating < 0 || c.Rating > 5 {
		return errors.New("rating must be between 1 and 5, or 0 for no rating")
	}
	return nil
}

// IsReply checks if the comment is a reply to another comment
func (c *Comment) IsReply() bool {
	return c.ParentCommentID != nil
}

// CanBeEditedBy checks if the comment can be edited by the given user
func (c *Comment) CanBeEditedBy(userID int64) bool {
	return c.UserID == userID && !c.IsDeleted
}

// CanBeDeletedBy checks if the comment can be deleted by the given user
func (c *Comment) CanBeDeletedBy(userID int64) bool {
	return c.UserID == userID && !c.IsDeleted
}

// SoftDelete marks the comment as deleted without removing it from database
func (c *Comment) SoftDelete() {
	c.IsDeleted = true
	c.Content = "[This comment has been deleted]"
	c.UpdatedAt = time.Now()
}

// TableName specifies the table name for GORM
func (Comment) TableName() string {
	return "comments"
}
