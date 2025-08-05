package repositories

import (
	"context"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
}

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error
	GetByID(ctx context.Context, id int64) (*domain.Project, error)
	GetByOwner(ctx context.Context, ownerID int64) ([]domain.Project, error)
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, filters domain.ProjectSearchFilters) ([]domain.Project, error)
	List(ctx context.Context, limit, offset int) ([]domain.Project, error)
	UpdateRating(ctx context.Context, projectID int64, averageRating float32, ratingCount int) error
}

// CommentRepository defines the interface for comment data access
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByID(ctx context.Context, id int64) (*domain.Comment, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]domain.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentID int64) ([]domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	GetWithUserInfo(ctx context.Context, projectID int64) ([]domain.CommentResponse, error)
}

// RatingRepository defines the interface for rating data access
type RatingRepository interface {
	Create(ctx context.Context, rating *domain.Rating) error
	GetByID(ctx context.Context, id int64) (*domain.Rating, error)
	GetByUserAndProject(ctx context.Context, userID, projectID int64) (*domain.Rating, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]domain.Rating, error)
	Update(ctx context.Context, rating *domain.Rating) error
	Delete(ctx context.Context, id int64) error
	GetProjectStats(ctx context.Context, projectID int64) (*domain.RatingStats, error)
	GetAverageRating(ctx context.Context, projectID int64) (float32, int, error)
}

// ProjectListRepository defines the interface for project list data access
type ProjectListRepository interface {
	Create(ctx context.Context, list *domain.ProjectList) error
	GetByID(ctx context.Context, id int64) (*domain.ProjectList, error)
	GetByUserID(ctx context.Context, userID int64) ([]domain.ProjectList, error)
	Update(ctx context.Context, list *domain.ProjectList) error
	Delete(ctx context.Context, id int64) error
	AddProject(ctx context.Context, item *domain.ProjectListItem) error
	RemoveProject(ctx context.Context, listID, projectID int64) error
	GetWithProjects(ctx context.Context, listID int64) (*domain.ProjectListResponse, error)
	GetUserListsWithProjects(ctx context.Context, userID int64) ([]domain.ProjectListResponse, error)
	IsProjectInList(ctx context.Context, listID, projectID int64) (bool, error)
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	User        UserRepository
	Project     ProjectRepository
	Comment     CommentRepository
	Rating      RatingRepository
	ProjectList ProjectListRepository
}
