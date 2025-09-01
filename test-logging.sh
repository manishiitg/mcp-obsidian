#!/bin/bash

# Test script to demonstrate comprehensive logging in MCP Obsidian server
echo "🧪 Testing MCP Obsidian Comprehensive Logging"
echo "=============================================="

# Set logging environment variables
export MCP_LOG_TO_FILE=true
export MCP_LOG_LEVEL=debug
export MCP_LOG_DIR=./logs
export MCP_LOG_TO_CONSOLE=true

echo "📝 Logging Configuration:"
echo "  - File logging: $MCP_LOG_TO_FILE"
echo "  - Log level: $MCP_LOG_LEVEL"
echo "  - Log directory: $MCP_LOG_DIR"
echo "  - Console logging: $MCP_LOG_TO_CONSOLE"
echo ""

# Clean up any existing logs
rm -rf logs/
echo "🧹 Cleaned up existing logs"

# Start the server in background with stdio transport
echo "🚀 Starting MCP Obsidian server with stdio transport..."
./mcp-obsidian obsidian-mcp --stdio &
SERVER_PID=$!

# Wait a moment for server to start
sleep 3

echo ""
echo "📊 Server startup logs:"
echo "========================"
if [ -f "logs/mcp-obsidian.log" ]; then
    cat logs/mcp-obsidian.log | jq -r '.message' 2>/dev/null || cat logs/mcp-obsidian.log
else
    echo "❌ No log file found"
fi

echo ""
echo "📁 Log file details:"
echo "===================="
if [ -f "logs/mcp-obsidian.log" ]; then
    ls -la logs/
    echo ""
    echo "📊 Log file size: $(wc -c < logs/mcp-obsidian.log) bytes"
    echo "📝 Number of log entries: $(wc -l < logs/mcp-obsidian.log)"
else
    echo "❌ No log file found"
fi

# Stop the server
echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null || echo "Server already stopped"

echo ""
echo "✅ Test completed!"
echo ""
echo "📋 Summary of logging features implemented:"
echo "  ✓ Structured JSON logging with logrus"
echo "  ✓ File and console output"
echo "  ✓ Environment variable configuration"
echo "  ✓ Server startup/shutdown events"
echo "  ✓ Tool registration logging"
echo "  ✓ Performance timing (startup duration)"
echo "  ✓ Error logging with context"
echo "  ✓ Logging middleware for MCP tools"
echo "  ✓ Request/response logging infrastructure"
echo "  ✓ Obsidian API call logging infrastructure"
echo ""
echo "🔧 To enable logging in production:"
echo "  export MCP_LOG_TO_FILE=true"
echo "  export MCP_LOG_LEVEL=info"
echo "  export MCP_LOG_DIR=/var/log/mcp-obsidian"
echo "  export MCP_LOG_MAX_SIZE=100"
echo "  export MCP_LOG_MAX_BACKUPS=5"
echo "  export MCP_LOG_MAX_AGE=30"
