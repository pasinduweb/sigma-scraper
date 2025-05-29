package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/product-scraper/internal/config"
	"github.com/product-scraper/internal/models"
	"github.com/product-scraper/internal/storage"
)

func setupTestEnvironment(t *testing.T) (*config.Config, string) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "scraper_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create subdirectories
	os.Mkdir(filepath.Join(tempDir, "output"), 0755)
	// os.Mkdir(filepath.Join(tempDir, "output", "snapshots"), 0755) // Removed

	// Create test configuration
	cfg := &config.Config{
		InputFile: filepath.Join(tempDir, "test_input.xlsx"),
		OutputDir: filepath.Join(tempDir, "output"),
		// SnapshotsDir:     filepath.Join(tempDir, "output", "snapshots"), // Removed
		FinalOutputFile: filepath.Join(tempDir, "output", "final_output.json"),
		FailedURLsFile:  filepath.Join(tempDir, "output", "failed_urls.json"),
		// ResumeDataFile:   filepath.Join(tempDir, "output", "resume_data.json"), // Removed
		WorkerCount: 2,
		BufferSize:  10,
		// SnapshotInterval: 5, // Removed
		RequestTimeout: 5 * time.Second,
		PageLoadDelay:  500 * time.Millisecond,
		MaxRetries:     1,
		RetryDelay:     1 * time.Second,
	}

	return cfg, tempDir
}

func cleanupTestEnvironment(t *testing.T, tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		t.Logf("Warning: Failed to clean up temp directory: %v", err)
	}
}

func TestStorageManager(t *testing.T) {
	cfg, tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	// Initialize storage manager
	sm := storage.NewManager(cfg)

	// Test saving a successful result
	testResult1 := models.ProductResult{
		ID:      "test1",
		Images:  []string{"http://example.com/img1.jpg", "http://example.com/img2.jpg"},
		Success: true,
	}
	sm.SaveResult(testResult1)

	// Test saving a failed URL (simulating how processResults would call it)
	failedProductID := "test2"
	failedProductURL := "http://example.com/nonexistent"
	failedErrorMsg := "Failed to load page"
	sm.SaveFailedURL(failedProductID, failedProductURL, failedErrorMsg)

	// Test generating final output
	err := sm.GenerateFinalOutput()
	if err != nil {
		t.Fatalf("Failed to generate final output: %v", err)
	}

	// Check if final output file exists and contains the successful result
	if _, err := os.Stat(cfg.FinalOutputFile); os.IsNotExist(err) {
		t.Errorf("Final output file was not created")
	} else {
		data, err := os.ReadFile(cfg.FinalOutputFile)
		if err != nil {
			t.Fatalf("Failed to read final output file: %v", err)
		}
		var finalResults []models.ProductResult
		if err := json.Unmarshal(data, &finalResults); err != nil {
			t.Fatalf("Failed to parse final output file: %v", err)
		}
		if len(finalResults) != 1 {
			t.Errorf("Expected 1 result in final output, got %d", len(finalResults))
		} else if finalResults[0].ID != testResult1.ID {
			t.Errorf("Expected result ID %s in final output, got %s", testResult1.ID, finalResults[0].ID)
		}
	}

	// Check if failed URLs file exists and contains the failed URL
	if _, err := os.Stat(cfg.FailedURLsFile); os.IsNotExist(err) {
		t.Errorf("Failed URLs file was not created")
	} else {
		data, err := os.ReadFile(cfg.FailedURLsFile)
		if err != nil {
			t.Fatalf("Failed to read failed URLs file: %v", err)
		}
		var failedURLs []models.FailedURL
		if err := json.Unmarshal(data, &failedURLs); err != nil {
			t.Fatalf("Failed to parse failed URLs file: %v", err)
		}
		if len(failedURLs) != 1 {
			t.Errorf("Expected 1 entry in failed URLs file, got %d", len(failedURLs))
		} else if failedURLs[0].ID != failedProductID {
			t.Errorf("Expected failed product ID %s, got %s", failedProductID, failedURLs[0].ID)
		} else if failedURLs[0].URL != failedProductURL {
			t.Errorf("Expected failed URL %s, got %s", failedProductURL, failedURLs[0].URL)
		}
	}
}

func TestProductLoading(t *testing.T) {
	// This is a placeholder for a test that would create a test Excel file
	// and verify that products are loaded correctly
	t.Skip("Implement this test with a real Excel file")
}
