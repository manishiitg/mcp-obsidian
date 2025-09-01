# Release v1.0.0 - Initial Release

## 🎉 Initial Release

This is the initial release of the MCP Obsidian Server, providing comprehensive tools for interacting with Obsidian notes and files through the Obsidian Local REST API.

## ✨ Key Features

### Core Functionality
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

### Advanced Features
- 🔄 **Nested Content**: Advanced nested content retrieval and manipulation
- 🚀 **Multiple Transports**: HTTP, SSE, and stdio support
- 🐳 **Docker Support**: Production-ready Docker containers
- 🎯 **Cursor Integration**: Easy integration with Cursor IDE

## 🛠️ Available Tools

### Core Tools (10)
1. `obsidian_test_connection` - Test connection to Obsidian API
2. `obsidian_list_files_in_vault` - List all files in the vault
3. `obsidian_list_files_in_dir` - List files in a specific directory
4. `obsidian_get_file_contents` - Get contents of a file
5. `obsidian_search` - Search for text in the vault
6. `obsidian_append_content` - Append content to a file
7. `obsidian_put_content` - Create or update a file
8. `obsidian_delete_file` - Delete a file or directory
9. `obsidian_patch_content` - Patch content in a file
10. `obsidian_search_json` - Perform complex JSON-based search

### Advanced Tools (10)
11. `obsidian_get_periodic_note` - Get or create a periodic note
12. `obsidian_create_periodic_note` - Create a new periodic note
13. `obsidian_get_recent_changes` - Get recent changes in the vault
14. `obsidian_get_tags` - Get all tags in the vault
15. `obsidian_get_frontmatter` - Get frontmatter from a file
16. `obsidian_set_frontmatter` - Set frontmatter for a file
17. `obsidian_get_block_reference` - Get a block reference from a file
18. `obsidian_discover_structure` - Discover and analyze markdown file structure
19. `obsidian_get_nested_content` - Get content using nested path selectors
20. `obsidian_read_content` - Read specific content using selectors

## 🚀 Quick Start

### Prerequisites
1. **Obsidian Local REST API Plugin**: Install and enable the [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin
2. **API Key**: Generate an API key from the plugin settings
3. **Go 1.24+**: Required to build and run the server

### Installation
```bash
# Clone the repository
git clone https://github.com/manishiitg/mcp-obsidian.git
cd mcp-obsidian

# Build the project
go build -o mcp-obsidian .

# Set environment variables
export OBSIDIAN_API_KEY="your-api-key-here"
export OBSIDIAN_HOST="127.0.0.1"
export OBSIDIAN_PORT="27124"
export OBSIDIAN_USE_HTTPS="true"

# Start the server
./mcp-obsidian obsidian-mcp --http-port 8080
```

## 🐳 Docker Support

Production-ready Docker containers are available:

```bash
# Deploy both services
./docker-deploy.sh start both

# Deploy HTTP only
./docker-deploy.sh start http

# Deploy SSE only
./docker-deploy.sh start sse
```

## 📚 Documentation

- [README.md](README.md) - General project documentation
- [OBSIDIAN_API_DOCUMENTATION.md](OBSIDIAN_API_DOCUMENTATION.md) - Detailed API documentation
- [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker deployment guide
- [CURSOR_SETUP.md](CURSOR_SETUP.md) - Cursor integration setup guide

## 🔧 Configuration

### Environment Variables
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OBSIDIAN_API_KEY` | ✅ | - | Your Obsidian Local REST API key |
| `OBSIDIAN_HOST` | ❌ | `127.0.0.1` | Obsidian API host |
| `OBSIDIAN_PORT` | ❌ | `27124` | Obsidian API port |
| `OBSIDIAN_VAULT_PATH` | ❌ | - | Path to your Obsidian vault |
| `OBSIDIAN_USE_HTTPS` | ❌ | `true` | Use HTTPS for API calls |
| `OBSIDIAN_PROTOCOL` | ❌ | - | Protocol to use (http/https) |

## 🎯 Use Cases

- **AI Assistant Integration**: Seamlessly integrate AI assistants with Obsidian vaults
- **Automated Note Management**: Automate note creation, updates, and organization
- **Content Discovery**: Advanced search and content retrieval capabilities
- **Workflow Automation**: Streamline note-taking and knowledge management workflows
- **Development Tools**: Integrate with development environments and IDEs

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) - The underlying API that makes this possible
- [mcp-go](https://github.com/mark3labs/mcp-go) - The Go MCP library
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Go](https://golang.org/) - The Go programming language

---

**Made with ❤️ for the Obsidian community**

