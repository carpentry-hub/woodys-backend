package services

import (
	"context"
	"fmt"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
)

// projectListService implements the ProjectListService interface
type projectListService struct {
	projectListRepo repositories.ProjectListRepository
	projectRepo     repositories.ProjectRepository
	userRepo        repositories.UserRepository
}

// NewProjectListService creates a new project list service
func NewProjectListService(projectListRepo repositories.ProjectListRepository, projectRepo repositories.ProjectRepository, userRepo repositories.UserRepository) ProjectListService {
	return &projectListService{
		projectListRepo: projectListRepo,
		projectRepo:     projectRepo,
		userRepo:        userRepo,
	}
}

// CreateProjectList creates a new project list
func (s *projectListService) CreateProjectList(ctx context.Context, req domain.ProjectListCreateRequest) (*domain.ProjectList, error) {
	// Validate that the user exists
	if _, err := s.userRepo.GetByID(ctx, req.UserID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	projectList := &domain.ProjectList{
		UserID:   req.UserID,
		Name:     req.Name,
		IsPublic: req.IsPublic,
	}

	if err := s.projectListRepo.Create(ctx, projectList); err != nil {
		return nil, fmt.Errorf("failed to create project list: %w", err)
	}

	return projectList, nil
}

// GetProjectListByID retrieves a project list by ID
func (s *projectListService) GetProjectListByID(ctx context.Context, id int64, userID int64) (*domain.ProjectListResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project list ID: %d", id)
	}

	projectList, err := s.projectListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project list not found: %w", err)
	}

	// Check if user can access this list (owner or public)
	if !projectList.CanBeAccessedBy(userID) {
		return nil, fmt.Errorf("unauthorized: user cannot access this project list")
	}

	response, err := s.projectListRepo.GetWithProjects(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project list with projects: %w", err)
	}

	return response, nil
}

// GetUserProjectLists retrieves all project lists for a user
func (s *projectListService) GetUserProjectLists(ctx context.Context, userID int64) ([]domain.ProjectListResponse, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Validate that the user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	lists, err := s.projectListRepo.GetUserListsWithProjects(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user project lists: %w", err)
	}

	return lists, nil
}

// UpdateProjectList updates an existing project list
func (s *projectListService) UpdateProjectList(ctx context.Context, id int64, req domain.ProjectListUpdateRequest, userID int64) (*domain.ProjectList, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project list ID: %d", id)
	}

	projectList, err := s.projectListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project list not found: %w", err)
	}

	if !projectList.CanBeEditedBy(userID) {
		return nil, fmt.Errorf("unauthorized: user cannot edit this project list")
	}

	// Update fields if provided
	if req.Name != nil {
		projectList.Name = *req.Name
	}
	if req.IsPublic != nil {
		projectList.IsPublic = *req.IsPublic
	}

	if err := s.projectListRepo.Update(ctx, projectList); err != nil {
		return nil, fmt.Errorf("failed to update project list: %w", err)
	}

	return projectList, nil
}

// DeleteProjectList deletes a project list
func (s *projectListService) DeleteProjectList(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid project list ID: %d", id)
	}

	projectList, err := s.projectListRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("project list not found: %w", err)
	}

	if !projectList.CanBeDeletedBy(userID) {
		return fmt.Errorf("unauthorized: user cannot delete this project list")
	}

	if err := s.projectListRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project list: %w", err)
	}

	return nil
}

// AddProjectToList adds a project to a project list
func (s *projectListService) AddProjectToList(ctx context.Context, listID, projectID, userID int64) error {
	if listID <= 0 {
		return fmt.Errorf("invalid project list ID: %d", listID)
	}
	if projectID <= 0 {
		return fmt.Errorf("invalid project ID: %d", projectID)
	}

	// Validate that the project list exists and user can edit it
	projectList, err := s.projectListRepo.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("project list not found: %w", err)
	}

	if !projectList.CanBeEditedBy(userID) {
		return fmt.Errorf("unauthorized: user cannot edit this project list")
	}

	// Validate that the project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Check if project is already in the list
	exists, err := s.projectListRepo.IsProjectInList(ctx, listID, projectID)
	if err != nil {
		return fmt.Errorf("failed to check if project is in list: %w", err)
	}
	if exists {
		return fmt.Errorf("project is already in the list")
	}

	item := &domain.ProjectListItem{
		ProjectListID: listID,
		ProjectID:     projectID,
	}

	if err := s.projectListRepo.AddProject(ctx, item); err != nil {
		return fmt.Errorf("failed to add project to list: %w", err)
	}

	return nil
}

// RemoveProjectFromList removes a project from a project list
func (s *projectListService) RemoveProjectFromList(ctx context.Context, listID, projectID, userID int64) error {
	if listID <= 0 {
		return fmt.Errorf("invalid project list ID: %d", listID)
	}
	if projectID <= 0 {
		return fmt.Errorf("invalid project ID: %d", projectID)
	}

	// Validate that the project list exists and user can edit it
	projectList, err := s.projectListRepo.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("project list not found: %w", err)
	}

	if !projectList.CanBeEditedBy(userID) {
		return fmt.Errorf("unauthorized: user cannot edit this project list")
	}

	if err := s.projectListRepo.RemoveProject(ctx, listID, projectID); err != nil {
		return fmt.Errorf("failed to remove project from list: %w", err)
	}

	return nil
}

// IsProjectInList checks if a project is in a project list
func (s *projectListService) IsProjectInList(ctx context.Context, listID, projectID int64) (bool, error) {
	if listID <= 0 {
		return false, fmt.Errorf("invalid project list ID: %d", listID)
	}
	if projectID <= 0 {
		return false, fmt.Errorf("invalid project ID: %d", projectID)
	}

	exists, err := s.projectListRepo.IsProjectInList(ctx, listID, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to check if project is in list: %w", err)
	}

	return exists, nil
}
