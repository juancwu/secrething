# Error Formatting

The Konbini validator provides robust error formatting capabilities to standardize validation error responses across your application.

## Overview

When validation fails, the validator returns a `ValidationErrors` type, which is a slice of `ValidationError` objects. These can be easily formatted into a consistent structure for API responses or user feedback.

## ValidationError Structure

Each validation error contains:

```go
type ValidationError struct {
    Field   string `json:"field"`   // Field path (e.g., "user.address.street")
    Message string `json:"message"` // Error message (e.g., "Street is required")
    Tag     string `json:"tag,omitempty"`    // Validation tag that failed (e.g., "required")
    Value   any    `json:"value,omitempty"`  // Invalid value that was provided
}
```

## Basic Formatting

The `ValidationErrors` type implements the `error` interface and provides an `AsMap()` method for simple formatting:

```go
err := validator.Validate(&myStruct)
if err != nil {
    validationErrors := err.(validator.ValidationErrors)
    
    // Get a flat map of field name -> error message
    errorMap := validationErrors.AsMap()
    
    // Example output:
    // {
    //   "name": "Name is required",
    //   "email": "Invalid email format"
    // }
}
```

## Structured Error Formatting

For more complex data structures, the validator provides the `FormatValidationErrors` function that handles:

- Nested objects with dot notation
- Arrays/slices with indexed notation
- Multi-dimensional arrays
- Arrays of objects with index information

```go
formattedErrors := validator.FormatValidationErrors(validationErrors)

// Create a standardized API response
response := map[string]interface{}{
    "code":    http.StatusBadRequest,
    "message": "Validation failed",
    "errors":  formattedErrors,
}

// Send the response
return c.JSON(http.StatusBadRequest, response)
```

## Formatting Examples

### Flat Structure

For a simple struct:

```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}
```

The formatted errors would look like:

```json
{
  "name": "Name is required",
  "email": "Invalid email format"
}
```

### Nested Objects

For nested structures:

```go
type Address struct {
    Street string `json:"street" validate:"required"`
    City   string `json:"city" validate:"required"`
}

type User struct {
    Name    string  `json:"name" validate:"required"`
    Address Address `json:"address" validate:"required"`
}
```

The formatted errors would look like:

```json
{
  "name": "Name is required",
  "address": {
    "street": "Street is required",
    "city": "City is required"
  }
}
```

### Arrays and Slices

For arrays and slices:

```go
type Item struct {
    Name     string  `json:"name" validate:"required"`
    Price    float64 `json:"price" validate:"required,gt=0"`
    Quantity int     `json:"quantity" validate:"required,gt=0"`
}

type Order struct {
    CustomerID string `json:"customerId" validate:"required"`
    Items      []Item `json:"items" validate:"required,dive"`
}
```

The formatted errors would show index information for each element:

```json
{
  "customerId": "Customer ID is required",
  "items": [
    {
      "index": 0,
      "field_errors": {
        "name": "Name is required",
        "price": "Price must be greater than zero"
      }
    },
    {
      "index": 1,
      "field_errors": {
        "quantity": "Quantity must be greater than zero"
      }
    }
  ]
}
```

## Implementation Details

The `FormatValidationErrors` function in `format.go` handles the transformation of validation errors into structured maps:

```go
func FormatValidationErrors(valErrors ValidationErrors) map[string]interface{} {
    fieldErrors := make(map[string]interface{})
    
    for _, validationErr := range valErrors {
        setNestedErrorWithArrays(fieldErrors, validationErr.Field, validationErr.Message)
    }
    
    return fieldErrors
}
```

The core algorithm recursively processes field paths and builds a nested structure:

1. Split the path into segments (e.g., "items[0].name" → ["items[0]", "name"])
2. Process each segment, handling array/slice notation
3. Create nested maps for object fields and arrays for indexed fields
4. Set the error message at the correct location

## Key Helper Functions

- **splitPathWithArrays**: Splits a path into segments
- **isArraySegment**: Checks if a segment contains array notation
- **parseArraySegment**: Extracts field name and index from array segment
- **ensureArrayExists**: Ensures an array exists and is large enough
- **setNestedErrorWithArrays**: Sets an error message in the nested structure

## Usage with Echo Framework

When using the Echo framework, you can integrate error formatting into your error handler:

```go
e.HTTPErrorHandler = func(err error, c echo.Context) {
    var validationErrors validator.ValidationErrors
    var statusCode int = http.StatusInternalServerError
    var message string = "Internal server error"
    
    if errors.As(err, &validationErrors) {
        // Handle validation errors
        statusCode = http.StatusBadRequest
        message = "Validation failed"
        
        // Format validation errors
        fieldErrors := validator.FormatValidationErrors(validationErrors)
        
        // Return standardized response
        c.JSON(statusCode, map[string]interface{}{
            "code":         statusCode,
            "message":      message,
            "field_errors": fieldErrors,
        })
        return
    }
    
    // Handle other types of errors
    // ...
}
```

## Best Practices

1. **Consistent formatting**: Use the same error format across your entire API
2. **Descriptive error messages**: Ensure messages clearly explain the issue
3. **Include error codes**: Consider adding error codes for programmatic handling
4. **Client-friendly**: Format errors in a way that's easy for clients to process
5. **Localization-ready**: Design your error system to support multiple languages

## Example: Complete API Error Response

```json
{
  "code": 400,
  "message": "Validation failed",
  "field_errors": {
    "user": {
      "name": "Name is required",
      "email": "Invalid email format",
      "address": {
        "street": "Street is required",
        "city": "City is required"
      }
    },
    "payment": {
      "cardNumber": "Invalid card number",
      "expiryDate": "Card has expired"
    },
    "items": [
      {
        "index": 0,
        "field_errors": {
          "name": "Item name is required",
          "price": "Price must be greater than zero"
        }
      },
      {
        "index": 2,
        "field_errors": {
          "quantity": "Quantity must be greater than zero"
        }
      }
    ]
  }
}
```

## Advanced: Custom Error Formatting

You can create custom formatting functions by extending the base implementation:

```go
func CustomFormatValidationErrors(valErrors validator.ValidationErrors) map[string]interface{} {
    // Start with standard formatting
    fieldErrors := validator.FormatValidationErrors(valErrors)
    
    // Add additional metadata
    result := map[string]interface{}{
        "validation_errors": fieldErrors,
        "error_count":       len(valErrors),
        "timestamp":         time.Now().Unix(),
    }
    
    return result
}
```