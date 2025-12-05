package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for storing user data in request context
	UserContextKey ContextKey = "user"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authUseCase *authUseCase.UseCase
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(uc *authUseCase.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: uc,
	}
}

// Authenticate validates JWT token and injects user context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.writeError(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.writeError(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.authUseCase.ValidateToken(tokenString)
		if err != nil {
			m.writeError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Inject user data into context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if the authenticated user has the required role
func (m *AuthMiddleware) RequireRole(role entity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				m.writeError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check role
			if claims.Role != role {
				m.writeError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission checks if the authenticated user has the required permission
func (m *AuthMiddleware) RequirePermission(permission Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				m.writeError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user's role has the required permission
			if !HasPermission(claims.Role, permission) {
				m.writeError(w, "Forbidden: insufficient permissions for this action", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth validates token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := m.authUseCase.ValidateToken(parts[1])
			if err == nil {
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
