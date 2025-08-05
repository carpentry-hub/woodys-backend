package domain

import (
	"errors"
	"strings"
	"time"
)

// ProjectList represents a user's collection of projects
type ProjectList struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null"`
	IsPublic  bool      `json:"is_public" gorm:"default:false"`

	// Relationships
	User  User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items []ProjectListItem `json:"items,omitempty" gorm:"foreignKey:ProjectListID"`
}

// ProjectListItem represents an item in a project list
type ProjectListItem struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	ProjectListID int64     `json:"project_list_id" gorm:"not null;index"`
	ProjectID     int64     `json:"project_id" gorm:"not null;index"`

	// Relationships
	ProjectList ProjectList `json:"project_list,omitempty" gorm:"foreignKey:ProjectListID"`
	Project     Project     `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectListCreateRequest represents the data needed to create a new project list
type ProjectListCreateRequest struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	IsPublic bool   `json:"is_public"`
}

// ProjectListUpdateRequest represents the data that can be updated for a project list
type ProjectListUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

// ProjectListResponse represents the project list data returned to clients
type ProjectListResponse struct {
	ID           int64                    `json:"id"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	UserID       int64                    `json:"user_id"`
	Name         string                   `json:"name"`
	IsPublic     bool                     `json:"is_public"`
	ProjectCount int                      `json:"project_count"`
	Projects     []ProjectListItemProject `json:"projects,omitempty"`
}

// ProjectListItemProject represents project data within a list
type ProjectListItemProject struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Portrait      string    `json:"portrait"`
	AverageRating float32   `json:"average_rating"`
	RatingCount   int       `json:"rating_count"`
	AddedAt       time.Time `json:"added_at"`
}

// AddProjectRequest represents the data needed to add a project to a list
type AddProjectRequest struct {
	ProjectListID int64 `json:"project_list_id" validate:"required"`
	ProjectID     int64 `json:"project_id" validate:"required"`
}

// Validate validates the project list data
func (pl *ProjectList) Validate() error {
	if strings.TrimSpace(pl.Name) == "" {
		return errors.New("name is required")
	}
	if len(pl.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	if pl.UserID <= 0 {
		return errors.New("valid user_id is required")
	}
	return nil
}

// CanBeAccessedBy checks if the project list can be accessed by the given user
func (pl *ProjectList) CanBeAccessedBy(userID int64) bool {
	return pl.IsPublic || pl.UserID == userID
}

// CanBeEditedBy checks if the project list can be edited by the given user
func (pl *ProjectList) CanBeEditedBy(userID int64) bool {
	return pl.UserID == userID
}

// CanBeDeletedBy checks if the project list can be deleted by the given user
func (pl *ProjectList) CanBeDeletedBy(userID int64) bool {
	return pl.UserID == userID
}

// Validate validates the project list item data
func (pli *ProjectListItem) Validate() error {
	if pli.ProjectListID <= 0 {
		return errors.New("valid project_list_id is required")
	}
	if pli.ProjectID <= 0 {
		return errors.New("valid project_id is required")
	}
	return nil
}

// TableName specifies the table name for GORM
func (ProjectList) TableName() string {
	return "project_lists"
}

// TableName specifies the table name for GORM
func (ProjectListItem) TableName() string {
	return "project_list_items"
}

// Index defines composite unique index for project list item
func (ProjectListItem) Indexes() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":    "idx_project_list_project",
			"columns": []string{"project_list_id", "project_id"},
			"unique":  true,
		},
	}
}
