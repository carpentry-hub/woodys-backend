package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// commentRepository implements the CommentRepository interface
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

// Create creates a new comment in the database
func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetByID retrieves a comment by its ID
func (r *commentRepository) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.WithContext(ctx).Preload("User").Preload("Project").First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	return &comment, nil
}

// GetByProjectID retrieves all comments for a specific project
func (r *commentRepository) GetByProjectID(ctx context.Context, projectID int64) ([]domain.Comment, error) {
	var comments []domain.Comment
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND parent_comment_id IS NULL AND is_deleted = false", projectID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments by project id: %w", err)
	}

	return comments, nil
}

// GetRepliesByCommentID retrieves all replies for a specific comment
func (r *commentRepository) GetRepliesByCommentID(ctx context.Context, commentID int64) ([]domain.Comment, error) {
	var replies []domain.Comment
	if err := r.db.WithContext(ctx).
		Where("parent_comment_id = ? AND is_deleted = false", commentID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	return replies, nil
}

// Update updates an existing comment
func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(comment)
	if result.Error != nil {
		return fmt.Errorf("failed to update comment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comment with id %d not found", comment.ID)
	}

	return nil
}

// Delete permanently deletes a comment by its ID
func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Comment{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete comment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comment with id %d not found", id)
	}

	return nil
}

// SoftDelete marks a comment as deleted without removing it from database
func (r *commentRepository) SoftDelete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&domain.Comment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"content":    "[This comment has been deleted]",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to soft delete comment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comment with id %d not found", id)
	}

	return nil
}

// GetWithUserInfo retrieves comments with user information for a project
func (r *commentRepository) GetWithUserInfo(ctx context.Context, projectID int64) ([]domain.CommentResponse, error) {
	var comments []domain.Comment

	// Get all non-deleted comments for the project with user info
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND parent_comment_id IS NULL AND is_deleted = false", projectID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments with user info: %w", err)
	}

	// Convert to response format
	var responses []domain.CommentResponse
	for _, comment := range comments {
		// Get reply count
		var replyCount int64
		r.db.WithContext(ctx).Model(&domain.Comment{}).
			Where("parent_comment_id = ? AND is_deleted = false", comment.ID).
			Count(&replyCount)

		response := domain.CommentResponse{
			ID:              comment.ID,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
			ProjectID:       comment.ProjectID,
			Content:         comment.Content,
			Rating:          comment.Rating,
			UserID:          comment.UserID,
			ParentCommentID: comment.ParentCommentID,
			Username:        comment.User.Username,
			UserReputation:  comment.User.Reputation,
			ReplyCount:      int(replyCount),
		}

		// Get replies if any
		if replyCount > 0 {
			replies, err := r.getRepliesWithUserInfo(ctx, comment.ID)
			if err == nil {
				response.Replies = replies
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// getRepliesWithUserInfo is a helper function to get replies with user info
func (r *commentRepository) getRepliesWithUserInfo(ctx context.Context, parentID int64) ([]domain.CommentResponse, error) {
	var replies []domain.Comment

	if err := r.db.WithContext(ctx).
		Where("parent_comment_id = ? AND is_deleted = false", parentID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("failed to get replies with user info: %w", err)
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
			ReplyCount:      0, // Replies don't have nested replies
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetCommentsByUser retrieves all comments made by a specific user
func (r *commentRepository) GetCommentsByUser(ctx context.Context, userID int64, limit, offset int) ([]domain.Comment, error) {
	var comments []domain.Comment

	query := r.db.WithContext(ctx).
		Where("user_id = ? AND is_deleted = false", userID).
		Preload("Project").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments by user: %w", err)
	}

	return comments, nil
}

// GetCommentCount returns the total number of comments for a project
func (r *commentRepository) GetCommentCount(ctx context.Context, projectID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Comment{}).
		Where("project_id = ? AND is_deleted = false", projectID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get comment count: %w", err)
	}

	return count, nil
}

// GetReplyCount returns the number of replies for a comment
func (r *commentRepository) GetReplyCount(ctx context.Context, commentID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Comment{}).
		Where("parent_comment_id = ? AND is_deleted = false", commentID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get reply count: %w", err)
	}

	return count, nil
}
