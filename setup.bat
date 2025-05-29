@echo off
echo Setting up Product Scraper environment...

REM Create necessary directories
echo Creating directory structure...
mkdir output 2>nul
mkdir data 2>nul
mkdir bin 2>nul

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go 1.21 or higher.
    exit /b 1
)

REM Check Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo Found Go version: %GO_VERSION%

REM Install dependencies
echo Installing Go dependencies...
go mod download

REM Check for Chrome
echo Checking for Chrome...
reg query "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe" >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Chrome not found in registry. The scraper requires Chrome or Chromium.
    echo Please install Chrome to use the scraper.
) else (
    echo Chrome found.
)

REM Build the application
echo Building application...
go build -o bin\scraper.exe .\cmd\scraper

echo.
echo Setup complete! You can now run the scraper with:
echo bin\scraper.exe
echo.
echo Make sure to place your input Excel file in the correct location
echo or set the INPUT_FILE environment variable.
echo.
echo For more information, see the README.md file.
