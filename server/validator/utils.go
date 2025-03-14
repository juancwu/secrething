package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

// getJSONFieldPath returns the JSON field path for a validation error
// For nested fields, it returns a dot-separated path like "profile.firstName"
// For array/slice fields, it returns indexed paths like "items[0].name"
func getJSONFieldPath(obj interface{}, fieldError govalidator.FieldError) string {
	// Build the namespace path based on JSON field names rather than struct field names
	namespace := fieldError.Namespace()
	parts := strings.Split(namespace, ".")

	// The first part is the type name, so we skip it
	parts = parts[1:]

	// Build a new path with JSON names
	var jsonParts []string
	currentObj := obj
	currentType := reflect.TypeOf(currentObj).Elem()
	currentValue := reflect.ValueOf(currentObj).Elem()

	for i, part := range parts {
		// Check if this part refers to an array/slice index
		indexMatch := regexp.MustCompile(`^(\w+)\[(\d+)\]$`).FindStringSubmatch(part)

		if len(indexMatch) == 3 {
			// This is an array/slice index reference like 'Items[0]'
			fieldName := indexMatch[1]
			indexStr := indexMatch[2]
			index, _ := strconv.Atoi(indexStr)

			// Find the struct field for the array/slice
			field, found := currentType.FieldByName(fieldName)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name for the array/slice field
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				jsonName = fieldName
			}

			// Add the field name and index to the path
			jsonParts = append(jsonParts, fmt.Sprintf("%s[%d]", jsonName, index))

			// Update currentType/currentValue for the next iteration
			fieldValue := currentValue.FieldByName(fieldName)
			if !fieldValue.IsValid() || index >= fieldValue.Len() {
				// If the field value is invalid or index is out of bounds, we can't continue
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}

			// Get the element at the specified index
			elemValue := fieldValue.Index(index)

			// Update currentType and currentValue based on the element type
			if elemValue.Kind() == reflect.Struct {
				currentType = elemValue.Type()
				currentValue = elemValue
			} else if elemValue.Kind() == reflect.Ptr && elemValue.Elem().Kind() == reflect.Struct {
				currentType = elemValue.Elem().Type()
				currentValue = elemValue.Elem()
			} else {
				// If the element is not a struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		} else {
			// Regular struct field
			field, found := currentType.FieldByName(part)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				// If there's no JSON tag or it's "-", use the original field name
				jsonParts = append(jsonParts, part)
			} else {
				jsonParts = append(jsonParts, jsonName)
			}

			// Update currentType and currentValue for the next iteration
			fieldValue := currentValue.FieldByName(part)

			if field.Type.Kind() == reflect.Struct {
				currentType = field.Type
				if fieldValue.IsValid() {
					currentValue = fieldValue
				}
			} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				currentType = field.Type.Elem()
				if fieldValue.IsValid() && !fieldValue.IsNil() {
					currentValue = fieldValue.Elem()
				}
			} else {
				// If the field is not a struct or pointer to struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		}
	}

	return strings.Join(jsonParts, ".")
}

// normalizePath standardizes field paths for consistent format handling
// - Handles array indices consistently: profile.addresses[0].street
// - Handles map keys: profile.metadata["key"].value
// - Also handles the new empty bracket notation: items[].name
// The path is normalized to ensure consistent lookup regardless of source format
func normalizePath(path string) string {
	// Trim leading/trailing spaces
	path = strings.TrimSpace(path)

	// Handle empty path case
	if path == "" {
		return ""
	}

	// Replace spaces around dots with just dots
	path = strings.ReplaceAll(path, " .", ".")
	path = strings.ReplaceAll(path, ". ", ".")

	// Split by dots to handle each segment
	segments := strings.Split(path, ".")
	normalizedSegments := make([]string, 0, len(segments))

	for _, segment := range segments {
		// Trim spaces from segment
		segment = strings.TrimSpace(segment)

		// Skip empty segments
		if segment == "" {
			continue
		}

		// Normalize array/map notation if present
		if strings.Contains(segment, "[") && strings.Contains(segment, "]") {
			// Extract the field name part (before the bracket)
			fieldName := segment
			if idx := strings.Index(segment, "["); idx > 0 {
				fieldName = strings.TrimSpace(segment[:idx])
			}

			// Extract all index/key parts but normalize them
			normalizedIndices := make([]string, 0)
			remaining := segment
			for strings.Contains(remaining, "[") && strings.Contains(remaining, "]") {
				start := strings.Index(remaining, "[")
				end := strings.Index(remaining, "]")

				if start >= 0 && end > start {
					// Extract the content between brackets
					indexContent := strings.TrimSpace(remaining[start+1 : end])

					// Create normalized index with no spaces
					var normalizedIndex string
					if indexContent == "" {
						// Keep [] notation for empty brackets
						normalizedIndex = "[]"
					} else {
						normalizedIndex = "[" + indexContent + "]"
					}

					normalizedIndices = append(normalizedIndices, normalizedIndex)

					// Move past this index/key
					if end+1 < len(remaining) {
						remaining = remaining[end+1:]
					} else {
						remaining = ""
					}
				} else {
					break
				}
			}

			// Reconstruct the segment with normalized field name and indices
			normalized := fieldName + strings.Join(normalizedIndices, "")
			normalizedSegments = append(normalizedSegments, normalized)
		} else {
			// Regular field segment without array/map notation
			normalizedSegments = append(normalizedSegments, segment)
		}
	}

	// Join segments with dots
	return strings.Join(normalizedSegments, ".")
}

// extractLeafName extracts just the leaf field name from a path
// For example:
// - "profile.firstName" returns "firstName"
// - "items[0].name" returns "name"
// - "data.points[0][1]" returns "points[0][1]"
func extractLeafName(path string) string {
	path = strings.TrimSpace(path)

	if path == "" {
		return ""
	}

	// If there are no dots, it's already a leaf
	if !strings.Contains(path, ".") {
		return path
	}

	// Split by dots and take the last segment
	segments := strings.Split(path, ".")
	return segments[len(segments)-1]
}

// matchArrayPath checks if a concrete path matches an array pattern with [] notation
// For example, "items[0].name" would match "items[].name"
func matchArrayPath(concretePath, patternPath string) bool {
	// Split both paths by dots
	concreteSegments := strings.Split(concretePath, ".")
	patternSegments := strings.Split(patternPath, ".")

	// Both must have the same number of segments
	if len(concreteSegments) != len(patternSegments) {
		return false
	}

	// Compare each segment
	for i, patternSegment := range patternSegments {
		concreteSegment := concreteSegments[i]

		// If pattern segment contains [] notation
		if strings.Contains(patternSegment, "[]") {
			// Extract field name (part before any brackets)
			patternField := patternSegment
			concreteField := concreteSegment

			if idx := strings.Index(patternSegment, "["); idx >= 0 {
				patternField = patternSegment[:idx]
			}

			if idx := strings.Index(concreteSegment, "["); idx >= 0 {
				concreteField = concreteSegment[:idx]
			}

			// Field names must match
			if patternField != concreteField {
				return false
			}

			// Create regex pattern to match the array notation
			// Replace [] with a regex that matches any numeric index [0], [1], etc.
			regexPattern := "^" + regexp.QuoteMeta(patternSegment)
			regexPattern = strings.ReplaceAll(regexPattern, "\\[\\]", "\\[(\\d+)\\]")

			// Check if the concrete segment matches the regex pattern
			matched, err := regexp.MatchString(regexPattern, concreteSegment)
			if err != nil || !matched {
				return false
			}
		} else if patternSegment != concreteSegment {
			// For non-array segments, direct comparison
			return false
		}
	}

	return true
}
