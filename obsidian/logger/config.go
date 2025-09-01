package logger

import (
	"os"
	"strconv"
)

// LoadLogConfigFromEnv loads logging configuration from environment variables
func LoadLogConfigFromEnv() *LogConfig {
	config := DefaultLogConfig()

	// Log level
	if level := os.Getenv("MCP_LOG_LEVEL"); level != "" {
		if isValidLogLevel(LogLevel(level)) {
			config.Level = LogLevel(level)
		}
	}

	// Log to file
	if logToFile := os.Getenv("MCP_LOG_TO_FILE"); logToFile != "" {
		if parsed, err := strconv.ParseBool(logToFile); err == nil {
			config.LogToFile = parsed
		}
	}

	// Log to console
	if logToConsole := os.Getenv("MCP_LOG_TO_CONSOLE"); logToConsole != "" {
		if parsed, err := strconv.ParseBool(logToConsole); err == nil {
			config.LogToConsole = parsed
		}
	}

	// Log directory
	if logDir := os.Getenv("MCP_LOG_DIR"); logDir != "" {
		config.LogDir = logDir
	}

	// Max file size (MB)
	if maxSize := os.Getenv("MCP_LOG_MAX_SIZE"); maxSize != "" {
		if parsed, err := strconv.Atoi(maxSize); err == nil && parsed > 0 {
			config.MaxSize = parsed
		}
	}

	// Max backups
	if maxBackups := os.Getenv("MCP_LOG_MAX_BACKUPS"); maxBackups != "" {
		if parsed, err := strconv.Atoi(maxBackups); err == nil && parsed >= 0 {
			config.MaxBackups = parsed
		}
	}

	// Max age (days)
	if maxAge := os.Getenv("MCP_LOG_MAX_AGE"); maxAge != "" {
		if parsed, err := strconv.Atoi(maxAge); err == nil && parsed >= 0 {
			config.MaxAge = parsed
		}
	}

	return config
}

// isValidLogLevel checks if the given log level is valid
func isValidLogLevel(level LogLevel) bool {
	validLevels := []LogLevel{
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarn,
		LogLevelError,
		LogLevelFatal,
		LogLevelPanic,
	}

	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}
	return false
}

// GetLogConfigSummary returns a summary of the current logging configuration
func GetLogConfigSummary() map[string]interface{} {
	logger := GetLogger()
	if logger == nil {
		return map[string]interface{}{
			"status": "logger_not_initialized",
		}
	}

	return map[string]interface{}{
		"level":          string(logger.config.Level),
		"log_to_file":    logger.config.LogToFile,
		"log_to_console": logger.config.LogToConsole,
		"log_dir":        logger.config.LogDir,
		"max_size_mb":    logger.config.MaxSize,
		"max_backups":    logger.config.MaxBackups,
		"max_age_days":   logger.config.MaxAge,
	}
}
