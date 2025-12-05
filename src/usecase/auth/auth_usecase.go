package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
)

// UseCase handles authentication business logic
type UseCase struct {
	userRepo    repository.UserRepository
	jwtProvider *auth.JWTProvider
}

// NewUseCase creates a new auth use case
func NewUseCase(userRepo repository.UserRepository, jwtProvider *auth.JWTProvider) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		jwtProvider: jwtProvider,
	}
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email    string
	Password string
	Name     string
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token     string      `json:"token"`
	UserID    uuid.UUID   `json:"user_id"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Role      entity.Role `json:"role"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Register creates a new user account
func (uc *UseCase) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	user := &entity.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Role:      entity.RoleCustomer, // Default role
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.jwtProvider.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Should match JWT expiration
	}, nil
}

// Login authenticates a user
func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := uc.jwtProvider.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (uc *UseCase) ValidateToken(tokenString string) (*auth.Claims, error) {
	return uc.jwtProvider.ValidateToken(tokenString)
}
