package database

import (
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
