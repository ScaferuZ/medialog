package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL         = "https://api.themoviedb.org/3"
	imageBaseURL    = "https://image.tmdb.org/t/p"
	defaultLanguage = "en-US"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Movie represents a TMDB movie result.
type Movie struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
	GenreIDs      []int   `json:"genre_ids"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
}

// MovieDetails extends Movie with full details.
type MovieDetails struct {
	Movie
	Genres              []Genre    `json:"genres"`
	Runtime             int        `json:"runtime"`
	Status              string     `json:"status"`
	Tagline             string     `json:"tagline"`
	Budget              int64      `json:"budget"`
	Revenue             int64      `json:"revenue"`
	ImdbID              string     `json:"imdb_id"`
	Homepage            string     `json:"homepage"`
	ProductionCountries []Country  `json:"production_countries"`
	SpokenLanguages     []Language `json:"spoken_languages"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Country struct {
	ISO3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

type Language struct {
	ISO639_1 string `json:"iso_639_1"`
	Name     string `json:"name"`
}

type SearchResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

var ErrMovieNotFound = errors.New("movie not found")

type apiErrorResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

// SearchMovies searches for movies by query.
func (c *Client) SearchMovies(query string, page int) (*SearchResponse, error) {
	return c.SearchMoviesWithContext(context.Background(), query, page)
}

// SearchMoviesWithContext searches for movies by query.
func (c *Client) SearchMoviesWithContext(ctx context.Context, query string, page int) (*SearchResponse, error) {
	if page < 1 {
		page = 1
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("language", defaultLanguage)

	var result SearchResponse
	_, err := c.getJSON(ctx, "/search/movie", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	return &result, nil
}

// GetMovieDetails gets full details for a movie.
func (c *Client) GetMovieDetails(tmdbID int64) (*MovieDetails, error) {
	return c.GetMovieDetailsWithContext(context.Background(), tmdbID)
}

// GetMovieDetailsWithContext gets full details for a movie.
func (c *Client) GetMovieDetailsWithContext(ctx context.Context, tmdbID int64) (*MovieDetails, error) {
	params := url.Values{}
	params.Set("language", defaultLanguage)

	var result MovieDetails
	statusCode, err := c.getJSON(ctx, fmt.Sprintf("/movie/%d", tmdbID), params, &result)
	if err != nil {
		if statusCode == http.StatusNotFound {
			return nil, ErrMovieNotFound
		}

		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}

	return &result, nil
}

// GetMovieByID gets movie by TMDB ID.
func (c *Client) GetMovieByID(tmdbID int64) (*MovieDetails, error) {
	return c.GetMovieDetails(tmdbID)
}

// GetPopularMovies gets popular movies.
func (c *Client) GetPopularMovies(page int) (*SearchResponse, error) {
	return c.GetPopularMoviesWithContext(context.Background(), page)
}

// GetPopularMoviesWithContext gets popular movies.
func (c *Client) GetPopularMoviesWithContext(ctx context.Context, page int) (*SearchResponse, error) {
	if page < 1 {
		page = 1
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("language", defaultLanguage)

	var result SearchResponse
	_, err := c.getJSON(ctx, "/movie/popular", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	return &result, nil
}

// GetImageURL returns a full TMDB image URL.
func GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}

	if size == "" {
		size = "w500"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return fmt.Sprintf("%s/%s%s", imageBaseURL, size, path)
}

// GetPosterURL returns a poster image URL.
func GetPosterURL(path string) string {
	return GetImageURL(path, "w500")
}

// GetBackdropURL returns a backdrop image URL.
func GetBackdropURL(path string) string {
	return GetImageURL(path, "original")
}

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, dest any) (int, error) {
	endpoint := c.buildURL(path, params)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, errors.New("TMDB request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, decodeAPIError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return resp.StatusCode, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp.StatusCode, nil
}

func (c *Client) buildURL(path string, params url.Values) string {
	if params == nil {
		params = url.Values{}
	}

	params = cloneValues(params)
	params.Set("api_key", c.apiKey)

	return fmt.Sprintf("%s%s?%s", baseURL, path, params.Encode())
}

func decodeAPIError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
	}

	var apiErr apiErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.StatusMessage != "" {
		return fmt.Errorf("TMDB API returned status %d: %s", resp.StatusCode, apiErr.StatusMessage)
	}

	return fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
}

func cloneValues(values url.Values) url.Values {
	cloned := url.Values{}
	for key, list := range values {
		copied := make([]string, len(list))
		copy(copied, list)
		cloned[key] = copied
	}

	return cloned
}
