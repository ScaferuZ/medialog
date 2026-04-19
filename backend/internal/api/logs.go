package api

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/db"
)

type LogsHandler struct {
	queries *db.Queries
}

func NewLogsHandler(queries *db.Queries) *LogsHandler {
	return &LogsHandler{queries: queries}
}

// RegisterRoutes registers log routes
func (h *LogsHandler) RegisterRoutes(router fiber.Router) {
	logs := router.Group("/logs")

	// Public routes
	logs.Get("/timeline", h.GetTimeline)

	// Protected routes
	logs.Get("/me", h.GetMyLogs)
	logs.Get("/me/:mediaId", h.GetLogForMedia)
	logs.Post("/", h.CreateLog)
	logs.Put("/:id", h.UpdateLog)
	logs.Delete("/:id", h.DeleteLog)
}

// GetTimeline returns activity feed from followed users
func (h *LogsHandler) GetTimeline(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.queries.GetTimeline(c.Context(), db.GetTimelineParams{
		FollowerID: userID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch timeline",
		})
	}

	return c.JSON(fiber.Map{
		"logs": logs,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetMyLogs returns current user's logs
func (h *LogsHandler) GetMyLogs(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.queries.ListLogsByUser(c.Context(), db.ListLogsByUserParams{
		UserID:  userID,
		Column2: status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch logs",
		})
	}

	return c.JSON(fiber.Map{
		"logs": logs,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetLogForMedia returns log for specific media
func (h *LogsHandler) GetLogForMedia(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	mediaIDStr := c.Params("mediaId")

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var mediaID pgtype.UUID
	if err := mediaID.Scan(mediaIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid media id",
		})
	}

	log, err := h.queries.GetLogByUserAndMedia(c.Context(), db.GetLogByUserAndMediaParams{
		UserID:  userID,
		MediaID: mediaID,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch log",
		})
	}

	return c.JSON(log)
}

// CreateLogRequest represents a log creation request
type CreateLogRequest struct {
	MediaID          string   `json:"media_id" validate:"required"`
	Status           string   `json:"status" validate:"required,oneof=planned in_progress completed dropped"`
	Rating           *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=10"`
	StartedAt        *string  `json:"started_at,omitempty"`
	CompletedAt      *string  `json:"completed_at,omitempty"`
	RewatchCount     int32    `json:"rewatch_count,omitempty"`
	Progress         *int32   `json:"progress,omitempty"`
	Total            *int32   `json:"total,omitempty"`
	Note             string   `json:"note,omitempty"`
	ContainsSpoilers bool     `json:"contains_spoilers,omitempty"`
}

// CreateLog creates a new log entry
func (h *LogsHandler) CreateLog(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var req CreateLogRequest
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

	var mediaID pgtype.UUID
	if err := mediaID.Scan(req.MediaID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid media id",
		})
	}

	// Check if log already exists for this user+media
	_, err := h.queries.GetLogByUserAndMedia(c.Context(), db.GetLogByUserAndMediaParams{
		UserID:  userID,
		MediaID: mediaID,
	})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "log already exists for this media",
		})
	}

	// Build params
	params := db.CreateLogParams{
		UserID:           userID,
		MediaID:          mediaID,
		Status:           req.Status,
		RewatchCount:     pgtype.Int4{Int32: req.RewatchCount, Valid: req.RewatchCount > 0},
		Note:             pgtype.Text{String: req.Note, Valid: req.Note != ""},
		ContainsSpoilers: pgtype.Bool{Bool: req.ContainsSpoilers, Valid: true},
	}

	if req.Rating != nil {
		ratingStr := strconv.FormatFloat(*req.Rating, 'f', 1, 64)
		var rating pgtype.Numeric
		if err := rating.Scan(ratingStr); err == nil {
			params.Rating = rating
		}
	}

	if req.StartedAt != nil {
		var startedAt pgtype.Date
		if err := startedAt.Scan(*req.StartedAt); err == nil {
			params.StartedAt = startedAt
		}
	}

	if req.CompletedAt != nil {
		var completedAt pgtype.Date
		if err := completedAt.Scan(*req.CompletedAt); err == nil {
			params.CompletedAt = completedAt
		}
	}

	if req.Progress != nil {
		params.Progress = pgtype.Int4{Int32: *req.Progress, Valid: true}
	}

	if req.Total != nil {
		params.Total = pgtype.Int4{Int32: *req.Total, Valid: true}
	}

	log, err := h.queries.CreateLog(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create log",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(log)
}

// UpdateLogRequest represents a log update request
type UpdateLogRequest struct {
	Status           *string  `json:"status,omitempty" validate:"omitempty,oneof=planned in_progress completed dropped"`
	Rating           *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=10"`
	StartedAt        *string  `json:"started_at,omitempty"`
	CompletedAt      *string  `json:"completed_at,omitempty"`
	RewatchCount     *int32   `json:"rewatch_count,omitempty"`
	Progress         *int32   `json:"progress,omitempty"`
	Total            *int32   `json:"total,omitempty"`
	Note             *string  `json:"note,omitempty"`
	ContainsSpoilers *bool    `json:"contains_spoilers,omitempty"`
}

// UpdateLog updates an existing log
func (h *LogsHandler) UpdateLog(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	logIDStr := c.Params("id")

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var logID pgtype.UUID
	if err := logID.Scan(logIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid log id",
		})
	}

	// Verify log exists and belongs to user
	existing, err := h.queries.GetLogByID(c.Context(), logID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch log",
		})
	}

	if existing.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "not authorized to update this log",
		})
	}

	var req UpdateLogRequest
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

	// Build update params - use zero values with Valid=false for optional fields
	params := db.UpdateLogParams{
		ID: logID,
	}

	if req.Status != nil {
		params.Status = *req.Status
	}

	if req.Rating != nil {
		ratingStr := strconv.FormatFloat(*req.Rating, 'f', 1, 64)
		var rating pgtype.Numeric
		if err := rating.Scan(ratingStr); err == nil {
			params.Rating = rating
		}
	}

	if req.StartedAt != nil {
		var startedAt pgtype.Date
		if err := startedAt.Scan(*req.StartedAt); err == nil {
			params.StartedAt = startedAt
		}
	}

	if req.CompletedAt != nil {
		var completedAt pgtype.Date
		if err := completedAt.Scan(*req.CompletedAt); err == nil {
			params.CompletedAt = completedAt
		}
	}

	if req.RewatchCount != nil {
		params.RewatchCount = pgtype.Int4{Int32: *req.RewatchCount, Valid: true}
	}

	if req.Progress != nil {
		params.Progress = pgtype.Int4{Int32: *req.Progress, Valid: true}
	}

	if req.Total != nil {
		params.Total = pgtype.Int4{Int32: *req.Total, Valid: true}
	}

	if req.Note != nil {
		params.Note = pgtype.Text{String: *req.Note, Valid: true}
	}

	if req.ContainsSpoilers != nil {
		params.ContainsSpoilers = pgtype.Bool{Bool: *req.ContainsSpoilers, Valid: true}
	}

	updated, err := h.queries.UpdateLog(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update log",
		})
	}

	return c.JSON(updated)
}

// DeleteLog deletes a log
func (h *LogsHandler) DeleteLog(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	logIDStr := c.Params("id")

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var logID pgtype.UUID
	if err := logID.Scan(logIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid log id",
		})
	}

	// Verify log exists and belongs to user
	existing, err := h.queries.GetLogByID(c.Context(), logID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch log",
		})
	}

	if existing.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "not authorized to delete this log",
		})
	}

	err = h.queries.DeleteLog(c.Context(), logID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete log",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
