package api

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/medialogg/backend/internal/db"
	"github.com/medialogg/backend/internal/tmdb"
)

type TMDBHandler struct {
	client  *tmdb.Client
	service *tmdb.Service
	queries *db.Queries
}

func NewTMDBHandler(client *tmdb.Client, queries *db.Queries) *TMDBHandler {
	service := tmdb.NewService(client, queries)
	return &TMDBHandler{
		client:  client,
		service: service,
		queries: queries,
	}
}

// RegisterRoutes registers TMDB routes
func (h *TMDBHandler) RegisterRoutes(router fiber.Router) {
	tmdb := router.Group("/tmdb")

	tmdb.Get("/search", h.SearchMovies)
	tmdb.Get("/movie/:id", h.GetMovieDetails)
	tmdb.Get("/popular", h.GetPopularMovies)
	tmdb.Post("/sync/:tmdbId", h.SyncMovie)
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query string `query:"q" validate:"required"`
	Page  int    `query:"page"`
}

// SearchMovies searches TMDB for movies
func (h *TMDBHandler) SearchMovies(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search query is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	result, err := h.client.SearchMovies(query, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to search movies",
		})
	}

	// Convert to response format with full image URLs
	type MovieResponse struct {
		ID            int64   `json:"id"`
		Title         string  `json:"title"`
		OriginalTitle string  `json:"original_title"`
		Overview      string  `json:"overview"`
		PosterURL     string  `json:"poster_url"`
		BackdropURL   string  `json:"backdrop_url"`
		ReleaseDate   string  `json:"release_date"`
		VoteAverage   float64 `json:"vote_average"`
		VoteCount     int     `json:"vote_count"`
	}

	var movies []MovieResponse
	for _, m := range result.Results {
		movies = append(movies, MovieResponse{
			ID:            m.ID,
			Title:         m.Title,
			OriginalTitle: m.OriginalTitle,
			Overview:      m.Overview,
			PosterURL:     tmdb.GetPosterURL(m.PosterPath),
			BackdropURL:   tmdb.GetBackdropURL(m.BackdropPath),
			ReleaseDate:   m.ReleaseDate,
			VoteAverage:   m.VoteAverage,
			VoteCount:     m.VoteCount,
		})
	}

	return c.JSON(fiber.Map{
		"movies":        movies,
		"page":          result.Page,
		"total_pages":   result.TotalPages,
		"total_results": result.TotalResults,
	})
}

// GetMovieDetails gets movie details from TMDB
func (h *TMDBHandler) GetMovieDetails(c *fiber.Ctx) error {
	tmdbID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid movie id",
		})
	}

	details, err := h.client.GetMovieDetails(tmdbID)
	if err != nil {
		if errors.Is(err, tmdb.ErrMovieNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "movie not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch movie details",
		})
	}

	// Convert to response format
	type GenreResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var genres []GenreResponse
	for _, g := range details.Genres {
		genres = append(genres, GenreResponse{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	return c.JSON(fiber.Map{
		"id":             details.ID,
		"title":          details.Title,
		"original_title": details.OriginalTitle,
		"overview":       details.Overview,
		"poster_url":     tmdb.GetPosterURL(details.PosterPath),
		"backdrop_url":   tmdb.GetBackdropURL(details.BackdropPath),
		"release_date":   details.ReleaseDate,
		"runtime":        details.Runtime,
		"status":         details.Status,
		"tagline":        details.Tagline,
		"budget":         details.Budget,
		"revenue":        details.Revenue,
		"imdb_id":        details.ImdbID,
		"homepage":       details.Homepage,
		"vote_average":   details.VoteAverage,
		"vote_count":     details.VoteCount,
		"genres":         genres,
	})
}

// GetPopularMovies gets popular movies from TMDB
func (h *TMDBHandler) GetPopularMovies(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	result, err := h.client.GetPopularMovies(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch popular movies",
		})
	}

	// Convert to response format with full image URLs
	type MovieResponse struct {
		ID          int64   `json:"id"`
		Title       string  `json:"title"`
		Overview    string  `json:"overview"`
		PosterURL   string  `json:"poster_url"`
		ReleaseDate string  `json:"release_date"`
		VoteAverage float64 `json:"vote_average"`
	}

	var movies []MovieResponse
	for _, m := range result.Results {
		movies = append(movies, MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Overview:    m.Overview,
			PosterURL:   tmdb.GetPosterURL(m.PosterPath),
			ReleaseDate: m.ReleaseDate,
			VoteAverage: m.VoteAverage,
		})
	}

	return c.JSON(fiber.Map{
		"movies":        movies,
		"page":          result.Page,
		"total_pages":   result.TotalPages,
		"total_results": result.TotalResults,
	})
}

// SyncRequest represents a sync request
type SyncRequest struct {
	Sync bool `json:"sync"`
}

// SyncMovie syncs a TMDB movie to local database
func (h *TMDBHandler) SyncMovie(c *fiber.Ctx) error {
	tmdbID, err := strconv.ParseInt(c.Params("tmdbId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid tmdb id",
		})
	}

	media, err := h.service.GetOrFetchMovie(c.Context(), tmdbID)
	if err != nil {
		if errors.Is(err, tmdb.ErrMovieNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "movie not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync movie",
		})
	}

	return c.JSON(fiber.Map{
		"message": "movie synced successfully",
		"media":   media,
	})
}
