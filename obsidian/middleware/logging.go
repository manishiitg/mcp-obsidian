package middleware

import (
	"context"
	"time"

	"mcp-obsidian/obsidian/logger"

	"github.com/mark3labs/mcp-go/mcp"
)

// LoggingMiddleware wraps an MCP tool handler to add comprehensive logging
func LoggingMiddleware(handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()

		// Generate a unique request ID for this call
		requestID := generateRequestID()
		ctx = context.WithValue(ctx, "request_id", requestID)

		// Log the incoming request
		logger.LogMCPRequest(ctx, "tool_call", req, map[string]interface{}{
			"request_id": requestID,
			"arguments":  req.Params.Arguments,
		})

		// Execute the actual handler
		result, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		logger.LogMCPResponse(ctx, "tool_call", result, duration, err, map[string]interface{}{
			"request_id": requestID,
			"arguments":  req.Params.Arguments,
		})

		// Log performance metrics
		logger.LogPerformance("mcp_tool_execution", duration, map[string]interface{}{
			"tool_name":  "tool_call",
			"request_id": requestID,
		})

		return result, err
	}
}

// generateRequestID generates a unique request identifier
func generateRequestID() string {
	// Simple timestamp-based ID for now
	// In production, you might want to use UUID or similar
	return time.Now().Format("20060102-150405.000000")
}
