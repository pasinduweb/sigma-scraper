package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/product-scraper/internal/config"
	"github.com/product-scraper/internal/models"
	"github.com/product-scraper/internal/scraper"
	"github.com/product-scraper/internal/storage"
	"github.com/product-scraper/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize storage manager
	storageManager := storage.NewManager(cfg)

	// Load products from Excel file
	log.Printf("Attempting to load Excel file from: %s", cfg.InputFile)
	products, err := utils.LoadProductsFromExcel(cfg.InputFile)
	if err != nil {
		log.Fatalf("Failed to load products from Excel: %v", err)
	}
	log.Printf("Loaded %d products from Excel file", len(products))

	startIndex := 0 // Resume logic removed

	// Initialize scraper
	scraperInstance := scraper.New(cfg)

	// Setup graceful shutdown
	ctx := setupGracefulShutdown()

	// Create channels for communication
	productChan := make(chan models.Product, cfg.BufferSize)
	resultChan := make(chan models.ProductResult, cfg.BufferSize)

	var wg sync.WaitGroup // Main WaitGroup

	// Start result processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		// cfg might not be needed by processResults if SnapshotInterval was its only use
		processResults(ctx, resultChan, storageManager) // cfg removed
	}()

	// Start product producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(productChan) // Producer closes productChan when done

		for i := startIndex; i < len(products); i++ {
			select {
			case <-ctx.Done():
				log.Println("Producer: Context done, stopping product feed.")
				return
			case productChan <- products[i]:
			}
		}
		log.Println("Producer: Finished sending all products.")
	}()

	// Start scraper workers and manage resultChan closure
	wg.Add(1) // Add the worker manager goroutine to the main WaitGroup
	go func() {
		defer wg.Done() // This goroutine is done when workers are done and resultChan is closed

		var workerWg sync.WaitGroup
		log.Printf("Starting %d scraper workers...", cfg.WorkerCount)
		for i := 0; i < cfg.WorkerCount; i++ {
			workerWg.Add(1)
			go func(workerID int) {
				defer workerWg.Done()
				scraperInstance.Worker(ctx, workerID, productChan, resultChan)
			}(i)
		}
		workerWg.Wait() // Wait for all workers to finish
		log.Println("All scraper workers finished.")
		close(resultChan) // Then close resultChan
		log.Println("Result channel closed.")
	}()

	log.Println("Main: Waiting for all goroutines to complete...")
	wg.Wait() // Wait for producer, processResults, and the worker manager (which closes resultChan)
	log.Println("Main: All goroutines completed.")

	// Clean up browser resources
	scraperInstance.Cleanup()

	// Generate final output
	log.Println("Main: Generating final output...")
	if err := storageManager.GenerateFinalOutput(); err != nil {
		log.Printf("Failed to generate final output: %v", err)
	}

	log.Println("Scraping completed successfully.")
}

func processResults(ctx context.Context, resultChan <-chan models.ProductResult, storage *storage.Manager) { // cfg *config.Config removed
	log.Println("ProcessResults: Started.")
	for {
		select {
		case <-ctx.Done():
			log.Println("ProcessResults: Context done, stopping.")
			return
		case result, ok := <-resultChan:
			if !ok {
				log.Println("ProcessResults: Result channel closed, stopping.")
				return
			}

			if result.Success {
				storage.SaveResult(result)
			} else {
				// Attempt to find the original product URL. This is a bit tricky as ProductResult doesn't store it.
				// For now, we'll use a placeholder if the URL isn't easily available.
				// This part might need adjustment based on how you want to log failed URLs.
				log.Printf("ProcessResults: Received failed result for ID %s: %s", result.ID, result.Error)
				storage.SaveFailedURL(result.ID, "N/A - URL not in ProductResult", result.Error)
			}
		}
	}
}

func setupGracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	return ctx
}
