package domain

import (
	"errors"
	"time"
)

// Rating represents a rating given to a project by a user
type Rating struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Value     int       `json:"value" gorm:"not null;check:value >= 1 AND value <= 5"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	ProjectID int64     `json:"project_id" gorm:"not null;index"`

	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// RatingCreateRequest represents the data needed to create a new rating
type RatingCreateRequest struct {
	Value     int   `json:"value" validate:"required,min=1,max=5"`
	UserID    int64 `json:"user_id" validate:"required"`
	ProjectID int64 `json:"project_id" validate:"required"`
}

// RatingUpdateRequest represents the data that can be updated for a rating
type RatingUpdateRequest struct {
	Value *int `json:"value,omitempty" validate:"omitempty,min=1,max=5"`
}

// RatingResponse represents the rating data returned to clients
type RatingResponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Value     int       `json:"value"`
	UserID    int64     `json:"user_id"`
	ProjectID int64     `json:"project_id"`
	Username  string    `json:"username"`
}

// RatingStats represents aggregated rating statistics for a project
type RatingStats struct {
	ProjectID     int64       `json:"project_id"`
	AverageRating float32     `json:"average_rating"`
	TotalRatings  int         `json:"total_ratings"`
	Distribution  map[int]int `json:"distribution"` // rating value -> count
}

// Validate validates the rating data
func (r *Rating) Validate() error {
	if r.Value < 1 || r.Value > 5 {
		return errors.New("rating value must be between 1 and 5")
	}
	if r.UserID <= 0 {
		return errors.New("valid user_id is required")
	}
	if r.ProjectID <= 0 {
		return errors.New("valid project_id is required")
	}
	return nil
}

// CanBeUpdatedBy checks if the rating can be updated by the given user
func (r *Rating) CanBeUpdatedBy(userID int64) bool {
	return r.UserID == userID
}

// CanBeDeletedBy checks if the rating can be deleted by the given user
func (r *Rating) CanBeDeletedBy(userID int64) bool {
	return r.UserID == userID
}

// TableName specifies the table name for GORM
func (Rating) TableName() string {
	return "ratings"
}

// Index defines composite unique index for user and project
func (Rating) Indexes() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":    "idx_user_project_rating",
			"columns": []string{"user_id", "project_id"},
			"unique":  true,
		},
	}
}
