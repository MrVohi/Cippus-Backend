package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cippus-backend/internal/models"
)

func Connect(dsn string) (*gorm.DB, error) {
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	err = db.AutoMigrate(models.Models...)
	if err != nil {
		return nil, err
	}

	return db, err
}
