package services

import (
	"context"
	"fmt"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
)

// commentService implements the CommentService interface
type commentService struct {
	commentRepo repositories.CommentRepository
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
}

// NewCommentService creates a new comment service
func NewCommentService(commentRepo repositories.CommentRepository, userRepo repositories.UserRepository, projectRepo repositories.ProjectRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// CreateComment creates a new comment
func (s *commentService) CreateComment(ctx context.Context, req domain.CommentCreateRequest) (*domain.Comment, error) {
	// Validate that the user exists
	if _, err := s.userRepo.GetByID(ctx, req.UserID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate that the project exists
	if _, err := s.projectRepo.GetByID(ctx, req.ProjectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	comment := &domain.Comment{
		ProjectID:       req.ProjectID,
		Content:         req.Content,
		Rating:          req.Rating,
		UserID:          req.UserID,
		ParentCommentID: req.ParentCommentID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetProjectComments retrieves all comments for a project
func (s *commentService) GetProjectComments(ctx context.Context, projectID int64) ([]domain.CommentResponse, error) {
	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	comments, err := s.commentRepo.GetWithUserInfo(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project comments: %w", err)
	}

	return comments, nil
}

// GetCommentReplies retrieves all replies for a comment
func (s *commentService) GetCommentReplies(ctx context.Context, commentID int64) ([]domain.CommentResponse, error) {
	if commentID <= 0 {
		return nil, fmt.Errorf("invalid comment ID: %d", commentID)
	}

	replies, err := s.commentRepo.GetRepliesByCommentID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	var responses []domain.CommentResponse
	for _, reply := range replies {
		response := domain.CommentResponse{
			ID:              reply.ID,
			CreatedAt:       reply.CreatedAt,
			UpdatedAt:       reply.UpdatedAt,
			ProjectID:       reply.ProjectID,
			Content:         reply.Content,
			Rating:          reply.Rating,
			UserID:          reply.UserID,
			ParentCommentID: reply.ParentCommentID,
			Username:        reply.User.Username,
			UserReputation:  reply.User.Reputation,
			ReplyCount:      0,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateComment updates an existing comment
func (s *commentService) UpdateComment(ctx context.Context, id int64, req domain.CommentUpdateRequest, userID int64) (*domain.Comment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid comment ID: %d", id)
	}

	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}

	if !comment.CanBeEditedBy(userID) {
		return nil, fmt.Errorf("unauthorized: user cannot edit this comment")
	}

	if req.Content != nil {
		comment.Content = *req.Content
	}
	if req.Rating != nil {
		comment.Rating = *req.Rating
	}

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (s *commentService) DeleteComment(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid comment ID: %d", id)
	}

	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	if !comment.CanBeDeletedBy(userID) {
		return fmt.Errorf("unauthorized: user cannot delete this comment")
	}

	if err := s.commentRepo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

// CreateReply creates a reply to an existing comment
func (s *commentService) CreateReply(ctx context.Context, parentID int64, req domain.CommentCreateRequest) (*domain.Comment, error) {
	if parentID <= 0 {
		return nil, fmt.Errorf("invalid parent comment ID: %d", parentID)
	}

	// Validate that the parent comment exists
	if _, err := s.commentRepo.GetByID(ctx, parentID); err != nil {
		return nil, fmt.Errorf("parent comment not found: %w", err)
	}

	// Set the parent comment ID
	req.ParentCommentID = &parentID

	return s.CreateComment(ctx, req)
}
