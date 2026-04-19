package api

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/config"
	"github.com/medialogg/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName *string   `json:"display_name,omitempty"`
	AvatarUrl   *string   `json:"avatar_url,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuthHandler struct {
	queries   *db.Queries
	jwtSecret string
}

func NewAuthHandler(queries *db.Queries, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[strings.ToLower(e.Field())] = validationErrorMessage(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": errors,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to process password",
		})
	}

	displayName := pgtype.Text{Valid: false}
	user, err := h.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username or email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	tokens, err := config.GenerateTokenPair(uuidToString(user.ID), user.Username, user.Email, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate tokens",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		User:         dbUserToResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[strings.ToLower(e.Field())] = validationErrorMessage(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": errors,
		})
	}

	user, err := h.queries.GetUserByUsername(context.Background(), req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve user",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	tokens, err := config.GenerateTokenPair(uuidToString(user.ID), user.Username, user.Email, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate tokens",
		})
	}

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		User:         dbUserToResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[strings.ToLower(e.Field())] = validationErrorMessage(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": errors,
		})
	}

	claims, err := config.ValidateToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired refresh token",
		})
	}

	userID, err := stringToPgUUID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	user, err := h.queries.GetUserByUsername(context.Background(), claims.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	_ = userID

	tokens, err := config.GenerateTokenPair(uuidToString(user.ID), user.Username, user.Email, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate tokens",
		})
	}

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		User:         dbUserToResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	username := c.Locals("username").(string)

	_ = userID

	user, err := h.queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dbUserToResponse(user))
}

func dbUserToResponse(user db.User) UserResponse {
	resp := UserResponse{
		ID:        uuidToString(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		IsPublic:  user.IsPublic,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}

	if user.DisplayName.Valid {
		resp.DisplayName = &user.DisplayName.String
	}
	if user.AvatarUrl.Valid {
		resp.AvatarUrl = &user.AvatarUrl.String
	}
	if user.Bio.Valid {
		resp.Bio = &user.Bio.String
	}

	return resp
}

func validationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return "must be at least " + e.Param() + " characters"
	case "max":
		return "must be at most " + e.Param() + " characters"
	case "email":
		return "must be a valid email address"
	case "username":
		return "must contain only letters, numbers, and underscores"
	default:
		return "invalid value"
	}
}
