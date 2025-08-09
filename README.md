# MCP Obsidian Server

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/mcp-obsidian)](https://goreportcard.com/report/github.com/yourusername/mcp-obsidian)
[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/mcp-obsidian.svg)](https://pkg.go.dev/github.com/yourusername/mcp-obsidian)

A Model Context Protocol (MCP) server that provides comprehensive tools for interacting with Obsidian notes and files through the Obsidian Local REST API. This server enables seamless integration between AI assistants and your Obsidian vault, allowing for advanced note management, content discovery, and automated workflows.

## ✨ Features

- 📁 **File Management**: List files in vault and directories
- 📄 **Content Operations**: Get, create, update, and delete file contents
- 🔍 **Advanced Search**: Simple text search and complex JSON-based search
- ✏️ **Content Editing**: Append, patch, and modify content with precision
- 🎯 **Markdown Discovery**: Discover and analyze markdown file structure
- 📖 **Content Reading**: Read specific content using various selectors
- 📝 **Heading Operations**: Extract headings and their content
- 📅 **Periodic Notes**: Get and create daily, weekly, monthly, quarterly, and yearly notes
- 🏷️ **Tags Management**: Get all tags in the vault
- 📊 **Recent Changes**: Track recent changes in the vault
- 📄 **Frontmatter**: Get and set frontmatter for files
- 🔗 **Block References**: Get block references from files
- 🔗 **Connection Testing**: Test connectivity to Obsidian API
- 🚀 **Multiple Transports**: HTTP, SSE, and stdio support
- 🐳 **Docker Support**: Production-ready Docker containers
- 🎯 **Cursor Integration**: Easy integration with Cursor IDE
- 🔄 **Nested Content**: Advanced nested content retrieval and manipulation

## 🚀 Quick Start

### Prerequisites

1. **Obsidian Local REST API Plugin**: Install and enable the [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin in your Obsidian vault
2. **API Key**: Generate an API key from the plugin settings
3. **Go 1.24+**: Required to build and run the server

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/mcp-obsidian.git
   cd mcp-obsidian
   ```

2. **Build the project**:
   ```bash
   go build -o mcp-obsidian .
   ```

3. **Set environment variables**:
   ```bash
   export OBSIDIAN_API_KEY="your-api-key-here"
   export OBSIDIAN_HOST="127.0.0.1"
   export OBSIDIAN_PORT="27124"
   export OBSIDIAN_USE_HTTPS="true"
   ```

4. **Start the server**:
   ```bash
   # HTTP transport (default)
   ./mcp-obsidian obsidian-mcp --http-port 8080
   
   # Or SSE transport
   ./mcp-obsidian obsidian-mcp --sse --sse-port 8081
   
   # Or both
   ./mcp-obsidian obsidian-mcp --both --http-port 8080 --sse-port 8081
   ```

## 📚 Documentation

- [README.md](README.md) - This file, general project documentation
- [OBSIDIAN_API_DOCUMENTATION.md](OBSIDIAN_API_DOCUMENTATION.md) - Detailed documentation of Obsidian Local REST API system endpoints and certificate management
- [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker deployment guide
- [CURSOR_SETUP.md](CURSOR_SETUP.md) - Cursor integration setup guide

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OBSIDIAN_API_KEY` | ✅ | - | Your Obsidian Local REST API key |
| `OBSIDIAN_HOST` | ❌ | `127.0.0.1` | Obsidian API host |
| `OBSIDIAN_PORT` | ❌ | `27124` | Obsidian API port |
| `OBSIDIAN_VAULT_PATH` | ❌ | - | Path to your Obsidian vault |
| `OBSIDIAN_USE_HTTPS` | ❌ | `true` | Use HTTPS for API calls |
| `OBSIDIAN_PROTOCOL` | ❌ | - | Protocol to use (http/https) |

### Example Configuration

```bash
# Required
export OBSIDIAN_API_KEY="your-api-key-here"

# Optional (with defaults)
export OBSIDIAN_HOST="127.0.0.1"
export OBSIDIAN_PORT="27124"
export OBSIDIAN_VAULT_PATH="/path/to/your/vault"
export OBSIDIAN_USE_HTTPS="true"
export OBSIDIAN_PROTOCOL="https"
```

## 🎯 Usage

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

## 🛠️ Available Tools

### Core Tools

| Tool | Description |
|------|-------------|
| `obsidian_test_connection` | Test connection to Obsidian API |
| `obsidian_list_files_in_vault` | List all files in the vault |
| `obsidian_list_files_in_dir` | List files in a specific directory |
| `obsidian_get_file_contents` | Get contents of a file |
| `obsidian_search` | Search for text in the vault |
| `obsidian_append_content` | Append content to a file |
| `obsidian_put_content` | Create or update a file |
| `obsidian_delete_file` | Delete a file or directory |
| `obsidian_patch_content` | Patch content in a file |
| `obsidian_search_json` | Perform complex JSON-based search |

### Advanced Tools

| Tool | Description |
|------|-------------|
| `obsidian_get_periodic_note` | Get or create a periodic note (daily, weekly, monthly, quarterly, yearly) |
| `obsidian_create_periodic_note` | Create a new periodic note |
| `obsidian_get_recent_changes` | Get recent changes in the vault |
| `obsidian_get_tags` | Get all tags in the vault |
| `obsidian_get_frontmatter` | Get frontmatter from a file |
| `obsidian_set_frontmatter` | Set frontmatter for a file |
| `obsidian_get_block_reference` | Get a block reference from a file |
| `obsidian_discover_structure` | Discover and analyze markdown file structure |
| `obsidian_get_nested_content` | Get content using nested path selectors |
| `obsidian_read_content` | Read specific content using selectors |

## 📖 Examples

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

## 🐳 Docker Deployment

This project includes production-ready Docker support with multiple deployment options:

- **HTTP Transport**: Port 8080 (default)
- **SSE Transport**: Port 8081 (default)
- **Both Services**: Simultaneous HTTP and SSE support

### Quick Docker Deployment

```bash
# Deploy both services
./docker-deploy.sh start both

# Deploy HTTP only
./docker-deploy.sh start http

# Deploy SSE only
./docker-deploy.sh start sse

# View logs
./docker-deploy.sh logs

# Check status
./docker-deploy.sh status
```

See [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) for comprehensive deployment instructions.

## 🏗️ Development

### Project Structure

```
mcp-obsidian/
├── cmd/
│   ├── root.go              # Root command
│   ├── obsidian-mcp.go      # Obsidian MCP command
│   └── comprehensive-target-test.go  # Comprehensive testing suite
├── obsidian/
│   ├── client/
│   │   └── client.go        # HTTP client for Obsidian API
│   ├── handlers/
│   │   └── obsidian.go      # MCP tool handlers
│   ├── types/
│   │   └── types.go         # Data types
│   └── prompts/
│       └── obsidian-comprehensive.md  # Comprehensive prompts
├── main.go                  # Entry point
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies checksum
├── Dockerfile.http          # HTTP transport Dockerfile
├── Dockerfile.sse           # SSE transport Dockerfile
├── docker-compose.http.yml  # HTTP service compose
├── docker-compose.sse.yml   # SSE service compose
├── docker-deploy.sh         # Deployment script
├── .dockerignore            # Docker ignore file
├── openapi.yaml             # OpenAPI specification
├── test-connection.sh       # Connection test script
├── cursor-mcp-obsidian.json # Cursor integration config
└── README.md               # This file
```

### Building

```bash
# Build for current platform
go build -o mcp-obsidian .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o mcp-obsidian-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o mcp-obsidian-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o mcp-obsidian-windows-amd64.exe .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./cmd -run TestComprehensiveTarget
```

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Add tests if applicable**
5. **Run tests**:
   ```bash
   go test ./...
   ```
6. **Commit your changes**:
   ```bash
   git commit -m "feat: add your feature description"
   ```
7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
8. **Submit a pull request**

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. **Fork and clone the repository**
2. **Install dependencies**:
   ```bash
   go mod download
   ```
3. **Set up environment variables** (see Configuration section)
4. **Run tests**:
   ```bash
   go test ./...
   ```
5. **Start development server**:
   ```bash
   go run main.go obsidian-mcp --http-port 8080
   ```

### Code Style

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for complex logic
- Write tests for new functionality
- Update documentation as needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) - The underlying API that makes this possible
- [mcp-go](https://github.com/mark3labs/mcp-go) - The Go MCP library
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Go](https://golang.org/) - The Go programming language

## 📞 Support

- 📧 **Email**: [your-email@example.com](mailto:your-email@example.com)
- 🐛 **Issues**: [GitHub Issues](https://github.com/yourusername/mcp-obsidian/issues)
- 📖 **Documentation**: [GitHub Wiki](https://github.com/yourusername/mcp-obsidian/wiki)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/yourusername/mcp-obsidian/discussions)

---

**Made with ❤️ for the Obsidian community**
