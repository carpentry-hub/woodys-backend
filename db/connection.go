// Package db proporciona las utilidades para la conexion con la base de datos
package db

import (
	"github.com/carpentry-hub/woodys-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB variable utilizada para referenciar a la base de datos
var DB *gorm.DB

// Connection realiza la conexion con la base de datos
func Connection(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}
