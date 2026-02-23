package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// IGDBClient fetches game metadata from IGDB using OAuth2 client credentials.
type IGDBClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
	mu           sync.Mutex
	accessToken  string
	tokenExpiry  time.Time
}

// NewIGDBClient creates a new IGDB client.
func NewIGDBClient(clientID, clientSecret string) *IGDBClient {
	return &IGDBClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
}

func (c *IGDBClient) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return nil
	}

	params := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://id.twitch.tv/oauth2/token", strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("get igdb token: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decode token: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

func (c *IGDBClient) doQuery(ctx context.Context, endpoint, body string) ([]byte, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.igdb.com/v4/"+endpoint, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("build igdb request: %w", err)
	}
	req.Header.Set("Client-ID", c.clientID)
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("igdb request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read igdb response: %w", err)
	}
	return data, nil
}

type igdbGame struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Summary      string `json:"summary"`
	FirstRelease int64  `json:"first_release_date"`
	Cover        *struct {
		URL string `json:"url"`
	} `json:"cover"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
	InvolvedCompanies []struct {
		Developer bool `json:"developer"`
		Company   struct {
			Name string `json:"name"`
		} `json:"company"`
	} `json:"involved_companies"`
}

func parseGame(g igdbGame, id string) *Result {
	yr := 0
	if g.FirstRelease > 0 {
		yr = time.Unix(g.FirstRelease, 0).Year()
	}
	coverURL := ""
	if g.Cover != nil {
		coverURL = "https:" + strings.Replace(g.Cover.URL, "t_thumb", "t_cover_big", 1)
	}
	genres := make([]string, 0, len(g.Genres))
	for _, genre := range g.Genres {
		genres = append(genres, genre.Name)
	}
	developer := ""
	for _, ic := range g.InvolvedCompanies {
		if ic.Developer {
			developer = ic.Company.Name
			break
		}
	}
	return &Result{
		ExternalID:  id,
		Title:       g.Name,
		Creator:     developer,
		Genres:      genres,
		CoverURL:    coverURL,
		ReleaseYear: yr,
		Overview:    g.Summary,
	}
}

// Search searches for games by title.
func (c *IGDBClient) Search(ctx context.Context, title string, year *int) ([]*Result, error) {
	if c.clientID == "" {
		return nil, fmt.Errorf("IGDB credentials not configured")
	}

	query := fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,cover.url,genres.name,involved_companies.developer,involved_companies.company.name; limit 5;`, title)
	data, err := c.doQuery(ctx, "games", query)
	if err != nil {
		return nil, err
	}

	var games []igdbGame
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, fmt.Errorf("decode igdb search: %w", err)
	}

	results := make([]*Result, 0, len(games))
	for _, g := range games {
		results = append(results, parseGame(g, strconv.Itoa(g.ID)))
	}
	return results, nil
}

// GetByID fetches a game by IGDB ID.
func (c *IGDBClient) GetByID(ctx context.Context, id string) (*Result, error) {
	if c.clientID == "" {
		return nil, fmt.Errorf("IGDB credentials not configured")
	}

	query := fmt.Sprintf(`fields name,summary,first_release_date,cover.url,genres.name,involved_companies.developer,involved_companies.company.name; where id=%s;`, id)
	data, err := c.doQuery(ctx, "games", query)
	if err != nil {
		return nil, err
	}

	var games []igdbGame
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, fmt.Errorf("decode igdb detail: %w", err)
	}
	if len(games) == 0 {
		return nil, fmt.Errorf("game not found")
	}

	return parseGame(games[0], id), nil
}
