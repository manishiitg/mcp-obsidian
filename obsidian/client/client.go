package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"mcp-obsidian/obsidian/types"
)

// ObsidianClient represents a client for the Obsidian REST API
type ObsidianClient struct {
	config     *types.ObsidianConfig
	httpClient *http.Client
	baseURL    string
}

// NewObsidianClient creates a new Obsidian client with the given configuration.
// The client is configured to handle self-signed certificates for local development.
// For production use, consider loading the certificate from /obsidian-local-rest-api.crt
// and configuring proper TLS settings.
func NewObsidianClient(config *types.ObsidianConfig) *ObsidianClient {
	timeout := time.Duration(config.Timeout) * time.Second

	// Create transport with SSL certificate verification disabled for local development
	// In production, you should load the certificate from the API endpoint:
	// GET https://127.0.0.1:27124/obsidian-local-rest-api.crt
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip SSL certificate verification for local development
		},
	}

	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	protocol := "https"
	if !config.UseHTTPS {
		protocol = "http"
	}

	baseURL := fmt.Sprintf("%s://%s:%s", protocol, config.Host, config.Port)

	return &ObsidianClient{
		config:     config,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// NewObsidianClientFromEnv creates a new ObsidianClient from environment variables
func NewObsidianClientFromEnv() (*ObsidianClient, error) {
	config := types.NewObsidianConfig()

	// Load configuration from environment variables
	if apiKey := os.Getenv("OBSIDIAN_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	} else {
		return nil, fmt.Errorf("OBSIDIAN_API_KEY environment variable is required")
	}

	if host := os.Getenv("OBSIDIAN_HOST"); host != "" {
		config.Host = host
	}

	if port := os.Getenv("OBSIDIAN_PORT"); port != "" {
		config.Port = port
	}

	if protocol := os.Getenv("OBSIDIAN_PROTOCOL"); protocol != "" {
		config.Protocol = protocol
		config.UseHTTPS = strings.ToLower(protocol) == "https"
	}

	if vaultPath := os.Getenv("OBSIDIAN_VAULT_PATH"); vaultPath != "" {
		config.VaultPath = vaultPath
	}

	if useHTTPS := os.Getenv("OBSIDIAN_USE_HTTPS"); useHTTPS != "" {
		if parsed, err := strconv.ParseBool(useHTTPS); err == nil {
			config.UseHTTPS = parsed
		}
	}

	return NewObsidianClient(config), nil
}

// getHeaders returns the headers needed for API requests
func (c *ObsidianClient) getHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.config.APIKey),
		"Content-Type":  "application/json",
	}
}

// makeRequest makes an HTTP request to the Obsidian API
func (c *ObsidianClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)

		var obsidianError types.ObsidianError
		if err := json.Unmarshal(bodyBytes, &obsidianError); err == nil {
			return nil, fmt.Errorf("API error %d: %s", obsidianError.ErrorCode, obsidianError.Message)
		}

		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// TestConnection tests the connection to the Obsidian API
func (c *ObsidianClient) TestConnection() error {
	resp, err := c.makeRequest("GET", "/vault/", nil)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("connection test failed with status: %d", resp.StatusCode)
	}

	return nil
}

// ListFilesInVault lists all files in the vault
func (c *ObsidianClient) ListFilesInVault() ([]types.FileInfo, error) {
	resp, err := c.makeRequest("GET", "/vault/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Files []string `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert string array to FileInfo array
	var files []types.FileInfo
	for _, filePath := range result.Files {
		fileInfo := types.FileInfo{
			Path: filePath,
			Name: filePath,
		}

		// Determine if it's a directory (ends with /)
		if strings.HasSuffix(filePath, "/") {
			fileInfo.Type = "directory"
			fileInfo.Name = strings.TrimSuffix(filePath, "/")
		} else {
			fileInfo.Type = "file"
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// ListFilesInDir lists files in a specific directory
func (c *ObsidianClient) ListFilesInDir(dirPath string) ([]types.FileInfo, error) {
	endpoint := fmt.Sprintf("/vault/%s/", url.PathEscape(dirPath))
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Files []string `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert string array to FileInfo array
	var files []types.FileInfo
	for _, filePath := range result.Files {
		fileInfo := types.FileInfo{
			Path: filePath,
			Name: filePath,
		}

		// Determine if it's a directory (ends with /)
		if strings.HasSuffix(filePath, "/") {
			fileInfo.Type = "directory"
			fileInfo.Name = strings.TrimSuffix(filePath, "/")
		} else {
			fileInfo.Type = "file"
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// GetFileContents gets the contents of a file
func (c *ObsidianClient) GetFileContents(filePath string) (string, error) {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

// Search performs a simple text search
func (c *ObsidianClient) Search(query string, contextLength int) ([]types.SearchResult, error) {
	endpoint := "/search/simple/"

	params := url.Values{}
	params.Set("query", query)
	params.Set("contextLength", strconv.Itoa(contextLength))

	fullURL := fmt.Sprintf("%s%s?%s", c.baseURL, endpoint, params.Encode())

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var results []types.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	return results, nil
}

// AppendContent appends content to a file
func (c *ObsidianClient) AppendContent(filePath, content string) error {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))

	req, err := http.NewRequest("POST", c.baseURL+endpoint, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "text/markdown"
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("append content failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// PutContent creates or updates a file
func (c *ObsidianClient) PutContent(filePath, content string) error {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))

	req, err := http.NewRequest("PUT", c.baseURL+endpoint, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "text/markdown"
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("put content failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteFile deletes a file or directory
func (c *ObsidianClient) DeleteFile(filePath string) error {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// PatchContent patches content in a file using the correct API format
func (c *ObsidianClient) PatchContent(filePath, operation, targetType, target, content string) error {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))

	req, err := http.NewRequest("PATCH", c.baseURL+endpoint, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "text/markdown"
	headers["Operation"] = operation
	headers["Target-Type"] = targetType

	// URL-encode the target as required by the API
	// The API expects the target to be URL-encoded, especially for headings with special characters
	encodedTarget := url.QueryEscape(target)
	headers["Target"] = encodedTarget

	// Add Target-Delimiter header for nested headings (default is "::")
	if targetType == "heading" && strings.Contains(target, "::") {
		headers["Target-Delimiter"] = "::"
	}

	// Add Trim-Target-Whitespace header to handle any whitespace issues
	headers["Trim-Target-Whitespace"] = "true"

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("patch content failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// SearchJSON performs a complex search using JsonLogic
func (c *ObsidianClient) SearchJSON(query map[string]interface{}) ([]types.SearchResult, error) {
	endpoint := "/search/"

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewReader(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "application/vnd.olrapi.jsonlogic+json"
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("JSON search failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var results []types.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	return results, nil
}

// GetPeriodicNote gets or creates a periodic note
func (c *ObsidianClient) GetPeriodicNote(period, date string) (*types.PeriodicNoteResponse, error) {
	// Parse the date to get year, month, day
	// Expected format: YYYY-MM-DD
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	year := parts[0]
	month := parts[1]
	day := parts[2]

	// Remove leading zeros from month and day
	month = strings.TrimLeft(month, "0")
	day = strings.TrimLeft(day, "0")

	endpoint := fmt.Sprintf("/periodic/%s/%s/%s/%s/", period, year, month, day)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.PeriodicNoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode periodic note response: %w", err)
	}

	return &result, nil
}

// CreatePeriodicNote creates a new periodic note
func (c *ObsidianClient) CreatePeriodicNote(period, date, content string) (*types.PeriodicNoteResponse, error) {
	// Parse the date to get year, month, day
	// Expected format: YYYY-MM-DD
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	year := parts[0]
	month := parts[1]
	day := parts[2]

	// Remove leading zeros from month and day
	month = strings.TrimLeft(month, "0")
	day = strings.TrimLeft(day, "0")

	endpoint := fmt.Sprintf("/periodic/%s/%s/%s/%s/", period, year, month, day)

	req, err := http.NewRequest("POST", c.baseURL+endpoint, strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "text/markdown"
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create periodic note failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result types.PeriodicNoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode periodic note response: %w", err)
	}

	return &result, nil
}

// GetRecentChanges gets recent changes in the vault
// Note: This endpoint is not available in the current version of the API
func (c *ObsidianClient) GetRecentChanges(limit, days int) (*types.RecentChangesResponse, error) {
	return nil, fmt.Errorf("recent changes endpoint is not available in this version of the API")
}

// GetTags gets all tags in the vault
// Note: This endpoint is not available in the current version of the API
func (c *ObsidianClient) GetTags() (*types.TagInfo, error) {
	return nil, fmt.Errorf("tags endpoint is not available in this version of the API")
}

// GetFrontmatter gets frontmatter from a file
func (c *ObsidianClient) GetFrontmatter(filePath string) (*types.FrontmatterResponse, error) {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Accept"] = "application/vnd.olrapi.note+json"
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get frontmatter failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Frontmatter map[string]interface{} `json:"frontmatter"`
		Tags        []string               `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode frontmatter response: %w", err)
	}

	return &types.FrontmatterResponse{
		Path: filePath,
		Data: result.Frontmatter,
	}, nil
}

// SetFrontmatter sets frontmatter for a file
func (c *ObsidianClient) SetFrontmatter(filePath, field, value string) error {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(filePath))

	// JSON encode the value
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	req, err := http.NewRequest("PATCH", c.baseURL+endpoint, bytes.NewReader(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	headers := c.getHeaders()
	headers["Content-Type"] = "application/json"
	headers["Operation"] = "replace"
	headers["Target-Type"] = "frontmatter"
	headers["Target"] = field
	headers["Create-Target-If-Missing"] = "true"

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("set frontmatter failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetBlockReference gets a block reference from a file
// Note: This endpoint is not available in the current version of the API
func (c *ObsidianClient) GetBlockReference(filePath, blockID string) (*types.BlockReferenceResponse, error) {
	return nil, fmt.Errorf("block reference endpoint is not available in this version of the API")
}
