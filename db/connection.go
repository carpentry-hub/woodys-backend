package db

import (
	"github.com/carpentry-hub/woodys-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	var error error
	DB, error = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return error
}
