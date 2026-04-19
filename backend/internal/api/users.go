package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/medialogg/backend/internal/db"
)

type UsersHandler struct {
	queries *db.Queries
}

func NewUsersHandler(queries *db.Queries) *UsersHandler {
	return &UsersHandler{queries: queries}
}

// RegisterRoutes registers user routes
func (h *UsersHandler) RegisterRoutes(router fiber.Router) {
	// Public profile
	router.Get("/users/:username", h.GetUserProfile)
	router.Get("/users/:username/stats", h.GetUserStats)
}

// UserProfileResponse represents a user's public profile
type UserProfileResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarUrl   *string `json:"avatar_url,omitempty"`
	IsPublic    bool    `json:"is_public"`
	CreatedAt   string  `json:"created_at"`
}

// GetUserProfile returns a user's public profile
func (h *UsersHandler) GetUserProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.queries.GetUserByUsername(c.UserContext(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch user profile",
		})
	}

	// TODO: Check if profile is private and viewer is not following

	resp := UserProfileResponse{
		ID:        uuidToString(user.ID),
		Username:  user.Username,
		IsPublic:  user.IsPublic,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02"),
	}

	if user.DisplayName.Valid {
		resp.DisplayName = &user.DisplayName.String
	}
	if user.Bio.Valid {
		resp.Bio = &user.Bio.String
	}
	if user.AvatarUrl.Valid {
		resp.AvatarUrl = &user.AvatarUrl.String
	}

	return c.JSON(resp)
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	CompletedCount  int64   `json:"completed_count"`
	InProgressCount int64   `json:"in_progress_count"`
	PlannedCount    int64   `json:"planned_count"`
	DroppedCount    int64   `json:"dropped_count"`
	AverageRating   float64 `json:"average_rating"`
	TotalMedia      int64   `json:"total_media"`
	FollowersCount  int64   `json:"followers_count"`
	FollowingCount  int64   `json:"following_count"`
}

// GetUserStats returns a user's statistics
func (h *UsersHandler) GetUserStats(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get user
	user, err := h.queries.GetUserByUsername(c.UserContext(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch user",
		})
	}

	// Get log stats
	logStats, err := h.queries.GetUserStats(c.UserContext(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch stats",
		})
	}

	// Get follower counts
	followersCount, err := h.queries.CountFollowers(c.UserContext(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch follower count",
		})
	}

	followingCount, err := h.queries.CountFollowing(c.UserContext(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch following count",
		})
	}

	return c.JSON(UserStatsResponse{
		CompletedCount:  logStats.CompletedCount,
		InProgressCount: logStats.InProgressCount,
		PlannedCount:    logStats.PlannedCount,
		DroppedCount:    logStats.DroppedCount,
		AverageRating:   logStats.AverageRating,
		TotalMedia:      logStats.TotalMedia,
		FollowersCount:  followersCount,
		FollowingCount:  followingCount,
	})
}
