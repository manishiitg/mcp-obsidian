package cmd

import (
	"fmt"
	"os"
	"time"

	obsidianHandlers "mcp-obsidian/obsidian/handlers"
	"mcp-obsidian/obsidian/logger"
	"mcp-obsidian/obsidian/middleware"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	obsidianUseSSE     bool
	obsidianPort       string
	obsidianSSEPort    string
	obsidianHTTPPort   string
	obsidianEnableBoth bool
	obsidianUseStdio   bool
)

// obsidianMcpCmd represents the obsidian-mcp command
var obsidianMcpCmd = &cobra.Command{
	Use:   "obsidian-mcp",
	Short: "Start the Obsidian MCP server",
	Long:  `Start an MCP server that provides Obsidian tools for managing notes, files, and content.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		// Initialize logging
		logConfig := logger.LoadLogConfigFromEnv()
		if err := logger.InitLogger(logConfig); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		// Log server startup
		logger.LogServerEvent("server_startup", "🚀 Starting Obsidian MCP server...", map[string]interface{}{
			"config": logger.GetLogConfigSummary(),
		})

		fmt.Fprintf(os.Stderr, "🚀 Starting Obsidian MCP server...\n")

		// Check environment
		checkObsidianMCPEnvironment()

		// Log server startup completion
		totalStartupTime := time.Since(startTime)
		logger.LogServerEvent("server_startup_complete", "Obsidian MCP Server startup completed", map[string]interface{}{
			"total_startup_time_ms": totalStartupTime.Milliseconds(),
			"total_startup_time":    totalStartupTime.String(),
			"transport":             getTransportType(),
		})

		// Create a new MCP server
		s := server.NewMCPServer(
			"Obsidian MCP Server 📝",
			"1.0.0",
			server.WithToolCapabilities(true),
			server.WithPromptCapabilities(true),
		)

		// Register Obsidian tools
		logger.LogInfo("Registering Obsidian tools", nil)
		registerObsidianTools(s)

		// Register Obsidian prompts from obsidian/prompts
		fmt.Fprintf(os.Stderr, "📝 Registering Obsidian prompts...\n")
		logger.LogInfo("Registering Obsidian prompts", nil)
		if err := obsidianHandlers.RegisterObsidianPrompts(s); err != nil {
			logger.LogError(err, "Failed to register Obsidian prompts", nil)
			fmt.Fprintf(os.Stderr, "⚠️  Failed to register Obsidian prompts: %v\n", err)
		} else {
			logger.LogInfo("Obsidian prompts registered successfully", nil)
			fmt.Fprintf(os.Stderr, "✅ Obsidian prompts registered successfully!\n")
		}

		if obsidianUseStdio {
			// Start stdio server
			fmt.Fprintf(os.Stderr, "📝 Starting Obsidian MCP Server with stdio transport\n")
			logger.LogServerEvent("server_started", "Obsidian MCP Server started with stdio transport", map[string]interface{}{
				"transport":       "stdio",
				"startup_time_ms": time.Since(startTime).Milliseconds(),
			})
			if err := server.ServeStdio(s); err != nil {
				logger.LogError(err, "Obsidian MCP stdio Server error", map[string]interface{}{
					"transport": "stdio",
				})
				fmt.Fprintf(os.Stderr, "❌ Obsidian MCP stdio Server error: %v\n", err)
				os.Exit(1)
			}
		} else if obsidianEnableBoth {
			// Start both SSE and StreamableHTTP servers simultaneously
			fmt.Fprintf(os.Stderr, "🚀 Starting Obsidian MCP Server with BOTH transports:\n")
			fmt.Fprintf(os.Stderr, "   📡 SSE Server on port %s (endpoint: /sse)\n", obsidianSSEPort)
			fmt.Fprintf(os.Stderr, "   🌐 StreamableHTTP Server on port %s (endpoint: /mcp)\n", obsidianHTTPPort)
			logger.LogServerEvent("server_starting", "Starting Obsidian MCP Server with BOTH transports", map[string]interface{}{
				"transport":       "both",
				"sse_port":        obsidianSSEPort,
				"http_port":       obsidianHTTPPort,
				"startup_time_ms": time.Since(startTime).Milliseconds(),
			})

			// Start SSE server in a goroutine
			go func() {
				sseServer := server.NewSSEServer(s,
					server.WithSSEEndpoint("/sse"),
				)
				fmt.Fprintf(os.Stderr, "📡 Starting SSE server on port %s...\n", obsidianSSEPort)
				logger.LogInfo("Starting SSE server", map[string]interface{}{
					"port": obsidianSSEPort,
				})
				if err := sseServer.Start(":" + obsidianSSEPort); err != nil {
					logger.LogError(err, "Obsidian MCP SSE Server error", map[string]interface{}{
						"transport": "sse",
						"port":      obsidianSSEPort,
					})
					fmt.Fprintf(os.Stderr, "❌ Obsidian MCP SSE Server error: %v\n", err)
				} else {
					logger.LogServerEvent("server_started", "Obsidian MCP SSE Server started successfully", map[string]interface{}{
						"transport": "sse",
						"port":      obsidianSSEPort,
					})
					fmt.Fprintf(os.Stderr, "✅ Obsidian MCP SSE Server started successfully on port %s\n", obsidianSSEPort)
				}
			}()

			// Start StreamableHTTP server in main thread (blocking)
			streamableServer := server.NewStreamableHTTPServer(s,
				server.WithEndpointPath("/mcp"),
			)
			fmt.Fprintf(os.Stderr, "🌐 Starting StreamableHTTP server on port %s...\n", obsidianHTTPPort)
			logger.LogInfo("Starting StreamableHTTP server", map[string]interface{}{
				"port": obsidianHTTPPort,
			})
			if err := streamableServer.Start(":" + obsidianHTTPPort); err != nil {
				logger.LogError(err, "Obsidian MCP StreamableHTTP Server error", map[string]interface{}{
					"transport": "streamable_http",
					"port":      obsidianHTTPPort,
				})
				fmt.Fprintf(os.Stderr, "❌ Obsidian MCP StreamableHTTP Server error: %v\n", err)
				os.Exit(1)
			} else {
				logger.LogServerEvent("server_started", "Obsidian MCP StreamableHTTP Server started successfully", map[string]interface{}{
					"transport": "streamable_http",
					"port":      obsidianHTTPPort,
				})
				fmt.Fprintf(os.Stderr, "✅ Obsidian MCP StreamableHTTP Server started successfully on port %s\n", obsidianHTTPPort)
			}
		} else if obsidianUseSSE {
			// Start SSE server only
			fmt.Fprintf(os.Stderr, "📡 Starting Obsidian MCP Server with SSE on port %s\n", obsidianSSEPort)
			logger.LogServerEvent("server_starting", "Starting Obsidian MCP Server with SSE transport", map[string]interface{}{
				"transport":       "sse",
				"port":            obsidianSSEPort,
				"startup_time_ms": time.Since(startTime).Milliseconds(),
			})

			sseServer := server.NewSSEServer(s,
				server.WithSSEEndpoint("/sse"),
			)

			if err := sseServer.Start(":" + obsidianSSEPort); err != nil {
				logger.LogError(err, "Obsidian MCP SSE Server error", map[string]interface{}{
					"transport": "sse",
					"port":      obsidianSSEPort,
				})
				fmt.Printf("Obsidian MCP SSE Server error: %v\n", err)
				os.Exit(1)
			} else {
				logger.LogServerEvent("server_started", "Obsidian MCP SSE Server started successfully", map[string]interface{}{
					"transport": "sse",
					"port":      obsidianSSEPort,
				})
			}
		} else {
			// Start StreamableHTTP server only
			fmt.Fprintf(os.Stderr, "🌐 Starting Obsidian MCP Server with StreamableHTTP on port %s\n", obsidianHTTPPort)
			logger.LogServerEvent("server_starting", "Starting Obsidian MCP Server with StreamableHTTP transport", map[string]interface{}{
				"transport":       "streamable_http",
				"port":            obsidianHTTPPort,
				"startup_time_ms": time.Since(startTime).Milliseconds(),
			})

			streamableServer := server.NewStreamableHTTPServer(s,
				server.WithEndpointPath("/mcp"),
			)

			if err := streamableServer.Start(":" + obsidianHTTPPort); err != nil {
				logger.LogError(err, "Obsidian MCP StreamableHTTP Server error", map[string]interface{}{
					"transport": "streamable_http",
					"port":      obsidianHTTPPort,
				})
				fmt.Printf("Obsidian MCP StreamableHTTP Server error: %v\n", err)
				os.Exit(1)
			} else {
				logger.LogServerEvent("server_started", "Obsidian MCP StreamableHTTP Server started successfully", map[string]interface{}{
					"transport": "streamable_http",
					"port":      obsidianHTTPPort,
				})
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(obsidianMcpCmd)

	// Flags
	obsidianMcpCmd.Flags().BoolVar(&obsidianUseSSE, "sse", false, "Use SSE transport")
	obsidianMcpCmd.Flags().StringVar(&obsidianPort, "port", "8080", "Port for HTTP transport")
	obsidianMcpCmd.Flags().StringVar(&obsidianSSEPort, "sse-port", "8081", "Port for SSE transport")
	obsidianMcpCmd.Flags().StringVar(&obsidianHTTPPort, "http-port", "8080", "Port for HTTP transport")
	obsidianMcpCmd.Flags().BoolVar(&obsidianEnableBoth, "both", false, "Enable both HTTP and SSE transports")
	obsidianMcpCmd.Flags().BoolVar(&obsidianUseStdio, "stdio", false, "Use stdio transport")
}

// registerObsidianTools registers all Obsidian-related tools
func registerObsidianTools(s *server.MCPServer) {
	fmt.Fprintf(os.Stderr, "🔧 Registering Obsidian tools...\n")
	logger.LogInfo("Registering Obsidian tools", nil)

	// Test connection tool
	testConnectionTool := mcp.NewTool("obsidian_test_connection",
		mcp.WithDescription("Test the connection to Obsidian API"),
	)
	s.AddTool(testConnectionTool, obsidianHandlers.TestConnection)

	// List files in vault tool
	listFilesInVaultTool := mcp.NewTool("obsidian_list_files_in_vault",
		mcp.WithDescription("List all files in the Obsidian vault"),
	)
	s.AddTool(listFilesInVaultTool, middleware.LoggingMiddleware(obsidianHandlers.ListFilesInVault))

	// List files in directory tool
	listFilesInDirTool := mcp.NewTool("obsidian_list_files_in_dir",
		mcp.WithDescription("List files in a specific directory"),
		mcp.WithString("dirpath", mcp.Required(), mcp.Description("Directory path to list files from")),
		mcp.WithString("max_depth", mcp.Description("Maximum depth to explore (default: 3, use 0 for unlimited)")),
	)
	s.AddTool(listFilesInDirTool, middleware.LoggingMiddleware(obsidianHandlers.ListFilesInDir))

	// Get file contents tool
	getFileContentsTool := mcp.NewTool("obsidian_get_file_contents",
		mcp.WithDescription("Get the contents of a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
	)
	s.AddTool(getFileContentsTool, middleware.LoggingMiddleware(obsidianHandlers.GetFileContents))

	// Search tool
	searchTool := mcp.NewTool("obsidian_search",
		mcp.WithDescription("Search for text in the vault"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("context_length", mcp.Description("Length of context around matches (default: 100)")),
	)
	s.AddTool(searchTool, obsidianHandlers.Search)

	// Append content tool
	appendContentTool := mcp.NewTool("obsidian_append_content",
		mcp.WithDescription("Append content to a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to append")),
	)
	s.AddTool(appendContentTool, obsidianHandlers.AppendContent)

	// Put content tool
	putContentTool := mcp.NewTool("obsidian_put_content",
		mcp.WithDescription("Create or update a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to write")),
	)
	s.AddTool(putContentTool, obsidianHandlers.PutContent)

	// Delete file tool
	deleteFileTool := mcp.NewTool("obsidian_delete_file",
		mcp.WithDescription("Delete a file or directory"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file or directory")),
		mcp.WithString("confirm", mcp.Description("Must be 'true' to confirm deletion")),
	)
	s.AddTool(deleteFileTool, obsidianHandlers.DeleteFile)

	// Patch content tool
	patchContentTool := mcp.NewTool("obsidian_patch_content",
		mcp.WithDescription(`Patch content in a file by targeting specific elements with surgical precision.

OVERVIEW:
This tool allows you to modify content in Obsidian files by targeting specific elements like headings, blocks, or frontmatter fields. It's designed for precise content updates without affecting the rest of the document.

TARGET TYPES:

1. HEADING TARGETS:
   - Use exact heading text as it appears in the file
   - Supports nested headings with "::" delimiter (e.g., "Section::Subsection")
   - Case-sensitive - must match exactly
   - Best practice: Use discover_structure first to get exact heading names

2. BLOCK TARGETS:
   - Use block IDs (e.g., "abc123", "def456")
   - Block IDs are unique identifiers for specific content blocks
   - Useful for targeting specific paragraphs or content sections

3. FRONTMATTER TARGETS:
   - Use field names (e.g., "status", "tags", "title", "last_updated")
   - Updates YAML frontmatter at the top of the file
   - Content should be the new field value (e.g., "in_progress" for status field)
   - For array fields like tags, use proper YAML format: "['tag1', 'tag2']"

OPERATIONS:
- append: Add content after the target element
- prepend: Add content before the target element  
- replace: Replace the target element's content entirely

EXAMPLES:
- Update status: target_type="frontmatter", target="status", content="completed"
- Add to heading: target_type="heading", target="Notes", operation="append", content="\n\nAdditional notes here"
- Replace block: target_type="block", target="abc123", operation="replace", content="New content"

TIPS:
- Always use discover_structure first to see available targets
- For frontmatter, use simple string values unless you know the field expects arrays/objects
- Test with small changes first to verify the target works correctly`),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file relative to vault root. Examples: 'Projects/task.md', 'Notes/meeting-notes.md', 'Tasks/AWS_Security_Audit/progress/completed_steps.md'")),
		mcp.WithString("operation", mcp.Required(), mcp.Description("Operation to perform: 'append' (add content after the target element), 'prepend' (add content before the target element), 'replace' (completely replace the target element's content)")),
		mcp.WithString("target_type", mcp.Required(), mcp.Description("Type of target element: 'heading' (target should be exact heading text), 'block' (target should be block ID like 'abc123'), 'frontmatter' (target should be field name like 'status' or 'tags')")),
		mcp.WithString("target", mcp.Required(), mcp.Description("Target identifier: For headings use exact heading text (case-sensitive), for blocks use block ID, for frontmatter use field name. Use discover_structure to find exact target names.")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to patch: For frontmatter use the new field value (e.g., 'completed' for status field), for headings/blocks use markdown content. Include newlines with \\n for proper formatting.")),
	)
	s.AddTool(patchContentTool, obsidianHandlers.PatchContent)

	// Discover markdown structure tool
	discoverStructureTool := mcp.NewTool("obsidian_discover_structure",
		mcp.WithDescription(`Discover and analyze the complete structure of a markdown file with patch-friendly targets.

OVERVIEW:
This tool provides a comprehensive analysis of a markdown file's structure, showing all available targets for patch operations. It's the essential first step before using obsidian_patch_content to ensure you have the correct target names.

WHAT IT SHOWS:

1. HEADING STRUCTURE:
   - All headings with their exact text (case-sensitive)
   - Nested hierarchy with proper indentation
   - Patch-friendly targets for each heading
   - Line numbers for reference

2. BLOCK STRUCTURE:
   - Available block IDs for targeting specific content
   - Block types and their locations
   - Unique identifiers for precise targeting

3. FRONTMATTER FIELDS:
   - All available frontmatter fields
   - Field names ready for patch operations
   - Current field values for reference

4. PATCH TARGETS:
   - Exact target strings to use with obsidian_patch_content
   - Properly formatted for copy-paste usage
   - Case-sensitive and whitespace-preserved

OUTPUT FORMATS:
- Structured JSON with hierarchical data
- Human-readable format with patch targets
- Simple list format for quick reference

BEST PRACTICES:
- Always run this before patch operations
- Use the exact target strings provided
- Check max_depth to control detail level
- Use this instead of reading full file contents

EXAMPLES:
- Discover all structure: max_depth=0 (unlimited)
- Quick overview: max_depth=2 (shallow)
- Detailed analysis: max_depth=5 (deep)

TIPS:
- Copy exact target strings from the output
- Pay attention to case sensitivity
- Use nested paths for complex documents
- Check frontmatter field names carefully`),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the markdown file relative to vault root. Examples: 'Projects/task.md', 'Notes/meeting-notes.md', 'Tasks/AWS_Security_Audit/progress/completed_steps.md'")),
		mcp.WithString("max_depth", mcp.Description("Maximum depth to explore: 0=unlimited, 1=top level only, 2-3=moderate depth, 4+=deep analysis. Default: 3")),
	)
	s.AddTool(discoverStructureTool, obsidianHandlers.DiscoverMarkdownStructure)

	// Get nested content tool
	getNestedContentTool := mcp.NewTool("obsidian_get_nested_content",
		mcp.WithDescription(`Get specific content from a markdown file using hierarchical path selectors.

OVERVIEW:
This tool extracts content from specific sections of a markdown file using nested path navigation. It's perfect for reading targeted content without loading entire files, making it ideal for large documents.

HOW IT WORKS:

1. PATH NAVIGATION:
   - Uses " -> " separator for nested paths
   - Navigates through heading hierarchy
   - Finds content under specific sections
   - Returns content with context

2. CONTENT EXTRACTION:
   - Extracts all content under the target heading
   - Includes subheadings and their content
   - Preserves markdown formatting
   - Shows hierarchical structure

3. PATCH TARGET DISPLAY:
   - Shows exact target string for patch operations
   - Provides copy-paste ready target names
   - Displays case-sensitive heading text
   - Includes nested path information

PATH FORMATS:

1. SIMPLE PATHS:
   - "Introduction" (top-level heading)
   - "Getting Started" (single heading)

2. NESTED PATHS:
   - "Troubleshooting -> Common Issues" (two levels)
   - "Setup -> Installation -> Dependencies" (three levels)
   - "Chapter 1 -> Section A -> Subsection 1" (deep nesting)

3. COMPLEX PATHS:
   - "AWS Security -> IAM Policies -> User Management"
   - "Project Status -> Completed Tasks -> Documentation"

BEST PRACTICES:
- Use discover_structure first to find available paths
- Copy exact path strings from discover_structure output
- Use this instead of reading full file contents
- Check path case sensitivity carefully

EXAMPLES:
- Get introduction: nested_path="Introduction"
- Get troubleshooting: nested_path="Troubleshooting -> Common Issues"
- Get deep content: nested_path="Chapter 1 -> Section A -> Subsection 1"

TIPS:
- Paths are case-sensitive
- Use " -> " (space-arrow-space) as separator
- Check discover_structure for exact path names
- Great for reading specific sections of large documents`),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the markdown file relative to vault root. Examples: 'Projects/task.md', 'Notes/meeting-notes.md', 'Tasks/AWS_Security_Audit/progress/completed_steps.md'")),
		mcp.WithString("nested_path", mcp.Required(), mcp.Description("Nested path using ' -> ' separator. Examples: 'Introduction', 'Troubleshooting -> Common Issues', 'Setup -> Installation -> Dependencies'. Use discover_structure to find exact path names.")),
	)
	s.AddTool(getNestedContentTool, obsidianHandlers.GetNestedContent)

	// Read markdown content tool
	readMarkdownContentTool := mcp.NewTool("obsidian_read_content",
		mcp.WithDescription("Read specific content from a markdown file using various selectors"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the markdown file")),
		mcp.WithString("selector_type", mcp.Required(), mcp.Description("Type of selector: heading, block, frontmatter")),
		mcp.WithString("query", mcp.Description("Search query or identifier")),
		mcp.WithString("level", mcp.Description("Heading level (1-6) for heading selectors")),
		mcp.WithString("exact", mcp.Description("Whether to use exact matching (true/false)")),
	)
	s.AddTool(readMarkdownContentTool, obsidianHandlers.ReadMarkdownContent)

	// JSON search tool
	searchJSONTool := mcp.NewTool("obsidian_search_json",
		mcp.WithDescription("Perform a complex search using JsonLogic"),
		mcp.WithString("query", mcp.Required(), mcp.Description("JSON query string")),
	)
	s.AddTool(searchJSONTool, obsidianHandlers.SearchJSON)

	// Periodic notes tools
	getPeriodicNoteTool := mcp.NewTool("obsidian_get_periodic_note",
		mcp.WithDescription("Get or create a periodic note (daily, weekly, monthly, quarterly, yearly)"),
		mcp.WithString("period", mcp.Required(), mcp.Description("Period type: daily, weekly, monthly, quarterly, yearly")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Date in YYYY-MM-DD format")),
	)
	s.AddTool(getPeriodicNoteTool, obsidianHandlers.GetPeriodicNote)

	createPeriodicNoteTool := mcp.NewTool("obsidian_create_periodic_note",
		mcp.WithDescription("Create a new periodic note"),
		mcp.WithString("period", mcp.Required(), mcp.Description("Period type: daily, weekly, monthly, quarterly, yearly")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Date in YYYY-MM-DD format")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content for the periodic note")),
	)
	s.AddTool(createPeriodicNoteTool, obsidianHandlers.CreatePeriodicNote)

	// Recent changes tool
	getRecentChangesTool := mcp.NewTool("obsidian_get_recent_changes",
		mcp.WithDescription("Get recent changes in the vault"),
		mcp.WithString("limit", mcp.Description("Number of changes to return (default: 10)")),
		mcp.WithString("days", mcp.Description("Number of days to look back (default: 7)")),
	)
	s.AddTool(getRecentChangesTool, obsidianHandlers.GetRecentChanges)

	// Tags tool
	getTagsTool := mcp.NewTool("obsidian_get_tags",
		mcp.WithDescription("Get all tags in the vault"),
	)
	s.AddTool(getTagsTool, obsidianHandlers.GetTags)

	// Frontmatter tools
	getFrontmatterTool := mcp.NewTool("obsidian_get_frontmatter",
		mcp.WithDescription("Get frontmatter from a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
	)
	s.AddTool(getFrontmatterTool, obsidianHandlers.GetFrontmatter)

	setFrontmatterTool := mcp.NewTool("obsidian_set_frontmatter",
		mcp.WithDescription("Set frontmatter for a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithString("data", mcp.Required(), mcp.Description("JSON string of frontmatter data")),
	)
	s.AddTool(setFrontmatterTool, obsidianHandlers.SetFrontmatter)

	// Block reference tool
	getBlockReferenceTool := mcp.NewTool("obsidian_get_block_reference",
		mcp.WithDescription("Get a block reference from a file"),
		mcp.WithString("filepath", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithString("block_id", mcp.Required(), mcp.Description("Block ID to retrieve")),
	)
	s.AddTool(getBlockReferenceTool, obsidianHandlers.GetBlockReference)

	fmt.Fprintf(os.Stderr, "✅ Obsidian tools registered successfully!\n")
	logger.LogInfo("Obsidian tools registered successfully", nil)
}

func checkObsidianMCPEnvironment() {
	requiredEnvVars := []string{"OBSIDIAN_API_KEY"}
	missingVars := []string{}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		logger.LogWarn("Missing required environment variables", map[string]interface{}{
			"missing_vars": missingVars,
		})
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Missing required environment variables: %v\n", missingVars)
		fmt.Fprintf(os.Stderr, "   Please set these environment variables before running the server.\n")
	}

	// Check if Obsidian is accessible
	fmt.Fprintf(os.Stderr, "🔍 Checking Obsidian connection...\n")
	logger.LogInfo("Checking Obsidian connection", nil)
}

// getTransportType returns a string describing the current transport configuration
func getTransportType() string {
	if obsidianUseStdio {
		return "stdio"
	} else if obsidianEnableBoth {
		return "both"
	} else if obsidianUseSSE {
		return "sse"
	} else {
		return "streamable_http"
	}
}
