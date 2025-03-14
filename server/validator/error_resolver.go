package validator

import (
	"strings"
)

// ErrorResolver handles translating validation errors to custom messages
type ErrorResolver struct {
	fieldErrors     map[string]map[string]string // map[field]map[tag]message
	arrayPathErrors map[string]map[string]string // map[arrayPath]map[tag]message for paths with [] notation
	leafFieldErrors map[string]map[string]string // map[leafName]map[tag]message for fallback
	defaultErrors   map[string]string            // map[tag]message
	defaultMessage  string
}

// NewErrorTranslator creates a new ErrorTranslator
func NewErrorTranslator() *ErrorResolver {
	return &ErrorResolver{
		fieldErrors:     make(map[string]map[string]string),
		arrayPathErrors: make(map[string]map[string]string),
		leafFieldErrors: make(map[string]map[string]string),
		defaultErrors:   make(map[string]string),
		defaultMessage:  "Invalid value",
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
// The field path is normalized for consistent lookup during validation
// If the field path contains empty bracket notation "[]", it's stored as an array path
func (t *ErrorResolver) SetFieldError(field, tag, message string) {
	// Normalize the field path for consistent lookup
	normalizedPath := normalizePath(field)

	// Check if this contains [] array notation
	if strings.Contains(normalizedPath, "[]") {
		// Store it as an array path pattern
		if _, ok := t.arrayPathErrors[normalizedPath]; !ok {
			t.arrayPathErrors[normalizedPath] = make(map[string]string)
		}
		t.arrayPathErrors[normalizedPath][tag] = message

		// Also store by leaf name for fallback
		leafName := extractLeafName(field)
		if leafName != "" && leafName != normalizedPath {
			if _, ok := t.leafFieldErrors[leafName]; !ok {
				t.leafFieldErrors[leafName] = make(map[string]string)
			}
			t.leafFieldErrors[leafName][tag] = message
		}

		return
	}

	// Store by normalized path
	if _, ok := t.fieldErrors[normalizedPath]; !ok {
		t.fieldErrors[normalizedPath] = make(map[string]string)
	}
	t.fieldErrors[normalizedPath][tag] = message

	// Also store by leaf name for fallback lookups
	leafName := extractLeafName(field)
	if leafName != "" && leafName != normalizedPath {
		if _, ok := t.leafFieldErrors[leafName]; !ok {
			t.leafFieldErrors[leafName] = make(map[string]string)
		}
		t.leafFieldErrors[leafName][tag] = message
	}
}

// SetDefaultError sets a default error message for a validation tag
func (t *ErrorResolver) SetDefaultError(tag, message string) {
	t.defaultErrors[tag] = message
}

// SetDefaultMessage sets the default error message for all validations
func (t *ErrorResolver) SetDefaultMessage(message string) {
	t.defaultMessage = message
}

// Translate translates a validation error to a custom message
func (t *ErrorResolver) Translate(field string, tag string) string {
	// Normalize the field path for consistent lookup
	normalizedPath := normalizePath(field)

	// First, check if there's a custom message for the normalized full field path and tag
	if fieldMessages, ok := t.fieldErrors[normalizedPath]; ok {
		if message, ok := fieldMessages[tag]; ok {
			return message
		}
	}

	// Second, check if there's a matching array path pattern
	for patternPath, tagMessages := range t.arrayPathErrors {
		// Skip patterns that don't have a message for this tag
		message, hasTagMessage := tagMessages[tag]
		if !hasTagMessage {
			continue
		}

		// Check if the normalized path matches this array pattern
		if matchArrayPath(normalizedPath, patternPath) {
			return message
		}
	}

	// Third, check if there's a leaf field name match
	leafName := extractLeafName(field)
	if leafName != "" && leafName != normalizedPath {
		// Check in the dedicated leaf field errors map first
		if fieldMessages, ok := t.leafFieldErrors[leafName]; ok {
			if message, ok := fieldMessages[tag]; ok {
				return message
			}
		}

		// For backward compatibility, also check in the main fieldErrors map
		if fieldMessages, ok := t.fieldErrors[leafName]; ok {
			if message, ok := fieldMessages[tag]; ok {
				return message
			}
		}
	}

	// Check if there's a default message for this tag
	if message, ok := t.defaultErrors[tag]; ok {
		return message
	}

	// Return the default message
	return t.defaultMessage
}

// Clone creates a copy of the ErrorTranslator
func (t *ErrorResolver) Clone() *ErrorResolver {
	clone := NewErrorTranslator()
	clone.defaultMessage = t.defaultMessage

	// Copy default errors
	for tag, msg := range t.defaultErrors {
		clone.defaultErrors[tag] = msg
	}

	// Copy field errors
	for field, tagMsgs := range t.fieldErrors {
		clone.fieldErrors[field] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.fieldErrors[field][tag] = msg
		}
	}

	// Copy array path errors
	for path, tagMsgs := range t.arrayPathErrors {
		clone.arrayPathErrors[path] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.arrayPathErrors[path][tag] = msg
		}
	}

	// Copy leaf field errors
	for field, tagMsgs := range t.leafFieldErrors {
		clone.leafFieldErrors[field] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.leafFieldErrors[field][tag] = msg
		}
	}

	return clone
}
