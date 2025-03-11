# Request Validator Middleware

This package provides middleware for validating request bodies in Echo-based APIs.

## Features

- Integrates with Echo framework
- Validates request bodies using validator.v10
- Returns validation errors in a structured format that matches request body fields
- Default error messages for common validation rules
- Support for custom error messages
- API for adding custom validation rules

## Usage

### Basic Usage

```go
// Register the middleware globally
e := echo.New()
e.Use(middleware.RequestValidator())

// In your handler
func createUser(c echo.Context) error {
    // Define request struct with validation tags
    type UserRequest struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
    }
    
    // Create and validate
    user := new(UserRequest)
    if err := middleware.Validate(c, user); err != nil {
        return err // Let the global error handler process this
    }
    
    // Process validated request...
    return c.JSON(200, map[string]string{"status": "Success"})
}
```

### Custom Error Messages

```go
// Add custom error messages for specific validation rules
customMessages := map[string]string{
    "name.required": "User name is required",
    "email.required": "Email address is required",
    "email.email": "Please provide a valid email address"
}

// Apply to a specific route group
userGroup := e.Group("/users")
userGroup.Use(middleware.WithCustomMessages(customMessages))
```

### Custom Validation Rules

```go
// Register a custom validator with a default message
middleware.RegisterCustomValidator("strong_password", validateStrongPassword, 
    "Password must be at least 8 characters with uppercase, lowercase, number, and special character")

// Example validator implementation
func validateStrongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    // Password validation logic here...
    return isValid
}
```

## Error Response Format

When validation fails, the response will have this structure:

```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": {
    "name": "This field is required",
    "email": "Invalid email format"
  },
  "req_id": "8f7d9fd32..."
}
```

This format follows your application's error handling convention. The errors object maps each field to its specific validation error message, making it easy for clients to display validation feedback next to the appropriate form fields.

## Available Validation Rules

This middleware uses the validation rules provided by the go-playground/validator.v10 package. 
Here are some common validation rules:

- `required`: Field is required
- `email`: Field must be a valid email
- `min=X`: Minimum value (for numbers) or length (for strings)
- `max=X`: Maximum value or length
- `len=X`: Exact length required
- `alphanum`: Must be alphanumeric
- `numeric`: Must be numeric
- `url`: Must be a valid URL
- `uuid`: Must be a valid UUID