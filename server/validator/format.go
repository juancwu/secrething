package validator

import (
	"strconv"
	"strings"
)

// FormatValidationErrors formats ValidationErrors into a standardized map structure
// This is used by error handlers to format validation errors consistently
// It supports:
//   - Nested field paths like "user.address.street"
//   - Array/slice indexing like "items[0].name"
//   - Multi-dimensional arrays of any depth like "matrix[0][1]" or "cube[0][1][2]"
//   - Arrays of structs like "profile.addresses[0].street" with index information
//     Example: {"profile": { "addresses": [{"index": 0, "field_errors": {"street": "error"}}]}}
func FormatValidationErrors(valErrors ValidationErrors) map[string]interface{} {
	fieldErrors := make(map[string]interface{})

	for _, validationErr := range valErrors {
		setNestedErrorWithArrays(fieldErrors, validationErr.Field, validationErr.Message)
	}

	return fieldErrors
}

// setNestedErrorWithArrays sets a validation error message in a nested structure
// It handles both nested object paths and array indices
// For arrays of structs, it preserves the struct/object representation inside array elements
func setNestedErrorWithArrays(container map[string]interface{}, path string, message string) {
	// Split the path into segments
	segments := splitPathWithArrays(path)
	if len(segments) == 0 {
		return
	}

	// Handle single-segment case
	if len(segments) == 1 {
		segment := segments[0]
		// Check if it's an array index notation
		if isArraySegment(segment) {
			fieldName, index := parseArraySegment(segment)
			// Create or ensure the array exists
			ensureArrayExists(container, fieldName, index)
			array := container[fieldName].([]interface{})
			array[index] = message
		} else {
			container[segment] = message
		}
		return
	}

	// Process the path segments
	current := container
	for i := 0; i < len(segments)-1; i++ {
		segment := segments[i]
		nextSegment := segments[i+1]

		if isArraySegment(segment) {
			// Handle array segment
			fieldName, index := parseArraySegment(segment)

			// Create or ensure the array exists
			ensureArrayExists(current, fieldName, index)
			array := current[fieldName].([]interface{})

			// Check if the next segment is an array or object field
			if isArraySegment(nextSegment) {
				// Next segment is another array index - we need a new nested array here
				if array[index] == nil {
					array[index] = make([]interface{}, 0)
				} else if _, ok := array[index].(map[string]interface{}); ok {
					// Already a map, leave it as is
				} else if _, ok := array[index].([]interface{}); !ok {
					// Not already an array, create a new array
					array[index] = make([]interface{}, 0)
				}
			} else {
				// Next segment is an object field - we need a map with field_errors structure
				if array[index] == nil {
					// For structured arrays with objects, create with index and field_errors
					errMap := make(map[string]interface{})
					errMap["index"] = index
					errMap["field_errors"] = make(map[string]interface{})
					array[index] = errMap
				} else if mapVal, ok := array[index].(map[string]interface{}); ok {
					// Already a map, ensure it has the right structure
					if _, hasIndex := mapVal["index"]; !hasIndex {
						mapVal["index"] = index
					}
					if _, hasErrors := mapVal["field_errors"]; !hasErrors {
						mapVal["field_errors"] = make(map[string]interface{})
					}
				} else {
					// Not already a map, create a new structured map
					errMap := make(map[string]interface{})
					errMap["index"] = index
					errMap["field_errors"] = make(map[string]interface{})
					array[index] = errMap
				}
			}

			// Update current to point to the field_errors map in the array element
			if mapVal, ok := array[index].(map[string]interface{}); ok {
				if fieldErrors, ok := mapVal["field_errors"].(map[string]interface{}); ok {
					current = fieldErrors
				} else {
					mapVal["field_errors"] = make(map[string]interface{})
					current = mapVal["field_errors"].(map[string]interface{})
				}
			} else {
				// Cannot continue - expected a map but got something else
				return
			}
		} else {
			// Handle regular field segment
			if _, exists := current[segment]; !exists {
				if isArraySegment(nextSegment) {
					// Next segment is an array, initialize an empty map to hold the array
					current[segment] = make(map[string]interface{})
				} else {
					// Next segment is another object field
					current[segment] = make(map[string]interface{})
				}
			} else if _, ok := current[segment].(map[string]interface{}); !ok {
				// Field exists but is not a map - convert to map and keep previous value
				existingValue := current[segment]
				newMap := make(map[string]interface{})
				newMap["_value"] = existingValue
				current[segment] = newMap
			}

			// Move to the next level
			current = current[segment].(map[string]interface{})
		}
	}

	// Set the final value
	lastSegment := segments[len(segments)-1]

	if isArraySegment(lastSegment) {
		fieldName, index := parseArraySegment(lastSegment)
		ensureArrayExists(current, fieldName, index)
		array := current[fieldName].([]interface{})

		// Create a structured error for primitive array elements
		errMap := make(map[string]interface{})
		errMap["index"] = index
		errMap["message"] = message
		array[index] = errMap
	} else {
		current[lastSegment] = message
	}
}

// splitPathWithArrays splits a field path into segments, preserving array notation
// Example: "profile.addresses[0].street" -> ["profile", "addresses[0]", "street"]
func splitPathWithArrays(path string) []string {
	// First we'll split by dots to get the main segments
	dotParts := strings.Split(path, ".")

	var segments []string

	for _, part := range dotParts {
		if part == "" {
			continue
		}
		segments = append(segments, part)
	}

	return segments
}

// isArraySegment checks if a path segment contains array notation
// Example: "items[0]" -> true, "name" -> false
func isArraySegment(segment string) bool {
	return strings.Contains(segment, "[") && strings.Contains(segment, "]")
}

// parseArraySegment extracts the field name and index from an array segment
// Example: "items[0]" -> "items", 0
func parseArraySegment(segment string) (string, int) {
	parts := strings.SplitN(segment, "[", 2)
	fieldName := parts[0]

	// Extract the index
	indexStr := strings.TrimSuffix(parts[1], "]")
	index, _ := strconv.Atoi(indexStr)

	return fieldName, index
}

// ensureArrayExists creates or ensures an array exists at the specified field and is large enough
// For arrays of structs, it also ensures each element has 'index' and 'field_errors' fields
func ensureArrayExists(container map[string]interface{}, fieldName string, index int) {
	// Create the array if it doesn't exist
	if _, exists := container[fieldName]; !exists {
		array := make([]interface{}, index+1)
		container[fieldName] = array
	} else if existingArray, ok := container[fieldName].([]interface{}); ok {
		// Expand the array if needed
		if len(existingArray) <= index {
			newArray := make([]interface{}, index+1)
			copy(newArray, existingArray)
			container[fieldName] = newArray
		}
	} else {
		// Not an array, create a new one
		container[fieldName] = make([]interface{}, index+1)
	}
}
