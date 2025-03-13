package validator

import (
	"strings"
)

// FormatValidationErrors formats ValidationErrors into a standardized map structure
// This is used by error handlers to format validation errors consistently
// It supports nested field paths like "user.address.street"
func FormatValidationErrors(valErrors ValidationErrors) map[string]interface{} {
	fieldErrors := make(map[string]interface{})

	for _, validationErr := range valErrors {
		// Check if the field path contains dots indicating nested structure
		if strings.Contains(validationErr.Field, ".") {
			setNestedField(fieldErrors, validationErr.Field, validationErr.Message)
		} else {
			fieldErrors[validationErr.Field] = validationErr.Message
		}
	}

	return fieldErrors
}

// setNestedField sets a value in a nested map structure based on a dot-separated path
// Example: "user.address.street" -> map[user]map[address]map[street]message
func setNestedField(m map[string]interface{}, path string, value string) {
	parts := strings.Split(path, ".")

	// Handle empty path case
	if len(parts) == 0 {
		return
	}

	lastIndex := len(parts) - 1

	// Handle single-part path (no dots)
	if lastIndex <= 0 {
		m[path] = value
		return
	}

	// Navigate to the final container
	current := m
	for _, part := range parts[:lastIndex] {
		// Skip empty parts
		if part == "" {
			continue
		}

		// If this part doesn't exist yet, create it
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}

		// If it's not a map, we can't continue (field name collision)
		next, ok := current[part].(map[string]interface{})
		if !ok {
			// Field exists but is not a map - convert scalar to map and keep previous value
			nextMap := make(map[string]interface{})
			nextMap["_value"] = current[part] // Store the scalar value under "_value" key
			current[part] = nextMap
			next = nextMap
		}

		current = next
	}

	// Set the final value if lastIndex part is not empty
	if parts[lastIndex] != "" {
		current[parts[lastIndex]] = value
	}
}
