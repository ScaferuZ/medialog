package api

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/db"
)

type MediaHandler struct {
	queries *db.Queries
}

func NewMediaHandler(queries *db.Queries) *MediaHandler {
	return &MediaHandler{queries: queries}
}

// RegisterRoutes registers media routes
func (h *MediaHandler) RegisterRoutes(router fiber.Router) {
	media := router.Group("/media")

	media.Get("/", h.ListMedia)
	media.Get("/search", h.SearchMedia)
	media.Get("/:id", h.GetMedia)
	media.Get("/:id/reviews", h.GetMediaReviews)
}

// ListMedia lists media with optional filtering
func (h *MediaHandler) ListMedia(c *fiber.Ctx) error {
	// Query params:
	// - type: filter by media type (film, anime, book, manga, game, doujin)
	// - limit: page size (default 20, max 50)
	// - offset: pagination offset (default 0)

	mediaType := c.Query("type")
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

	media, err := h.queries.ListMedia(c.UserContext(), db.ListMediaParams{
		Column1: mediaType,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch media",
		})
	}

	media = ensureSlice(media)

	return c.JSON(fiber.Map{
		"media": media,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// SearchMedia searches media by title
func (h *MediaHandler) SearchMedia(c *fiber.Ctx) error {
	// Query params:
	// - q: search query (required)
	// - type: filter by media type
	// - limit: page size (default 20, max 50)
	// - offset: pagination offset

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search query is required",
		})
	}

	mediaType := c.Query("type")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 20
	}

	media, err := h.queries.SearchMedia(c.UserContext(), db.SearchMediaParams{
		Query:   query,
		Column2: mediaType,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to search media",
		})
	}

	media = ensureSlice(media)

	return c.JSON(fiber.Map{
		"media": media,
		"query": query,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetMedia gets a single media by ID
func (h *MediaHandler) GetMedia(c *fiber.Ctx) error {
	idStr := c.Params("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid media id",
		})
	}

	media, err := h.queries.GetMediaByID(c.UserContext(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "media not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch media",
		})
	}

	// Get genres for this media
	genres, _ := h.queries.GetMediaGenres(c.UserContext(), media.ID)
	genres = ensureSlice(genres)

	return c.JSON(fiber.Map{
		"media":  media,
		"genres": genres,
	})
}

// GetMediaReviews gets reviews for a media
func (h *MediaHandler) GetMediaReviews(c *fiber.Ctx) error {
	idStr := c.Params("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid media id",
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

	reviews, err := h.queries.ListReviewsByMedia(c.UserContext(), db.ListReviewsByMediaParams{
		MediaID: id,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch reviews",
		})
	}

	reviews = ensureSlice(reviews)
	activityLogs, logsErr := h.queries.ListLogsByMedia(c.UserContext(), db.ListLogsByMediaParams{
		MediaID: id,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if logsErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch reviews",
		})
	}

	reviewRows := make([]fiber.Map, 0, len(reviews))
	existingByLogID := make(map[string]struct{}, len(reviews))
	for _, review := range reviews {
		logID := uuidToString(review.LogID)
		if logID != "" {
			existingByLogID[logID] = struct{}{}
		}
		reviewRows = append(reviewRows, fiber.Map{
			"ID":               uuidToString(review.ID),
			"UserID":           uuidToString(review.UserID),
			"MediaID":          uuidToString(review.MediaID),
			"LogID":            logID,
			"Title":            review.Title,
			"Content":          review.Content,
			"Rating":           review.Rating,
			"ContainsSpoilers": review.ContainsSpoilers,
			"CreatedAt":        review.CreatedAt,
			"UpdatedAt":        review.UpdatedAt,
			"Username":         review.Username,
			"DisplayName":      review.DisplayName,
		})
	}

	for _, log := range activityLogs {
		logID := uuidToString(log.ID)
		if _, exists := existingByLogID[logID]; exists {
			continue
		}
		if strings.TrimSpace(log.Note.String) == "" {
			continue
		}

		reviewRows = append(reviewRows, fiber.Map{
			"ID":               logID,
			"UserID":           uuidToString(log.UserID),
			"MediaID":          uuidToString(log.MediaID),
			"LogID":            logID,
			"Title":            nil,
			"Content":          log.Note.String,
			"Rating":           log.Rating,
			"ContainsSpoilers": log.ContainsSpoilers,
			"CreatedAt":        log.CreatedAt,
			"UpdatedAt":        log.UpdatedAt,
			"Username":         log.Username,
			"DisplayName":      log.DisplayName,
		})
	}

	return c.JSON(fiber.Map{
		"reviews": reviewRows,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}
