package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
)

type AuthHandler struct {
	authUseCase authUseCase.AuthService
}

func NewAuthHandler(uc authUseCase.AuthService) *AuthHandler {
	return &AuthHandler{
		authUseCase: uc,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role,omitempty" example:"customer"` // Optional: customer (default) or admin
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account. Public registration creates customer accounts. Creating admin accounts requires admin authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - Admin authentication required for admin role"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - Only admins can create admin accounts"
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Security check: Only authenticated admins can create admin accounts
	if req.Role == "admin" || req.Role == string(entity.RoleAdmin) {
		claims, err := middleware.GetUserFromContext(r)
		if err != nil {
			// Not authenticated
			respondError(w, http.StatusUnauthorized, "Only authenticated admin users can create admin accounts")
			return
		}
		if claims.Role != entity.RoleAdmin {
			// Authenticated but not admin
			respondError(w, http.StatusForbidden, "Only admin users can create admin accounts")
			return
		}
	}

	authReq := authUseCase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}

	response, err := h.authUseCase.Register(r.Context(), authReq)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "Password is required")
		return
	}

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
