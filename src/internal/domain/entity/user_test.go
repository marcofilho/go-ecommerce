package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestUser_SetPassword(t *testing.T) {
	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	password := "validPassword123"
	err := user.SetPassword(password)

	if err != nil {
		t.Errorf("SetPassword() error = %v, want nil", err)
	}

	if user.PasswordHash == "" {
		t.Error("SetPassword() did not set PasswordHash")
	}

	if user.PasswordHash == password {
		t.Error("SetPassword() stored plain password instead of hash")
	}
}

func TestUser_SetPassword_TooShort(t *testing.T) {
	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := user.SetPassword("short")

	if err == nil {
		t.Error("SetPassword() with short password should return error")
	}
}

func TestUser_CheckPassword_Valid(t *testing.T) {
	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	password := "validPassword123"
	user.SetPassword(password)

	if !user.CheckPassword(password) {
		t.Error("CheckPassword() returned false for valid password")
	}
}

func TestUser_CheckPassword_Invalid(t *testing.T) {
	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	password := "validPassword123"
	user.SetPassword(password)

	if user.CheckPassword("wrongPassword") {
		t.Error("CheckPassword() returned true for invalid password")
	}
}

func TestUser_Validate_Success(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         RoleCustomer,
		Active:       true,
	}

	err := user.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestUser_Validate_EmptyEmail(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         RoleCustomer,
	}

	err := user.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty email")
	}
}

func TestUser_Validate_ShortName(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "A",
		Role:         RoleCustomer,
	}

	err := user.Validate()
	if err == nil {
		t.Error("Validate() should return error for name shorter than 2 characters")
	}
}

func TestUser_Validate_EmptyName(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "",
		Role:         RoleCustomer,
	}

	err := user.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty name")
	}
}

func TestUser_Validate_InvalidRole(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         Role("invalid"),
	}

	err := user.Validate()
	if err == nil {
		t.Error("Validate() should return error for invalid role")
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"Admin role", RoleAdmin, true},
		{"Customer role", RoleCustomer, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			if got := user.IsAdmin(); got != tt.want {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		active bool
		want   bool
	}{
		{"Active user", true, true},
		{"Inactive user", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Active: tt.active}
			if got := user.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
