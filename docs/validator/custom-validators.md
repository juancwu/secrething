# Custom Validators

The Konbini validator package allows you to extend its validation capabilities with custom validation rules. This document covers how to use the built-in custom validators and how to create your own.

## Built-in Custom Validators

### Password Validator

The package includes a built-in password validator that checks for strong passwords:

```go
type User struct {
    Password string `json:"password" validate:"required,password"`
}
```

The `password` validation tag ensures passwords meet these requirements:
- At least 8 characters long
- Contains at least one uppercase letter
- Contains at least one lowercase letter
- Contains at least one digit
- Contains at least one special character from `!@#$%^&*()-_=+[]{}|;:'",.<>/?`

#### Implementation

The password validator is implemented in `password.go`:

```go
func validatePassword(fl govalidator.FieldLevel) bool {
    password := fl.Field().String()
    
    if len(password) < 8 {
        return false
    }
    
    var (
        hasUpper   bool
        hasLower   bool
        hasDigit   bool
        hasSpecial bool
    )
    
    specialChars := `!@#$%^&*()-_=+[]{}|;:'",.<>/?`
    
    for _, char := range password {
        switch {
        case 'A' <= char && char <= 'Z':
            hasUpper = true
        case 'a' <= char && char <= 'z':
            hasLower = true
        case '0' <= char && char <= '9':
            hasDigit = true
        case strings.ContainsRune(specialChars, char):
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasDigit && hasSpecial
}
```

#### Custom Error Message

You can set a custom error message for password validation:

```go
validator.Translator().SetDefaultError("password", "Password must meet security requirements")
validator.Translator().SetFieldError("user.password", "password", "User password must be secure")
```

## Creating Custom Validators

You can create your own custom validators to extend the validation capabilities.

### Basic Custom Validator

First, define a validation function that returns a boolean:

```go
// validateIsbn checks if a string is a valid ISBN (example)
func validateIsbn(fl govalidator.FieldLevel) bool {
    isbn := fl.Field().String()
    
    // Remove dashes and spaces
    isbn = strings.ReplaceAll(isbn, "-", "")
    isbn = strings.ReplaceAll(isbn, " ", "")
    
    // Check length (ISBN-10 or ISBN-13)
    if len(isbn) != 10 && len(isbn) != 13 {
        return false
    }
    
    // Example ISBN-10 check (simplified)
    if len(isbn) == 10 {
        sum := 0
        for i := 0; i < 9; i++ {
            digit := int(isbn[i] - '0')
            sum += digit * (10 - i)
        }
        
        // Check digit can be 'X' for 10
        var checkDigit int
        if isbn[9] == 'X' || isbn[9] == 'x' {
            checkDigit = 10
        } else {
            checkDigit = int(isbn[9] - '0')
        }
        
        return (sum+checkDigit) % 11 == 0
    }
    
    // Example ISBN-13 check (simplified)
    if len(isbn) == 13 {
        sum := 0
        for i := 0; i < 12; i++ {
            digit := int(isbn[i] - '0')
            if i % 2 == 0 {
                sum += digit
            } else {
                sum += digit * 3
            }
        }
        
        checkDigit := int(isbn[12] - '0')
        return (10 - (sum % 10)) % 10 == checkDigit
    }
    
    return false
}
```

### Registering a Custom Validator

Register your custom validator with the validator instance:

```go
func initValidator() *validator.CustomValidator {
    v := validator.NewCustomValidator()
    
    // Get access to the underlying go-playground validator
    playgroundValidator := v.GetUnderlyingValidator()
    
    // Register your custom validator
    playgroundValidator.RegisterValidation("isbn", validateIsbn)
    
    // Set custom error message
    v.Translator().SetDefaultError("isbn", "Must be a valid ISBN number")
    
    return v
}
```

### Using the Custom Validator

Use your custom validator in struct tags:

```go
type Book struct {
    Title       string `json:"title" validate:"required"`
    Author      string `json:"author" validate:"required"`
    ISBN        string `json:"isbn" validate:"required,isbn"`
    PublishYear int    `json:"publishYear" validate:"required,gt=1900"`
}
```

### Custom Validators with Parameters

You can create validators that accept parameters:

```go
// validateCreditCard validates different types of credit cards
// Usage: `validate:"credit_card=visa|mastercard"`
func validateCreditCard(fl govalidator.FieldLevel) bool {
    cardNumber := fl.Field().String()
    cardTypes := strings.Split(fl.Param(), "|")
    
    // Remove spaces and dashes
    cardNumber = strings.ReplaceAll(cardNumber, " ", "")
    cardNumber = strings.ReplaceAll(cardNumber, "-", "")
    
    // Check if all digits
    for _, c := range cardNumber {
        if c < '0' || c > '9' {
            return false
        }
    }
    
    // Basic Luhn algorithm check
    valid := luhnCheck(cardNumber)
    if !valid {
        return false
    }
    
    // Check card type based on prefix and length
    for _, cardType := range cardTypes {
        switch strings.ToLower(cardType) {
        case "visa":
            if strings.HasPrefix(cardNumber, "4") && (len(cardNumber) == 13 || len(cardNumber) == 16) {
                return true
            }
        case "mastercard":
            if strings.HasPrefix(cardNumber, "5") && len(cardNumber) == 16 {
                return true
            }
        case "amex":
            if (strings.HasPrefix(cardNumber, "34") || strings.HasPrefix(cardNumber, "37")) && len(cardNumber) == 15 {
                return true
            }
        }
    }
    
    return false
}

// Register the parameterized validator
playgroundValidator.RegisterValidation("credit_card", validateCreditCard)
```

### Validator Registration in Different Packages

If you're adding validators from a different package, make a helper function:

```go
// RegisterCustomValidators registers all custom validators with the provided validator
func RegisterCustomValidators(v *validator.CustomValidator) {
    playgroundValidator := v.GetUnderlyingValidator()
    
    // Register custom validations
    playgroundValidator.RegisterValidation("isbn", validateIsbn)
    playgroundValidator.RegisterValidation("credit_card", validateCreditCard)
    
    // Set default error messages
    v.Translator().SetDefaultError("isbn", "Must be a valid ISBN number")
    v.Translator().SetDefaultError("credit_card", "Must be a valid credit card number of the specified type(s)")
}

// In your main package
func main() {
    v := validator.NewCustomValidator()
    customvalidators.RegisterCustomValidators(v)
    // ...
}
```

## Best Practices

1. **Keep validators focused**: Each validator should check one specific thing
2. **Provide clear error messages**: Make error messages descriptive
3. **Optimize for performance**: Validation runs on every request
4. **Test extensively**: Add unit tests for all custom validators
5. **Document validation rules**: Add comments explaining the validation logic
6. **Consider security implications**: For sensitive validations like passwords
7. **Use parameters for flexibility**: Create configurable validators where appropriate

## Examples

### Email Domain Validator

```go
// validateEmailDomain checks if an email uses an allowed domain
// Usage: `validate:"email,email_domain=company.com|approved-vendor.com"`
func validateEmailDomain(fl govalidator.FieldLevel) bool {
    email := fl.Field().String()
    allowedDomains := strings.Split(fl.Param(), "|")
    
    // Split email at @ symbol
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false // Not a valid email format
    }
    
    domain := parts[1]
    
    // Check if domain is in allowed list
    for _, allowedDomain := range allowedDomains {
        if domain == allowedDomain {
            return true
        }
    }
    
    return false
}
```

### File Extension Validator

```go
// validateFileExt checks if a filename has an allowed extension
// Usage: `validate:"file_ext=jpg|png|gif"`
func validateFileExt(fl govalidator.FieldLevel) bool {
    filename := fl.Field().String()
    allowedExts := strings.Split(fl.Param(), "|")
    
    // Get the file extension
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return false // No extension
    }
    
    ext := strings.ToLower(parts[len(parts)-1])
    
    // Check if extension is allowed
    for _, allowedExt := range allowedExts {
        if ext == strings.ToLower(allowedExt) {
            return true
        }
    }
    
    return false
}
```

### Phone Number Validator

```go
// validatePhone checks if a string is a valid phone number in a specific format
// Usage: `validate:"phone=us"`
func validatePhone(fl govalidator.FieldLevel) bool {
    phone := fl.Field().String()
    format := fl.Param()
    
    // Remove common formatting characters
    phone = strings.ReplaceAll(phone, " ", "")
    phone = strings.ReplaceAll(phone, "-", "")
    phone = strings.ReplaceAll(phone, "(", "")
    phone = strings.ReplaceAll(phone, ")", "")
    phone = strings.ReplaceAll(phone, "+", "")
    
    switch format {
    case "us":
        // US phone numbers: 10 digits
        if len(phone) != 10 {
            return false
        }
        for _, c := range phone {
            if c < '0' || c > '9' {
                return false
            }
        }
        return true
        
    case "international":
        // International format: at least 8 digits, all numeric
        if len(phone) < 8 {
            return false
        }
        for _, c := range phone {
            if c < '0' || c > '9' {
                return false
            }
        }
        return true
    }
    
    return false
}
```