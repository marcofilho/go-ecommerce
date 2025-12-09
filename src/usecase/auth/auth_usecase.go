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

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type UseCase struct {
	userRepo    repository.UserRepository
	jwtProvider auth.TokenProvider
}

func NewUseCase(userRepo repository.UserRepository, jwtProvider auth.TokenProvider) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		jwtProvider: jwtProvider,
	}
}

type RegisterRequest struct {
	Email    string
	Password string
	Name     string
	Role     string
}

type LoginRequest struct {
	Email    string
	Password string
}

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
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("Email already registered")
	}

	role := entity.RoleCustomer
	if req.Role != "" {
		if req.Role == string(entity.RoleAdmin) {
			role = entity.RoleAdmin
		} else if req.Role == string(entity.RoleCustomer) {
			role = entity.RoleCustomer
		} else {
			return nil, errors.New("Invalid role. Must be 'customer' or 'admin'")
		}
	}

	user := &entity.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Role:      role,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

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

func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Invalid credentials")
	}

	if !user.IsActive() {
		return nil, errors.New("Account is inactive")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("Invalid credentials")
	}

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

func (uc *UseCase) ValidateToken(tokenString string) (*auth.Claims, error) {
	return uc.jwtProvider.ValidateToken(tokenString)
}
