package database

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

// ComputeHash computes a hash for a given password and salt combination.
func ComputeHash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 32*1024, 4, 32)
}

// User represents a single Okra user.
type User struct {
	gorm.Model
	Username       string `gorm:"unique;index"`
	HashedPassword []byte
	PasswordSalt   []byte
}

// CheckPassword checks a provided password against the stored password for a given user.
func (user *User) CheckPassword(password string) bool {
	return bytes.Equal(user.HashedPassword, ComputeHash(password, user.PasswordSalt))
}

// NewUser creates a new user.
func NewUser(username, password string) (*User, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error generating password salt: %w", err)
	}

	hash := ComputeHash(password, salt)

	instance := User{
		Username:       username,
		HashedPassword: hash,
		PasswordSalt:   salt,
	}

	tx := Instance.Create(&instance)
	if tx.Error != nil {
		return nil, fmt.Errorf("error inserting user into database: %w", tx.Error)
	}

	log.Info().Str("username", instance.Username).Uint("id", instance.ID).Msg("New user created")

	return &instance, nil
}
