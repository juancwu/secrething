# Konbini Server Architecture

This directory contains the server-side code for the Konbini application.

## Directory Structure

The server code is organized as follows:

```
/server
  /config         # Application configuration
  /cookies        # Cookie-related utilities
  /database       # Database connection wrappers
  /db             # Generated SQL code and database models
  /errors         # Error types and handling
  /handlers       # HTTP request handlers
  /helpers        # Utility functions for request handling
  /middleware     # HTTP middleware
  /observability  # Monitoring, logging, and error reporting
  /routes         # API route definitions
  /utils          # General utility functions
```

## Components

### Config

Configuration loading and validation, including environment variables management.

### Cookies

Utilities for handling HTTP cookies securely.

### Database

Database connection wrappers and utilities. The actual SQL models and queries are in the `/db` directory.

### Errors

Custom error types and error handling utilities that provide standardized error responses.

### Handlers

HTTP request handlers organized by functional domain.

### Helpers

Helper functions for things like parsing query parameters.

### Middleware

HTTP middleware for cross-cutting concerns like:
- Request validation
- Error handling
- Logging
- Request metrics

### Observability

Tools for monitoring, logging, and error reporting (e.g., Sentry integration).

### Routes

API route definitions and route registration.

### Utils

General utility functions used across the application.

## Development

When adding new functionality:

1. Create handlers in the appropriate domain directory under `/handlers`
2. Register routes in the `/routes` package
3. Add database queries to the appropriate file in `/db/queries` and regenerate with `sqlc`
4. Add middleware in the `/middleware` directory

For more comprehensive development instructions, see the project README.md.