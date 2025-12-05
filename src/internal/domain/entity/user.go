package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Role         Role      `gorm:"type:varchar(50);not null;default:customer"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Validate validates user fields
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}

	if u.Role != RoleAdmin && u.Role != RoleCustomer {
		return errors.New("invalid role")
	}

	return nil
}

// SetPassword hashes and sets the user password
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Active
}
