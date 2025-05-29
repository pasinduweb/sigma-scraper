#!/bin/bash
# Setup script for Product Scraper

echo "Setting up Product Scraper environment..."

# Create necessary directories
echo "Creating directory structure..."
mkdir -p output data bin

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_VERSION_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_VERSION_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_VERSION_MAJOR" -lt 1 ] || ([ "$GO_VERSION_MAJOR" -eq 1 ] && [ "$GO_VERSION_MINOR" -lt 21 ]); then
    echo "Warning: Recommended Go version is 1.21 or higher. You have $GO_VERSION"
fi

# Install dependencies
echo "Installing Go dependencies..."
go mod download

# Check for Chrome/Chromium
echo "Checking for Chrome/Chromium..."
if ! command -v google-chrome &> /dev/null && ! command -v chromium &> /dev/null && ! command -v chromium-browser &> /dev/null; then
    echo "Warning: Chrome or Chromium not found. The scraper requires a Chrome-compatible browser."
    echo "Please install Chrome or Chromium to use the scraper."
else
    echo "Chrome/Chromium found."
fi

# Build the application
echo "Building application..."
go build -o bin/scraper ./cmd/scraper

echo ""
echo "Setup complete! You can now run the scraper with:"
echo "./bin/scraper"
echo ""
echo "Make sure to place your input Excel file in the correct location"
echo "or set the INPUT_FILE environment variable."
echo ""
echo "For more information, see the README.md file."
