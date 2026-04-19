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

func (h *SocialHandler) RegisterPublicUserRoutes(router fiber.Router) {
	router.Get("/:username/followers", h.GetFollowers)
	router.Get("/:username/following", h.GetFollowing)
}

func (h *SocialHandler) RegisterProtectedUserRoutes(router fiber.Router) {
	router.Post("/:username/follow", h.FollowUser)
	router.Delete("/:username/follow", h.UnfollowUser)
}

func (h *SocialHandler) RegisterProtectedLogRoutes(router fiber.Router) {
	router.Post("/:id/like", h.LikeLog)
	router.Delete("/:id/like", h.UnlikeLog)
}

func (h *SocialHandler) RegisterProtectedReviewRoutes(router fiber.Router) {
	router.Post("/:id/like", h.LikeReview)
	router.Delete("/:id/like", h.UnlikeReview)
}

func (h *SocialHandler) FollowUser(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(string)
	username := c.Params("username")

	targetUser, err := h.queries.GetUserByUsername(c.UserContext(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if uuidToString(targetUser.ID) == followerID {
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

	isFollowing, _ := h.queries.IsFollowing(c.UserContext(), db.IsFollowingParams{
		FollowerID:  followerUUID,
		FollowingID: targetUser.ID,
	})

	if isFollowing {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already following this user",
		})
	}

	_, err = h.queries.CreateFollow(c.UserContext(), db.CreateFollowParams{
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

	targetUser, err := h.queries.GetUserByUsername(c.UserContext(), username)
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

	err = h.queries.DeleteFollow(c.UserContext(), db.DeleteFollowParams{
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

	targetUser, err := h.queries.GetUserByUsername(c.UserContext(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	followers, err := h.queries.ListFollowers(c.UserContext(), db.ListFollowersParams{
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
	response = ensureSlice(response)

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

	targetUser, err := h.queries.GetUserByUsername(c.UserContext(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	following, err := h.queries.ListFollowing(c.UserContext(), db.ListFollowingParams{
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
	response = ensureSlice(response)

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

	_, err = h.queries.GetLogByID(c.UserContext(), logUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "log not found",
		})
	}

	_, err = h.queries.GetLikeByLog(c.UserContext(), db.GetLikeByLogParams{
		UserID: userUUID,
		LogID:  logUUID,
	})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already liked this log",
		})
	}

	_, err = h.queries.CreateLike(c.UserContext(), db.CreateLikeParams{
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

	err = h.queries.DeleteLike(c.UserContext(), db.DeleteLikeParams{
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

	_, err = h.queries.GetReviewByID(c.UserContext(), reviewUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "review not found",
		})
	}

	_, err = h.queries.GetLikeByReview(c.UserContext(), db.GetLikeByReviewParams{
		UserID:   userUUID,
		ReviewID: reviewUUID,
	})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "already liked this review",
		})
	}

	_, err = h.queries.CreateLike(c.UserContext(), db.CreateLikeParams{
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

	err = h.queries.DeleteLike(c.UserContext(), db.DeleteLikeParams{
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
