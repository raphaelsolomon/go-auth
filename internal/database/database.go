package database

import (
	"golang/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	databaseName := config.GetEnv("DB_NAME")
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	return db
}
