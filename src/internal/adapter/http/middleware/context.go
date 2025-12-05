package middleware

import (
	"errors"
	"net/http"

	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
)

// GetUserFromContext retrieves the authenticated user claims from request context
func GetUserFromContext(r *http.Request) (*auth.Claims, error) {
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return claims, nil
}
