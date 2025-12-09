package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
)

// mockAuthService is a mock implementation of AuthService for testing
type mockAuthService struct {
	registerFunc      func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error)
	loginFunc         func(ctx context.Context, req authUseCase.LoginRequest) (*authUseCase.AuthResponse, error)
	validateTokenFunc func(tokenString string) (*auth.Claims, error)
}

func (m *mockAuthService) Register(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil, errors.New("Not implemented")
}

func (m *mockAuthService) Login(ctx context.Context, req authUseCase.LoginRequest) (*authUseCase.AuthResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return nil, errors.New("Not implemented")
}

func (m *mockAuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(tokenString)
	}
	return nil, errors.New("Not implemented")
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
			return &authUseCase.AuthResponse{
				Token:     "test-token",
				UserID:    uuid.New(),
				Email:     req.Email,
				Name:      req.Name,
				Role:      entity.RoleCustomer,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Register() status = %d, want %d", w.Code, http.StatusCreated)
	}

	var response authUseCase.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Email != reqBody.Email {
		t.Errorf("Register() email = %s, want %s", response.Email, reqBody.Email)
	}

	if response.Token == "" {
		t.Error("Register() returned empty token")
	}
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Register() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Register_UseCaseError(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
			return nil, errors.New("Email already registered")
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Register() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := &mockAuthService{
		loginFunc: func(ctx context.Context, req authUseCase.LoginRequest) (*authUseCase.AuthResponse, error) {
			return &authUseCase.AuthResponse{
				Token:     "test-token",
				UserID:    uuid.New(),
				Email:     req.Email,
				Name:      "Test User",
				Role:      entity.RoleCustomer,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Login() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response authUseCase.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Email != reqBody.Email {
		t.Errorf("Login() email = %s, want %s", response.Email, reqBody.Email)
	}

	if response.Token == "" {
		t.Error("Login() returned empty token")
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Login() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := &mockAuthService{
		loginFunc: func(ctx context.Context, req authUseCase.LoginRequest) (*authUseCase.AuthResponse, error) {
			return nil, errors.New("Invalid credentials")
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Login() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_InactiveAccount(t *testing.T) {
	mockService := &mockAuthService{
		loginFunc: func(ctx context.Context, req authUseCase.LoginRequest) (*authUseCase.AuthResponse, error) {
			return nil, errors.New("Account is inactive")
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := LoginRequest{
		Email:    "inactive@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Login() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Register_AdminWithoutAuth(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "admin@example.com",
		Password: "password123",
		Name:     "Admin User",
		Role:     "admin",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Register() status = %d, want %d for admin creation without auth", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Register_AdminWithCustomerAuth(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "admin@example.com",
		Password: "password123",
		Name:     "Admin User",
		Role:     "admin",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	customerClaims := &auth.Claims{
		UserID: uuid.New(),
		Email:  "customer@example.com",
		Role:   entity.RoleCustomer,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, customerClaims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Register() status = %d, want %d for admin creation by customer", w.Code, http.StatusForbidden)
	}
}

func TestAuthHandler_Register_AdminWithAdminAuth(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
			return &authUseCase.AuthResponse{
				Token:     "test-token",
				UserID:    uuid.New(),
				Email:     req.Email,
				Name:      req.Name,
				Role:      entity.RoleAdmin,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "newadmin@example.com",
		Password: "password123",
		Name:     "New Admin",
		Role:     "admin",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	adminClaims := &auth.Claims{
		UserID: uuid.New(),
		Email:  "admin@example.com",
		Role:   entity.RoleAdmin,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, adminClaims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Register() status = %d, want %d for admin creation by admin", w.Code, http.StatusCreated)
	}

	var response authUseCase.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Role != entity.RoleAdmin {
		t.Errorf("Register() role = %s, want %s", response.Role, entity.RoleAdmin)
	}
}

func TestAuthHandler_Register_CustomerWithoutRole(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
			return &authUseCase.AuthResponse{
				Token:     "test-token",
				UserID:    uuid.New(),
				Email:     req.Email,
				Name:      req.Name,
				Role:      entity.RoleCustomer,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "customer@example.com",
		Password: "password123",
		Name:     "Customer User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Register() status = %d, want %d", w.Code, http.StatusCreated)
	}

	var response authUseCase.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Role != entity.RoleCustomer {
		t.Errorf("Register() role = %s, want %s", response.Role, entity.RoleCustomer)
	}
}

func TestAuthHandler_Register_InvalidRole(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
			return nil, errors.New("Invalid role. Must be 'customer' or 'admin'")
		},
	}

	handler := NewAuthHandler(mockService)

	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "superuser",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Register() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
