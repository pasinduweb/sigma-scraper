package models

import "time"

// Product represents a single product to be scraped
type Product struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

// ProductResult represents the result of scraping a product
type ProductResult struct {
	ID      string   `json:"id"`
	Images  []string `json:"images"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

// FailedURL represents a failed scraping attempt
type FailedURL struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
	Attempts  int       `json:"attempts"`
}

// ResumeData represents the state for resuming interrupted scraping
type ResumeData struct {
	LastProcessedIndex int       `json:"last_processed_index"`
	Timestamp          time.Time `json:"timestamp"`
}
