package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Project represents a woodworking project in the system
type Project struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Owner         int64          `json:"owner" gorm:"not null;index"`
	Title         string         `json:"title" gorm:"not null"`
	AverageRating float32        `json:"average_rating" gorm:"default:0"`
	RatingCount   int            `json:"rating_count" gorm:"default:0"`
	Materials     pq.StringArray `json:"materials" gorm:"type:varchar[]"`
	Tools         pq.StringArray `json:"tools" gorm:"type:varchar[]"`
	Description   string         `json:"description" gorm:"type:text"`
	Style         pq.StringArray `json:"style" gorm:"type:varchar[]"`
	Environment   pq.StringArray `json:"environment" gorm:"type:text[]"`
	Portrait      string         `json:"portrait"`
	Images        pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Tutorial      string         `json:"tutorial" gorm:"type:text"`
	TimeToBuild   int            `json:"time_to_build" gorm:"default:0"` // in minutes
}

// ProjectCreateRequest represents the data needed to create a new project
type ProjectCreateRequest struct {
	Owner       int64          `json:"owner" validate:"required"`
	Title       string         `json:"title" validate:"required,min=3,max=200"`
	Materials   pq.StringArray `json:"materials"`
	Tools       pq.StringArray `json:"tools"`
	Description string         `json:"description" validate:"max=5000"`
	Style       pq.StringArray `json:"style"`
	Environment pq.StringArray `json:"environment"`
	Portrait    string         `json:"portrait"`
	Images      pq.StringArray `json:"images"`
	Tutorial    string         `json:"tutorial" validate:"max=10000"`
	TimeToBuild int            `json:"time_to_build" validate:"min=0"`
}

// ProjectUpdateRequest represents the data that can be updated for a project
type ProjectUpdateRequest struct {
	Title       *string         `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Materials   *pq.StringArray `json:"materials,omitempty"`
	Tools       *pq.StringArray `json:"tools,omitempty"`
	Description *string         `json:"description,omitempty" validate:"omitempty,max=5000"`
	Style       *pq.StringArray `json:"style,omitempty"`
	Environment *pq.StringArray `json:"environment,omitempty"`
	Portrait    *string         `json:"portrait,omitempty"`
	Images      *pq.StringArray `json:"images,omitempty"`
	Tutorial    *string         `json:"tutorial,omitempty" validate:"omitempty,max=10000"`
	TimeToBuild *int            `json:"time_to_build,omitempty" validate:"omitempty,min=0"`
}

// ProjectSearchFilters represents search filters for projects
type ProjectSearchFilters struct {
	Style          string  `json:"style"`
	Environment    string  `json:"environment"`
	MaxTimeToBuild int     `json:"max_time_to_build"`
	Materials      string  `json:"materials"`
	Tools          string  `json:"tools"`
	MinRating      float32 `json:"min_rating"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}

// Validate validates the project data
func (p *Project) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if len(p.Title) < 3 || len(p.Title) > 200 {
		return errors.New("title must be between 3 and 200 characters")
	}
	if p.Owner <= 0 {
		return errors.New("valid owner is required")
	}
	if len(p.Description) > 5000 {
		return errors.New("description cannot exceed 5000 characters")
	}
	if len(p.Tutorial) > 10000 {
		return errors.New("tutorial cannot exceed 10000 characters")
	}
	if p.TimeToBuild < 0 {
		return errors.New("time to build cannot be negative")
	}
	return nil
}

// UpdateAverageRating updates the average rating based on new rating
func (p *Project) UpdateAverageRating(newRating int, isUpdate bool, oldRating int) {
	if isUpdate {
		// Remove old rating and add new one
		if p.RatingCount > 0 {
			total := p.AverageRating * float32(p.RatingCount)
			total = total - float32(oldRating) + float32(newRating)
			p.AverageRating = total / float32(p.RatingCount)
		}
	} else {
		// Add new rating
		total := p.AverageRating * float32(p.RatingCount)
		p.RatingCount++
		total += float32(newRating)
		p.AverageRating = total / float32(p.RatingCount)
	}
}

// TableName specifies the table name for GORM
func (Project) TableName() string {
	return "projects"
}
