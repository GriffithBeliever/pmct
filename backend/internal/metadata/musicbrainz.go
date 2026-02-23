package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MusicBrainzClient fetches music metadata from MusicBrainz (no key, 1 req/s).
type MusicBrainzClient struct {
	httpClient *http.Client
	baseURL    string
	lastReq    time.Time
}

// NewMusicBrainzClient creates a new MusicBrainz client.
func NewMusicBrainzClient() *MusicBrainzClient {
	return &MusicBrainzClient{
		httpClient: &http.Client{},
		baseURL:    "https://musicbrainz.org/ws/2",
	}
}

type mbSearchResponse struct {
	Releases []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Date         string `json:"date"`
		ArtistCredit []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"artist-credit"`
		Genres []struct {
			Name string `json:"name"`
		} `json:"genres"`
	} `json:"releases"`
}

// Search searches for releases by title.
func (c *MusicBrainzClient) Search(ctx context.Context, title string, year *int) ([]*Result, error) {
	if time.Since(c.lastReq) < time.Second {
		time.Sleep(time.Second - time.Since(c.lastReq))
	}
	c.lastReq = time.Now()

	params := url.Values{
		"query": {fmt.Sprintf("release:%s", url.QueryEscape(title))},
		"fmt":   {"json"},
		"limit": {"5"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/release?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "EMS/1.0 (ems@example.com)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz search: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var data mbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode mb response: %w", err)
	}

	results := make([]*Result, 0, len(data.Releases))
	for _, r := range data.Releases {
		yr := 0
		if len(r.Date) >= 4 {
			yr, _ = strconv.Atoi(r.Date[:4]) //nolint:errcheck
		}
		creator := ""
		if len(r.ArtistCredit) > 0 {
			creator = r.ArtistCredit[0].Artist.Name
		}
		genres := make([]string, 0, len(r.Genres))
		for _, g := range r.Genres {
			genres = append(genres, g.Name)
		}
		results = append(results, &Result{
			ExternalID:  r.ID,
			Title:       r.Title,
			Creator:     creator,
			Genres:      genres,
			ReleaseYear: yr,
		})
	}
	return results, nil
}

// GetByID fetches a release by MusicBrainz ID.
func (c *MusicBrainzClient) GetByID(ctx context.Context, id string) (*Result, error) {
	if time.Since(c.lastReq) < time.Second {
		time.Sleep(time.Second - time.Since(c.lastReq))
	}
	c.lastReq = time.Now()

	params := url.Values{"inc": {"artist-credits genres"}, "fmt": {"json"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/release/"+id+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "EMS/1.0 (ems@example.com)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz get: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var release struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Date  string `json:"date"`
		ArtistCredit []struct {
			Artist struct{ Name string `json:"name"` } `json:"artist"`
		} `json:"artist-credit"`
		Genres []struct{ Name string `json:"name"` } `json:"genres"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode mb release: %w", err)
	}

	yr := 0
	if len(release.Date) >= 4 {
		yr, _ = strconv.Atoi(release.Date[:4]) //nolint:errcheck
	}
	creator := ""
	if len(release.ArtistCredit) > 0 {
		creator = release.ArtistCredit[0].Artist.Name
	}
	genres := make([]string, 0, len(release.Genres))
	for _, g := range release.Genres {
		genres = append(genres, g.Name)
	}

	return &Result{
		ExternalID:  id,
		Title:       release.Title,
		Creator:     creator,
		Genres:      genres,
		ReleaseYear: yr,
	}, nil
}
