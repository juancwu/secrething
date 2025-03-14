# Validation Contexts

Validation contexts provide a powerful way to create request-specific validation rules without modifying your base validator. This is particularly useful for different HTTP endpoints that require different validation behavior for the same data types.

## Overview

A validation context is a lightweight copy of your base validator with its own error translator. It allows you to:

- Set request-specific error messages
- Override default validation behavior for specific fields
- Keep your base validator clean and reusable
- Apply context-specific rules only where needed

## Creating a Validation Context

Start with a base validator and create a context from it:

```go
// Create your base validator
baseValidator := validator.NewCustomValidator()

// Create a validation context
context := validator.NewValidationContext(baseValidator)
```

The context will contain a clone of the base validator's translator, so any changes you make to the context won't affect the base validator.

## Customizing a Context

You can customize a validation context using a fluent API:

```go
// Chain method calls for concise setup
context.SetFieldError("username", "required", "Please enter your username")
       .SetFieldError("email", "email", "Invalid email address format")
       .SetDefaultError("required", "This field must be provided")
       .SetDefaultMessage("Validation error occurred")
```

## Using a Context for Validation

Once configured, use the context to validate your structs:

```go
// Validate directly with context
err := context.Validate(&myStruct)

// Or with Echo framework
err := validator.BindAndValidateWithContext(c, &myStruct, context)
```

## Practical Examples

### Different Validation Rules by Endpoint

```go
// Base validator with common rules
baseValidator := validator.NewCustomValidator()
baseValidator.Translator().SetDefaultError("required", "Field is required")
baseValidator.Translator().SetDefaultError("email", "Invalid email")

// Context for user registration
registerContext := validator.NewValidationContext(baseValidator)
registerContext.SetFieldError("password", "password", "Password must be secure with uppercase, lowercase, numbers, and special characters")
registerContext.SetFieldError("email", "email", "Please provide a valid email address for account verification")

// Context for user login
loginContext := validator.NewValidationContext(baseValidator)
loginContext.SetFieldError("password", "required", "Password is required to log in")
loginContext.SetFieldError("email", "email", "Please enter a valid email address")

// Context for password reset
resetContext := validator.NewValidationContext(baseValidator)
resetContext.SetFieldError("email", "email", "Please enter the email address associated with your account")
```

### Localization Support

Create contexts for different languages:

```go
// Create base validator
baseValidator := validator.NewCustomValidator()

// English context
enContext := validator.NewValidationContext(baseValidator)
enContext.SetDefaultError("required", "This field is required")
       .SetDefaultError("email", "Please enter a valid email address")
       .SetFieldError("password", "password", "Password must be secure")

// Spanish context
esContext := validator.NewValidationContext(baseValidator)
esContext.SetDefaultError("required", "Este campo es obligatorio")
       .SetDefaultError("email", "Por favor ingrese un correo electrónico válido")
       .SetFieldError("password", "password", "La contraseña debe ser segura")

// Select context based on user language preference
func getValidationContext(locale string) *validator.ValidationContext {
    switch locale {
    case "es":
        return esContext
    default:
        return enContext
    }
}

// Use in handler
func createUser(c echo.Context) error {
    locale := c.Request().Header.Get("Accept-Language")
    validationContext := getValidationContext(locale)
    
    user := new(User)
    if err := validator.BindAndValidateWithContext(c, user, validationContext); err != nil {
        return err
    }
    
    // Process valid user...
    return c.JSON(http.StatusCreated, user)
}
```

### Different Rules by Role

Apply different validation rules based on user roles:

```go
// Base validator
baseValidator := validator.NewCustomValidator()

// Admin context - more permissive
adminContext := validator.NewValidationContext(baseValidator)
adminContext.SetFieldError("accessLevel", "range", "Access level must be between 0 and 10")

// User context - more restrictive
userContext := validator.NewValidationContext(baseValidator)
userContext.SetFieldError("accessLevel", "range", "Access level must be between 0 and 5")

// Choose context based on user role
func getContextByRole(role string) *validator.ValidationContext {
    switch role {
    case "admin":
        return adminContext
    default:
        return userContext
    }
}

// In handler
func updatePermissions(c echo.Context) error {
    user := getCurrentUser(c)
    context := getContextByRole(user.Role)
    
    permissions := new(Permissions)
    if err := validator.BindAndValidateWithContext(c, permissions, context); err != nil {
        return err
    }
    
    // Process valid permissions...
    return c.JSON(http.StatusOK, permissions)
}
```

## Pattern Error Support in Contexts

Validation contexts also support wildcard pattern errors:

```go
// Create context for order validation
orderContext := validator.NewValidationContext(baseValidator)
orderContext.SetPatternError("items[*].name", "required", "All items must have a name")
orderContext.SetPatternError("items[*].quantity", "gt", "All quantities must be greater than zero")
```

## Implementation Details

Validation contexts maintain their own copy of the error translator, allowing independent customization:

```go
// NewValidationContext creates a new validation context with a cloned validator
func NewValidationContext(baseValidator *CustomValidator) *ValidationContext {
    cloned := baseValidator.Clone()
    return &ValidationContext{
        validator:  cloned,
        translator: cloned.translator,
    }
}
```

The `ValidationContext` struct is designed for fluent API usage:

```go
// ValidationContext stores context-specific validation settings
type ValidationContext struct {
    validator  *CustomValidator
    translator *ErrorTranslator
}

// SetFieldError sets a custom error message for a specific field and validation tag
func (vc *ValidationContext) SetFieldError(field, tag, message string) *ValidationContext {
    vc.translator.SetFieldError(field, tag, message)
    return vc
}

// Additional methods follow the same pattern...
```

## Best Practices

1. **Create a solid base validator** with sensible default messages
2. **Create contexts for specific endpoints** or use cases
3. **Chain method calls** for concise, readable context configuration
4. **Organize contexts by feature area** for better maintainability
5. **Consider performance** - contexts clone the translator so avoid creating them inside request handlers
6. **Reuse contexts** where possible - create them once at startup

## Advanced: Context Factory Pattern

For applications with many contexts, consider using a factory pattern:

```go
// ContextFactory manages validation contexts
type ContextFactory struct {
    baseValidator *validator.CustomValidator
    contexts      map[string]*validator.ValidationContext
}

// NewContextFactory creates a new context factory
func NewContextFactory(baseValidator *validator.CustomValidator) *ContextFactory {
    return &ContextFactory{
        baseValidator: baseValidator,
        contexts:      make(map[string]*validator.ValidationContext),
    }
}

// RegisterContext creates and registers a named validation context
func (f *ContextFactory) RegisterContext(name string, configurator func(*validator.ValidationContext)) {
    context := validator.NewValidationContext(f.baseValidator)
    configurator(context)
    f.contexts[name] = context
}

// GetContext retrieves a registered validation context by name
func (f *ContextFactory) GetContext(name string) *validator.ValidationContext {
    if context, exists := f.contexts[name]; exists {
        return context
    }
    
    // Return default context if named context doesn't exist
    return validator.NewValidationContext(f.baseValidator)
}

// Usage
factory := NewContextFactory(baseValidator)

// Register contexts
factory.RegisterContext("login", func(vc *validator.ValidationContext) {
    vc.SetFieldError("email", "email", "Please enter a valid email address")
    vc.SetFieldError("password", "required", "Password is required")
})

factory.RegisterContext("register", func(vc *validator.ValidationContext) {
    vc.SetFieldError("email", "email", "Please provide a valid email for verification")
    vc.SetFieldError("password", "password", "Create a strong password")
    vc.SetFieldError("username", "required", "Choose a username")
})

// In handlers
func loginHandler(c echo.Context) error {
    context := factory.GetContext("login")
    credentials := new(LoginCredentials)
    
    if err := validator.BindAndValidateWithContext(c, credentials, context); err != nil {
        return err
    }
    
    // Process valid login...
}
```