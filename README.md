# MCP Obsidian Server

A Model Context Protocol (MCP) server that provides tools for interacting with Obsidian notes and files through the Obsidian Local REST API.

## Documentation

- [README.md](README.md) - This file, general project documentation
- [OBSIDIAN_API_DOCUMENTATION.md](OBSIDIAN_API_DOCUMENTATION.md) - Detailed documentation of Obsidian Local REST API system endpoints and certificate management
- [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker deployment guide
- [CURSOR_SETUP.md](CURSOR_SETUP.md) - Cursor integration setup guide

## Features

- 📁 **File Management**: List files in vault and directories
- 📄 **Content Operations**: Get, create, update, and delete file contents
- 🔍 **Search**: Simple text search and complex JSON-based search
- ✏️ **Content Editing**: Append, patch, and modify content
- 🎯 **Markdown Discovery**: Discover and analyze markdown file structure
- 📖 **Content Reading**: Read specific content using various selectors
- 📝 **Heading Operations**: Extract headings and their content
- 📅 **Periodic Notes**: Get and create daily, weekly, monthly, quarterly, and yearly notes
- 🏷️ **Tags**: Get all tags in the vault
- 📊 **Recent Changes**: Track recent changes in the vault
- 📄 **Frontmatter**: Get and set frontmatter for files
- 🔗 **Block References**: Get block references from files
- 🔗 **Connection Testing**: Test connectivity to Obsidian API
- 🚀 **Multiple Transports**: HTTP, SSE, and stdio support
- 🐳 **Docker Support**: Production-ready Docker containers
- 🎯 **Cursor Integration**: Easy integration with Cursor IDE

## Prerequisites

1. **Obsidian Local REST API Plugin**: Install and enable the [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin in your Obsidian vault
2. **API Key**: Generate an API key from the plugin settings
3. **Go 1.24+**: Required to build and run the server

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd mcp-obsidian
```

2. Build the project:
```bash
go build
```

## Configuration

Set the following environment variables:

```bash
export OBSIDIAN_API_KEY="your-api-key-here"
export OBSIDIAN_HOST="127.0.0.1"
export OBSIDIAN_PORT="27124"
export OBSIDIAN_VAULT_PATH="/path/to/your/vault"
export OBSIDIAN_USE_HTTPS="true"
```

### Environment Variables

- `OBSIDIAN_API_KEY` (required): Your Obsidian Local REST API key
- `OBSIDIAN_HOST` (default: "127.0.0.1"): Obsidian API host
- `OBSIDIAN_PORT` (default: "27124"): Obsidian API port
- `OBSIDIAN_VAULT_PATH` (optional): Path to your Obsidian vault
- `OBSIDIAN_USE_HTTPS` (default: "true"): Use HTTPS for API calls
- `OBSIDIAN_PROTOCOL` (optional): Protocol to use (http/https)

## Usage

### Start the Server

#### HTTP Transport (Default)
```bash
./mcp-obsidian obsidian-mcp --http-port 8080
```

#### SSE Transport
```bash
./mcp-obsidian obsidian-mcp --sse --sse-port 8081
```

#### Both HTTP and SSE
```bash
./mcp-obsidian obsidian-mcp --both --http-port 8080 --sse-port 8081
```

#### stdio Transport
```bash
./mcp-obsidian obsidian-mcp --stdio
```

### Cursor Integration

For easy integration with Cursor IDE, use the provided JSON configuration:

1. **Build the application**:
   ```bash
   go build -o mcp-obsidian .
   ```

2. **Get the absolute path** to the executable:
   ```bash
   # On macOS/Linux
   realpath mcp-obsidian
   
   # Or use pwd + filename
   echo "$(pwd)/mcp-obsidian"
   ```

3. **Configure environment variables** in `cursor-mcp-obsidian.json`:
   - Replace `/absolute/path/to/mcp-obsidian` with the actual absolute path from step 2
   - Update the `env` section with your actual values:
   ```json
   {
     "mcpServers": {
       "obsidian": {
         "command": "/Users/yourusername/path/to/mcp-obsidian/mcp-obsidian",
         "args": ["obsidian-mcp", "--stdio"],
         "env": {
           "OBSIDIAN_API_KEY": "your-actual-api-key-here",
           "OBSIDIAN_HOST": "127.0.0.1",
           "OBSIDIAN_PORT": "27124",
           "OBSIDIAN_USE_HTTPS": "true",
           "OBSIDIAN_PROTOCOL": "https"
         }
       }
     }
   }
   ```

4. **Add to Cursor**:
   - Open Cursor
   - Go to Settings > Extensions > MCP Servers
   - Add the `cursor-mcp-obsidian.json` file
   - Restart Cursor

See [CURSOR_SETUP.md](CURSOR_SETUP.md) for detailed setup instructions.

### Available Tools

#### Core Tools
1. **`obsidian_test_connection`**: Test connection to Obsidian API
2. **`obsidian_list_files_in_vault`**: List all files in the vault
3. **`obsidian_list_files_in_dir`**: List files in a specific directory
4. **`obsidian_get_file_contents`**: Get contents of a file
5. **`obsidian_search`**: Search for text in the vault
6. **`obsidian_append_content`**: Append content to a file
7. **`obsidian_put_content`**: Create or update a file
8. **`obsidian_delete_file`**: Delete a file or directory
9. **`obsidian_patch_content`**: Patch content in a file
10. **`obsidian_search_json`**: Perform complex JSON-based search

#### Advanced Tools
11. **`obsidian_get_periodic_note`**: Get or create a periodic note (daily, weekly, monthly, quarterly, yearly)
12. **`obsidian_create_periodic_note`**: Create a new periodic note
13. **`obsidian_get_recent_changes`**: Get recent changes in the vault
14. **`obsidian_get_tags`**: Get all tags in the vault
15. **`obsidian_get_frontmatter`**: Get frontmatter from a file
16. **`obsidian_set_frontmatter`**: Set frontmatter for a file
17. **`obsidian_get_block_reference`**: Get a block reference from a file

## Examples

### Test Connection
```bash
# Test if the server can connect to Obsidian
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_test_connection", "arguments": {}}}'
```

### List Files in Vault
```bash
# List all files in the vault
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_list_files_in_vault", "arguments": {}}}'
```

### Search for Content
```bash
# Search for "todo" in the vault
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_search", "arguments": {"query": "todo"}}}'
```

### Create a New File
```bash
# Create a new note
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_put_content", "arguments": {"filepath": "test.md", "content": "# Test Note\n\nThis is a test note."}}}'
```

### Get Periodic Note
```bash
# Get today's daily note
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_get_periodic_note", "arguments": {"period": "daily", "date": "2025-01-09"}}}'
```

### Get Recent Changes
```bash
# Get recent changes in the last 7 days
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_get_recent_changes", "arguments": {"limit": "10", "days": "7"}}}'
```

### Get Tags
```bash
# Get all tags in the vault
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "obsidian_get_tags", "arguments": {}}}'
```

## Development

### Project Structure

```
mcp-obsidian/
├── cmd/
│   ├── root.go              # Root command
│   └── obsidian-mcp.go      # Obsidian MCP command
├── obsidian/
│   ├── client/
│   │   └── client.go        # HTTP client for Obsidian API
│   ├── handlers/
│   │   └── obsidian.go      # MCP tool handlers
│   └── types/
│       └── types.go         # Data types
├── main.go                  # Entry point
└── README.md               # This file
```

### Building

```bash
go build -o mcp-obsidian
```

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) - The underlying API that makes this possible
- [mcp-go](https://github.com/mark3labs/mcp-go) - The Go MCP library
- [Cobra](https://github.com/spf13/cobra) - CLI framework
