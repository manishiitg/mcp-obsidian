package main

import (
	"log"
	"os"
	"path/filepath"

	"mcp-obsidian/cmd"

	"github.com/joho/godotenv"
)

func main() {
	// Try loading .env files from a few common locations
	tryPaths := []string{".env"}

	// Directory where the executable resides (e.g., repo root if built there)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		tryPaths = append(tryPaths, filepath.Join(exeDir, ".env"))
		// One level up from executable directory (useful when binary lives in build/ or similar)
		tryPaths = append(tryPaths, filepath.Join(exeDir, "..", ".env"))
	}

	loaded := false
	for _, p := range tryPaths {
		if err := godotenv.Load(p); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		log.Printf("No .env file found in paths: %v", tryPaths)
	}

	cmd.Execute()
}
