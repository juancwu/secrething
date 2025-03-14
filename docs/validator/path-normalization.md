# Path Normalization

The Konbini validator includes a path normalization system that ensures consistent handling of field paths, especially for nested structures, arrays, and maps.

## Overview

Path normalization standardizes different representations of the same field path to ensure error message lookup works correctly regardless of formatting variations.

### Benefits

- **Consistent lookup**: Ensures custom error messages can be found during validation
- **Flexible syntax**: Allows developers to use their preferred formatting style
- **Reliable matching**: Handles spaces, array indices, and nested paths consistently

## How Path Normalization Works

The path normalization process:

1. Trims whitespace from the path
2. Handles whitespace around dots and brackets
3. Standardizes array/map notation
4. Preserves the semantic structure of the path

### Examples

All of these paths are normalized to the same internal representation:

```
user.profile.addresses[0].street
user . profile . addresses[0].street
user.profile.addresses [ 0 ].street
```

## Usage

Path normalization happens automatically when:

1. Setting custom error messages:
   ```go
   validator.SetFieldError("user.addresses [ 0 ].city", "required", "City is required")
   ```

2. Looking up error messages during validation:
   ```go
   message := validator.Translate("user.addresses[0].city", "required")
   ```

## Implementation Details

### normalizePath Function

The `normalizePath` function handles standardization of field paths:

```go
// normalizePath standardizes field paths for consistent format handling
// - Handles array indices consistently: profile.addresses[0].street
// - Handles map keys: profile.metadata["key"].value
func normalizePath(path string) string {
    // Implementation...
}
```

### Path Format Rules

- Dot (`.`) separates field segments
- Square brackets (`[]`) enclose array indices or map keys
- Whitespace around dots and brackets is removed
- Array indices remain numeric
- Map keys preserve their format, including quotes

## Integration with Other Features

Path normalization works in conjunction with:

- **Leaf name extraction**: For extracting the final field name from a path
- **Wildcard pattern matching**: For handling wildcards in array indices
- **Error message lookup**: For finding the best matching error message

## Best Practices

- Use consistent path formatting in your code for readability
- When setting custom messages, use the most natural format
- For complex paths with arrays and maps, use a clear, readable format

## Example: Complete Path Normalization 

```go
// These different formats are all normalized to the same internal representation
validator.SetFieldError("user.addresses[0].contacts[\"primary\"].email", "email", "Invalid email format")
validator.SetFieldError("user . addresses [ 0 ] . contacts [ \"primary\" ] . email", "email", "Invalid email format")
```

During validation, regardless of how the path is formatted in the validation error, the custom message will be correctly applied due to path normalization.
