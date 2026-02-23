// Package ai provides Claude-powered features for the media tracker.
package ai

// Recommendation represents a suggested media item.
type Recommendation struct {
	Title       string `json:"title"`
	MediaType   string `json:"media_type"`
	Creator     string `json:"creator"`
	Reason      string `json:"reason"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year,omitempty"`
}

// NLSearchResult is a structured filter parsed from natural language.
type NLSearchResult struct {
	Query   string         `json:"query"`
	Filters map[string]any `json:"filters"`
}

// MoodResult holds mood-based discovery suggestions.
type MoodResult struct {
	Mood        string `json:"mood"`
	Interpretation string `json:"interpretation,omitempty"`
	FromCollection []map[string]any `json:"from_collection,omitempty"`
	NewSuggestions []map[string]any `json:"new_suggestions,omitempty"`
}

// DuplicateResult flags potential duplicate items.
type DuplicateResult struct {
	IsDuplicate bool   `json:"is_duplicate"`
	Reason      string `json:"reason,omitempty"`
	MatchTitle  string `json:"match_title,omitempty"`
}
