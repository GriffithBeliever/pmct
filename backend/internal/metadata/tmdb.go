package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// TMDBClient fetches movie metadata from The Movie Database.
type TMDBClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewTMDBClient creates a new TMDB metadata client.
func NewTMDBClient(apiKey string) *TMDBClient {
	return &TMDBClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    "https://api.themoviedb.org/3",
	}
}

type tmdbSearchResponse struct {
	Results []struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Overview    string  `json:"overview"`
		ReleaseDate string  `json:"release_date"`
		PosterPath  string  `json:"poster_path"`
		VoteAverage float64 `json:"vote_average"`
	} `json:"results"`
}

type tmdbMovieDetail struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
	Genres      []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Credits struct {
		Crew []struct {
			Job  string `json:"job"`
			Name string `json:"name"`
		} `json:"crew"`
	} `json:"credits"`
}

// Search searches for movies by title.
func (c *TMDBClient) Search(ctx context.Context, title string, year *int) ([]*Result, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	params := url.Values{"api_key": {c.apiKey}, "query": {title}}
	if year != nil {
		params.Set("year", strconv.Itoa(*year))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/search/movie?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tmdb search: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var data tmdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode tmdb response: %w", err)
	}

	results := make([]*Result, 0, len(data.Results))
	for _, r := range data.Results {
		yr := 0
		if len(r.ReleaseDate) >= 4 {
			yr, _ = strconv.Atoi(r.ReleaseDate[:4]) //nolint:errcheck
		}
		coverURL := ""
		if r.PosterPath != "" {
			coverURL = "https://image.tmdb.org/t/p/w500" + r.PosterPath
		}
		results = append(results, &Result{
			ExternalID:  strconv.Itoa(r.ID),
			Title:       r.Title,
			CoverURL:    coverURL,
			ReleaseYear: yr,
			Overview:    r.Overview,
		})
	}
	return results, nil
}

// GetByID fetches a movie by TMDB ID with full credits.
func (c *TMDBClient) GetByID(ctx context.Context, id string) (*Result, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	params := url.Values{"api_key": {c.apiKey}, "append_to_response": {"credits"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/movie/"+id+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tmdb get: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var movie tmdbMovieDetail
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("decode tmdb detail: %w", err)
	}

	yr := 0
	if len(movie.ReleaseDate) >= 4 {
		yr, _ = strconv.Atoi(movie.ReleaseDate[:4]) //nolint:errcheck
	}
	coverURL := ""
	if movie.PosterPath != "" {
		coverURL = "https://image.tmdb.org/t/p/w500" + movie.PosterPath
	}
	genres := make([]string, 0, len(movie.Genres))
	for _, g := range movie.Genres {
		genres = append(genres, g.Name)
	}
	creator := ""
	for _, crew := range movie.Credits.Crew {
		if crew.Job == "Director" {
			creator = crew.Name
			break
		}
	}

	return &Result{
		ExternalID:  id,
		Title:       movie.Title,
		Creator:     creator,
		Genres:      genres,
		CoverURL:    coverURL,
		ReleaseYear: yr,
		Overview:    movie.Overview,
	}, nil
}
