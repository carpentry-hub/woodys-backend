package domain

import (
	"errors"
	"strings"
	"time"
)

// User represents a user in the system
type User struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Username       string    `json:"username" gorm:"uniqueIndex;not null"`
	Email          string    `json:"email" gorm:"uniqueIndex;not null"`
	Reputation     float32   `json:"reputation" gorm:"default:0"`
	ProfilePicture int64     `json:"profile_picture" gorm:"default:0"`
	FirebaseUID    string    `json:"firebase_uid" gorm:"uniqueIndex;not null"`
}

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	FirebaseUID string `json:"firebase_uid" validate:"required"`
}

// UserUpdateRequest represents the data that can be updated for a user
type UserUpdateRequest struct {
	Username       *string  `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Reputation     *float32 `json:"reputation,omitempty"`
	ProfilePicture *int64   `json:"profile_picture,omitempty"`
}

// Validate validates the user data
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username is required")
	}
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(u.FirebaseUID) == "" {
		return errors.New("firebase_uid is required")
	}
	return nil
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
