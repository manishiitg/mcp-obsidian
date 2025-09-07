package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"mcp-obsidian/obsidian/client"
	"mcp-obsidian/obsidian/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ListFilesInVault lists all files in the vault
func ListFilesInVault(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	files, err := obsidianClient.ListFilesInVault()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list files in vault: %v", err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Files in Obsidian Vault:\n\n")

	if len(files) == 0 {
		fmt.Fprintf(&buf, "No files found in the vault.\n")
	} else {
		for i, file := range files {
			icon := ""
			if file.Type == "directory" {
				icon = ""
			}
			fmt.Fprintf(&buf, "%d. %s %s\n", i+1, icon, file.Path)
			if file.Size > 0 {
				fmt.Fprintf(&buf, "   Size: %d bytes\n", file.Size)
			}
			if !file.ModifiedTime.IsZero() {
				fmt.Fprintf(&buf, "   Modified: %s\n", file.ModifiedTime.Format("2006-01-02 15:04:05"))
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// ListFilesInDir lists files in a specific directory
func ListFilesInDir(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dirPath, err := req.RequireString("dirpath")
	if err != nil {
		return mcp.NewToolResultError("dirpath parameter required"), nil
	}

	// Parse max_depth parameter (default: 3 for recursive listing)
	maxDepth := 3
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if depthStr, ok := args["max_depth"].(string); ok {
			if depth, err := strconv.Atoi(depthStr); err == nil {
				maxDepth = depth
			}
		}
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	files, err := obsidianClient.ListFilesInDir(dirPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list files in directory %s: %v", dirPath, err)), nil
	}

	// Convert to JSON structure
	fileList := make([]map[string]interface{}, 0, len(files))

	for _, file := range files {
		fileInfo := map[string]interface{}{
			"name": file.Name,
			"type": file.Type,
			"path": file.Path,
		}

		// Add file-specific information
		if file.Type == "file" {
			if file.Size > 0 {
				fileInfo["size"] = file.Size
			}
			if !file.ModifiedTime.IsZero() {
				fileInfo["modified"] = file.ModifiedTime.Format("2006-01-02 15:04:05")
			}
		}

		// If it's a directory and max_depth > 1, recursively list contents
		if file.Type == "directory" && maxDepth > 1 {
			// For now, we'll just indicate it's a directory with potential children
			// In a full implementation, we'd recursively call ListFilesInDir
			fileInfo["has_children"] = true
			fileInfo["children_depth"] = maxDepth - 1
		}

		fileList = append(fileList, fileInfo)
	}

	// Create JSON structure
	structureData := map[string]interface{}{
		"directory":   dirPath,
		"max_depth":   maxDepth,
		"total_items": len(files),
		"items":       fileList,
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(structureData, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// GetFileContents gets the contents of a file
func GetFileContents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "File: %s\n\n", filePath)
	fmt.Fprintf(&buf, "%s", content)

	return mcp.NewToolResultText(buf.String()), nil
}

// Search performs a simple text search
func Search(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError("query parameter required"), nil
	}

	contextLength := 100
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if contextLengthStr, ok := args["context_length"].(string); ok {
			if parsed, err := strconv.Atoi(contextLengthStr); err == nil {
				contextLength = parsed
			}
		}
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	results, err := obsidianClient.Search(query, contextLength)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search for '%s': %v", query, err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Search results for: '%s'\n\n", query)

	if len(results) == 0 {
		fmt.Fprintf(&buf, "No results found.\n")
	} else {
		for i, result := range results {
			fmt.Fprintf(&buf, "%d. %s (Score: %.2f)\n", i+1, result.Filename, result.Score)
			if result.Matches != nil {
				for _, match := range result.Matches {
					fmt.Fprintf(&buf, "   Context: %s\n", match.Context)
				}
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// AppendContent appends content to a file
func AppendContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	content, err := req.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("content parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	err = obsidianClient.AppendContent(filePath, content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to append content to %s: %v", filePath, err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Successfully appended content to %s\n\n", filePath)
	fmt.Fprintf(&buf, "Content appended:\n%s", content)

	return mcp.NewToolResultText(buf.String()), nil
}

// PutContent creates or updates a file
func PutContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	content, err := req.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("content parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	err = obsidianClient.PutContent(filePath, content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to put content to %s: %v", filePath, err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Successfully created/updated %s", filePath)

	return mcp.NewToolResultText(buf.String()), nil
}

// DeleteFile deletes a file or directory
func DeleteFile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	confirm := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if confirmStr, ok := args["confirm"].(string); ok {
			if parsed, err := strconv.ParseBool(confirmStr); err == nil {
				confirm = parsed
			}
		}
	}

	if !confirm {
		return mcp.NewToolResultError("confirm must be set to true to delete a file"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	err = obsidianClient.DeleteFile(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete %s: %v", filePath, err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "🗑️ Successfully deleted %s", filePath)

	return mcp.NewToolResultText(buf.String()), nil
}

// PatchContent patches content in a file
func PatchContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	operation, err := req.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError("operation parameter required"), nil
	}

	targetType, err := req.RequireString("target_type")
	if err != nil {
		return mcp.NewToolResultError("target_type parameter required"), nil
	}

	target, err := req.RequireString("target")
	if err != nil {
		return mcp.NewToolResultError("target parameter required"), nil
	}

	content, err := req.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("content parameter required"), nil
	}

	// Validate target type
	validTargetTypes := map[string]bool{
		"heading":     true,
		"block":       true,
		"frontmatter": true,
	}
	if !validTargetTypes[targetType] {
		return mcp.NewToolResultError(fmt.Sprintf("invalid target_type: %s. Must be one of: heading, block, frontmatter", targetType)), nil
	}

	// Validate operation
	validOperations := map[string]bool{
		"append":  true,
		"prepend": true,
		"replace": true,
	}
	if !validOperations[operation] {
		return mcp.NewToolResultError(fmt.Sprintf("invalid operation: %s. Must be one of: append, prepend, replace", operation)), nil
	}

	// Handle target formatting based on type
	switch targetType {
	case "heading":
		// Convert nested path format from " -> " to "::" for Obsidian API
		if strings.Contains(target, " -> ") {
			target = strings.ReplaceAll(target, " -> ", "::")
		}
		// Clean up any extra whitespace but preserve the original format
		target = strings.TrimSpace(target)
	case "block":
		// Block references should be clean block IDs
		target = strings.TrimSpace(target)
	case "frontmatter":
		// Frontmatter field names should be clean
		target = strings.TrimSpace(target)
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	// Validate that the target exists before patching
	if targetType == "heading" {
		fileContent, err := obsidianClient.GetFileContents(filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for validation: %v", err)), nil
		}

		elements := parseMarkdownElements(fileContent)
		nestedElements := buildNestedStructure(elements)

		// Check if the target heading exists with exact matching
		targetExists := false
		var exactMatch string
		for _, element := range nestedElements {
			if element.Element.Type == "heading" {
				// Try exact match first
				if element.Element.Title == target {
					targetExists = true
					exactMatch = element.Element.Title
					break
				}
				// Fallback to case-insensitive match
				if strings.EqualFold(element.Element.Title, target) {
					targetExists = true
					exactMatch = element.Element.Title
					break
				}
			}
		}

		if !targetExists {
			// Try to find similar headings
			var similarHeadings []string
			for _, element := range nestedElements {
				if element.Element.Type == "heading" {
					similarHeadings = append(similarHeadings, element.Element.Title)
				}
			}

			return mcp.NewToolResultError(fmt.Sprintf("target heading '%s' not found in file. Available headings: %v", target, similarHeadings)), nil
		}

		// Use the exact match from the file if we found one
		if exactMatch != "" && exactMatch != target {
			target = exactMatch
		}
	}

	err = obsidianClient.PatchContent(filePath, operation, targetType, target, content)
	if err != nil {
		// Provide more detailed error information for debugging
		errorMsg := fmt.Sprintf("failed to patch content in %s: %v\n\nDebug info:\n- File: %s\n- Operation: %s\n- Target Type: %s\n- Target: '%s'\n- Content length: %d chars",
			filePath, err, filePath, operation, targetType, target, len(content))
		return mcp.NewToolResultError(errorMsg), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "✅ Successfully patched content in %s\n\n", filePath)
	fmt.Fprintf(&buf, "📝 Operation: %s\n", operation)
	fmt.Fprintf(&buf, "🎯 Target Type: %s\n", targetType)
	fmt.Fprintf(&buf, "🎯 Target: %s\n", target)
	fmt.Fprintf(&buf, "📄 Content:\n%s", content)

	return mcp.NewToolResultText(buf.String()), nil
}

// SearchJSON performs a complex search using JsonLogic
func SearchJSON(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	queryStr, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError("query parameter required"), nil
	}

	var query map[string]interface{}
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid JSON query: %v", err)), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	results, err := obsidianClient.SearchJSON(query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to perform JSON search: %v", err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "JSON Search Results\n\n")

	if len(results) == 0 {
		fmt.Fprintf(&buf, "No results found.\n")
	} else {
		for i, result := range results {
			fmt.Fprintf(&buf, "%d. %s (Score: %.2f)\n", i+1, result.Filename, result.Score)
			if result.Matches != nil {
				for _, match := range result.Matches {
					fmt.Fprintf(&buf, "   Context: %s\n", match.Context)
				}
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// TestConnection tests the connection to Obsidian
func TestConnection(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	err = obsidianClient.TestConnection()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("connection test failed: %v", err)), nil
	}

	// Get configuration from environment variables for display
	host := "127.0.0.1"
	if envHost := os.Getenv("OBSIDIAN_HOST"); envHost != "" {
		host = envHost
	}

	port := "27124"
	if envPort := os.Getenv("OBSIDIAN_PORT"); envPort != "" {
		port = envPort
	}

	protocol := "https"
	if envProtocol := os.Getenv("OBSIDIAN_PROTOCOL"); envProtocol != "" {
		protocol = envProtocol
	}

	vaultPath := os.Getenv("OBSIDIAN_VAULT_PATH")

	var buf strings.Builder
	fmt.Fprintf(&buf, "Successfully connected to Obsidian!\n\n")
	fmt.Fprintf(&buf, "Configuration:\n")
	fmt.Fprintf(&buf, "  Host: %s\n", host)
	fmt.Fprintf(&buf, "  Port: %s\n", port)
	fmt.Fprintf(&buf, "  Protocol: %s\n", protocol)
	if vaultPath != "" {
		fmt.Fprintf(&buf, "  Vault Path: %s\n", vaultPath)
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// GetPeriodicNote gets or creates a periodic note
func GetPeriodicNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	period, err := req.RequireString("period")
	if err != nil {
		return mcp.NewToolResultError("period parameter required"), nil
	}

	date, err := req.RequireString("date")
	if err != nil {
		return mcp.NewToolResultError("date parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	note, err := obsidianClient.GetPeriodicNote(period, date)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get periodic note: %v", err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Periodic Note: %s (%s)\n\n", period, date)
	fmt.Fprintf(&buf, "Path: %s\n", note.Path)
	fmt.Fprintf(&buf, "Exists: %t\n", note.Exists)
	if note.Content != "" {
		fmt.Fprintf(&buf, "\nContent:\n%s", note.Content)
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// CreatePeriodicNote creates a new periodic note
func CreatePeriodicNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	period, err := req.RequireString("period")
	if err != nil {
		return mcp.NewToolResultError("period parameter required"), nil
	}

	date, err := req.RequireString("date")
	if err != nil {
		return mcp.NewToolResultError("date parameter required"), nil
	}

	content, err := req.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("content parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	note, err := obsidianClient.CreatePeriodicNote(period, date, content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create periodic note: %v", err)), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Created Periodic Note: %s (%s)\n\n", period, date)
	fmt.Fprintf(&buf, "Path: %s\n", note.Path)
	fmt.Fprintf(&buf, "Exists: %t\n", note.Exists)

	return mcp.NewToolResultText(buf.String()), nil
}

// GetRecentChanges gets recent changes in the vault
func GetRecentChanges(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError("Recent changes endpoint is not available in this version of the Obsidian Local REST API. Use the search functionality to find recent files."), nil
}

// GetTags gets all tags in the vault
func GetTags(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError("Tags endpoint is not available in this version of the Obsidian Local REST API. Tags can be accessed through the note JSON response when getting file contents."), nil
}

// GetFrontmatter gets frontmatter from a file
func GetFrontmatter(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filepath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	result, err := obsidianClient.GetFrontmatter(filepath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get frontmatter: %v", err)), nil
	}

	jsonData, err := json.Marshal(result.Data)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal frontmatter: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Frontmatter for %s: %s", filepath, string(jsonData))), nil
}

// SetFrontmatter sets frontmatter for a file
func SetFrontmatter(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filepath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	field, err := req.RequireString("field")
	if err != nil {
		return mcp.NewToolResultError("field parameter required"), nil
	}

	value, err := req.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError("value parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	err = obsidianClient.SetFrontmatter(filepath, field, value)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to set frontmatter: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Frontmatter field '%s' set to '%s' for file %s", field, value, filepath)), nil
}

// GetBlockReference gets a block reference from a file
func GetBlockReference(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError("Block reference endpoint is not available in this version of the Obsidian Local REST API. Use the PATCH endpoint with Target-Type: block to work with block references."), nil
}

// GetHeadings gets all headings from a markdown file
func GetHeadings(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	headings := parseHeadings(content)

	var buf strings.Builder
	fmt.Fprintf(&buf, "Headings in file: %s\n\n", filePath)

	if len(headings) == 0 {
		fmt.Fprintf(&buf, "No headings found in the file.\n")
	} else {
		for i, heading := range headings {
			indent := strings.Repeat("  ", heading.Level-1)
			fmt.Fprintf(&buf, "%d. %s%s (Level %d, Line %d)\n", i+1, indent, heading.Title, heading.Level, heading.Line)
		}
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// GetHeadingContent gets the content under a specific heading
func GetHeadingContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	headingTitle, err := req.RequireString("heading")
	if err != nil {
		return mcp.NewToolResultError("heading parameter required"), nil
	}

	exact := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if exactStr, ok := args["exact"].(string); ok {
			if parsed, err := strconv.ParseBool(exactStr); err == nil {
				exact = parsed
			}
		}
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	headings := parseHeadings(content)

	var targetHeading *types.HeadingInfo
	for _, heading := range headings {
		if exact {
			if heading.Title == headingTitle {
				targetHeading = &heading
				break
			}
		} else {
			if strings.Contains(strings.ToLower(heading.Title), strings.ToLower(headingTitle)) {
				targetHeading = &heading
				break
			}
		}
	}

	if targetHeading == nil {
		var buf strings.Builder
		fmt.Fprintf(&buf, "Heading '%s' not found in file: %s\n\n", headingTitle, filePath)
		fmt.Fprintf(&buf, "Available headings:\n")
		for i, heading := range headings {
			indent := strings.Repeat("  ", heading.Level-1)
			fmt.Fprintf(&buf, "%d. %s%s %s\n", i+1, indent, strings.Repeat("#", heading.Level), heading.Title)
		}
		return mcp.NewToolResultText(buf.String()), nil
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "Content under heading '%s' in file: %s\n\n", targetHeading.Title, filePath)
	fmt.Fprintf(&buf, "%s", targetHeading.Content)

	return mcp.NewToolResultText(buf.String()), nil
}

// parseHeadings parses markdown content and extracts headings with their content
func parseHeadings(content string) []types.HeadingInfo {
	var headings []types.HeadingInfo
	lines := strings.Split(content, "\n")

	var currentHeading *types.HeadingInfo
	var currentContent strings.Builder

	for i, line := range lines {
		lineNum := i + 1
		trimmedLine := strings.TrimSpace(line)

		// Check if this is a heading
		if strings.HasPrefix(trimmedLine, "#") {
			// If we have a current heading, save it
			if currentHeading != nil {
				currentHeading.Content = strings.TrimSpace(currentContent.String())
				headings = append(headings, *currentHeading)
			}

			// Parse the new heading
			level := 0
			title := trimmedLine
			for strings.HasPrefix(title, "#") {
				level++
				title = strings.TrimPrefix(title, "#")
			}
			title = strings.TrimSpace(title)

			// Create new heading
			currentHeading = &types.HeadingInfo{
				Level: level,
				Title: title,
				Line:  lineNum,
			}

			// Reset content builder
			currentContent.Reset()
		} else if currentHeading != nil {
			// Add line to current heading's content
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		}
	}

	// Don't forget the last heading
	if currentHeading != nil {
		currentHeading.Content = strings.TrimSpace(currentContent.String())
		headings = append(headings, *currentHeading)
	}

	return headings
}

// selectElements selects elements based on the given criteria
func selectElements(elements []types.MarkdownElement, selectorType, query string, level int, exact bool) []types.MarkdownElement {
	var selected []types.MarkdownElement

	for _, element := range elements {
		switch selectorType {
		case "heading":
			if element.Type == "heading" {
				if level == 0 || element.Level == level {
					if query == "" || matchesQuery(element.Title, query, exact) {
						selected = append(selected, element)
					}
				}
			}
		case "section":
			if element.Type == "heading" {
				if query == "" || matchesQuery(element.Title, query, exact) {
					selected = append(selected, element)
				}
			}
		case "list":
			if element.Type == "list" {
				if query == "" || matchesQuery(element.Content, query, exact) {
					selected = append(selected, element)
				}
			}
		case "code":
			if element.Type == "code_block" {
				if query == "" || matchesQuery(element.Content, query, exact) {
					selected = append(selected, element)
				}
			}
		case "table":
			if element.Type == "table" {
				if query == "" || matchesQuery(element.Content, query, exact) {
					selected = append(selected, element)
				}
			}
		case "all":
			if query == "" || matchesQuery(element.Title, query, exact) || matchesQuery(element.Content, query, exact) {
				selected = append(selected, element)
			}
		}
	}

	return selected
}

// matchesQuery checks if content matches the query
func matchesQuery(content, query string, exact bool) bool {
	if exact {
		return strings.EqualFold(content, query)
	}
	return strings.Contains(strings.ToLower(content), strings.ToLower(query))
}

// parseMarkdownElements parses markdown content and extracts all elements
func parseMarkdownElements(content string) []types.MarkdownElement {
	var elements []types.MarkdownElement
	lines := strings.Split(content, "\n")

	var currentElement *types.MarkdownElement
	var currentContent strings.Builder
	var inCodeBlock bool

	for i, line := range lines {
		lineNum := i + 1
		trimmedLine := strings.TrimSpace(line)

		// Check for code block start/end
		if strings.HasPrefix(trimmedLine, "```") {
			if inCodeBlock {
				// End of code block
				if currentElement != nil {
					currentElement.Content = strings.TrimSpace(currentContent.String())
					elements = append(elements, *currentElement)
				}
				inCodeBlock = false
				currentElement = nil
				currentContent.Reset()
			} else {
				// Start of code block
				if currentElement != nil {
					currentElement.Content = strings.TrimSpace(currentContent.String())
					elements = append(elements, *currentElement)
				}
				language := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "```"))
				currentElement = &types.MarkdownElement{
					Type:       "code_block",
					Line:       lineNum,
					Attributes: map[string]string{"language": language},
				}
				inCodeBlock = true
				currentContent.Reset()
			}
			continue
		}

		if inCodeBlock {
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
			continue
		}

		// Check for headings
		if strings.HasPrefix(trimmedLine, "#") {
			if currentElement != nil {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
			}

			level := 0
			title := trimmedLine
			for strings.HasPrefix(title, "#") {
				level++
				title = strings.TrimPrefix(title, "#")
			}
			title = strings.TrimSpace(title)

			currentElement = &types.MarkdownElement{
				Type:  "heading",
				Level: level,
				Title: title,
				Line:  lineNum,
			}
			currentContent.Reset()
			continue
		}

		// Check for lists
		if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") || strings.HasPrefix(trimmedLine, "+ ") {
			if currentElement != nil && currentElement.Type != "list" {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
			}

			if currentElement == nil || currentElement.Type != "list" {
				currentElement = &types.MarkdownElement{
					Type: "list",
					Line: lineNum,
				}
				currentContent.Reset()
			}

			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(trimmedLine)
			continue
		}

		// Check for tables
		if strings.Contains(trimmedLine, "|") && strings.Count(trimmedLine, "|") >= 2 {
			if currentElement != nil && currentElement.Type != "table" {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
			}

			if currentElement == nil || currentElement.Type != "table" {
				currentElement = &types.MarkdownElement{
					Type: "table",
					Line: lineNum,
				}
				currentContent.Reset()
			}

			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(trimmedLine)
			continue
		}

		// Check for blockquotes
		if strings.HasPrefix(trimmedLine, ">") {
			if currentElement != nil && currentElement.Type != "blockquote" {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
			}

			if currentElement == nil || currentElement.Type != "blockquote" {
				currentElement = &types.MarkdownElement{
					Type: "blockquote",
					Line: lineNum,
				}
				currentContent.Reset()
			}

			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(trimmedLine)
			continue
		}

		// Check for links and images
		if strings.Contains(trimmedLine, "[") && strings.Contains(trimmedLine, "]") {
			if currentElement != nil {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
			}

			elementType := "link"
			if strings.Contains(trimmedLine, "![") {
				elementType = "image"
			}

			currentElement = &types.MarkdownElement{
				Type:  elementType,
				Line:  lineNum,
				Title: extractLinkTitle(trimmedLine),
			}
			currentContent.Reset()
			continue
		}

		// Regular paragraph content
		if trimmedLine != "" {
			if currentElement == nil || currentElement.Type != "paragraph" {
				if currentElement != nil {
					currentElement.Content = strings.TrimSpace(currentContent.String())
					elements = append(elements, *currentElement)
				}
				currentElement = &types.MarkdownElement{
					Type: "paragraph",
					Line: lineNum,
				}
				currentContent.Reset()
			}

			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		} else {
			// Empty line - end current paragraph
			if currentElement != nil && currentElement.Type == "paragraph" {
				currentElement.Content = strings.TrimSpace(currentContent.String())
				elements = append(elements, *currentElement)
				currentElement = nil
				currentContent.Reset()
			}
		}
	}

	// Don't forget the last element
	if currentElement != nil {
		currentElement.Content = strings.TrimSpace(currentContent.String())
		elements = append(elements, *currentElement)
	}

	return elements
}

// buildNestedStructure builds a nested structure from flat elements
func buildNestedStructure(elements []types.MarkdownElement) []types.NestedElement {
	var nestedElements []types.NestedElement
	var headingStack []*types.NestedElement

	for _, element := range elements {
		nestedElement := types.NestedElement{
			Element: element,
			Level:   element.Level,
			Path:    []string{fmt.Sprintf("%s:%s", element.Type, element.Title)},
		}

		if element.Type == "heading" {
			// This is a heading - find the appropriate parent heading
			for len(headingStack) > 0 && headingStack[len(headingStack)-1].Element.Level >= element.Level {
				headingStack = headingStack[:len(headingStack)-1]
			}

			if len(headingStack) > 0 {
				// Add to parent heading's children
				parent := headingStack[len(headingStack)-1]
				parent.Children = append(parent.Children, nestedElement)
				nestedElement.Path = append(parent.Path, nestedElement.Path[0])
				// Add to stack
				headingStack = append(headingStack, &parent.Children[len(parent.Children)-1])
			} else {
				// Top-level heading
				nestedElements = append(nestedElements, nestedElement)
				// Add to stack
				headingStack = append(headingStack, &nestedElements[len(nestedElements)-1])
			}

		} else {
			// This is content - add to the current heading or top level
			if len(headingStack) > 0 {
				// Add to the most recent heading's children
				currentHeading := headingStack[len(headingStack)-1]
				currentHeading.Children = append(currentHeading.Children, nestedElement)
				nestedElement.Path = append(currentHeading.Path, nestedElement.Path[0])
			} else {
				// No current heading - add to top level
				nestedElements = append(nestedElements, nestedElement)
			}
		}
	}

	return nestedElements
}

// selectNestedElements selects elements based on nested structure
func selectNestedElements(nestedElements []types.NestedElement, selector *types.NestedSelector) []types.NestedElement {
	var selected []types.NestedElement

	for _, element := range nestedElements {
		if matchesNestedPath(element, selector.Path, 0, selector.Depth) {
			if selector.Include {
				selected = append(selected, element)
			}
			if selector.Children {
				selected = append(selected, getChildrenRecursive(element)...)
			}
		} else {
			// Check children
			children := selectNestedElementsRecursive(element.Children, selector)
			selected = append(selected, children...)
		}
	}

	return selected
}

// selectNestedElementsRecursive recursively selects nested elements
func selectNestedElementsRecursive(elements []types.NestedElement, selector *types.NestedSelector) []types.NestedElement {
	var selected []types.NestedElement

	for _, element := range elements {
		if matchesNestedPath(element, selector.Path, 0, selector.Depth) {
			if selector.Include {
				selected = append(selected, element)
			}
			if selector.Children {
				selected = append(selected, getChildrenRecursive(element)...)
			}
		} else {
			// Check children
			children := selectNestedElementsRecursive(element.Children, selector)
			selected = append(selected, children...)
		}
	}

	return selected
}

// matchesNestedPath checks if an element matches the nested path
func matchesNestedPath(element types.NestedElement, path []string, currentIndex, maxDepth int) bool {
	if maxDepth > 0 && len(element.Path) > maxDepth {
		return false
	}

	if currentIndex >= len(path) {
		return true
	}

	if currentIndex >= len(element.Path) {
		return false
	}

	pathPart := path[currentIndex]
	elementPath := element.Path[currentIndex]

	// Parse path part (e.g., "heading:Introduction" or "list")
	parts := strings.SplitN(pathPart, ":", 2)
	selectorType := parts[0]
	selectorQuery := ""
	if len(parts) > 1 {
		selectorQuery = parts[1]
	}

	// Parse element path part
	elementParts := strings.SplitN(elementPath, ":", 2)
	elementType := elementParts[0]
	elementQuery := ""
	if len(elementParts) > 1 {
		elementQuery = elementParts[1]
	}

	// Check type match
	if selectorType != elementType && selectorType != "all" {
		return false
	}

	// Check query match if provided
	if selectorQuery != "" {
		if !strings.Contains(strings.ToLower(elementQuery), strings.ToLower(selectorQuery)) {
			return false
		}
	}

	// Check if this is the last path part
	if currentIndex == len(path)-1 {
		return true
	}

	// Check children
	for _, child := range element.Children {
		if matchesNestedPath(child, path, currentIndex+1, maxDepth) {
			return true
		}
	}

	return false
}

// getChildrenRecursive gets all children recursively
func getChildrenRecursive(element types.NestedElement) []types.NestedElement {
	var children []types.NestedElement
	children = append(children, element.Children...)

	for _, child := range element.Children {
		children = append(children, getChildrenRecursive(child)...)
	}

	return children
}

// DiscoverMarkdownStructure discovers and returns the complete structure of a markdown file
func DiscoverMarkdownStructure(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	// Parse max_depth parameter (default: 3)
	maxDepth := 3
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if depthStr, ok := args["max_depth"].(string); ok {
			if depth, err := strconv.Atoi(depthStr); err == nil {
				maxDepth = depth
			}
		}
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	elements := parseMarkdownElements(content)
	nestedElements := buildNestedStructure(elements)

	// Create simplified JSON structure
	structureData := map[string]interface{}{
		"filepath": filePath,
		"headings": buildSimpleHeadingsJSON(nestedElements, maxDepth),
	}

	// Convert to JSON with proper encoding
	jsonData, err := json.MarshalIndent(structureData, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal JSON: %v", err)), nil
	}

	// Convert to string and replace Unicode escapes with actual characters
	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, "\\u0026", "&")
	jsonString = strings.ReplaceAll(jsonString, "\\u003c", "<")
	jsonString = strings.ReplaceAll(jsonString, "\\u003e", ">")

	return mcp.NewToolResultText(jsonString), nil
}

// printPatchTargets prints patch-friendly target information
func printPatchTargets(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		if element.Element.Type == "heading" {
			// Create patch-friendly target
			target := createPatchTarget(element)
			if target != "" {
				fmt.Fprintf(buf, "%s• %s: `%s`\n", indent, element.Element.Title, target)
			}

			// Add nested targets for children
			if len(element.Children) > 0 {
				printPatchTargets(buf, element.Children, depth+1)
			}
		}
	}
}

// createPatchTarget creates a patch-friendly target string for the given element
//
// API Request Structure for PATCH /vault/{filepath}:
//
// Headers:
// - Content-Type: text/markdown
// - Operation: append|prepend|replace
// - Target-Type: heading|block|frontmatter
// - Target: URL-encoded target string (e.g., "Heading%201::Subheading%201.1")
// - Target-Delimiter: :: (for nested headings)
// - Trim-Target-Whitespace: true
//
// Body: Content to insert
//
// Examples:
// - Single heading: "Troubleshooting"
// - Nested heading: "Troubleshooting::Common Issues"
// - Deep nested: "Quick Summary::Key Points About the Certificate File"
//
// Note: This function returns the exact original heading text to ensure
// compatibility with the Obsidian API patch operations.
func createPatchTarget(element types.NestedElement) string {
	if element.Element.Type != "heading" {
		return ""
	}

	// Build the target path by using the element's path
	var pathParts []string

	// Extract heading titles from the path
	for _, pathItem := range element.Path {
		if strings.HasPrefix(pathItem, "heading:") {
			title := strings.TrimPrefix(pathItem, "heading:")
			if title != "" {
				// Use the original title, not cleaned version
				pathParts = append(pathParts, title)
			}
		}
	}

	// If no path parts found, use the current element's title
	if len(pathParts) == 0 {
		return element.Element.Title
	}

	// Join with the delimiter expected by the API (->)
	if len(pathParts) == 1 {
		return pathParts[0]
	} else if len(pathParts) > 1 {
		return strings.Join(pathParts, " -> ")
	}

	return element.Element.Title
}

// printNestedStructure recursively prints the nested structure with proper indentation
func printNestedStructure(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)
		icon := getElementIcon(element.Element.Type)

		fmt.Fprintf(buf, "%s%s %s", indent, icon, element.Element.Type)

		if element.Element.Title != "" {
			fmt.Fprintf(buf, ": %s", element.Element.Title)
		}

		if element.Element.Level > 1 {
			fmt.Fprintf(buf, " (Level %d)", element.Element.Level)
		}

		fmt.Fprintf(buf, " (Line %d)", element.Element.Line)

		// Add descriptive text based on element type
		description := getElementDescription(element.Element)
		if description != "" {
			fmt.Fprintf(buf, " - %s", description)
		}

		// Add full content (not just preview)
		if element.Element.Content != "" {
			content := getContentPreview(element.Element)
			if content != "" {
				fmt.Fprintf(buf, "\n%s    Content:\n%s```\n%s\n%s```",
					indent,
					indent,
					content,
					indent)
			}
		}

		fmt.Fprintf(buf, "\n")

		// Print children recursively
		if len(element.Children) > 0 {
			printNestedStructure(buf, element.Children, depth+1)
		}
	}
}

// printNestedStructureStructured prints nested structure in LLM-friendly structured format
func printNestedStructureStructured(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		fmt.Fprintf(buf, "%s- type: \"%s\"\n", indent, element.Element.Type)
		fmt.Fprintf(buf, "%s  line: %d\n", indent, element.Element.Line)

		if element.Element.Title != "" {
			fmt.Fprintf(buf, "%s  title: \"%s\"\n", indent, element.Element.Title)
		}

		if element.Element.Level > 1 {
			fmt.Fprintf(buf, "%s  level: %d\n", indent, element.Element.Level)
		}

		// Add descriptive text based on element type
		description := getElementDescription(element.Element)
		if description != "" {
			fmt.Fprintf(buf, "%s  description: \"%s\"\n", indent, description)
		}

		// Add content preview (truncated for structured output)
		if element.Element.Content != "" {
			content := getContentPreview(element.Element)
			if content != "" {
				// Truncate content for structured output
				if len(content) > 200 {
					content = content[:200] + "..."
				}
				fmt.Fprintf(buf, "%s  content_preview: \"%s\"\n", indent, strings.ReplaceAll(content, "\n", "\\n"))
			}
		}

		// Print children recursively
		if len(element.Children) > 0 {
			fmt.Fprintf(buf, "%s  children:\n", indent)
			printNestedStructureStructured(buf, element.Children, depth+1)
		}
	}
}

// buildHeadingTargetsJSON builds JSON structure for heading targets
func buildHeadingTargetsJSON(elements []types.NestedElement, maxDepth int) []map[string]interface{} {
	var headings []map[string]interface{}

	for _, element := range elements {
		if element.Element.Type == "heading" {
			target := createPatchTarget(element)
			if target != "" {
				heading := map[string]interface{}{
					"title":  element.Element.Title,
					"target": target,
					"level":  element.Element.Level,
					"line":   element.Element.Line,
				}

				// Add children if they exist and within depth limit
				if len(element.Children) > 0 && (maxDepth == 0 || element.Level < maxDepth) {
					heading["children"] = buildHeadingTargetsJSON(element.Children, maxDepth)
				}

				headings = append(headings, heading)
			}
		}
	}

	return headings
}

// buildBlockTargetsJSON builds JSON structure for block targets
func buildBlockTargetsJSON(elements []types.NestedElement, maxDepth int) []map[string]interface{} {
	var blocks []map[string]interface{}

	for _, element := range elements {
		if element.Element.Type == "code_block" || element.Element.Type == "table" || element.Element.Type == "paragraph" {
			blockID := extractBlockID(element.Element)
			if blockID != "" {
				content := getContentPreview(element.Element)
				block := map[string]interface{}{
					"type":     element.Element.Type,
					"block_id": blockID,
					"line":     element.Element.Line,
				}

				if content != "" {
					block["content"] = content
				}

				blocks = append(blocks, block)
			}
		}
	}

	return blocks
}

// buildFrontmatterTargetsJSON builds JSON structure for frontmatter targets
func buildFrontmatterTargetsJSON(elements []types.NestedElement, maxDepth int) []map[string]interface{} {
	var frontmatter []map[string]interface{}

	for _, element := range elements {
		if element.Element.Type == "frontmatter" {
			fields := extractFrontmatterFields(element.Element.Content)
			for fieldName, fieldValue := range fields {
				frontmatter = append(frontmatter, map[string]interface{}{
					"field": fieldName,
					"value": fieldValue,
					"line":  element.Element.Line,
				})
			}
		}
	}

	return frontmatter
}

// buildDetailedStructureJSON builds JSON structure for detailed structure
func buildDetailedStructureJSON(elements []types.NestedElement, maxDepth int) []map[string]interface{} {
	var structure []map[string]interface{}

	for _, element := range elements {
		item := map[string]interface{}{
			"type": element.Element.Type,
			"line": element.Element.Line,
		}

		if element.Element.Title != "" {
			item["title"] = element.Element.Title
		}

		if element.Element.Level > 1 {
			item["level"] = element.Element.Level
		}

		description := getElementDescription(element.Element)
		if description != "" {
			item["description"] = description
		}

		if element.Element.Content != "" {
			content := getContentPreview(element.Element)
			if content != "" {
				item["content"] = content
			}
		}

		// Add children if they exist and within depth limit
		if len(element.Children) > 0 && (maxDepth == 0 || element.Level < maxDepth) {
			item["children"] = buildDetailedStructureJSON(element.Children, maxDepth)
		}

		structure = append(structure, item)
	}

	return structure
}

// buildSimpleHeadingsJSON builds a simplified JSON structure for headings
func buildSimpleHeadingsJSON(elements []types.NestedElement, maxDepth int) []string {
	var headings []string

	for _, element := range elements {
		if element.Element.Type == "heading" {
			target := createPatchTarget(element)
			if target != "" {
				headings = append(headings, target)
			}

			// Add children if they exist and within depth limit
			if len(element.Children) > 0 && (maxDepth == 0 || element.Element.Level < maxDepth) {
				children := buildSimpleHeadingsJSON(element.Children, maxDepth)
				headings = append(headings, children...)
			}
		}
	}

	return headings
}

// getContentPreview returns the full content of the element (no longer truncated)
func getContentPreview(element types.MarkdownElement) string {
	if element.Content == "" {
		return ""
	}

	// Return the full content, cleaned up but not truncated
	content := strings.TrimSpace(element.Content)
	if content == "" {
		return ""
	}

	// Clean up the content but keep newlines and formatting
	// Only remove excessive whitespace
	content = strings.ReplaceAll(content, "  ", " ")

	return content
}

// getElementDescription returns a descriptive text for the element type
func getElementDescription(element types.MarkdownElement) string {
	switch element.Type {
	case "heading":
		if element.Level == 1 {
			return "Main document title"
		} else if element.Level == 2 {
			return "Major section"
		} else if element.Level == 3 {
			return "Subsection"
		} else if element.Level == 4 {
			return "Sub-subsection"
		} else if element.Level == 5 {
			return "Minor subsection"
		} else {
			return "Deep subsection"
		}
	case "paragraph":
		return "Text content block"
	case "list":
		return "List of items"
	case "code_block":
		language := element.Attributes["language"]
		if language != "" {
			return fmt.Sprintf("Code block in %s language", language)
		}
		return "Code block"
	case "table":
		return "Data table"
	case "blockquote":
		return "Quoted text"
	case "link":
		return "External or internal link"
	case "image":
		return "Image or media file"
	default:
		return "Content element"
	}
}

// extractHeadings extracts all headings from the elements
func extractHeadings(elements []types.MarkdownElement) []types.MarkdownElement {
	var headings []types.MarkdownElement
	for _, element := range elements {
		if element.Type == "heading" {
			headings = append(headings, element)
		}
	}
	return headings
}

// selectNestedElementsForRead selects elements from nested structure for reading
func selectNestedElementsForRead(nestedElements []types.NestedElement, selectorType, query string, level int, exact bool) []types.NestedElement {
	var selected []types.NestedElement

	for _, element := range nestedElements {
		// Check if this element matches the selector
		if matchesSelector(element, selectorType, query, level, exact) {
			selected = append(selected, element)
		}

		// Recursively check children
		if len(element.Children) > 0 {
			children := selectNestedElementsForRead(element.Children, selectorType, query, level, exact)
			selected = append(selected, children...)
		}
	}

	return selected
}

// matchesSelector checks if an element matches the given selector criteria
func matchesSelector(element types.NestedElement, selectorType, query string, level int, exact bool) bool {
	switch selectorType {
	case "heading":
		if element.Element.Type == "heading" {
			if level == 0 || element.Element.Level == level {
				if query == "" || matchesQuery(element.Element.Title, query, exact) {
					return true
				}
			}
		}
	case "block":
		// Look for elements that can have block references
		if element.Element.Type == "code_block" || element.Element.Type == "table" || element.Element.Type == "paragraph" {
			if query == "" {
				// If no query, return all block-eligible elements
				return true
			}
			// Check if the block ID matches the query
			blockID := extractBlockID(element.Element)
			if blockID != "" && matchesQuery(blockID, query, exact) {
				return true
			}
			// Also check content for matches
			if matchesQuery(element.Element.Content, query, exact) {
				return true
			}
		}
	case "frontmatter":
		if element.Element.Type == "frontmatter" {
			if query == "" {
				// If no query, return all frontmatter elements
				return true
			}
			// Check if any frontmatter field matches the query
			fields := extractFrontmatterFields(element.Element.Content)
			for fieldName, fieldValue := range fields {
				if matchesQuery(fieldName, query, exact) || matchesQuery(fieldValue, query, exact) {
					return true
				}
			}
		}
	}
	return false
}

func ReadMarkdownContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	selectorType, err := req.RequireString("selector_type")
	if err != nil {
		return mcp.NewToolResultError("selector_type parameter required"), nil
	}

	// Validate selector type - only support the three target types
	validSelectorTypes := map[string]bool{
		"heading":     true,
		"block":       true,
		"frontmatter": true,
	}
	if !validSelectorTypes[selectorType] {
		return mcp.NewToolResultError(fmt.Sprintf("invalid selector_type: %s. Must be one of: heading, block, frontmatter", selectorType)), nil
	}

	query := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if q, ok := args["query"].(string); ok {
			query = q
		}
	}

	level := 0
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if l, ok := args["level"].(string); ok {
			if parsed, err := strconv.Atoi(l); err == nil {
				level = parsed
			}
		}
	}

	exact := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if exactStr, ok := args["exact"].(string); ok {
			if parsed, err := strconv.ParseBool(exactStr); err == nil {
				exact = parsed
			}
		}
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	elements := parseMarkdownElements(content)
	nestedElements := buildNestedStructure(elements)
	selectedElements := selectNestedElementsForRead(nestedElements, selectorType, query, level, exact)

	var buf strings.Builder
	fmt.Fprintf(&buf, "Content from %s (Selector: %s", filePath, selectorType)
	if query != "" {
		fmt.Fprintf(&buf, ", Query: %s", query)
	}
	if level > 0 {
		fmt.Fprintf(&buf, ", Level: %d", level)
	}
	fmt.Fprintf(&buf, ")\n\n")

	if len(selectedElements) == 0 {
		fmt.Fprintf(&buf, "No content found matching the selector.\n\n")
		fmt.Fprintf(&buf, "Available elements:\n")
		printNestedStructure(&buf, nestedElements, 0)
	} else {
		fmt.Fprintf(&buf, "Found %d matching element(s):\n\n", len(selectedElements))
		for i, element := range selectedElements {
			icon := getElementIcon(element.Element.Type)
			fmt.Fprintf(&buf, "--- Element %d: %s %s ---\n", i+1, icon, element.Element.Type)
			if element.Element.Title != "" {
				fmt.Fprintf(&buf, "Title: %s\n", element.Element.Title)
			}
			if element.Element.Level > 1 {
				fmt.Fprintf(&buf, "Level: %d\n", element.Element.Level)
			}
			if element.Element.Content != "" {
				fmt.Fprintf(&buf, "Content:\n```markdown\n%s\n```\n", element.Element.Content)
			}

			// If this is a heading, also show its nested content
			if element.Element.Type == "heading" && len(element.Children) > 0 {
				fmt.Fprintf(&buf, "Nested content:\n")
				printNestedStructure(&buf, element.Children, 1)
			}

			fmt.Fprintf(&buf, "\n")
		}
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateStructureSummary generates a summary of the markdown structure
func generateStructureSummary(elements []types.MarkdownElement) string {
	headingCount := 0
	paragraphCount := 0
	listCount := 0
	codeCount := 0
	tableCount := 0

	for _, element := range elements {
		switch element.Type {
		case "heading":
			headingCount++
		case "paragraph":
			paragraphCount++
		case "list":
			listCount++
		case "code_block":
			codeCount++
		case "table":
			tableCount++
		}
	}

	parts := []string{}
	if headingCount > 0 {
		parts = append(parts, fmt.Sprintf("%d heading(s)", headingCount))
	}
	if paragraphCount > 0 {
		parts = append(parts, fmt.Sprintf("%d paragraph(s)", paragraphCount))
	}
	if listCount > 0 {
		parts = append(parts, fmt.Sprintf("%d list(s)", listCount))
	}
	if codeCount > 0 {
		parts = append(parts, fmt.Sprintf("%d code block(s)", codeCount))
	}
	if tableCount > 0 {
		parts = append(parts, fmt.Sprintf("%d table(s)", tableCount))
	}

	if len(parts) == 0 {
		return "Empty file"
	}

	return strings.Join(parts, ", ")
}

// getElementIcon returns an appropriate icon for the element type
func getElementIcon(elementType string) string {
	switch elementType {
	case "heading":
		return ""
	case "paragraph":
		return ""
	case "list":
		return ""
	case "code_block":
		return ""
	case "table":
		return ""
	case "blockquote":
		return ""
	case "link":
		return ""
	case "image":
		return ""
	default:
		return ""
	}
}

// extractLinkTitle extracts the title from a markdown link or image
func extractLinkTitle(line string) string {
	// Simple extraction - look for [title](url) pattern
	if strings.Contains(line, "[") && strings.Contains(line, "]") {
		start := strings.Index(line, "[")
		end := strings.Index(line, "]")
		if start >= 0 && end > start {
			return strings.TrimSpace(line[start+1 : end])
		}
	}
	return ""
}

// getHeadingLevelDescription returns a description for a heading level
func getHeadingLevelDescription(level int) string {
	switch level {
	case 1:
		return "Main document title"
	case 2:
		return "Major section"
	case 3:
		return "Subsection"
	case 4:
		return "Sub-subsection"
	case 5:
		return "Minor subsection"
	default:
		return "Deep subsection"
	}
}

// GetNestedContent gets content from a markdown file using nested path selectors
func GetNestedContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, err := req.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError("filepath parameter required"), nil
	}

	nestedPath, err := req.RequireString("nested_path")
	if err != nil {
		return mcp.NewToolResultError("nested_path parameter required"), nil
	}

	obsidianClient, err := client.NewObsidianClientFromEnv()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create Obsidian client: %v", err)), nil
	}

	content, err := obsidianClient.GetFileContents(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get file contents for %s: %v", filePath, err)), nil
	}

	elements := parseMarkdownElements(content)
	nestedElements := buildNestedStructure(elements)

	// Parse the nested path
	pathParts := strings.Split(nestedPath, " -> ")
	for i, part := range pathParts {
		pathParts[i] = strings.TrimSpace(part)
	}

	// Find the content using the nested path
	foundContent := findNestedContent(nestedElements, pathParts)

	var buf strings.Builder
	fmt.Fprintf(&buf, "Nested Content for: %s\n", nestedPath)
	fmt.Fprintf(&buf, "File: %s\n\n", filePath)

	if foundContent == nil {
		fmt.Fprintf(&buf, "No content found for path: %s\n\n", nestedPath)
		fmt.Fprintf(&buf, "Available paths:\n")
		printAvailablePaths(&buf, nestedElements, 0)
	} else {
		fmt.Fprintf(&buf, "Found content:\n\n")
		if foundContent.Element.Type == "heading" {
			fmt.Fprintf(&buf, "📝 PATCH TARGET: Use this exact target for patch_content operations:\n")
			fmt.Fprintf(&buf, "Target: \"%s\"\n\n", foundContent.Element.Title)
		}
		printNestedContent(&buf, *foundContent, 0)
	}

	return mcp.NewToolResultText(buf.String()), nil
}

// findNestedContent recursively searches for content using a nested path
func findNestedContent(elements []types.NestedElement, pathParts []string) *types.NestedElement {
	if len(pathParts) == 0 {
		return nil
	}

	currentPart := strings.TrimSpace(pathParts[0])
	remainingParts := pathParts[1:]

	// Search through all elements recursively
	for _, element := range elements {
		if element.Element.Type == "heading" {
			title := strings.TrimSpace(element.Element.Title)

			// Try exact match first (case-insensitive, with emojis)
			if strings.EqualFold(title, currentPart) {
				if len(remainingParts) == 0 {
					// Found the target element
					return &element
				} else {
					// Continue searching in children
					if found := findNestedContent(element.Children, remainingParts); found != nil {
						return found
					}
				}
			}

			// Clean both strings for comparison (remove emojis and special chars)
			cleanCurrentPart := removeEmojisAndSpecialChars(currentPart)
			cleanTitle := removeEmojisAndSpecialChars(title)

			// Try exact match with cleaned strings (case-insensitive)
			if strings.EqualFold(cleanTitle, cleanCurrentPart) {
				if len(remainingParts) == 0 {
					// Found the target element
					return &element
				} else {
					// Continue searching in children
					if found := findNestedContent(element.Children, remainingParts); found != nil {
						return found
					}
				}
			} else if strings.Contains(strings.ToLower(cleanTitle), strings.ToLower(cleanCurrentPart)) {
				// Try partial match (case-insensitive)
				if len(remainingParts) == 0 {
					// Found the target element
					return &element
				} else {
					// Continue searching in children
					if found := findNestedContent(element.Children, remainingParts); found != nil {
						return found
					}
				}
			} else if strings.Contains(strings.ToLower(cleanCurrentPart), strings.ToLower(cleanTitle)) {
				// Try reverse partial match (case-insensitive)
				if len(remainingParts) == 0 {
					// Found the target element
					return &element
				} else {
					// Continue searching in children
					if found := findNestedContent(element.Children, remainingParts); found != nil {
						return found
					}
				}
			}
		}

		// Also search in children of this element (for nested structures)
		if len(element.Children) > 0 {
			if found := findNestedContent(element.Children, pathParts); found != nil {
				return found
			}
		}
	}

	return nil
}

// removeEmojisAndSpecialChars removes emojis and special characters from a string
func removeEmojisAndSpecialChars(s string) string {
	var result strings.Builder
	for _, r := range s {
		// Keep only ASCII characters (32-126) and basic whitespace
		if r >= 32 && r <= 126 {
			result.WriteRune(r)
		}
	}
	return strings.ToLower(strings.TrimSpace(result.String()))
}

// printAvailablePaths prints all available paths in the nested structure
func printAvailablePaths(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		if element.Element.Type == "heading" {
			indent := strings.Repeat("  ", depth)
			fmt.Fprintf(buf, "%s%s\n", indent, element.Element.Title)

			if len(element.Children) > 0 {
				printAvailablePaths(buf, element.Children, depth+1)
			}
		}
	}
}

// printNestedContent prints the content of a nested element and its children
func printNestedContent(buf *strings.Builder, element types.NestedElement, depth int) {
	indent := strings.Repeat("  ", depth)

	fmt.Fprintf(buf, "%s%s", indent, element.Element.Type)
	if element.Element.Title != "" {
		fmt.Fprintf(buf, ": %s", element.Element.Title)
	}
	if element.Element.Level > 1 {
		fmt.Fprintf(buf, " (Level %d)", element.Element.Level)
	}
	fmt.Fprintf(buf, " (Line %d)\n", element.Element.Line)

	// Add full content (not just preview)
	if element.Element.Content != "" {
		content := getContentPreview(element.Element)
		if content != "" {
			fmt.Fprintf(buf, "%s  Content:\n%s  ```\n%s\n%s  ```\n", indent, indent, content, indent)
		}
	}

	// Print children
	if len(element.Children) > 0 {
		for _, child := range element.Children {
			printNestedContent(buf, child, depth+1)
		}
	}
}

// printHeadingTargets prints patch-friendly heading targets
func printHeadingTargets(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		if element.Element.Type == "heading" {
			// Create patch-friendly target
			target := createPatchTarget(element)
			if target != "" {
				fmt.Fprintf(buf, "%s• %s: `%s`\n", indent, element.Element.Title, target)
			}

			// Add nested targets for children
			if len(element.Children) > 0 {
				printHeadingTargets(buf, element.Children, depth+1)
			}
		}
	}
}

// printHeadingTargetsStructured prints heading targets in LLM-friendly structured format
func printHeadingTargetsStructured(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		if element.Element.Type == "heading" {
			// Create patch-friendly target
			target := createPatchTarget(element)
			if target != "" {
				fmt.Fprintf(buf, "%s- title: \"%s\"\n", indent, element.Element.Title)
				fmt.Fprintf(buf, "%s  target: \"%s\"\n", indent, target)
				fmt.Fprintf(buf, "%s  level: %d\n", indent, element.Element.Level)
				fmt.Fprintf(buf, "%s  line: %d\n", indent, element.Element.Line)
			}

			// Add nested targets for children
			if len(element.Children) > 0 {
				printHeadingTargetsStructured(buf, element.Children, depth+1)
			}
		}
	}
}

// printBlockTargets prints patch-friendly block targets
func printBlockTargets(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		// Look for elements that can have block references
		if element.Element.Type == "code_block" || element.Element.Type == "table" || element.Element.Type == "paragraph" {
			// Try to extract block ID from the content
			blockID := extractBlockID(element.Element)
			if blockID != "" {
				content := getContentPreview(element.Element)
				fmt.Fprintf(buf, "%s• %s (Block ID: `%s`): %s\n", indent, element.Element.Type, blockID, content)
			}
		}

		// Check children recursively
		if len(element.Children) > 0 {
			printBlockTargets(buf, element.Children, depth+1)
		}
	}
}

// printBlockTargetsStructured prints block targets in LLM-friendly structured format
func printBlockTargetsStructured(buf *strings.Builder, elements []types.NestedElement, depth int) {
	for _, element := range elements {
		indent := strings.Repeat("  ", depth)

		// Look for elements that can have block references
		if element.Element.Type == "code_block" || element.Element.Type == "table" || element.Element.Type == "paragraph" {
			// Try to extract block ID from the content
			blockID := extractBlockID(element.Element)
			if blockID != "" {
				content := getContentPreview(element.Element)
				fmt.Fprintf(buf, "%s- type: \"%s\"\n", indent, element.Element.Type)
				fmt.Fprintf(buf, "%s  block_id: \"%s\"\n", indent, blockID)
				fmt.Fprintf(buf, "%s  line: %d\n", indent, element.Element.Line)
				if content != "" {
					fmt.Fprintf(buf, "%s  content_preview: \"%s\"\n", indent, strings.ReplaceAll(content, "\n", "\\n"))
				}
			}
		}

		// Check children recursively
		if len(element.Children) > 0 {
			printBlockTargetsStructured(buf, element.Children, depth+1)
		}
	}
}

// printFrontmatterTargets prints patch-friendly frontmatter targets
func printFrontmatterTargets(buf *strings.Builder, elements []types.NestedElement, depth int) {
	// Frontmatter is typically at the top level, so we'll look for it in the root elements
	for _, element := range elements {
		if element.Element.Type == "frontmatter" {
			// Extract frontmatter fields
			fields := extractFrontmatterFields(element.Element.Content)
			for fieldName := range fields {
				fmt.Fprintf(buf, "• %s: `%s`\n", fieldName, fieldName)
			}
		}
	}
}

// printFrontmatterTargetsStructured prints frontmatter targets in LLM-friendly structured format
func printFrontmatterTargetsStructured(buf *strings.Builder, elements []types.NestedElement, depth int) {
	// Frontmatter is typically at the top level, so we'll look for it in the root elements
	for _, element := range elements {
		if element.Element.Type == "frontmatter" {
			// Extract frontmatter fields
			fields := extractFrontmatterFields(element.Element.Content)
			for fieldName, fieldValue := range fields {
				fmt.Fprintf(buf, "- field: \"%s\"\n", fieldName)
				fmt.Fprintf(buf, "  value: \"%s\"\n", fieldValue)
				fmt.Fprintf(buf, "  line: %d\n", element.Element.Line)
			}
		}
	}
}

// extractBlockID tries to extract a block ID from an element
func extractBlockID(element types.MarkdownElement) string {
	// Look for block reference patterns like ^2d9b4a
	lines := strings.Split(element.Content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "^") && len(line) > 1 {
			// Extract the block ID (alphanumeric characters after ^)
			blockID := strings.TrimPrefix(line, "^")
			// Clean up any trailing characters
			if idx := strings.IndexAny(blockID, " \t\n\r"); idx != -1 {
				blockID = blockID[:idx]
			}
			return blockID
		}
	}
	return ""
}

// extractFrontmatterFields extracts frontmatter field names from content
func extractFrontmatterFields(content string) map[string]string {
	fields := make(map[string]string)

	// Simple frontmatter parsing - look for YAML-like key-value pairs
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				fieldName := strings.TrimSpace(parts[0])
				if fieldName != "" {
					fields[fieldName] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return fields
}

// RegisterObsidianPrompts dynamically loads every markdown file under
// obsidian/prompts and registers it as an MCP prompt.
//
// The prompt name is the filename (without extension) with spaces replaced by
// underscores and lower-cased.
func RegisterObsidianPrompts(s *server.MCPServer) error {
	fmt.Fprintf(os.Stderr, "[RegisterObsidianPrompts] Starting prompt registration...\n")
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	searchDirs := []string{
		filepath.Join(exeDir, "obsidian", "prompts"),
		filepath.Join(exeDir, "..", "obsidian", "prompts"),
		filepath.Join(cwd, "obsidian", "prompts"),
	}

	var promptsDir string
	var entries []fs.DirEntry
	for _, dir := range searchDirs {
		if e, err := os.ReadDir(dir); err == nil {
			promptsDir = dir
			entries = e
			break
		}
	}

	if promptsDir == "" {
		return fmt.Errorf("failed to locate prompts directory in any of: %v", searchDirs)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(promptsDir, entry.Name())
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Failed to read prompt file %s: %v\n", filePath, err)
			continue
		}

		content := strings.TrimSpace(string(contentBytes))

		base := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		promptName := strings.ReplaceAll(strings.ToLower(base), " ", "_")

		prompt := mcp.Prompt{
			Name:        promptName,
			Description: fmt.Sprintf("Prompt loaded from %s", entry.Name()),
			Arguments:   []mcp.PromptArgument{},
		}

		s.AddPrompt(prompt, func(c context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return mcp.NewGetPromptResult(
				base,
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(content)),
				},
			), nil
		})

		fmt.Fprintf(os.Stderr, "✅ Registered prompt: %s from %s\n", promptName, entry.Name())
	}

	return nil
}
