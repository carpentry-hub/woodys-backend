package services

import (
	"context"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req domain.UserCreateRequest) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, req domain.UserUpdateRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetUserProjects(ctx context.Context, userID int64) ([]domain.Project, error)
}

// ProjectService defines the interface for project business logic
type ProjectService interface {
	CreateProject(ctx context.Context, req domain.ProjectCreateRequest) (*domain.Project, error)
	GetProjectByID(ctx context.Context, id int64) (*domain.Project, error)
	UpdateProject(ctx context.Context, id int64, req domain.ProjectUpdateRequest) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int64, userID int64) error
	SearchProjects(ctx context.Context, filters domain.ProjectSearchFilters) ([]domain.Project, error)
	ListProjects(ctx context.Context, limit, offset int) ([]domain.Project, error)
	GetProjectsByOwner(ctx context.Context, ownerID int64) ([]domain.Project, error)
	GetPopularProjects(ctx context.Context, limit, offset int) ([]domain.Project, error)
	GetRecentProjects(ctx context.Context, limit, offset int) ([]domain.Project, error)
}

// CommentService defines the interface for comment business logic
type CommentService interface {
	CreateComment(ctx context.Context, req domain.CommentCreateRequest) (*domain.Comment, error)
	GetProjectComments(ctx context.Context, projectID int64) ([]domain.CommentResponse, error)
	GetCommentReplies(ctx context.Context, commentID int64) ([]domain.CommentResponse, error)
	UpdateComment(ctx context.Context, id int64, req domain.CommentUpdateRequest, userID int64) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int64, userID int64) error
	CreateReply(ctx context.Context, parentID int64, req domain.CommentCreateRequest) (*domain.Comment, error)
}

// RatingService defines the interface for rating business logic
type RatingService interface {
	CreateRating(ctx context.Context, req domain.RatingCreateRequest) (*domain.Rating, error)
	UpdateRating(ctx context.Context, userID, projectID int64, req domain.RatingUpdateRequest) (*domain.Rating, error)
	GetProjectRatings(ctx context.Context, projectID int64) ([]domain.RatingResponse, error)
	GetUserRating(ctx context.Context, userID, projectID int64) (*domain.Rating, error)
	DeleteRating(ctx context.Context, userID, projectID int64) error
	GetProjectRatingStats(ctx context.Context, projectID int64) (*domain.RatingStats, error)
}

// ProjectListService defines the interface for project list business logic
type ProjectListService interface {
	CreateProjectList(ctx context.Context, req domain.ProjectListCreateRequest) (*domain.ProjectList, error)
	GetProjectListByID(ctx context.Context, id int64, userID int64) (*domain.ProjectListResponse, error)
	GetUserProjectLists(ctx context.Context, userID int64) ([]domain.ProjectListResponse, error)
	UpdateProjectList(ctx context.Context, id int64, req domain.ProjectListUpdateRequest, userID int64) (*domain.ProjectList, error)
	DeleteProjectList(ctx context.Context, id int64, userID int64) error
	AddProjectToList(ctx context.Context, listID, projectID, userID int64) error
	RemoveProjectFromList(ctx context.Context, listID, projectID, userID int64) error
	IsProjectInList(ctx context.Context, listID, projectID int64) (bool, error)
}

// Services aggregates all service interfaces
type Services struct {
	User        UserService
	Project     ProjectService
	Comment     CommentService
	Rating      RatingService
	ProjectList ProjectListService
}
