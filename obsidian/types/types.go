package types

import "time"

// ObsidianError represents an error response from the Obsidian API
type ObsidianError struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

// FileInfo represents a file in the Obsidian vault
type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Type         string    `json:"type"` // "file" or "directory"
	Size         int64     `json:"size,omitempty"`
	ModifiedTime time.Time `json:"modifiedTime,omitempty"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
}

// DirectoryInfo represents a directory in the Obsidian vault
type DirectoryInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	ModifiedTime time.Time `json:"modifiedTime,omitempty"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
}

// SearchResult represents a search result from Obsidian
type SearchResult struct {
	Filename string        `json:"filename"`
	Score    float64       `json:"score"`
	Matches  []SearchMatch `json:"matches"`
}

// SearchMatch represents a match within a search result
type SearchMatch struct {
	Context       string   `json:"context"`
	MatchPosition MatchPos `json:"match_position"`
}

// MatchPos represents the position of a match in the content
type MatchPos struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// MarkdownElement represents any element in a markdown file
type MarkdownElement struct {
	Type       string            `json:"type"`       // "heading", "paragraph", "list", "code_block", "table", "blockquote", "link", "image"
	Level      int               `json:"level"`      // For headings (1-6), list items (1+)
	Content    string            `json:"content"`    // The actual content
	Title      string            `json:"title"`      // For headings, links, images
	Line       int               `json:"line"`       // Line number where element starts
	Attributes map[string]string `json:"attributes"` // Additional attributes (e.g., language for code blocks)
	Children   []MarkdownElement `json:"children"`   // For nested elements
}

// MarkdownStructure represents the complete structure of a markdown file
type MarkdownStructure struct {
	Filepath string            `json:"filepath"`
	Elements []MarkdownElement `json:"elements"`
	Summary  string            `json:"summary"` // Brief summary of the file structure
}

// ContentSelector represents a way to select specific content
type ContentSelector struct {
	Type    string            `json:"type"`    // "heading", "section", "list", "code", "table", "all", "nested"
	Query   string            `json:"query"`   // Search query or identifier
	Level   int               `json:"level"`   // For headings (1-6)
	Exact   bool              `json:"exact"`   // Exact match or partial
	Filters map[string]string `json:"filters"` // Additional filters
	Nested  *NestedSelector   `json:"nested"`  // For nested structure selection
}

// NestedSelector represents a nested selection path
type NestedSelector struct {
	Path     []string `json:"path"`     // Path of selectors (e.g., ["heading:Introduction", "list", "code"])
	Depth    int      `json:"depth"`    // Maximum depth to search (0 = unlimited)
	Include  bool     `json:"include"`  // Whether to include parent content
	Children bool     `json:"children"` // Whether to include children
}

// NestedElement represents an element with its nested structure
type NestedElement struct {
	Element  MarkdownElement `json:"element"`
	Children []NestedElement `json:"children"`
	Level    int             `json:"level"`
	Path     []string        `json:"path"`
}

// ContentRequest represents a request for specific content
type ContentRequest struct {
	Filepath string          `json:"filepath"`
	Selector ContentSelector `json:"selector"`
	Format   string          `json:"format"` // "markdown", "text", "json"
}

// ContentResponse represents a response with selected content
type ContentResponse struct {
	Filepath string            `json:"filepath"`
	Selector ContentSelector   `json:"selector"`
	Content  []MarkdownElement `json:"content"`
	Summary  string            `json:"summary"`
}

// PatchRequest represents a request to patch content
type PatchRequest struct {
	Operation  string `json:"operation"`  // "append", "prepend", "replace"
	TargetType string `json:"targetType"` // "heading", "block", "frontmatter"
	Target     string `json:"target"`
	Content    string `json:"content"`
}

// PeriodicNoteRequest represents a request for periodic notes
type PeriodicNoteRequest struct {
	Period string `json:"period"` // "daily", "weekly", "monthly", "quarterly", "yearly"
	Type   string `json:"type"`   // "content" or "metadata"
	Date   string `json:"date"`   // Date in YYYY-MM-DD format
}

// PeriodicNoteResponse represents a response for periodic notes
type PeriodicNoteResponse struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Exists  bool   `json:"exists"`
}

// RecentChangesRequest represents a request for recent changes
type RecentChangesRequest struct {
	Limit int `json:"limit"`
	Days  int `json:"days"`
}

// RecentChangesResponse represents a response for recent changes
type RecentChangesResponse struct {
	Changes []ChangeInfo `json:"changes"`
}

// ChangeInfo represents information about a change
type ChangeInfo struct {
	Path         string    `json:"path"`
	Type         string    `json:"type"` // "created", "modified", "deleted"
	ModifiedTime time.Time `json:"modifiedTime"`
	Size         int64     `json:"size,omitempty"`
}

// TagInfo represents tag information
type TagInfo struct {
	Tag   string   `json:"tag"`
	Count int      `json:"count"`
	Files []string `json:"files"`
}

// FrontmatterRequest represents a request for frontmatter operations
type FrontmatterRequest struct {
	Path   string                 `json:"path"`
	Action string                 `json:"action"` // "get", "set", "update"
	Data   map[string]interface{} `json:"data,omitempty"`
}

// FrontmatterResponse represents a response for frontmatter operations
type FrontmatterResponse struct {
	Path string                 `json:"path"`
	Data map[string]interface{} `json:"data"`
}

// BlockReferenceRequest represents a request for block references
type BlockReferenceRequest struct {
	Path string `json:"path"`
	ID   string `json:"id"`
}

// BlockReferenceResponse represents a response for block references
type BlockReferenceResponse struct {
	Path    string `json:"path"`
	ID      string `json:"id"`
	Content string `json:"content"`
	Exists  bool   `json:"exists"`
}

// ObsidianConfig represents the configuration for Obsidian client
type ObsidianConfig struct {
	APIKey    string
	Host      string
	Port      string
	Protocol  string
	VaultPath string
	UseHTTPS  bool
	Timeout   int
	VerifySSL bool
}

// NewObsidianConfig creates a new ObsidianConfig with defaults
func NewObsidianConfig() *ObsidianConfig {
	return &ObsidianConfig{
		Host:      "127.0.0.1",
		Port:      "27124",
		Protocol:  "https",
		UseHTTPS:  true,
		Timeout:   30,
		VerifySSL: false,
	}
}

// HeadingInfo represents a heading in a markdown file
type HeadingInfo struct {
	Level   int    `json:"level"`   // Heading level (1-6)
	Title   string `json:"title"`   // Heading title
	Content string `json:"content"` // Content under this heading
	Line    int    `json:"line"`    // Line number where heading starts
}

// HeadingResponse represents a response with heading information
type HeadingResponse struct {
	Filepath string        `json:"filepath"`
	Headings []HeadingInfo `json:"headings"`
}

// HeadingContentRequest represents a request for heading content
type HeadingContentRequest struct {
	Filepath string `json:"filepath"`
	Heading  string `json:"heading"` // Heading title to search for
	Exact    bool   `json:"exact"`   // Whether to match exact heading title
}
