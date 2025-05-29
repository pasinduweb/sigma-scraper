# Sigma Scraper

A fast and scalable web scraper built with Go for extracting product data from e-commerce sites. It handles modern JavaScript-heavy websites using a headless browser approach with chromedp. Built to process large data sets with strong error handling and efficient concurrency.

## Features

-   **High Performance**: Concurrent processing with configurable worker pools
-   **Robust Error Handling**: Retry mechanisms with exponential backoff
-   **Resource Efficient**: Optimized memory usage and browser context management
-   **Production Ready**: Docker support, monitoring, and graceful shutdown

## Project Structure

```
├── bin/                  # Compiled application binary (e.g., scraper.exe, scraper)
├── cmd/
│   └── scraper/          # Main application entry point
├── data/                 # Directory for input data or persistent storage (if used)
├── internal/
│   ├── config/           # Configuration management
│   ├── models/           # Data structures
│   ├── scraper/          # Core scraping logic
│   ├── storage/          # Data storage and persistence
│   └── utils/            # Utility functions
├── output/               # Default directory for scraper output files
├── test/                 # Test files
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Makefile              # Build automation for Linux/macOS
├── README.md             # This file
├── .env                  # Optional: Environment variables file
├── samples.xlsx          # Example input Excel file
├── setup.bat             # Setup script for Windows
└── setup.sh              # Setup script for Linux/macOS
```

## Requirements

-   Go 1.21 or higher
-   Chrome/Chromium (for headless browser automation)
-   Excel input file with product IDs and URLs

## Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/pasinduweb/sigma-scraper.git
    cd sigma-scraper
    ```

2.  **Setup and Install Dependencies:**

    -   **For Linux/macOS:**
        Run the setup script. This will create necessary directories, check for Go and Chrome/Chromium, install dependencies, and build the application.

        ```bash
        ./setup.sh
        ```

        Alternatively, if you only want to install dependencies using Make (you will need to build separately):

        ```bash
        make install-deps
        ```

    -   **For Windows:**
        Run the setup batch file. This will create necessary directories, check for Go and Chrome, install dependencies, and build the application.
        ```cmd
        .\setup.bat
        ```

## Configuration

Configuration is managed by environment variables. The application follows this order of precedence:

1.  Environment variables set directly in your terminal session. These will always take the highest priority.
2.  Variables loaded from a `.env` file located in the project root.
3.  Default values defined within the system if a variable is not found in either of the above locations.

To use a `.env` file, create one in the root of the project. For example:

```env
INPUT_FILE="my_products.xlsx"
WORKER_COUNT=10
```

The following environment variables are supported:

-   `INPUT_FILE`: Path to Excel file with product data (default: "samples.xlsx" located in the project root directory).
-   `OUTPUT_DIR`: Directory for output files (default: "output" in the project root).
-   `WORKER_COUNT`: Number of concurrent workers (default: 5)
-   `BUFFER_SIZE`: Channel buffer size (default: 100)
-   `REQUEST_TIMEOUT_SECONDS`: Browser request timeout in seconds (default: 30)
-   `PAGE_LOAD_DELAY_MS`: Delay after page load in milliseconds (default: 1000)
-   `MAX_RETRIES`: Maximum retry attempts for failed pages (default: 3)
-   `RETRY_DELAY_SECONDS`: Delay between retries in seconds (default: 5)

## Usage

-   **For Linux/macOS:**

    -   **If you used `./setup.sh` for installation:**
        The application is already built. Run it from the project root directory:
        ```bash
        ./bin/scraper
        ```
    -   **If you prefer using `make` (after running `make install-deps` or `./setup.sh`):**
        Build the scraper (if not already built by `./setup.sh`):
        ```bash
        make build
        ```
        Run the scraper:
        ```bash
        make run
        ```

-   **For Windows:**

    -   **After running `.\setup.bat` for installation:**
        The application is already built. Run it from the project root directory:
        ```cmd
        .\bin\scraper.exe
        ```
    -   **Note:** The `make` commands (`make build`, `make run`) are primarily for Linux/macOS environments. While `make` can be installed on Windows, the recommended method for this project on Windows is to use `.\setup.bat` for setup and then directly execute `.\bin\scraper.exe`.

## Input Format

The scraper expects an Excel file with at least two columns:

-   Product ID (labeled as "id", "product_id", "productid", or "sku")
-   Product URL (labeled as "link", "url", "product_link", or "pdp_url")

## Output Format

The scraper generates the following output:

-   `output/final_output.json`: Final results with product IDs and image URLs
-   `output/failed_urls.json`: List of URLs that failed to scrape

## License

This project is licensed under the MIT License

See the [LICENSE](./LICENSE) file for detailed licensing terms.
