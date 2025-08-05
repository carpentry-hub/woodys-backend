package services

import (
	"context"
	"fmt"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
)

// projectService implements the ProjectService interface
type projectService struct {
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repositories.ProjectRepository, userRepo repositories.UserRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateProject creates a new project
func (s *projectService) CreateProject(ctx context.Context, req domain.ProjectCreateRequest) (*domain.Project, error) {
	// Validate that the owner exists
	if _, err := s.userRepo.GetByID(ctx, req.Owner); err != nil {
		return nil, fmt.Errorf("owner not found: %w", err)
	}

	project := &domain.Project{
		Owner:         req.Owner,
		Title:         req.Title,
		Materials:     req.Materials,
		Tools:         req.Tools,
		Description:   req.Description,
		Style:         req.Style,
		Environment:   req.Environment,
		Portrait:      req.Portrait,
		Images:        req.Images,
		Tutorial:      req.Tutorial,
		TimeToBuild:   req.TimeToBuild,
		AverageRating: 0,
		RatingCount:   0,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetProjectByID retrieves a project by ID
func (s *projectService) GetProjectByID(ctx context.Context, id int64) (*domain.Project, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", id)
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// UpdateProject updates an existing project
func (s *projectService) UpdateProject(ctx context.Context, id int64, req domain.ProjectUpdateRequest) (*domain.Project, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", id)
	}

	// Get existing project
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Update fields if provided
	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Materials != nil {
		project.Materials = *req.Materials
	}
	if req.Tools != nil {
		project.Tools = *req.Tools
	}
	if req.Style != nil {
		project.Style = *req.Style
	}
	if req.Environment != nil {
		project.Environment = *req.Environment
	}
	if req.Portrait != nil {
		project.Portrait = *req.Portrait
	}
	if req.Images != nil {
		project.Images = *req.Images
	}
	if req.Tutorial != nil {
		project.Tutorial = *req.Tutorial
	}
	if req.TimeToBuild != nil {
		project.TimeToBuild = *req.TimeToBuild
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *projectService) DeleteProject(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid project ID: %d", id)
	}

	// Check if project exists and user is the owner
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if project.Owner != userID {
		return fmt.Errorf("unauthorized: user %d is not the owner of project %d", userID, id)
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// SearchProjects searches for projects based on filters
func (s *projectService) SearchProjects(ctx context.Context, filters domain.ProjectSearchFilters) ([]domain.Project, error) {
	projects, err := s.projectRepo.Search(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}

	return projects, nil
}

// ListProjects retrieves a list of projects with pagination
func (s *projectService) ListProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	projects, err := s.projectRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}

// GetProjectsByOwner retrieves all projects owned by a specific user
func (s *projectService) GetProjectsByOwner(ctx context.Context, ownerID int64) ([]domain.Project, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("invalid owner ID: %d", ownerID)
	}

	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, ownerID); err != nil {
		return nil, fmt.Errorf("owner not found: %w", err)
	}

	projects, err := s.projectRepo.GetByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects by owner: %w", err)
	}

	return projects, nil
}

// GetPopularProjects retrieves popular projects
func (s *projectService) GetPopularProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	// This is a placeholder implementation
	// In a real scenario, you might implement this in the repository
	projects, err := s.projectRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular projects: %w", err)
	}

	return projects, nil
}

// GetRecentProjects retrieves recently created projects
func (s *projectService) GetRecentProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	projects, err := s.projectRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent projects: %w", err)
	}

	return projects, nil
}
