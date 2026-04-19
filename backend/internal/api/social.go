package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/db"
)

type SocialHandler struct {
	queries *db.Queries
}

func NewSocialHandler(queries *db.Queries) *SocialHandler {
	return &SocialHandler{queries: queries}
}

func (h *SocialHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/users/:username/follow", h.FollowUser)
	router.Delete("/users/:username/follow", h.UnfollowUser)
	router.Get("/users/:username/followers", h.GetFollowers)
	router.Get("/users/:username/following", h.GetFollowing)
	router.Post("/logs/:id/like", h.LikeLog)
	router.Delete("/logs/:id/like", h.UnlikeLog)
	router.Post("/reviews/:id/like", h.LikeReview)
	router.Delete("/reviews/:id/like", h.UnlikeReview)
}

func (h *SocialHandler) FollowUser(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(string)
	username := c.Params("username")

	targetUser, err := h.queries.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if targetUser.ID.String() == followerID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot follow yourself",
		})
	}

	followerUUID, err := stringToPgUUID(followerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid follower ID",
		})
	}

	isFollowing, _ := h.queries.IsFollowing(c.Context(), db.IsFollowingParams{
		FollowerID:  followerUUID,
		FollowingID: targetUser.ID,
	})

	if isFollowing {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already following this user",
		})
	}

	_, err = h.queries.CreateFollow(c.Context(), db.CreateFollowParams{
		FollowerID:  followerUUID,
		FollowingID: targetUser.ID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to follow user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "now following user",
	})
}

func (h *SocialHandler) UnfollowUser(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(string)
	username := c.Params("username")

	targetUser, err := h.queries.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	followerUUID, err := stringToPgUUID(followerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid follower ID",
		})
	}

	err = h.queries.DeleteFollow(c.Context(), db.DeleteFollowParams{
		FollowerID:  followerUUID,
		FollowingID: targetUser.ID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unfollow user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "unfollowed user",
	})
}

func (h *SocialHandler) GetFollowers(c *fiber.Ctx) error {
	username := c.Params("username")

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 50 {
		limit = 50
	}

	targetUser, err := h.queries.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	followers, err := h.queries.ListFollowers(c.Context(), db.ListFollowersParams{
		FollowingID: targetUser.ID,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch followers",
		})
	}

	response := make([]UserResponse, 0, len(followers))
	for _, follower := range followers {
		response = append(response, dbUserToResponse(follower))
	}

	return c.JSON(fiber.Map{
		"followers": response,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *SocialHandler) GetFollowing(c *fiber.Ctx) error {
	username := c.Params("username")

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 50 {
		limit = 50
	}

	targetUser, err := h.queries.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	following, err := h.queries.ListFollowing(c.Context(), db.ListFollowingParams{
		FollowerID: targetUser.ID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch following",
		})
	}

	response := make([]UserResponse, 0, len(following))
	for _, user := range following {
		response = append(response, dbUserToResponse(user))
	}

	return c.JSON(fiber.Map{
		"following": response,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *SocialHandler) LikeLog(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	logID := c.Params("id")

	userUUID, err := stringToPgUUID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	logUUID, err := stringToPgUUID(logID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid log ID",
		})
	}

	_, err = h.queries.GetLogByID(c.Context(), logUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "log not found",
		})
	}

	_, err = h.queries.GetLikeByLog(c.Context(), db.GetLikeByLogParams{
		UserID: userUUID,
		LogID:  logUUID,
	})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already liked this log",
		})
	}

	_, err = h.queries.CreateLike(c.Context(), db.CreateLikeParams{
		UserID:   userUUID,
		LogID:    logUUID,
		ReviewID: pgtype.UUID{},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to like log",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "log liked",
	})
}

func (h *SocialHandler) UnlikeLog(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	logID := c.Params("id")

	userUUID, err := stringToPgUUID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	logUUID, err := stringToPgUUID(logID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid log ID",
		})
	}

	err = h.queries.DeleteLike(c.Context(), db.DeleteLikeParams{
		UserID:   userUUID,
		LogID:    logUUID,
		ReviewID: pgtype.UUID{},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unlike log",
		})
	}

	return c.JSON(fiber.Map{
		"message": "log unliked",
	})
}

func (h *SocialHandler) LikeReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	reviewID := c.Params("id")

	userUUID, err := stringToPgUUID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	reviewUUID, err := stringToPgUUID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid review ID",
		})
	}

	_, err = h.queries.GetReviewByID(c.Context(), reviewUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "review not found",
		})
	}

	_, err = h.queries.GetLikeByReview(c.Context(), db.GetLikeByReviewParams{
		UserID:   userUUID,
		ReviewID: reviewUUID,
	})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already liked this review",
		})
	}

	_, err = h.queries.CreateLike(c.Context(), db.CreateLikeParams{
		UserID:   userUUID,
		LogID:    pgtype.UUID{},
		ReviewID: reviewUUID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to like review",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "review liked",
	})
}

func (h *SocialHandler) UnlikeReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	reviewID := c.Params("id")

	userUUID, err := stringToPgUUID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	reviewUUID, err := stringToPgUUID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid review ID",
		})
	}

	err = h.queries.DeleteLike(c.Context(), db.DeleteLikeParams{
		UserID:   userUUID,
		LogID:    pgtype.UUID{},
		ReviewID: reviewUUID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unlike review",
		})
	}

	return c.JSON(fiber.Map{
		"message": "review unliked",
	})
}
