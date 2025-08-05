package services

import (
	"context"
	"fmt"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
)

// ratingService implements the RatingService interface
type ratingService struct {
	ratingRepo  repositories.RatingRepository
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

// NewRatingService creates a new rating service
func NewRatingService(ratingRepo repositories.RatingRepository, projectRepo repositories.ProjectRepository, userRepo repositories.UserRepository) RatingService {
	return &ratingService{
		ratingRepo:  ratingRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateRating creates a new rating
func (s *ratingService) CreateRating(ctx context.Context, req domain.RatingCreateRequest) (*domain.Rating, error) {
	// Validate that the user exists
	if _, err := s.userRepo.GetByID(ctx, req.UserID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate that the project exists
	if _, err := s.projectRepo.GetByID(ctx, req.ProjectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Check if user has already rated this project
	if _, err := s.ratingRepo.GetByUserAndProject(ctx, req.UserID, req.ProjectID); err == nil {
		return nil, fmt.Errorf("user has already rated this project")
	}

	rating := &domain.Rating{
		Value:     req.Value,
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	}

	if err := s.ratingRepo.Create(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to create rating: %w", err)
	}

	// Update project's average rating
	if err := s.updateProjectRating(ctx, req.ProjectID); err != nil {
		// Log the error but don't fail the rating creation
		fmt.Printf("Warning: failed to update project rating: %v\n", err)
	}

	return rating, nil
}

// UpdateRating updates an existing rating
func (s *ratingService) UpdateRating(ctx context.Context, userID, projectID int64, req domain.RatingUpdateRequest) (*domain.Rating, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	rating, err := s.ratingRepo.GetByUserAndProject(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	if !rating.CanBeUpdatedBy(userID) {
		return nil, fmt.Errorf("unauthorized: user cannot update this rating")
	}

	if req.Value != nil {
		rating.Value = *req.Value
	}

	if err := s.ratingRepo.Update(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	// Update project's average rating
	if err := s.updateProjectRating(ctx, projectID); err != nil {
		// Log the error but don't fail the rating update
		fmt.Printf("Warning: failed to update project rating: %v\n", err)
	}

	return rating, nil
}

// GetProjectRatings retrieves all ratings for a project
func (s *ratingService) GetProjectRatings(ctx context.Context, projectID int64) ([]domain.RatingResponse, error) {
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	ratings, err := s.ratingRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project ratings: %w", err)
	}

	var responses []domain.RatingResponse
	for _, rating := range ratings {
		response := domain.RatingResponse{
			ID:        rating.ID,
			CreatedAt: rating.CreatedAt,
			UpdatedAt: rating.UpdatedAt,
			Value:     rating.Value,
			UserID:    rating.UserID,
			ProjectID: rating.ProjectID,
			Username:  rating.User.Username,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetUserRating retrieves a user's rating for a specific project
func (s *ratingService) GetUserRating(ctx context.Context, userID, projectID int64) (*domain.Rating, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	rating, err := s.ratingRepo.GetByUserAndProject(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	return rating, nil
}

// DeleteRating deletes a rating
func (s *ratingService) DeleteRating(ctx context.Context, userID, projectID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}
	if projectID <= 0 {
		return fmt.Errorf("invalid project ID: %d", projectID)
	}

	rating, err := s.ratingRepo.GetByUserAndProject(ctx, userID, projectID)
	if err != nil {
		return fmt.Errorf("rating not found: %w", err)
	}

	if !rating.CanBeDeletedBy(userID) {
		return fmt.Errorf("unauthorized: user cannot delete this rating")
	}

	if err := s.ratingRepo.Delete(ctx, rating.ID); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	// Update project's average rating
	if err := s.updateProjectRating(ctx, projectID); err != nil {
		// Log the error but don't fail the rating deletion
		fmt.Printf("Warning: failed to update project rating: %v\n", err)
	}

	return nil
}

// GetProjectRatingStats retrieves rating statistics for a project
func (s *ratingService) GetProjectRatingStats(ctx context.Context, projectID int64) (*domain.RatingStats, error) {
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	stats, err := s.ratingRepo.GetProjectStats(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project rating stats: %w", err)
	}

	return stats, nil
}

// updateProjectRating updates the project's average rating and rating count
func (s *ratingService) updateProjectRating(ctx context.Context, projectID int64) error {
	avgRating, count, err := s.ratingRepo.GetAverageRating(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to calculate average rating: %w", err)
	}

	if err := s.projectRepo.UpdateRating(ctx, projectID, avgRating, count); err != nil {
		return fmt.Errorf("failed to update project rating: %w", err)
	}

	return nil
}
