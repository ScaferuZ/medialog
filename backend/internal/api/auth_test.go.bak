package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/api"
	"github.com/medialogg/backend/internal/config"
	"github.com/medialogg/backend/internal/db"
	"github.com/medialogg/backend/internal/middleware"
	testutil "github.com/medialogg/backend/internal/test"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "integration-test-secret"

var authTestDB *testutil.TestDatabase

type errorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var err error
	authTestDB, err = testutil.NewTestDatabase(ctx)
	if err != nil {
		log.Fatalf("failed to initialize auth test database: %v", err)
	}

	code := m.Run()

	closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Minute)
	defer closeCancel()
	if err := authTestDB.Close(closeCtx); err != nil {
		log.Printf("failed to close auth test database: %v", err)
	}

	os.Exit(code)
}

func TestAuthRegister(t *testing.T) {
	t.Run("successful registration returns 201 with tokens", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
			"username": "new_user",
			"email":    "new_user@example.com",
			"password": "password123",
		}, "")

		if status != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, status, string(body))
		}

		resp := decodeJSON[api.AuthResponse](t, body)
		if resp.User.Username != "new_user" {
			t.Fatalf("expected username %q, got %q", "new_user", resp.User.Username)
		}
		if resp.User.Email != "new_user@example.com" {
			t.Fatalf("expected email %q, got %q", "new_user@example.com", resp.User.Email)
		}
		if resp.AccessToken == "" {
			t.Fatal("expected access token to be present")
		}
		if resp.RefreshToken == "" {
			t.Fatal("expected refresh token to be present")
		}
		if resp.ExpiresIn <= 0 {
			t.Fatalf("expected positive expires_in, got %d", resp.ExpiresIn)
		}

		storedUser, err := authTestDB.Queries.GetUserByUsername(context.Background(), "new_user")
		if err != nil {
			t.Fatalf("failed to fetch stored user: %v", err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte("password123")); err != nil {
			t.Fatalf("expected stored password hash to match submitted password: %v", err)
		}
	})

	t.Run("duplicate username returns 409", func(t *testing.T) {
		app := newAuthTestApp(t)
		createTestUser(t, "taken_user", "taken_user@example.com", "password123")

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
			"username": "taken_user",
			"email":    "other@example.com",
			"password": "password123",
		}, "")

		if status != http.StatusConflict {
			t.Fatalf("expected status %d, got %d: %s", http.StatusConflict, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "username or email already exists" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		app := newAuthTestApp(t)
		createTestUser(t, "existing_user", "duplicate@example.com", "password123")

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
			"username": "other_user",
			"email":    "duplicate@example.com",
			"password": "password123",
		}, "")

		if status != http.StatusConflict {
			t.Fatalf("expected status %d, got %d: %s", http.StatusConflict, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "username or email already exists" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("invalid username format returns 400", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
			"username": "bad-user!",
			"email":    "bad-user@example.com",
			"password": "password123",
		}, "")

		if status != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Details["username"] != "must contain only letters, numbers, and underscores" {
			t.Fatalf("unexpected username validation message: %q", resp.Details["username"])
		}
	})

	t.Run("short password returns 400", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
			"username": "valid_user",
			"email":    "valid_user@example.com",
			"password": "short",
		}, "")

		if status != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Details["password"] != "must be at least 8 characters" {
			t.Fatalf("unexpected password validation message: %q", resp.Details["password"])
		}
	})

	t.Run("missing fields returns 400", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{}, "")

		if status != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Details["username"] != "this field is required" {
			t.Fatalf("unexpected username validation message: %q", resp.Details["username"])
		}
		if resp.Details["email"] != "this field is required" {
			t.Fatalf("unexpected email validation message: %q", resp.Details["email"])
		}
		if resp.Details["password"] != "this field is required" {
			t.Fatalf("unexpected password validation message: %q", resp.Details["password"])
		}
	})
}

func TestAuthLogin(t *testing.T) {
	t.Run("successful login returns 200 with tokens", func(t *testing.T) {
		app := newAuthTestApp(t)
		createTestUser(t, "login_user", "login_user@example.com", "password123")

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
			"username": "login_user",
			"password": "password123",
		}, "")

		if status != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, status, string(body))
		}

		resp := decodeJSON[api.AuthResponse](t, body)
		if resp.User.Username != "login_user" {
			t.Fatalf("expected username %q, got %q", "login_user", resp.User.Username)
		}
		if resp.AccessToken == "" {
			t.Fatal("expected access token to be present")
		}
		if resp.RefreshToken == "" {
			t.Fatal("expected refresh token to be present")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)
		createTestUser(t, "wrong_password_user", "wrong_password_user@example.com", "password123")

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
			"username": "wrong_password_user",
			"password": "not_the_password",
		}, "")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "invalid credentials" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("non-existent user returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
			"username": "missing_user",
			"password": "password123",
		}, "")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "invalid credentials" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("missing fields returns 400", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{}, "")

		if status != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Details["username"] != "this field is required" {
			t.Fatalf("unexpected username validation message: %q", resp.Details["username"])
		}
		if resp.Details["password"] != "this field is required" {
			t.Fatalf("unexpected password validation message: %q", resp.Details["password"])
		}
	})
}

func TestAuthRefresh(t *testing.T) {
	t.Run("valid refresh token returns new access token", func(t *testing.T) {
		app := newAuthTestApp(t)
		user := createTestUser(t, "refresh_user", "refresh_user@example.com", "password123")
		tokens := generateTestTokens(t, user)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/refresh", map[string]string{
			"refresh_token": tokens.RefreshToken,
		}, "")

		if status != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, status, string(body))
		}

		resp := decodeJSON[api.AuthResponse](t, body)
		if resp.User.Username != "refresh_user" {
			t.Fatalf("expected username %q, got %q", "refresh_user", resp.User.Username)
		}
		if resp.AccessToken == "" {
			t.Fatal("expected refreshed access token to be present")
		}
		if resp.RefreshToken == "" {
			t.Fatal("expected refreshed refresh token to be present")
		}

		claims, err := config.ValidateToken(resp.AccessToken, testJWTSecret)
		if err != nil {
			t.Fatalf("expected refreshed access token to be valid: %v", err)
		}
		if claims.Username != user.Username {
			t.Fatalf("expected refreshed access token username %q, got %q", user.Username, claims.Username)
		}
	})

	t.Run("invalid refresh token returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/refresh", map[string]string{
			"refresh_token": "not-a-valid-token",
		}, "")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "invalid or expired refresh token" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("expired refresh token returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)
		user := createTestUser(t, "expired_refresh_user", "expired_refresh_user@example.com", "password123")
		expiredRefreshToken := generateExpiredRefreshToken(t, user)

		status, body := doJSONRequest(t, app, http.MethodPost, "/api/auth/refresh", map[string]string{
			"refresh_token": expiredRefreshToken,
		}, "")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "invalid or expired refresh token" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})
}

func TestAuthMe(t *testing.T) {
	t.Run("valid token returns user profile", func(t *testing.T) {
		app := newAuthTestApp(t)
		user := createTestUser(t, "profile_user", "profile_user@example.com", "password123")
		tokens := generateTestTokens(t, user)

		status, body := doJSONRequest(t, app, http.MethodGet, "/api/auth/me", nil, tokens.AccessToken)

		if status != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, status, string(body))
		}

		resp := decodeJSON[api.UserResponse](t, body)
		if resp.ID != user.ID.String() {
			t.Fatalf("expected user id %q, got %q", user.ID.String(), resp.ID)
		}
		if resp.Username != user.Username {
			t.Fatalf("expected username %q, got %q", user.Username, resp.Username)
		}
		if resp.Email != user.Email {
			t.Fatalf("expected email %q, got %q", user.Email, resp.Email)
		}
	})

	t.Run("missing token returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodGet, "/api/auth/me", nil, "")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "missing authorization header" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		app := newAuthTestApp(t)

		status, body := doJSONRequest(t, app, http.MethodGet, "/api/auth/me", nil, "not-a-valid-token")

		if status != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, status, string(body))
		}

		resp := decodeJSON[errorResponse](t, body)
		if resp.Error != "invalid or expired token" {
			t.Fatalf("unexpected error message: %q", resp.Error)
		}
	})
}

func newAuthTestApp(t *testing.T) *fiber.App {
	t.Helper()

	if err := authTestDB.Reset(context.Background()); err != nil {
		t.Fatalf("failed to reset auth test database: %v", err)
	}

	t.Cleanup(func() {
		if err := authTestDB.Reset(context.Background()); err != nil {
			t.Errorf("failed to clean auth test database: %v", err)
		}
	})

	app := fiber.New()
	handler := api.NewAuthHandler(authTestDB.Queries, testJWTSecret)

	apiGroup := app.Group("/api")
	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
	authGroup.Get("/me", middleware.AuthMiddleware(testJWTSecret), handler.Me)

	return app
}

func doJSONRequest(t *testing.T, app *fiber.App, method, path string, body any, bearerToken string) (int, []byte) {
	t.Helper()

	var requestBody io.Reader = http.NoBody
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		requestBody = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return resp.StatusCode, responseBody
}

func decodeJSON[T any](t *testing.T, body []byte) T {
	t.Helper()

	var response T
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to decode JSON response %q: %v", string(body), err)
	}

	return response
}

func createTestUser(t *testing.T, username, email, password string) db.User {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash test password: %v", err)
	}

	user, err := authTestDB.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  pgtype.Text{},
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

func generateTestTokens(t *testing.T, user db.User) *config.TokenPair {
	t.Helper()

	tokens, err := config.GenerateTokenPair(user.ID, user.Username, user.Email, testJWTSecret)
	if err != nil {
		t.Fatalf("failed to generate test tokens: %v", err)
	}

	return tokens
}

func generateExpiredRefreshToken(t *testing.T, user db.User) string {
	t.Helper()

	claims := config.Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("failed to generate expired refresh token: %v", err)
	}

	return tokenString
}
