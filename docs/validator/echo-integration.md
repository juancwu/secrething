# Echo Framework Integration

The Konbini validator provides seamless integration with the [Echo](https://echo.labstack.com/) web framework, making request validation straightforward and consistent.

## Setup

### Registering the Validator

First, register the validator with your Echo instance:

```go
import (
    "github.com/juancwu/konbini/server/validator"
    "github.com/labstack/echo/v4"
)

func setupEcho() *echo.Echo {
    e := echo.New()
    
    // Create and register the validator
    v := validator.NewCustomValidator()
    e.Validator = v
    
    // Configure custom error messages as needed
    v.Translator().SetFieldError("email", "email", "Please enter a valid email address")
    v.Translator().SetDefaultError("required", "This field is required")
    
    return e
}
```

## Basic Usage

### BindAndValidate Helper

For simple cases, use the `BindAndValidate` helper:

```go
func createUser(c echo.Context) error {
    user := new(User)
    
    if err := validator.BindAndValidate(c, user); err != nil {
        // Handle validation errors
        return err // Echo will invoke your error handler
    }
    
    // Process valid user data
    return c.JSON(http.StatusCreated, user)
}
```

The `BindAndValidate` function:

1. Binds the request body to your struct using Echo's `Bind` method
2. Validates the struct using the registered validator
3. Returns any validation errors in a standardized format

## Context-Specific Validation

For cases where you need different validation rules for different endpoints:

### Creating a Validation Context

```go
// Create base validator
baseValidator := validator.NewCustomValidator()

// Create a context for registration validation
registerContext := validator.NewValidationContext(baseValidator)
registerContext.SetFieldError("password", "password", "Password must contain uppercase, lowercase, numbers, and special characters")
registerContext.SetFieldError("email", "email", "Please use a valid email address for registration")

// Create a context for login validation
loginContext := validator.NewValidationContext(baseValidator)
loginContext.SetFieldError("email", "email", "Invalid login credentials")
```

### Using a Validation Context

```go
func registerUser(c echo.Context) error {
    user := new(User)
    
    if err := validator.BindAndValidateWithContext(c, user, registerContext); err != nil {
        return err
    }
    
    // Process valid registration
    return c.JSON(http.StatusCreated, user)
}

func loginUser(c echo.Context) error {
    credentials := new(LoginCredentials)
    
    if err := validator.BindAndValidateWithContext(c, credentials, loginContext); err != nil {
        return err
    }
    
    // Process valid login
    return c.JSON(http.StatusOK, map[string]interface{}{"token": "..."})
}
```

## Error Handling

To provide consistent error responses, configure Echo's error handler:

```go
e.HTTPErrorHandler = func(err error, c echo.Context) {
    var validationErrors validator.ValidationErrors
    var statusCode int = http.StatusInternalServerError
    var message string = "Internal server error"
    
    if errors.As(err, &validationErrors) {
        // Handle validation errors
        statusCode = http.StatusBadRequest
        message = "Validation failed"
        
        // Use the FormatValidationErrors function for consistent formatting
        fieldErrors := validator.FormatValidationErrors(validationErrors)
        
        c.JSON(statusCode, map[string]interface{}{
            "code":         statusCode,
            "message":      message,
            "field_errors": fieldErrors,
        })
        return
    }
    
    // Handle other error types
    if httpErr, ok := err.(*echo.HTTPError); ok {
        statusCode = httpErr.Code
        message = fmt.Sprintf("%v", httpErr.Message)
    }
    
    c.JSON(statusCode, map[string]interface{}{
        "code":    statusCode,
        "message": message,
    })
}
```

## Example: Complete Request Handler

```go
type CreateProductRequest struct {
    Name        string   `json:"name" validate:"required"`
    Description string   `json:"description" validate:"required"`
    Price       float64  `json:"price" validate:"required,gt=0"`
    Categories  []string `json:"categories" validate:"required,dive,required"`
    SKU         string   `json:"sku" validate:"required,alphanum"`
}

func setupProductRoutes(e *echo.Echo) {
    // Create a validator with custom messages
    v := validator.NewCustomValidator()
    v.Translator().SetFieldError("name", "required", "Product name is required")
    v.Translator().SetFieldError("price", "gt", "Price must be greater than zero")
    v.Translator().SetPatternError("categories[*]", "required", "All categories must be specified")
    
    // Create a product context with specific rules
    productContext := validator.NewValidationContext(v)
    productContext.SetFieldError("sku", "alphanum", "SKU must contain only letters and numbers")
    
    // Register routes
    products := e.Group("/products")
    products.POST("", createProduct(productContext))
}

func createProduct(vc *validator.ValidationContext) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Create request struct
        req := new(CreateProductRequest)
        
        // Bind and validate with context
        if err := validator.BindAndValidateWithContext(c, req, vc); err != nil {
            return err
        }
        
        // Process valid product creation
        product := saveProduct(req)
        
        // Return success response
        return c.JSON(http.StatusCreated, map[string]interface{}{
            "message": "Product created successfully",
            "product": product,
        })
    }
}

func saveProduct(req *CreateProductRequest) *Product {
    // Implementation of product saving logic
    return &Product{
        ID:          generateID(),
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Categories:  req.Categories,
        SKU:         req.SKU,
        CreatedAt:   time.Now(),
    }
}
```

## Best Practices

1. **Reuse validation contexts** where possible to minimize duplication
2. **Group related validations** in shared contexts
3. **Use meaningful error messages** that guide users to correct their input
4. **Structure request structs** to match exactly what each endpoint needs
5. **Return consistent error responses** using the validation error formatter
6. **Document your validation rules** for frontend developers