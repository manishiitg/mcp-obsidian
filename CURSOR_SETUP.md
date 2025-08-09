# Cursor Integration Setup

This guide explains how to integrate the MCP Obsidian Server with Cursor.

## Quick Setup

1. **Build the application** (if not already built):
   ```bash
   go build -o mcp-obsidian .
   ```

2. **Get the absolute path** to the executable:
   ```bash
   # On macOS/Linux
   realpath mcp-obsidian
   
   # Or use pwd + filename
   echo "$(pwd)/mcp-obsidian"
   
   # On Windows (PowerShell)
   (Get-Item mcp-obsidian.exe).FullName
   ```

3. **Configure environment variables**:
   Edit the `cursor-mcp-obsidian.json` file and update:
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

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OBSIDIAN_API_KEY` | ✅ Yes | - | Your Obsidian Local REST API key |
| `OBSIDIAN_HOST` | ❌ No | `127.0.0.1` | Obsidian API host |
| `OBSIDIAN_PORT` | ❌ No | `27124` | Obsidian API port |
| `OBSIDIAN_USE_HTTPS` | ❌ No | `true` | Use HTTPS for API calls |
| `OBSIDIAN_PROTOCOL` | ❌ No | `https` | Protocol to use (http/https) |

## Usage

Once configured, you can use the Obsidian MCP server in Cursor with commands like:

- `obsidian_test_connection` - Test connection to Obsidian API
- `obsidian_list_files_in_vault` - List all files in the vault
- `obsidian_get_file_contents` - Get contents of a file
- `obsidian_search` - Search for text in the vault
- `obsidian_append_content` - Append content to a file
- `obsidian_put_content` - Create or update a file
- `obsidian_delete_file` - Delete a file or directory
- `obsidian_discover_structure` - Discover markdown file structure
- `obsidian_get_nested_content` - Get content using nested path selectors
- `obsidian_read_content` - Read specific content using selectors

## Troubleshooting

1. **Connection issues**: Make sure Obsidian Local REST API plugin is running
2. **Permission issues**: Ensure the `mcp-obsidian` binary is executable (`chmod +x mcp-obsidian`)
3. **API key issues**: Verify your API key is correct in the JSON file
4. **Path issues**: Make sure the absolute path in the JSON file is correct
5. **Path not found**: Double-check the absolute path using `realpath mcp-obsidian` or `which mcp-obsidian`

## Notes

- The server uses stdio transport for seamless integration with Cursor
- Environment variables are passed through the JSON configuration
- No Docker required - runs directly as a local binary
- Supports all Obsidian MCP tools and features
- **Important**: Always use absolute paths for the executable in the JSON file
