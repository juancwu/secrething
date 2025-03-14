# Wildcard Pattern Matching

The Konbini validator includes a powerful wildcard pattern matching system for validation error messages, allowing you to set custom messages for array or slice elements without specifying exact indices.

## Overview

Wildcard pattern matching uses the special `[*]` syntax to represent "any index" within a field path. This allows creating generic error messages that apply to any element in an array or slice.

## Key Features

- **Wildcard syntax**: Use `[*]` to match any array index
- **Deep nesting support**: Apply patterns at any level of nested structures
- **Mixed specific/wildcard indices**: Combine specific indices with wildcards
- **Precedence handling**: Clear priority order for message resolution

## Basic Usage

### Setting Pattern Error Messages

```go
validator := validator.NewCustomValidator()

// Set a message for any item's name field
validator.Translator().SetPatternError("items[*].name", "required", "Every item must have a name")

// Set a message for any item's price
validator.Translator().SetPatternError("items[*].price", "gt", "Price must be greater than zero")
```

### Complex Patterns

```go
// Deeply nested wildcards
validator.Translator().SetPatternError(
    "departments[*].teams[*].members[*].email", 
    "email", 
    "All team member emails must be valid"
)

// Mixed specific and wildcard indices
validator.Translator().SetPatternError(
    "matrix[0][*][2]", 
    "gte", 
    "Values at position [0][*][2] must be greater than or equal to zero"
)
```

## Pattern Matching Rules

The pattern matching algorithm follows these rules:

1. **Path length**: The pattern must have the same number of segments as the field path
2. **Field name matching**: The field names (before brackets) must match exactly
3. **Array index matching**: 
   - `[*]` in the pattern matches any numeric index in the field path
   - Specific indices (e.g., `[0]`) must match exactly
4. **Wildcards per segment**: A segment can have multiple indices with individual wildcards

## Error Message Precedence

When multiple error message definitions could apply to a field, the validator uses the following precedence:

1. **Exact path match**: e.g., `items[0].name` (highest priority)
2. **Pattern match**: e.g., `items[*].name`
3. **Leaf field name**: e.g., `name`
4. **Default tag message**: e.g., default message for `required`
5. **Generic default message**: e.g., "Invalid value" (lowest priority)

## Examples

### Basic Array Validation

```go
type Item struct {
    Name     string  `json:"name" validate:"required"`
    Price    float64 `json:"price" validate:"required,gt=0"`
    Quantity int     `json:"quantity" validate:"required,gt=0"`
}

type Order struct {
    Items []Item `json:"items" validate:"required,dive"`
}

// Set wildcard pattern messages
validator.Translator().SetPatternError("items[*].name", "required", "Every item needs a name")
validator.Translator().SetPatternError("items[*].price", "gt", "All prices must be positive")

// Create an order with validation errors
order := Order{
    Items: []Item{
        {
            Name:     "", // Required error - will get wildcard message
            Price:    0,  // GT error - will get wildcard message
            Quantity: 5,
        },
        {
            Name:     "", // Required error - will get the same wildcard message
            Price:    -1, // GT error - will get the same wildcard message
            Quantity: 3,
        },
    },
}

// Validate
err := validator.Validate(&order)
```

### Demonstrating Precedence

```go
// Set messages with different specificity levels
validator.Translator().SetDefaultError("required", "Generic required field message")
validator.Translator().SetFieldError("name", "required", "Name is required (leaf)")
validator.Translator().SetPatternError("items[*].name", "required", "Every item needs a name (pattern)")
validator.Translator().SetFieldError("items[0].name", "required", "First item name is required (exact)")

// Create an order with two invalid items
order := Order{
    Items: []Item{
        { Name: "" }, // First item - will get "First item name is required (exact)"
        { Name: "" }, // Second item - will get "Every item needs a name (pattern)"
    },
}
```

## Validation Context Support

Wildcard patterns can be scoped to specific validation contexts:

```go
// Create a validation context
ctx := validator.NewValidationContext(baseValidator)
ctx.SetPatternError("items[*].name", "required", "Context: Item name is required")
ctx.SetPatternError("items[*].price", "gt", "Context: Item price must be positive")

// Validate with context
err := ctx.Validate(&order)
```

## Implementation Details

### Key Functions

- **SetPatternError**: Stores a pattern-based custom error message
- **matchPattern**: Finds a matching pattern for a field path
- **matchesPattern**: Compares segments to determine if a pattern matches
- **splitPathForPatternMatch**: Splits a path into segments for matching
- **extractIndices**: Extracts array indices from a path segment
- **isNumericIndex**: Checks if an index is numeric (not a wildcard)

### Pattern Storage

Pattern errors are stored in the `patternErrors` map within the `ErrorTranslator` struct, similar to field-specific errors.

## Best Practices

- Use wildcard patterns for consistency across array elements
- For special cases, combine with exact path messages (higher precedence)
- Keep patterns clear and aligned with your data structure
- Consider using contexts for request-specific pattern messages
