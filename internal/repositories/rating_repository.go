package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// ratingRepository implements the RatingRepository interface
type ratingRepository struct {
	db *gorm.DB
}

// NewRatingRepository creates a new rating repository
func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{
		db: db,
	}
}

// Create creates a new rating in the database
func (r *ratingRepository) Create(ctx context.Context, rating *domain.Rating) error {
	if err := rating.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(rating).Error; err != nil {
		return fmt.Errorf("failed to create rating: %w", err)
	}

	return nil
}

// GetByID retrieves a rating by its ID
func (r *ratingRepository) GetByID(ctx context.Context, id int64) (*domain.Rating, error) {
	var rating domain.Rating
	if err := r.db.WithContext(ctx).Preload("User").Preload("Project").First(&rating, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rating with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get rating by id: %w", err)
	}

	return &rating, nil
}

// GetByUserAndProject retrieves a rating by user and project
func (r *ratingRepository) GetByUserAndProject(ctx context.Context, userID, projectID int64) (*domain.Rating, error) {
	var rating domain.Rating
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		First(&rating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rating not found for user %d and project %d", userID, projectID)
		}
		return nil, fmt.Errorf("failed to get rating by user and project: %w", err)
	}

	return &rating, nil
}

// GetByProjectID retrieves all ratings for a specific project
func (r *ratingRepository) GetByProjectID(ctx context.Context, projectID int64) ([]domain.Rating, error) {
	var ratings []domain.Rating
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		return nil, fmt.Errorf("failed to get ratings by project id: %w", err)
	}

	return ratings, nil
}

// Update updates an existing rating
func (r *ratingRepository) Update(ctx context.Context, rating *domain.Rating) error {
	if err := rating.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(rating)
	if result.Error != nil {
		return fmt.Errorf("failed to update rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("rating with id %d not found", rating.ID)
	}

	return nil
}

// Delete deletes a rating by its ID
func (r *ratingRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Rating{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("rating with id %d not found", id)
	}

	return nil
}

// GetProjectStats returns rating statistics for a project
func (r *ratingRepository) GetProjectStats(ctx context.Context, projectID int64) (*domain.RatingStats, error) {
	var stats domain.RatingStats
	stats.ProjectID = projectID
	stats.Distribution = make(map[int]int)

	// Get average rating and total count
	var result struct {
		AvgRating    float32
		TotalRatings int64
	}

	if err := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Select("COALESCE(AVG(value), 0) as avg_rating, COUNT(*) as total_ratings").
		Where("project_id = ?", projectID).
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}

	stats.AverageRating = result.AvgRating
	stats.TotalRatings = int(result.TotalRatings)

	// Get distribution of ratings
	var distributions []struct {
		Value int
		Count int
	}

	if err := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Select("value, COUNT(*) as count").
		Where("project_id = ?", projectID).
		Group("value").
		Scan(&distributions).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	// Initialize all rating values with 0
	for i := 1; i <= 5; i++ {
		stats.Distribution[i] = 0
	}

	// Fill in actual counts
	for _, dist := range distributions {
		stats.Distribution[dist.Value] = dist.Count
	}

	return &stats, nil
}

// GetAverageRating returns the average rating and count for a project
func (r *ratingRepository) GetAverageRating(ctx context.Context, projectID int64) (float32, int, error) {
	var result struct {
		AvgRating float32
		Count     int64
	}

	if err := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Select("COALESCE(AVG(value), 0) as avg_rating, COUNT(*) as count").
		Where("project_id = ?", projectID).
		Scan(&result).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return result.AvgRating, int(result.Count), nil
}

// GetUserRatings retrieves all ratings made by a specific user
func (r *ratingRepository) GetUserRatings(ctx context.Context, userID int64, limit, offset int) ([]domain.Rating, error) {
	var ratings []domain.Rating

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Project").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&ratings).Error; err != nil {
		return nil, fmt.Errorf("failed to get user ratings: %w", err)
	}

	return ratings, nil
}

// GetTopRatedProjects returns projects with highest average ratings
func (r *ratingRepository) GetTopRatedProjects(ctx context.Context, limit int) ([]struct {
	ProjectID     int64
	AverageRating float32
	RatingCount   int
}, error) {
	var results []struct {
		ProjectID     int64
		AverageRating float32
		RatingCount   int
	}

	query := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Select("project_id, AVG(value) as average_rating, COUNT(*) as rating_count").
		Group("project_id").
		Having("COUNT(*) >= 3"). // Only include projects with at least 3 ratings
		Order("average_rating DESC, rating_count DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get top rated projects: %w", err)
	}

	return results, nil
}

// HasUserRatedProject checks if a user has already rated a project
func (r *ratingRepository) HasUserRatedProject(ctx context.Context, userID, projectID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user rated project: %w", err)
	}

	return count > 0, nil
}

// GetRatingTrends returns rating trends over time for a project
func (r *ratingRepository) GetRatingTrends(ctx context.Context, projectID int64, days int) ([]struct {
	Date          string
	AverageRating float32
	Count         int
}, error) {
	var results []struct {
		Date          string
		AverageRating float32
		Count         int
	}

	query := r.db.WithContext(ctx).Model(&domain.Rating{}).
		Select("DATE(created_at) as date, AVG(value) as average_rating, COUNT(*) as count").
		Where("project_id = ? AND created_at >= NOW() - INTERVAL ? DAY", projectID, days).
		Group("DATE(created_at)").
		Order("date ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating trends: %w", err)
	}

	return results, nil
}
