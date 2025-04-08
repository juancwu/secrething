#!/bin/bash

set -e

# Use the passed CWD or default to the current directory
CWD=${CWD:-$(pwd)}
cd "$CWD" || exit 1

echo "Setting up Konbini local development environment..."

# Create necessary directories
echo "Creating required directories..."
mkdir -p "$CWD/.local"

# Check and install required tools
echo "Checking for required tools..."

# Check Go
if ! command -v go &> /dev/null; then
  echo "Error: Go is required but not installed."
  echo "Please install Go 1.24.0 from https://golang.org/doc/install"
  exit 1
else
  GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
  echo "Go version $GO_VERSION installed"
  if [[ "$GO_VERSION" < "1.24.0" ]]; then
    echo "Warning: Your Go version is older than 1.24.0, consider upgrading."
  fi
fi

# Check Air
if ! command -v air &> /dev/null; then
  echo "Installing Air v1.61.5..."
  go install github.com/cosmtrek/air@v1.61.5
else
  AIR_VERSION=$(air -v | grep -oP 'v\d+\.\d+\.\d+')
  echo "Air $AIR_VERSION installed"
fi

# Check Goose
if ! command -v goose &> /dev/null; then
  echo "Installing Goose v3.24.0..."
  go install github.com/pressly/goose/v3/cmd/goose@v3.24.0
else
  GOOSE_VERSION=$(goose -version | awk '{print $3}')
  echo "Goose $GOOSE_VERSION installed"
fi

# Check sqlc
if ! command -v sqlc &> /dev/null; then
  echo "Installing sqlc v1.28.0..."
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0
else
  SQLC_VERSION=$(sqlc version)
  echo "sqlc $SQLC_VERSION installed"
fi

# Check for Turso
if ! command -v turso &> /dev/null; then
  echo "Installing Turso CLI..."
  curl -sSfL https://get.turso.tech/install.sh | bash
else
  TURSO_VERSION=$(turso -v)
  echo "Turso CLI $TURSO_VERSION installed"
fi

# Download dependencies
echo "Downloading Go dependencies..."
go mod download

# Generate SQL code
echo "Generating SQL code with sqlc..."
sqlc generate

# Set up local database
echo "Setting up local database..."
if [ ! -f "$CWD/.local/local.db" ]; then
  echo "Creating new local database..."
  turso dev new --db-file "$CWD/.local/local.db"
fi

# Run migrations
echo "Running database migrations..."
chmod +x "$CWD/.scripts/migrate.sh"
"$CWD/.scripts/migrate.sh" up

echo "Setup complete! You can now run the development server with:"
echo "  gopack run dev:server"
echo ""
echo "Available commands:"
echo "  gopack run build:cli   - Build the CLI"
echo "  gopack run dev:server  - Run the development server with hot reload"
echo "  gopack run migrate     - Manage database migrations"
echo "  gopack run test        - Run tests"
echo ""
echo "Happy coding!"

# Make scripts executable if they aren't already
chmod +x "$CWD/.scripts/migrate.sh" "$CWD/.scripts/run_tests.sh"
