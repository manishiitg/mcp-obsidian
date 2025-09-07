#!/bin/bash

# Test script for Obsidian MCP Server

echo "🔍 Testing Obsidian MCP Server connection..."

# Set environment variables for testing
export OBSIDIAN_API_KEY="0488e4465eb68ba7cef5c38e35709e83e2579d67845d21be99edf1936ede1e48"
export OBSIDIAN_HOST="127.0.0.1"
export OBSIDIAN_PORT="27124"
export OBSIDIAN_VAULT_PATH="/Users/mipl/Library/CloudStorage/GoogleDrive-manish.prakash@citymall.live/My Drive/obsidian/drive"
export OBSIDIAN_USE_HTTPS="true"

echo "📋 Environment variables set:"
echo "  OBSIDIAN_API_KEY: ${OBSIDIAN_API_KEY:0:8}..."
echo "  OBSIDIAN_HOST: $OBSIDIAN_HOST"
echo "  OBSIDIAN_PORT: $OBSIDIAN_PORT"
echo "  OBSIDIAN_VAULT_PATH: $OBSIDIAN_VAULT_PATH"
echo "  OBSIDIAN_USE_HTTPS: $OBSIDIAN_USE_HTTPS"

echo ""
echo "🧪 Running comprehensive tests..."

# Run the test command
if ./mcp-obsidian test-obsidian; then
    echo ""
    echo "✅ All tests passed! MCP server is working correctly."
    echo ""
    echo "🚀 You can now start the MCP server with:"
    echo "   ./mcp-obsidian obsidian-mcp --stdio"
    echo ""
    echo "📝 Or start it with HTTP transport:"
    echo "   ./mcp-obsidian obsidian-mcp --port 8080"
    echo ""
    echo "🔗 The server will be available at:"
    echo "   - Stdio transport: Use with MCP clients"
    echo "   - HTTP transport: http://localhost:8080/mcp"
    echo "   - SSE transport: http://localhost:8080/sse"
else
    echo ""
    echo "❌ Some tests failed. Please check the output above for details."
    echo ""
    echo "🔧 Troubleshooting tips:"
    echo "   1. Make sure Obsidian is running with the Local REST API plugin enabled"
    echo "   2. Check that the API key is correct"
    echo "   3. Verify the host and port settings"
    echo "   4. Ensure the vault path is accessible"
    exit 1
fi
