package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
)

// userService implements the UserService interface
type userService struct {
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, projectRepo repositories.ProjectRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(ctx context.Context, req domain.UserCreateRequest) (*domain.User, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user with email already exists
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Check if user with Firebase UID already exists
	if _, err := s.userRepo.GetByFirebaseUID(ctx, req.FirebaseUID); err == nil {
		return nil, fmt.Errorf("user with firebase_uid %s already exists", req.FirebaseUID)
	}

	// Create user domain object
	user := &domain.User{
		Username:       strings.TrimSpace(req.Username),
		Email:          strings.ToLower(strings.TrimSpace(req.Email)),
		FirebaseUID:    strings.TrimSpace(req.FirebaseUID),
		Reputation:     0.0,
		ProfilePicture: 0,
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByFirebaseUID retrieves a user by Firebase UID
func (s *userService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	if strings.TrimSpace(firebaseUID) == "" {
		return nil, fmt.Errorf("firebase_uid cannot be empty")
	}

	user, err := s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by firebase_uid: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, id int64, req domain.UserUpdateRequest) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if err := s.validateUsername(username); err != nil {
			return nil, fmt.Errorf("invalid username: %w", err)
		}
		user.Username = username
	}

	if req.Reputation != nil {
		if *req.Reputation < 0 {
			return nil, fmt.Errorf("reputation cannot be negative")
		}
		user.Reputation = *req.Reputation
	}

	if req.ProfilePicture != nil {
		if *req.ProfilePicture < 0 {
			return nil, fmt.Errorf("profile_picture cannot be negative")
		}
		user.ProfilePicture = *req.ProfilePicture
	}

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// TODO: In a real application, you might want to:
	// 1. Soft delete instead of hard delete
	// 2. Handle related data (projects, comments, ratings)
	// 3. Send notifications
	// 4. Archive user data

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with pagination
func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	// Set default limit if not provided
	if limit == 0 {
		limit = 50
	}

	// Enforce maximum limit
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// GetUserProjects retrieves all projects owned by a user
func (s *userService) GetUserProjects(ctx context.Context, userID int64) ([]domain.Project, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	projects, err := s.projectRepo.GetByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user projects: %w", err)
	}

	return projects, nil
}

// validateCreateRequest validates the user creation request
func (s *userService) validateCreateRequest(req domain.UserCreateRequest) error {
	if err := s.validateUsername(req.Username); err != nil {
		return err
	}

	if err := s.validateEmail(req.Email); err != nil {
		return err
	}

	if strings.TrimSpace(req.FirebaseUID) == "" {
		return fmt.Errorf("firebase_uid is required")
	}

	return nil
}

// validateUsername validates the username
func (s *userService) validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is required")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return fmt.Errorf("username cannot exceed 50 characters")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}

	return nil
}

// validateEmail validates the email format (basic validation)
func (s *userService) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}

	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// CalculateUserReputation calculates user reputation based on their activity
func (s *userService) CalculateUserReputation(ctx context.Context, userID int64) (float32, error) {
	// This is a placeholder for reputation calculation logic
	// In a real application, you might calculate based on:
	// - Number of projects created
	// - Average rating of user's projects
	// - Number of helpful comments
	// - Community engagement metrics

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	projects, err := s.projectRepo.GetByOwner(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user projects: %w", err)
	}

	// Simple reputation calculation based on project ratings
	var totalRating float32
	var projectCount int

	for _, project := range projects {
		if project.RatingCount > 0 {
			totalRating += project.AverageRating
			projectCount++
		}
	}

	reputation := user.Reputation
	if projectCount > 0 {
		avgProjectRating := totalRating / float32(projectCount)
		reputation = avgProjectRating * float32(projectCount) * 0.1 // Simple formula
	}

	return reputation, nil
}
