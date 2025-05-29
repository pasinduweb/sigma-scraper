package scraper

import (
	"context"
	"fmt"
	"log"
	"strings" // Added for message filtering
	"time"

	"github.com/chromedp/chromedp"
	"github.com/product-scraper/internal/config"
	"github.com/product-scraper/internal/models"
)

type Scraper struct {
	config *config.Config
}

func New(cfg *config.Config) *Scraper {
	return &Scraper{
		config: cfg,
	}
}

// Worker processes products from the productChan and sends results to resultChan
func (s *Scraper) Worker(ctx context.Context, workerID int, productChan <-chan models.Product, resultChan chan<- models.ProductResult) {
	log.Printf("Worker %d started", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping due to context done", workerID)
			return
		case product, ok := <-productChan:
			if !ok {
				log.Printf("Worker %d stopping due to closed channel", workerID)
				return
			}

			// Process the product
			result := s.ScrapeProduct(ctx, workerID, product)

			// Send result back
			select {
			case <-ctx.Done():
				return
			case resultChan <- result:
			}
		}
	}
}

// ScrapeProduct scrapes a single product and returns the result
func (s *Scraper) ScrapeProduct(ctx context.Context, workerID int, product models.Product) models.ProductResult {
	result := models.ProductResult{
		ID:      product.ID,
		Images:  make([]string, 0),
		Success: false,
	}

	// Implement retry logic
	var lastErr error
	for attempt := 0; attempt < s.config.MaxRetries; attempt++ {
		// Check if parent context is done before attempting
		select {
		case <-ctx.Done():
			result.Error = "parent context canceled"
			return result
		default:
			// Continue with scraping
		}

		if attempt > 0 {
			log.Printf("Worker %d: Retrying product %s (attempt %d/%d)",
				workerID, product.ID, attempt+1, s.config.MaxRetries)
			time.Sleep(s.config.RetryDelay)
		}

		// Create a completely fresh browser context for each request like the working script
		images, err := s.scrapeWithFreshContext(ctx, product.Link)

		if err == nil {
			result.Images = images
			result.Success = true
			return result
		}

		lastErr = err
		log.Printf("Worker %d: Error scraping product %s: %v", workerID, product.ID, err)
	}

	if lastErr != nil {
		result.Error = lastErr.Error()
	} else {
		result.Error = "failed after maximum retries"
	}
	return result
}

// scrapeWithFreshContext creates a fresh browser context for each request
func (s *Scraper) scrapeWithFreshContext(parentCtx context.Context, url string) ([]string, error) {
	// Check if parent context is already canceled before starting
	select {
	case <-parentCtx.Done():
		return nil, fmt.Errorf("parent context canceled before starting scrape")
	default:
	}

	// Create a completely fresh browser context using background context
	// This detaches browser operations from the parent context hierarchy
	// exactly like the working single script approach

	// Custom error logger to filter out specific benign messages
	customErrorFunc := func(s string, i ...interface{}) {
		msg := fmt.Sprintf(s, i...)
		if strings.Contains(msg, "could not unmarshal event") || strings.Contains(msg, "unknown ClientNavigationReason value") {
			// Optionally log these to a debug log or just ignore
			// log.Printf("Chromedp Info (suppressed): %s", msg)
			return
		}
		log.Printf("Chromedp Error: %s", msg)
	}

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Ensure headless is explicitly set
		// Add other allocator options if needed
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer allocCancel()

	// Pass the custom error logger to the context
	browserCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithErrorf(customErrorFunc))
	defer cancel()

	// Use a shorter timeout like the working script (15s instead of 45s)
	requestTimeout := 15 * time.Second
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, requestTimeout)
	defer timeoutCancel()

	var imageURLs []string

	err := chromedp.Run(timeoutCtx,
		// Navigate to the page
		chromedp.Navigate(url),
		// Wait for the product images container to be visible
		chromedp.WaitVisible("#js-product-images-container", chromedp.ByID),
		// Give a little more time for everything to load - use exact same delay as working script
		chromedp.Sleep(2*time.Second),
		// Extract all image URLs
		chromedp.Evaluate(`
            Array.from(document.querySelectorAll('[data-slide-id="zoom"]'))
                .map(el => el.getAttribute('href'))
                .filter(href => href && href.startsWith('https://'))
        `, &imageURLs),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape URL %s: %v", url, err)
	}
	return imageURLs, nil
}

// Cleanup releases resources - now simplified since we don't maintain browser pool
func (s *Scraper) Cleanup() {
	log.Println("Cleanup completed - using fresh contexts per request")
}
