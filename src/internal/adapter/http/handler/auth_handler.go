package handler

import (
	"encoding/json"
	"net/http"

	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUseCase *authUseCase.UseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(uc *authUseCase.UseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: uc,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} authUseCase.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}

	// Register user
	authReq := authUseCase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	response, err := h.authUseCase.Register(r.Context(), authReq)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} authUseCase.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}

	// Login user
	authReq := authUseCase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.authUseCase.Login(r.Context(), authReq)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}
