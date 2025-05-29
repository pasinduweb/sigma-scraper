package utils

import (
	"fmt"
	"strings"

	"github.com/product-scraper/internal/models"
	"github.com/xuri/excelize/v2"
)

// LoadProductsFromExcel loads products from an Excel file
func LoadProductsFromExcel(filename string) ([]models.Product, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing Excel file: %v\n", err)
		}
	}()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %v", err)
	}

	// Check if the sheet has any rows
	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file has no data rows")
	}

	// Find column indices for ID and Link
	headers := rows[0]
	idCol, linkCol := -1, -1

	for i, header := range headers {
		header = strings.ToLower(strings.TrimSpace(header))
		switch header {
		case "id", "product_id", "productid", "sku":
			idCol = i
		case "link", "url", "product_link", "pdp_url":
			linkCol = i
		}
	}

	if idCol == -1 || linkCol == -1 {
		return nil, fmt.Errorf("required columns (ID and Link) not found in Excel file")
	}

	// Parse data rows
	var products []models.Product

	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// Skip empty rows
		if len(row) <= idCol || len(row) <= linkCol {
			continue
		}

		id := strings.TrimSpace(row[idCol])
		link := strings.TrimSpace(row[linkCol])

		// Skip rows without ID or link
		if id == "" || link == "" {
			continue
		}

		products = append(products, models.Product{
			ID:   id,
			Link: link,
		})
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no valid products found in Excel file")
	}

	return products, nil
}

// FilterFailedURLs filters products to only include those that previously failed
func FilterFailedURLs(products []models.Product, failedURLs []models.FailedURL) []models.Product {
	// Create map of failed product IDs for quick lookup
	failedMap := make(map[string]bool)
	for _, failed := range failedURLs {
		failedMap[failed.ID] = true
	}

	// Filter products
	filtered := make([]models.Product, 0)
	for _, product := range products {
		if failedMap[product.ID] {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

// ChunkProducts splits a slice of products into chunks of the specified size
func ChunkProducts(products []models.Product, chunkSize int) [][]models.Product {
	if chunkSize <= 0 {
		return [][]models.Product{products}
	}

	chunks := make([][]models.Product, 0, (len(products)+chunkSize-1)/chunkSize)

	for i := 0; i < len(products); i += chunkSize {
		end := i + chunkSize
		if end > len(products) {
			end = len(products)
		}
		chunks = append(chunks, products[i:end])
	}

	return chunks
}
