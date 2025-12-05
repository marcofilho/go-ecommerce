package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type Claims struct {
	UserID uuid.UUID   `json:"user_id"`
	Email  string      `json:"email"`
	Role   entity.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secretKey       string
	expirationHours int
}

func NewJWTProvider(secretKey string, expirationHours int) *JWTProvider {
	return &JWTProvider{
		secretKey:       secretKey,
		expirationHours: expirationHours,
	}
}

// GenerateToken generates a new JWT token for a user
func (p *JWTProvider) GenerateToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(p.expirationHours) * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-ecommerce",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (p *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(p.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}
