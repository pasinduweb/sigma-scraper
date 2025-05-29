package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/product-scraper/internal/config"
	"github.com/product-scraper/internal/models"
)

type Manager struct {
	config     *config.Config
	mutex      sync.RWMutex
	failedURLs []models.FailedURL
	allResults []models.ProductResult
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:     cfg,
		failedURLs: make([]models.FailedURL, 0),
		allResults: make([]models.ProductResult, 0),
	}
}

// SaveResult directly appends a successful result to the in-memory list
func (m *Manager) SaveResult(result models.ProductResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.allResults = append(m.allResults, result)
	log.Printf("Stored result for product ID: %s. Total results: %d", result.ID, len(m.allResults))
}

// SaveFailedURL saves a failed URL
func (m *Manager) SaveFailedURL(id string, url string, errMsg string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Look for existing entry
	var found bool
	for i, failed := range m.failedURLs {
		if failed.ID == id && failed.URL == url {
			// Update existing entry
			m.failedURLs[i].Error = errMsg
			m.failedURLs[i].Timestamp = time.Now()
			m.failedURLs[i].Attempts++
			found = true
			break
		}
	}

	// Add new entry if not found
	if !found {
		m.failedURLs = append(m.failedURLs, models.FailedURL{
			ID:        id,
			URL:       url,
			Error:     errMsg,
			Timestamp: time.Now(),
			Attempts:  1,
		})
	}
}

// GenerateFinalOutput combines all snapshots into a single output file
func (m *Manager) GenerateFinalOutput() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Save any remaining failed URLs
	m.saveFailedURLs()

	// Save final output from the in-memory slice
	if err := saveToJSON(m.config.FinalOutputFile, m.allResults); err != nil {
		return fmt.Errorf("failed to save final output: %v", err)
	}

	log.Printf("Final output saved with %d records", len(m.allResults))
	return nil
}

// saveFailedURLs saves the failed URLs to a file
func (m *Manager) saveFailedURLs() {
	if err := saveToJSON(m.config.FailedURLsFile, m.failedURLs); err != nil {
		log.Printf("Failed to save failed URLs: %v", err)
	}
}

// saveToJSON saves data to a JSON file
func saveToJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
