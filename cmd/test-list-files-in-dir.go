package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	obsidianHandlers "mcp-obsidian/obsidian/handlers"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

// testListFilesInDirCmd represents the test-list-files-in-dir command
var testListFilesInDirCmd = &cobra.Command{
	Use:   "test-list-files-in-dir",
	Short: "Test the obsidian_list_files_in_dir tool specifically",
	Long:  `Test the obsidian_list_files_in_dir tool to debug 40400 errors and directory listing issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "🔍 Testing obsidian_list_files_in_dir tool...\n\n")

		// Set environment variables for testing
		setTestEnvironment()

		// Run specific tests
		runListFilesInDirTests()
	},
}

func init() {
	rootCmd.AddCommand(testListFilesInDirCmd)
}

// runListFilesInDirTests runs specific tests for the list files in dir tool
func runListFilesInDirTests() {
	ctx := context.Background()

	fmt.Fprintf(os.Stderr, "🧪 Testing: Tasks/ Directory\n")
	fmt.Fprintf(os.Stderr, "============================\n")
	testListFilesInDir(ctx, "Tasks/", "Tasks/ directory")

	fmt.Fprintf(os.Stderr, "\n🧪 Testing: Tasks/AWS_Security_Audit/ Directory\n")
	fmt.Fprintf(os.Stderr, "===============================================\n")
	testListFilesInDir(ctx, "Tasks/AWS_Security_Audit/", "Tasks/AWS_Security_Audit/ directory")

	fmt.Fprintf(os.Stderr, "\n🎉 Test Completed!\n")
}

// testListFilesInDir tests listing files in a specific directory
func testListFilesInDir(ctx context.Context, dirPath, description string) {
	fmt.Fprintf(os.Stderr, "  📁 Testing: %s\n", description)
	fmt.Fprintf(os.Stderr, "  📍 Directory: '%s'\n", dirPath)

	// Create request
	req := createMockRequestForDirTest(map[string]interface{}{
		"dirpath": dirPath,
	})

	// Call the handler
	result, err := obsidianHandlers.ListFilesInDir(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "    ❌ FAILED: %v\n", err)
	} else if result.IsError {
		// Try to get the error message from the result
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				fmt.Fprintf(os.Stderr, "    ❌ FAILED: %s\n", textContent.Text)
			} else {
				fmt.Fprintf(os.Stderr, "    ❌ FAILED: Unknown error format\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "    ❌ FAILED: No error message available\n")
		}
	} else {
		fmt.Fprintf(os.Stderr, "    ✅ SUCCESS: Directory listing completed\n")
		// Show first few lines of result
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				lines := splitLines(textContent.Text)
				if len(lines) > 0 {
					fmt.Fprintf(os.Stderr, "    📄 Result preview:\n")
					for i, line := range lines {
						if i >= 5 { // Show only first 5 lines
							fmt.Fprintf(os.Stderr, "    ... (truncated)\n")
							break
						}
						fmt.Fprintf(os.Stderr, "      %s\n", line)
					}
				}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
}

// testListFilesInDirWithDepth tests listing files with specific max depth
func testListFilesInDirWithDepth(ctx context.Context, dirPath string, maxDepth int, description string) {
	fmt.Fprintf(os.Stderr, "  📁 Testing: %s\n", description)
	fmt.Fprintf(os.Stderr, "  📍 Directory: '%s', Max Depth: %d\n", dirPath, maxDepth)

	// Create request
	req := createMockRequestForDirTest(map[string]interface{}{
		"dirpath":   dirPath,
		"max_depth": fmt.Sprintf("%d", maxDepth),
	})

	// Call the handler
	result, err := obsidianHandlers.ListFilesInDir(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "    ❌ FAILED: %v\n", err)
	} else if result.IsError {
		// Try to get the error message from the result
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				fmt.Fprintf(os.Stderr, "    ❌ FAILED: %s\n", textContent.Text)
			} else {
				fmt.Fprintf(os.Stderr, "    ❌ FAILED: Unknown error format\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "    ❌ FAILED: No error message available\n")
		}
	} else {
		fmt.Fprintf(os.Stderr, "    ✅ SUCCESS: Directory listing completed\n")
	}
	fmt.Fprintf(os.Stderr, "\n")
}

// createMockRequest creates a mock MCP request for testing
func createMockRequestForDirTest(args map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: args,
		},
	}
}

// splitLines splits a string into lines
func splitLines(s string) []string {
	var lines []string
	var current strings.Builder

	for _, char := range s {
		if char == '\n' {
			lines = append(lines, current.String())
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	return lines
}
