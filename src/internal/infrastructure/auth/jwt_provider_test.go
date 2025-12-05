package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

func TestNewJWTProvider(t *testing.T) {
	provider := NewJWTProvider("test-secret-key", 24)

	if provider == nil {
		t.Fatal("NewJWTProvider() returned nil")
	}

	if provider.expirationHours != 24 {
		t.Errorf("NewJWTProvider() expirationHours = %d, want 24", provider.expirationHours)
	}
}

func TestJWTProvider_GenerateToken(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	user := &entity.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleCustomer,
	}

	token, err := provider.GenerateToken(user)

	if err != nil {
		t.Fatalf("GenerateToken() error = %v, want nil", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestJWTProvider_ValidateToken_Success(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	user := &entity.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleCustomer,
	}

	token, err := provider.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := provider.ValidateToken(token)

	if err != nil {
		t.Fatalf("ValidateToken() error = %v, want nil", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("ValidateToken() UserID = %s, want %s", claims.UserID, user.ID)
	}

	if claims.Email != user.Email {
		t.Errorf("ValidateToken() Email = %s, want %s", claims.Email, user.Email)
	}

	if claims.Role != user.Role {
		t.Errorf("ValidateToken() Role = %s, want %s", claims.Role, user.Role)
	}
}

func TestJWTProvider_ValidateToken_InvalidToken(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	_, err := provider.ValidateToken("invalid.token.here")

	if err == nil {
		t.Error("ValidateToken() should return error for invalid token")
	}
}

func TestJWTProvider_ValidateToken_ExpiredToken(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 0)

	user := &entity.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleCustomer,
	}

	token, err := provider.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	time.Sleep(2 * time.Second)

	_, err = provider.ValidateToken(token)

	if err == nil {
		t.Error("ValidateToken() should return error for expired token")
	}
}

func TestJWTProvider_ValidateToken_WrongSecret(t *testing.T) {
	provider1 := NewJWTProvider("secret-key-one", 24)

	user := &entity.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleCustomer,
	}

	token, err := provider1.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	provider2 := NewJWTProvider("secret-key-two", 24)

	_, err = provider2.ValidateToken(token)

	if err == nil {
		t.Error("ValidateToken() should return error for token signed with different secret")
	}
}

func TestJWTProvider_GenerateToken_AdminRole(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	user := &entity.User{
		ID:    uuid.New(),
		Email: "admin@example.com",
		Name:  "Admin User",
		Role:  entity.RoleAdmin,
	}

	token, err := provider.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := provider.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.Role != entity.RoleAdmin {
		t.Errorf("ValidateToken() Role = %s, want %s", claims.Role, entity.RoleAdmin)
	}
}

func TestJWTProvider_ValidateToken_InvalidSigningMethod(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	// Create a token with RSA signing method instead of HMAC
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.invalid"

	_, err := provider.ValidateToken(tokenString)

	if err == nil {
		t.Error("ValidateToken() should return error for token with invalid signing method")
	}
}

func TestJWTProvider_ValidateToken_MalformedToken(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	_, err := provider.ValidateToken("not.a.valid.jwt.token")

	if err == nil {
		t.Error("ValidateToken() should return error for malformed token")
	}
}

func TestJWTProvider_ValidateToken_EmptyToken(t *testing.T) {
	provider := NewJWTProvider("test-secret-key-for-jwt", 24)

	_, err := provider.ValidateToken("")

	if err == nil {
		t.Error("ValidateToken() should return error for empty token")
	}
}
