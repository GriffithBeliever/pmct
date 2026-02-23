package ai

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/ems/internal/media"
)

//go:embed prompts/*.txt
var promptsFS embed.FS

// Service provides AI-powered features for the media tracker.
type Service struct {
	client  *Client
	cache   *LRUCache
	prompts map[string]string
}

// NewService creates a new AI Service.
func NewService(client *Client, cache *LRUCache) (*Service, error) {
	s := &Service{
		client:  client,
		cache:   cache,
		prompts: make(map[string]string),
	}
	if err := s.loadPrompts(); err != nil {
		return nil, fmt.Errorf("load prompts: %w", err)
	}
	return s, nil
}

func (s *Service) loadPrompts() error {
	names := []string{"recommendations", "insights", "nl_search", "mood_discovery", "duplicate_detection"}
	for _, name := range names {
		content, err := promptsFS.ReadFile("prompts/" + name + ".txt")
		if err != nil {
			return fmt.Errorf("read prompt %s: %w", name, err)
		}
		s.prompts[name] = string(content)
	}
	return nil
}

func collectionSummary(items []*media.Item) string {
	if len(items) == 0 {
		return "Empty collection."
	}
	var sb strings.Builder
	for _, item := range items {
		fmt.Fprintf(&sb, "- [%s] %s by %s (%s) - Status: %s\n",
			item.MediaType, item.Title, item.Creator,
			strings.Join(item.Genre, ", "), item.Status)
	}
	return sb.String()
}

// Recommend returns AI-generated recommendations based on the user's collection.
func (s *Service) Recommend(ctx context.Context, userID uuid.UUID, items []*media.Item) ([]Recommendation, error) {
	key, err := CollectionKey(userID, items)
	if err != nil {
		return nil, err
	}
	cacheKey := "rec:" + key

	if cached, ok := s.cache.Get(cacheKey); ok {
		if recs, ok := cached.([]Recommendation); ok {
			return recs, nil
		}
	}

	prompt := strings.ReplaceAll(s.prompts["recommendations"],
		"{{COLLECTION_SUMMARY}}", collectionSummary(items))

	result, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("recommendations: %w", err)
	}

	result = extractJSON(result)
	var recs []Recommendation
	if err := json.Unmarshal([]byte(result), &recs); err != nil {
		return nil, fmt.Errorf("parse recommendations: %w", err)
	}

	s.cache.Set(cacheKey, recs)
	return recs, nil
}

// StreamInsights streams AI-generated collection insights via the out channel.
func (s *Service) StreamInsights(ctx context.Context, items []*media.Item, out chan<- string) error {
	prompt := strings.ReplaceAll(s.prompts["insights"],
		"{{COLLECTION_SUMMARY}}", collectionSummary(items))
	return s.client.StreamComplete(ctx, prompt, out)
}

// NLSearch parses a natural language query into structured filters.
func (s *Service) NLSearch(ctx context.Context, query string) (*NLSearchResult, error) {
	prompt := strings.ReplaceAll(s.prompts["nl_search"], "{{QUERY}}", query)

	result, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("nl_search: %w", err)
	}

	result = extractJSON(result)
	var filters map[string]any
	if err := json.Unmarshal([]byte(result), &filters); err != nil {
		return nil, fmt.Errorf("parse nl_search: %w", err)
	}

	return &NLSearchResult{Query: query, Filters: filters}, nil
}

// MoodDiscovery returns items matching a mood description.
func (s *Service) MoodDiscovery(ctx context.Context, mood string, items []*media.Item) (*MoodResult, error) {
	prompt := strings.ReplaceAll(s.prompts["mood_discovery"], "{{MOOD}}", mood)
	prompt = strings.ReplaceAll(prompt, "{{COLLECTION_SUMMARY}}", collectionSummary(items))

	result, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("mood discovery: %w", err)
	}

	result = extractJSON(result)
	var moodResult MoodResult
	if err := json.Unmarshal([]byte(result), &moodResult); err != nil {
		return nil, fmt.Errorf("parse mood result: %w", err)
	}
	moodResult.Mood = mood
	return &moodResult, nil
}

// DetectDuplicates checks if a new item might be a duplicate.
func (s *Service) DetectDuplicates(ctx context.Context, title, mediaType, creator string, existing []*media.Item) (*DuplicateResult, error) {
	var itemList strings.Builder
	for _, item := range existing {
		fmt.Fprintf(&itemList, "- %s by %s (%s)\n", item.Title, item.Creator, item.MediaType)
	}

	prompt := s.prompts["duplicate_detection"]
	prompt = strings.ReplaceAll(prompt, "{{NEW_TITLE}}", title)
	prompt = strings.ReplaceAll(prompt, "{{MEDIA_TYPE}}", mediaType)
	prompt = strings.ReplaceAll(prompt, "{{CREATOR}}", creator)
	prompt = strings.ReplaceAll(prompt, "{{EXISTING_ITEMS}}", itemList.String())

	result, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("duplicate detection: %w", err)
	}

	result = extractJSON(result)
	var dupResult DuplicateResult
	if err := json.Unmarshal([]byte(result), &dupResult); err != nil {
		return nil, fmt.Errorf("parse duplicate result: %w", err)
	}
	return &dupResult, nil
}

// extractJSON extracts a JSON object or array from a string that may have surrounding text.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	for i, c := range s {
		if c == '{' || c == '[' {
			s = s[i:]
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '}' || s[i] == ']' {
			s = s[:i+1]
			break
		}
	}
	return s
}
