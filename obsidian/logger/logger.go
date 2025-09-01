package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
	LogLevelPanic LogLevel = "panic"
)

// LogConfig holds the logging configuration
type LogConfig struct {
	Level        LogLevel `json:"level"`
	LogToFile    bool     `json:"log_to_file"`
	LogToConsole bool     `json:"log_to_console"`
	LogDir       string   `json:"log_dir"`
	MaxSize      int      `json:"max_size"`    // Max file size in MB
	MaxBackups   int      `json:"max_backups"` // Max number of backup files
	MaxAge       int      `json:"max_age"`     // Max age in days
}

// DefaultLogConfig returns the default logging configuration
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:        LogLevelInfo,
		LogToFile:    true,
		LogToConsole: true,
		LogDir:       "logs",
		MaxSize:      100, // 100 MB
		MaxBackups:   5,
		MaxAge:       30, // 30 days
	}
}

// Logger is the main logger instance
type Logger struct {
	logger *logrus.Logger
	config *LogConfig
}

// Global logger instance
var globalLogger *Logger

// InitLogger initializes the global logger with the given configuration
func InitLogger(config *LogConfig) error {
	if config == nil {
		config = DefaultLogConfig()
	}

	// Create logs directory if it doesn't exist
	if config.LogToFile {
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// Create new logrus logger
	logrusLogger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logrusLogger.SetLevel(level)

	// Set formatter
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set output
	var outputs []io.Writer

	if config.LogToConsole {
		outputs = append(outputs, os.Stderr)
	}

	if config.LogToFile {
		logFile := filepath.Join(config.LogDir, "mcp-obsidian.log")
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		outputs = append(outputs, file)
	}

	if len(outputs) > 0 {
		logrusLogger.SetOutput(io.MultiWriter(outputs...))
	}

	globalLogger = &Logger{
		logger: logrusLogger,
		config: config,
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Initialize with default config if not already initialized
		if err := InitLogger(DefaultLogConfig()); err != nil {
			panic(fmt.Sprintf("failed to initialize logger: %v", err))
		}
	}
	return globalLogger
}

// LogMCPRequest logs an MCP tool request
func (l *Logger) LogMCPRequest(ctx context.Context, toolName string, request interface{}, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":       "mcp_request",
		"tool_name":  toolName,
		"timestamp":  time.Now().UTC(),
		"request_id": getRequestID(ctx),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	// Add request data (sanitized if needed)
	if request != nil {
		fields["request_data"] = sanitizeRequestData(request)
	}

	l.logger.WithFields(fields).Info("MCP Tool Request")
}

// LogMCPResponse logs an MCP tool response
func (l *Logger) LogMCPResponse(ctx context.Context, toolName string, response interface{}, duration time.Duration, err error, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":        "mcp_response",
		"tool_name":   toolName,
		"timestamp":   time.Now().UTC(),
		"duration_ms": duration.Milliseconds(),
		"request_id":  getRequestID(ctx),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	// Add response data (sanitized if needed)
	if response != nil {
		fields["response_data"] = sanitizeResponseData(response)
	}

	// Add error information if present
	if err != nil {
		fields["error"] = err.Error()
		fields["success"] = false
		l.logger.WithFields(fields).Error("MCP Tool Response Error")
	} else {
		fields["success"] = true
		l.logger.WithFields(fields).Info("MCP Tool Response Success")
	}
}

// LogObsidianAPICall logs an Obsidian API call
func (l *Logger) LogObsidianAPICall(ctx context.Context, method, endpoint string, request interface{}, response interface{}, duration time.Duration, err error) {
	fields := logrus.Fields{
		"type":        "obsidian_api_call",
		"method":      method,
		"endpoint":    endpoint,
		"timestamp":   time.Now().UTC(),
		"duration_ms": duration.Milliseconds(),
		"request_id":  getRequestID(ctx),
	}

	// Add request/response data (sanitized)
	if request != nil {
		fields["request_data"] = sanitizeRequestData(request)
	}
	if response != nil {
		fields["response_data"] = sanitizeResponseData(response)
	}

	// Add error information if present
	if err != nil {
		fields["error"] = err.Error()
		fields["success"] = false
		l.logger.WithFields(fields).Error("Obsidian API Call Error")
	} else {
		fields["success"] = true
		l.logger.WithFields(fields).Info("Obsidian API Call Success")
	}
}

// LogServerEvent logs server lifecycle events
func (l *Logger) LogServerEvent(eventType, message string, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "server_event",
		"event":     eventType,
		"timestamp": time.Now().UTC(),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Info(message)
}

// LogError logs an error with stack trace
func (l *Logger) LogError(err error, message string, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "error",
		"timestamp": time.Now().UTC(),
		"error":     err.Error(),
	}

	// Add stack trace
	if _, file, line, ok := runtime.Caller(1); ok {
		fields["file"] = file
		fields["line"] = line
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Error(message)
}

// LogPerformance logs performance metrics
func (l *Logger) LogPerformance(operation string, duration time.Duration, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":        "performance",
		"operation":   operation,
		"duration_ms": duration.Milliseconds(),
		"timestamp":   time.Now().UTC(),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Info("Performance Metric")
}

// LogInfo logs an info message
func (l *Logger) LogInfo(message string, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "info",
		"timestamp": time.Now().UTC(),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Info(message)
}

// LogDebug logs a debug message
func (l *Logger) LogDebug(message string, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "debug",
		"timestamp": time.Now().UTC(),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Debug(message)
}

// LogWarn logs a warning message
func (l *Logger) LogWarn(message string, metadata map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "warn",
		"timestamp": time.Now().UTC(),
	}

	// Add metadata if provided
	for k, v := range metadata {
		fields[k] = v
	}

	l.logger.WithFields(fields).Warn(message)
}

// Helper functions

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return "no-context"
	}

	// Try to get request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}

	return "unknown"
}

// sanitizeRequestData removes sensitive information from request data
func sanitizeRequestData(data interface{}) interface{} {
	// For now, just return the data as-is
	// In the future, this could filter out sensitive fields like API keys
	return data
}

// sanitizeResponseData removes sensitive information from response data
func sanitizeResponseData(data interface{}) interface{} {
	// For now, just return the data as-is
	// In the future, this could filter out sensitive fields
	return data
}

// Convenience functions for global logger

// LogMCPRequest logs an MCP tool request using the global logger
func LogMCPRequest(ctx context.Context, toolName string, request interface{}, metadata map[string]interface{}) {
	GetLogger().LogMCPRequest(ctx, toolName, request, metadata)
}

// LogMCPResponse logs an MCP tool response using the global logger
func LogMCPResponse(ctx context.Context, toolName string, response interface{}, duration time.Duration, err error, metadata map[string]interface{}) {
	GetLogger().LogMCPResponse(ctx, toolName, response, duration, err, metadata)
}

// LogObsidianAPICall logs an Obsidian API call using the global logger
func LogObsidianAPICall(ctx context.Context, method, endpoint string, request interface{}, response interface{}, duration time.Duration, err error) {
	GetLogger().LogObsidianAPICall(ctx, method, endpoint, request, response, duration, err)
}

// LogServerEvent logs a server event using the global logger
func LogServerEvent(eventType, message string, metadata map[string]interface{}) {
	GetLogger().LogServerEvent(eventType, message, metadata)
}

// LogError logs an error using the global logger
func LogError(err error, message string, metadata map[string]interface{}) {
	GetLogger().LogError(err, message, metadata)
}

// LogPerformance logs performance metrics using the global logger
func LogPerformance(operation string, duration time.Duration, metadata map[string]interface{}) {
	GetLogger().LogPerformance(operation, duration, metadata)
}

// LogInfo logs an info message using the global logger
func LogInfo(message string, metadata map[string]interface{}) {
	GetLogger().LogInfo(message, metadata)
}

// LogDebug logs a debug message using the global logger
func LogDebug(message string, metadata map[string]interface{}) {
	GetLogger().LogDebug(message, metadata)
}

// LogWarn logs a warning message using the global logger
func LogWarn(message string, metadata map[string]interface{}) {
	GetLogger().LogWarn(message, metadata)
}
