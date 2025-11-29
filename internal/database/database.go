package database

import (
	"golang/internal/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the database connection once
func InitDatabase() error {
	databaseName := config.GetEnv("DB_NAME")
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		return err
	}

	// Get underlying SQL database for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connected successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// ConnectDatabase is deprecated, use GetDB() instead
// Kept for backward compatibility
func ConnectDatabase() *gorm.DB {
	return DB
}
