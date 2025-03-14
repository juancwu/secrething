# Konbini Validator Documentation

Welcome to the comprehensive documentation for the Konbini Validator package. This powerful validation system is designed to provide flexible, customizable validation for Go applications with first-class support for the Echo web framework.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Getting Started](#getting-started)
- [Core Concepts](#core-concepts)
- [Documentation Sections](#documentation-sections)
- [Examples](#examples)
- [API Reference](#api-reference)

## Overview

The Konbini Validator extends the popular [go-playground/validator](https://github.com/go-playground/validator) package with enhanced features for web applications:

- Custom error messages with multiple levels of specificity
- Normalized field paths for consistent handling of nested structures
- Wildcard pattern matching for validating arrays and slices
- Context-specific validation rules
- Standardized error formatting
- Direct integration with Echo framework

## Key Features

### Custom Error Messages

Define validation error messages at multiple levels:

```go
v := validator.NewCustomValidator()

// Field-specific messages
v.Translator().SetFieldError("email", "email", "Please enter a valid email address")

// Default messages for validation tags
v.Translator().SetDefaultError("required", "This field is required")

// Generic fallback message
v.Translator().SetDefaultMessage("Invalid input")
```

### Validation Contexts

Create request-specific validation rules without modifying your base validator:

```go
// Base validator with common rules
baseValidator := validator.NewCustomValidator()

// Context for registration endpoint
registerContext := validator.NewValidationContext(baseValidator)
registerContext.SetFieldError("password", "password", "Password must be secure")
```

### Wildcard Pattern Matching

Set validation messages for array elements without specifying exact indices:

```go
// Apply to any item in an array
v.Translator().SetPatternError("items[*].name", "required", "Every item must have a name")

// Deeply nested wildcards
v.Translator().SetPatternError("departments[*].teams[*].members[*].email", "email", 
                               "All email addresses must be valid")
```

### Path Normalization

Consistently handle field paths regardless of formatting variations:

```go
// These are treated as equivalent
v.Translator().SetFieldError("user.profile.addresses[0].street", "required", "Street is required")
v.Translator().SetFieldError("user . profile . addresses [ 0 ] . street", "required", "Street is required")
```

### Structured Error Formatting

Format validation errors into consistent, nested API responses:

```go
formattedErrors := validator.FormatValidationErrors(validationErrors)
response := map[string]interface{}{
    "code":    400,
    "message": "Validation failed",
    "errors":  formattedErrors,
}
```

### Echo Framework Integration

Seamless integration with the Echo web framework:

```go
// Register with Echo
e := echo.New()
e.Validator = validator.NewCustomValidator()

// Use in handlers
func createUser(c echo.Context) error {
    user := new(User)
    
    if err := validator.BindAndValidate(c, user); err != nil {
        return err
    }
    
    // Process valid user...
}
```

## Getting Started

### Installation

The validator package is part of the Konbini framework:

```bash
go get github.com/juancwu/konbini
```

### Basic Usage

```go
import "github.com/juancwu/konbini/server/validator"

// Create a validator
v := validator.NewCustomValidator()

// Set custom error messages
v.Translator().SetFieldError("name", "required", "Please enter your name")
v.Translator().SetDefaultError("email", "Invalid email format")

// Define a struct with validation tags
type User struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
}

// Validate a struct
user := User{
    Name:     "",
    Email:    "not-an-email",
    Password: "weak",
}

if err := v.Validate(&user); err != nil {
    validationErrors := err.(validator.ValidationErrors)
    fmt.Printf("Validation errors: %+v\n", validator.FormatValidationErrors(validationErrors))
}
```

## Core Concepts

### CustomValidator

The main validator instance that wraps go-playground/validator and adds enhanced features.

### ErrorTranslator

Handles translation of validation errors into custom messages with support for field-specific, pattern, and default messages.

### ValidationContext

Provides request-specific validation rules for different use cases.

### Path Normalization

Standardizes field paths to ensure consistent error message lookup.

### Wildcard Pattern Matching

Uses special `[*]` syntax to match any array index in field paths.

### Structured Error Formatting

Formats validation errors into nested objects that match your data structure.

## Documentation Sections

Detailed documentation is available for each feature:

- [Path Normalization](path-normalization.md) - How path normalization works
- [Wildcard Pattern Matching](wildcard-patterns.md) - Using wildcards in validation
- [Validation Contexts](validation-contexts.md) - Creating endpoint-specific rules
- [Error Formatting](error-formatting.md) - Structured API error responses
- [Custom Validators](custom-validators.md) - Creating your own validation rules
- [Echo Integration](echo-integration.md) - Using with Echo framework

## Examples

See the [example_validator_usage.go](/example_validator_usage.go) file for comprehensive examples of validator usage.

## API Reference

### Core Types

- **CustomValidator**: Main validator instance
  - `NewCustomValidator()`: Creates a new validator
  - `Validate(interface{}) error`: Validates a struct
  - `Translator() *ErrorTranslator`: Gets the translator instance
  - `Clone() *CustomValidator`: Creates a copy of the validator

- **ErrorTranslator**: Handles error message translations
  - `SetFieldError(field, tag, message)`: Sets a field-specific message
  - `SetPatternError(pattern, tag, message)`: Sets a wildcard pattern message
  - `SetDefaultError(tag, message)`: Sets a default message for a tag
  - `SetDefaultMessage(message)`: Sets a generic fallback message

- **ValidationContext**: Provides context-specific validation rules
  - `NewValidationContext(baseValidator)`: Creates a new context
  - `SetFieldError(field, tag, message)`: Sets a context-specific field message
  - `SetPatternError(pattern, tag, message)`: Sets a context-specific pattern message
  - `Validate(interface{}) error`: Validates using context-specific rules

- **ValidationErrors**: Collection of validation errors
  - `Error() string`: Implements the error interface
  - `AsMap() map[string]interface{}`: Converts to a simple map

### Key Functions

- `FormatValidationErrors(validationErrors)`: Formats errors into a structured map
- `BindAndValidate(c, i)`: Binds and validates an Echo request
- `BindAndValidateWithContext(c, i, vc)`: Binds and validates with a context