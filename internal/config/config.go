package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv" // Import the godotenv package
)

type Config struct {
	// Input/Output
	InputFile       string
	OutputDir       string
	FinalOutputFile string
	FailedURLsFile  string

	// Scraping settings
	WorkerCount    int
	BufferSize     int
	RequestTimeout time.Duration
	PageLoadDelay  time.Duration

	// Browser settings
	BrowserFlags []string
	MaxRetries   int
	RetryDelay   time.Duration
}

// Load loads configuration from environment variables or uses defaults
func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		// Log a warning if .env file is not found or fails to load, but don't make it fatal
		// as environment variables might be set directly.
		log.Printf("Warning: Could not load .env file: %v. Proceeding with environment variables or defaults.", err)
	}

	// Get the executable directory
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not determine executable path: %v", err)
		exePath = ""
	}

	// Get the project root directory (parent of the bin directory)
	projectRoot := ""
	if exePath != "" {
		exeDir := filepath.Dir(exePath)
		// If the directory name is 'bin', go one level up
		if filepath.Base(exeDir) == "bin" {
			projectRoot = filepath.Dir(exeDir)
			log.Printf("Determined project root as: %s", projectRoot)
		}
	}
	cfg := &Config{
		// Default values
		InputFile:       getEnv("INPUT_FILE", filepath.Join(projectRoot, "samples.xlsx")),
		OutputDir:       getEnv("OUTPUT_DIR", "output"),
		FinalOutputFile: getEnv("FINAL_OUTPUT_FILE", "output/final_output.json"),
		FailedURLsFile:  getEnv("FAILED_URLS_FILE", "output/failed_urls.json"),
		// Match working script settings exactly
		WorkerCount: getEnvInt("WORKER_COUNT", 5),
		BufferSize:  getEnvInt("BUFFER_SIZE", 100),
		// Use same timeout as working script
		RequestTimeout: time.Duration(getEnvInt("REQUEST_TIMEOUT_SECONDS", 15)) * time.Second,
		PageLoadDelay:  time.Duration(getEnvInt("PAGE_LOAD_DELAY_MS", 2000)) * time.Millisecond,

		BrowserFlags: getBrowserFlags(),
		// Match working script retry settings
		MaxRetries: getEnvInt("MAX_RETRIES", 3),
		// Match working script retry delay
		RetryDelay: time.Duration(getEnvInt("RETRY_DELAY_SECONDS", 2)) * time.Second,
	}

	// Create directories if they don't exist
	createDirectoryIfNotExists(cfg.OutputDir)

	return cfg
}

// Helper functions
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid value for %s, using default: %d", key, fallback)
	}
	return fallback
}

func getBrowserFlags() []string {
	defaultFlags := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
	}

	if customFlagsStr := getEnv("BROWSER_FLAGS", ""); customFlagsStr != "" {
		// Parse custom flags (comma separated)
		// This is simplified; you might want a more robust parsing
		return append(defaultFlags, customFlagsStr)
	}

	return defaultFlags
}

func createDirectoryIfNotExists(path string) {
	// Check if path is absolute, if not and we're in the bin directory, make it relative to project root
	if !filepath.IsAbs(path) {
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			// If we're in the bin directory, assume we need to go up one level for project root
			if filepath.Base(exeDir) == "bin" {
				projectRoot := filepath.Dir(exeDir)
				path = filepath.Join(projectRoot, path)
				log.Printf("Creating directory at: %s", path)
			}
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", path, err)
		}
	}
}
