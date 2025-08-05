package repositories

import (
	"gorm.io/gorm"
)

// NewRepositories creates and returns all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Project:     NewProjectRepository(db),
		Comment:     NewCommentRepository(db),
		Rating:      NewRatingRepository(db),
		ProjectList: NewProjectListRepository(db),
	}
}
