package repositories

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// projectRepository implements the ProjectRepository interface
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

// Create creates a new project in the database
func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	if err := project.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetByID retrieves a project by its ID
func (r *projectRepository) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.WithContext(ctx).First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	}

	return &project, nil
}

// GetByOwner retrieves all projects owned by a specific user
func (r *projectRepository) GetByOwner(ctx context.Context, ownerID int64) ([]domain.Project, error) {
	var projects []domain.Project
	if err := r.db.WithContext(ctx).Where("owner = ?", ownerID).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to get projects by owner: %w", err)
	}

	return projects, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	if err := project.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(project)
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", project.ID)
	}

	return nil
}

// Delete deletes a project by its ID
func (r *projectRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Project{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", id)
	}

	return nil
}

// Search searches for projects based on filters
func (r *projectRepository) Search(ctx context.Context, filters domain.ProjectSearchFilters) ([]domain.Project, error) {
	var projects []domain.Project
	query := r.db.WithContext(ctx).Model(&domain.Project{})

	// Apply filters
	if filters.Style != "" {
		query = query.Where("? = ANY(style)", filters.Style)
	}

	if filters.Environment != "" {
		query = query.Where("? = ANY(environment)", filters.Environment)
	}

	if filters.Materials != "" {
		query = query.Where("? = ANY(materials)", filters.Materials)
	}

	if filters.Tools != "" {
		query = query.Where("? = ANY(tools)", filters.Tools)
	}

	if filters.MaxTimeToBuild > 0 {
		query = query.Where("time_to_build <= ?", filters.MaxTimeToBuild)
	}

	if filters.MinRating > 0 {
		query = query.Where("average_rating >= ?", filters.MinRating)
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	} else {
		query = query.Limit(50) // Default limit
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Order by rating and creation date
	query = query.Order("average_rating DESC, created_at DESC")

	if err := query.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}

	return projects, nil
}

// List retrieves a list of projects with pagination
func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	var projects []domain.Project

	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(50) // Default limit
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}

// UpdateRating updates the average rating and rating count for a project
func (r *projectRepository) UpdateRating(ctx context.Context, projectID int64, averageRating float32, ratingCount int) error {
	result := r.db.WithContext(ctx).Model(&domain.Project{}).
		Where("id = ?", projectID).
		Updates(map[string]interface{}{
			"average_rating": averageRating,
			"rating_count":   ratingCount,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update project rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", projectID)
	}

	return nil
}

// SearchByTitle searches for projects by title (case-insensitive)
func (r *projectRepository) SearchByTitle(ctx context.Context, title string, limit, offset int) ([]domain.Project, error) {
	var projects []domain.Project

	query := r.db.WithContext(ctx).
		Where("LOWER(title) LIKE LOWER(?)", "%"+strings.ToLower(title)+"%").
		Order("average_rating DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(50)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to search projects by title: %w", err)
	}

	return projects, nil
}

// GetPopularProjects retrieves projects ordered by rating and rating count
func (r *projectRepository) GetPopularProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	var projects []domain.Project

	query := r.db.WithContext(ctx).
		Where("rating_count > 0").
		Order("average_rating DESC, rating_count DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(20)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular projects: %w", err)
	}

	return projects, nil
}

// GetRecentProjects retrieves recently created projects
func (r *projectRepository) GetRecentProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	var projects []domain.Project

	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(20)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent projects: %w", err)
	}

	return projects, nil
}
