# Konbini Validator

A powerful, flexible validation package for Go applications, built on top of [go-playground/validator](https://github.com/go-playground/validator/).

## Features

- **Custom error messages** for fields, patterns, and validation tags
- **Path normalization** for consistent field path handling
- **Wildcard pattern matching** for array/slice validation messages
- **Validation contexts** for request-specific validation rules
- **Nested field support** for complex data structures
- **Echo framework integration** for seamless HTTP request validation
- **Consistent error formatting** for standardized API responses
- **Built-in password validation** for secure applications

## Installation

The validator package is part of the Konbini framework:

```bash
go get github.com/juancwu/konbini
```

## Quick Start

```go
import "github.com/juancwu/konbini/server/validator"

// Create a user struct with validation tags
type User struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,gte=18"`
    Password string `json:"password" validate:"required,password"`
}

// Create a validator
v := validator.NewCustomValidator()

// Set custom error messages
v.Translator().SetFieldError("name", "required", "Please enter your name")
v.Translator().SetFieldError("email", "email", "Please enter a valid email address")
v.Translator().SetFieldError("age", "gte", "You must be at least 18 years old")

// Validate a user
user := User{
    Name:     "",
    Email:    "not-an-email",
    Age:      16,
    Password: "weak",
}

err := v.Validate(&user)
if err != nil {
    validationErrors := err.(validator.ValidationErrors)
    
    // Format errors for API response
    formattedErrors := validator.FormatValidationErrors(validationErrors)
    fmt.Printf("Validation errors: %+v\n", formattedErrors)
}
```

## Core Components

### CustomValidator

The main validator instance that wraps go-playground/validator with enhanced features:

```go
v := validator.NewCustomValidator()
err := v.Validate(myStruct)
```

### ErrorTranslator

Handles mapping validation errors to custom messages:

```go
translator := v.Translator()
translator.SetFieldError("name", "required", "Name is required")
translator.SetDefaultError("email", "Invalid email format")
translator.SetDefaultMessage("Invalid input")
```

### ValidationContext

Creates request-specific validation rules:

```go
ctx := validator.NewValidationContext(baseValidator)
ctx.SetFieldError("password", "min", "Password must be at least 8 characters")
```

### ValidationErrors

Collection of validation errors with formatting capabilities:

```go
errors := err.(validator.ValidationErrors)
formatted := errors.AsMap()
```

## Path Normalization

The validator normalizes field paths to ensure consistent handling:

```go
// These are all treated as equivalent
v.Translator().SetFieldError("user.profile.addresses[0].street", "required", "Street is required")
v.Translator().SetFieldError("user . profile . addresses [ 0 ] . street", "required", "Street is required")
```

## Wildcard Pattern Matching

Set messages for array/slice elements without specifying exact indices:

```go
// Applies to all items
v.Translator().SetPatternError("items[*].name", "required", "Every item must have a name")

// Deeply nested wildcards
v.Translator().SetPatternError("departments[*].teams[*].members[*].email", "email", "All emails must be valid")
```

## Custom Validation Rules

The validator includes built-in custom validation rules:

### Password Validation

Validates that passwords meet security requirements:
- At least 8 characters long
- Contains uppercase and lowercase letters
- Contains digits and special characters

```go
type User struct {
    Password string `json:"password" validate:"required,password"`
}
```

## Error Formatting

Format validation errors into structured responses:

```go
formattedErrors := validator.FormatValidationErrors(validationErrors)
response := map[string]interface{}{
    "code":    400,
    "message": "Validation failed",
    "errors":  formattedErrors,
}
```

The formatter handles:
- Nested objects: `user.profile.firstName`
- Array indices: `items[0].name`
- Multi-dimensional arrays: `matrix[0][1]`
- Arrays of objects with index information: 
  ```json
  {
    "items": [
      {"index": 0, "field_errors": {"name": "Name is required"}}
    ]
  }
  ```

## Echo Framework Integration

Seamless integration with the Echo web framework:

```go
// Register validator with Echo
e := echo.New()
e.Validator = validator.NewCustomValidator()

// In request handlers
func createUser(c echo.Context) error {
    user := new(User)
    
    if err := validator.BindAndValidate(c, user); err != nil {
        return err // Echo will handle the error
    }
    
    // Process valid user
    return c.JSON(http.StatusCreated, user)
}

// With validation context
func signupUser(c echo.Context) error {
    user := new(User)
    
    if err := validator.BindAndValidateWithContext(c, user, signupContext); err != nil {
        return err
    }
    
    // Process valid user
    return c.JSON(http.StatusCreated, user)
}
```

## Advanced Usage Examples

For more advanced usage examples, check:
- [Example Validator Usage](/example_validator_usage.go)
- [Path Normalization Documentation](/docs/validator/path-normalization.md)
- [Wildcard Patterns Documentation](/docs/validator/wildcard-patterns.md)

## API Reference

### Core Types
- `CustomValidator`: Main validator instance
- `ErrorTranslator`: Handles error message translations
- `ValidationContext`: Context-specific validation rules
- `ValidationError`: Single validation error
- `ValidationErrors`: Collection of validation errors

### Main Functions
- `NewCustomValidator()`: Creates a new validator instance
- `NewValidationContext()`: Creates a context with custom rules
- `FormatValidationErrors()`: Formats errors into structured map
- `BindAndValidate()`: Binds and validates an Echo request
- `BindAndValidateWithContext()`: Binds and validates with a context

### Error Message Customization
- `SetFieldError(field, tag, message)`: Set message for specific field and tag
- `SetPatternError(pattern, tag, message)`: Set message for wildcard pattern
- `SetDefaultError(tag, message)`: Set default message for validation tag
- `SetDefaultMessage(message)`: Set generic fallback message

## License

This package is part of the Konbini framework and is licensed under the terms of the included LICENSE file.