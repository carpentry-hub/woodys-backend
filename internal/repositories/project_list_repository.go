package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// projectListRepository implements the ProjectListRepository interface
type projectListRepository struct {
	db *gorm.DB
}

// NewProjectListRepository creates a new project list repository
func NewProjectListRepository(db *gorm.DB) ProjectListRepository {
	return &projectListRepository{
		db: db,
	}
}

// Create creates a new project list in the database
func (r *projectListRepository) Create(ctx context.Context, list *domain.ProjectList) error {
	if err := list.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(list).Error; err != nil {
		return fmt.Errorf("failed to create project list: %w", err)
	}

	return nil
}

// GetByID retrieves a project list by its ID
func (r *projectListRepository) GetByID(ctx context.Context, id int64) (*domain.ProjectList, error) {
	var list domain.ProjectList
	if err := r.db.WithContext(ctx).Preload("User").First(&list, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project list with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get project list by id: %w", err)
	}

	return &list, nil
}

// GetByUserID retrieves all project lists owned by a specific user
func (r *projectListRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.ProjectList, error) {
	var lists []domain.ProjectList
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&lists).Error; err != nil {
		return nil, fmt.Errorf("failed to get project lists by user id: %w", err)
	}

	return lists, nil
}

// Update updates an existing project list
func (r *projectListRepository) Update(ctx context.Context, list *domain.ProjectList) error {
	if err := list.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(list)
	if result.Error != nil {
		return fmt.Errorf("failed to update project list: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project list with id %d not found", list.ID)
	}

	return nil
}

// Delete deletes a project list by its ID
func (r *projectListRepository) Delete(ctx context.Context, id int64) error {
	// Start a transaction to delete both the list and its items
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all project list items first
		if err := tx.Where("project_list_id = ?", id).Delete(&domain.ProjectListItem{}).Error; err != nil {
			return fmt.Errorf("failed to delete project list items: %w", err)
		}

		// Delete the project list
		result := tx.Delete(&domain.ProjectList{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete project list: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("project list with id %d not found", id)
		}

		return nil
	})
}

// AddProject adds a project to a project list
func (r *projectListRepository) AddProject(ctx context.Context, item *domain.ProjectListItem) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if project is already in the list
	exists, err := r.IsProjectInList(ctx, item.ProjectListID, item.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to check if project exists in list: %w", err)
	}
	if exists {
		return fmt.Errorf("project %d is already in list %d", item.ProjectID, item.ProjectListID)
	}

	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("failed to add project to list: %w", err)
	}

	return nil
}

// RemoveProject removes a project from a project list
func (r *projectListRepository) RemoveProject(ctx context.Context, listID, projectID int64) error {
	result := r.db.WithContext(ctx).
		Where("project_list_id = ? AND project_id = ?", listID, projectID).
		Delete(&domain.ProjectListItem{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove project from list: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project %d not found in list %d", projectID, listID)
	}

	return nil
}

// GetWithProjects retrieves a project list with its projects
func (r *projectListRepository) GetWithProjects(ctx context.Context, listID int64) (*domain.ProjectListResponse, error) {
	var list domain.ProjectList
	if err := r.db.WithContext(ctx).Preload("User").First(&list, listID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project list with id %d not found", listID)
		}
		return nil, fmt.Errorf("failed to get project list: %w", err)
	}

	// Get projects in the list
	var items []struct {
		ProjectID     int64   `json:"project_id"`
		Title         string  `json:"title"`
		Portrait      string  `json:"portrait"`
		AverageRating float32 `json:"average_rating"`
		RatingCount   int     `json:"rating_count"`
		AddedAt       string  `json:"added_at"`
	}

	if err := r.db.WithContext(ctx).
		Table("project_list_items pli").
		Select("pli.project_id, p.title, p.portrait, p.average_rating, p.rating_count, pli.created_at as added_at").
		Joins("JOIN projects p ON p.id = pli.project_id").
		Where("pli.project_list_id = ?", listID).
		Order("pli.created_at DESC").
		Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get projects in list: %w", err)
	}

	// Convert to response format
	response := &domain.ProjectListResponse{
		ID:           list.ID,
		CreatedAt:    list.CreatedAt,
		UpdatedAt:    list.UpdatedAt,
		UserID:       list.UserID,
		Name:         list.Name,
		IsPublic:     list.IsPublic,
		ProjectCount: len(items),
	}

	// Convert items to response format
	for _, item := range items {
		project := domain.ProjectListItemProject{
			ID:            item.ProjectID,
			Title:         item.Title,
			Portrait:      item.Portrait,
			AverageRating: item.AverageRating,
			RatingCount:   item.RatingCount,
		}
		response.Projects = append(response.Projects, project)
	}

	return response, nil
}

// GetUserListsWithProjects retrieves all project lists for a user with their projects
func (r *projectListRepository) GetUserListsWithProjects(ctx context.Context, userID int64) ([]domain.ProjectListResponse, error) {
	var lists []domain.ProjectList
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&lists).Error; err != nil {
		return nil, fmt.Errorf("failed to get user project lists: %w", err)
	}

	var responses []domain.ProjectListResponse
	for _, list := range lists {
		// Get project count for each list
		var projectCount int64
		if err := r.db.WithContext(ctx).Model(&domain.ProjectListItem{}).
			Where("project_list_id = ?", list.ID).
			Count(&projectCount).Error; err != nil {
			return nil, fmt.Errorf("failed to get project count for list %d: %w", list.ID, err)
		}

		response := domain.ProjectListResponse{
			ID:           list.ID,
			CreatedAt:    list.CreatedAt,
			UpdatedAt:    list.UpdatedAt,
			UserID:       list.UserID,
			Name:         list.Name,
			IsPublic:     list.IsPublic,
			ProjectCount: int(projectCount),
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// IsProjectInList checks if a project is already in a project list
func (r *projectListRepository) IsProjectInList(ctx context.Context, listID, projectID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.ProjectListItem{}).
		Where("project_list_id = ? AND project_id = ?", listID, projectID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if project is in list: %w", err)
	}

	return count > 0, nil
}

// GetPublicLists retrieves all public project lists
func (r *projectListRepository) GetPublicLists(ctx context.Context, limit, offset int) ([]domain.ProjectListResponse, error) {
	var lists []domain.ProjectList

	query := r.db.WithContext(ctx).Where("is_public = true").Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&lists).Error; err != nil {
		return nil, fmt.Errorf("failed to get public project lists: %w", err)
	}

	var responses []domain.ProjectListResponse
	for _, list := range lists {
		// Get project count for each list
		var projectCount int64
		if err := r.db.WithContext(ctx).Model(&domain.ProjectListItem{}).
			Where("project_list_id = ?", list.ID).
			Count(&projectCount).Error; err != nil {
			return nil, fmt.Errorf("failed to get project count for list %d: %w", list.ID, err)
		}

		response := domain.ProjectListResponse{
			ID:           list.ID,
			CreatedAt:    list.CreatedAt,
			UpdatedAt:    list.UpdatedAt,
			UserID:       list.UserID,
			Name:         list.Name,
			IsPublic:     list.IsPublic,
			ProjectCount: int(projectCount),
		}

		responses = append(responses, response)
	}

	return responses, nil
}
