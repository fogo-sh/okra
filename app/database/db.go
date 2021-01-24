package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Instance is a shared, global Gorm instance. Initialize with CreateInstance.
var Instance *gorm.DB

// CreateInstance creates a global Gorm instance.
func CreateInstance(connectionString string) error {
	var err error
	Instance, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening db connection: %w", err)
	}

	return nil
}

// Migrate automatically creates the necessary database structure for all models.
func Migrate() {
	if Instance == nil {
		panic(errors.New("Attempted to migrate uninstantiated database - ensure you call database.CreateInstance before making any database calls"))
	}
	Instance.AutoMigrate(&User{})
}
