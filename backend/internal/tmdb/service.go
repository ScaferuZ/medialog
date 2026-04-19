package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medialogg/backend/internal/db"
)

// Service handles syncing TMDB data to the local database.
type Service struct {
	client  *Client
	queries *db.Queries
}

func NewService(client *Client, queries *db.Queries) *Service {
	return &Service{
		client:  client,
		queries: queries,
	}
}

// SearchAndSync searches TMDB and optionally syncs results to the local database.
func (s *Service) SearchAndSync(ctx context.Context, query string, sync bool) ([]db.Medium, error) {
	result, err := s.client.SearchMoviesWithContext(ctx, query, 1)
	if err != nil {
		return nil, err
	}

	media := make([]db.Medium, 0, len(result.Results))
	for _, movie := range result.Results {
		if !sync {
			media = append(media, s.tmdbMovieToMedium(movie))
			continue
		}

			syncedMedia, err := s.syncMovie(ctx, movie)
			if err != nil {
				return nil, err
			}

			media = append(media, *syncedMedia)
	}

	return media, nil
}

// GetOrFetchMovie gets a movie from the local database or fetches it from TMDB.
func (s *Service) GetOrFetchMovie(ctx context.Context, tmdbID int64) (*db.Medium, error) {
	media, err := s.queries.GetMediaByTMDBID(ctx, tmdbIDText(tmdbID))
	if err == nil {
		return &media, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	details, err := s.client.GetMovieDetailsWithContext(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return s.syncMovieDetails(ctx, details)
}

func (s *Service) syncMovie(ctx context.Context, movie Movie) (*db.Medium, error) {
	media, err := s.queries.GetMediaByTMDBID(ctx, tmdbIDText(movie.ID))
	if err == nil {
		return &media, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	details, err := s.client.GetMovieDetailsWithContext(ctx, movie.ID)
	if err != nil {
		return nil, err
	}

	return s.createMovieDetails(ctx, details)
}

func (s *Service) syncMovieDetails(ctx context.Context, details *MovieDetails) (*db.Medium, error) {
	if details == nil {
		return nil, fmt.Errorf("movie details are required")
	}

	media, err := s.queries.GetMediaByTMDBID(ctx, tmdbIDText(details.ID))
	if err == nil {
		return &media, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return s.createMovieDetails(ctx, details)
}

func (s *Service) createMovieDetails(ctx context.Context, details *MovieDetails) (*db.Medium, error) {
	if details == nil {
		return nil, fmt.Errorf("movie details are required")
	}

	metadataJSON, err := buildMovieDetailsMetadata(details)
	if err != nil {
		return nil, err
	}

	media, err := s.queries.CreateMedia(ctx, db.CreateMediaParams{
		Type:          "film",
		Title:         details.Title,
		OriginalTitle: textFromString(details.OriginalTitle),
		Description:   textFromString(details.Overview),
		CoverImage:    textFromString(GetPosterURL(details.PosterPath)),
		ReleaseDate:   dateFromString(details.ReleaseDate),
		Metadata:      metadataJSON,
		TmdbID:        tmdbIDText(details.ID),
	})
	if err != nil {
		return nil, err
	}

	if err := s.syncGenres(ctx, details.Genres, media.ID); err != nil {
		return nil, err
	}

	return &media, nil
}

func (s *Service) syncGenres(ctx context.Context, genres []Genre, mediaID pgtype.UUID) error {
	for _, genre := range genres {
		if err := s.syncGenre(ctx, genre, mediaID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) syncGenre(ctx context.Context, genre Genre, mediaID pgtype.UUID) error {
	name := strings.TrimSpace(genre.Name)
	if name == "" {
		return nil
	}

	dbGenre, err := s.queries.GetGenreByName(ctx, name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		dbGenre, err = s.queries.CreateGenre(ctx, name)
		if err != nil {
			if !isUniqueViolation(err) {
				return err
			}

			dbGenre, err = s.queries.GetGenreByName(ctx, name)
			if err != nil {
				return err
			}
		}
	}

	err = s.queries.AddMediaGenre(ctx, db.AddMediaGenreParams{
		MediaID: mediaID,
		GenreID: dbGenre.ID,
	})
	if err != nil && isUniqueViolation(err) {
		return nil
	}

	return err
}

func (s *Service) tmdbMovieToMedium(movie Movie) db.Medium {
	metadataJSON, _ := buildMovieMetadata(movie)

	return db.Medium{
		Type:          "film",
		Title:         movie.Title,
		OriginalTitle: textFromString(movie.OriginalTitle),
		Description:   textFromString(movie.Overview),
		CoverImage:    textFromString(GetPosterURL(movie.PosterPath)),
		ReleaseDate:   dateFromString(movie.ReleaseDate),
		Metadata:      metadataJSON,
		TmdbID:        tmdbIDText(movie.ID),
	}
}

func buildMovieMetadata(movie Movie) ([]byte, error) {
	return json.Marshal(map[string]any{
		"vote_average": movie.VoteAverage,
		"vote_count":   movie.VoteCount,
		"popularity":   movie.Popularity,
	})
}

func buildMovieDetailsMetadata(details *MovieDetails) ([]byte, error) {
	return json.Marshal(map[string]any{
		"runtime":      details.Runtime,
		"status":       details.Status,
		"tagline":      details.Tagline,
		"budget":       details.Budget,
		"revenue":      details.Revenue,
		"imdb_id":      details.ImdbID,
		"homepage":     details.Homepage,
		"vote_average": details.VoteAverage,
		"vote_count":   details.VoteCount,
		"popularity":   details.Popularity,
	})
}

func tmdbIDText(tmdbID int64) pgtype.Text {
	return textFromString(fmt.Sprintf("%d", tmdbID))
}

func textFromString(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{String: value, Valid: true}
}

func dateFromString(value string) pgtype.Date {
	if value == "" {
		return pgtype.Date{}
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}
	}

	return pgtype.Date{Time: parsed, Valid: true}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}
