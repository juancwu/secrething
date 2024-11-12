package tag

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Parse a struct with the 'errormsg' tag that defines messages for the validator.FieldError.
// Returns the message based on the 'errormsg' format and if the validator.FieldError can be
// matched or not. Otherwise, it returns an empty string.
func ParseErrorMsgTag(structType reflect.Type, fieldError validator.FieldError) string {
	var field reflect.StructField
	var found bool

	// regex to match all square brackets
	re, err := regexp.Compile(`\[(\d+)\]`)
	if err != nil {
		log.Error().Err(err).Str("struct_name", structType.Name()).Msg("Failed to compiled regex to replace square brackets.")
		return ""
	}
	// remove any square brackets, they come up when there is a field of type slice.
	fields := strings.Split(re.ReplaceAllString(fieldError.StructNamespace(), ""), ".")
	for i, fieldName := range fields {
		// continue since this is the struct name
		if i == 0 {
			continue
		}
		field, found = structType.FieldByName(fieldName)
		if !found {
			return ""
		}
		// get the new struct
		structType = field.Type
		if structType.Kind() == reflect.Slice {
			structType = structType.Elem()
		}
	}

	errormsg := field.Tag.Get("errormsg")

	validationTags := strings.Split(errormsg, ",")
	if len(validationTags) == 1 {
		parts := strings.Split(validationTags[0], "=")
		if len(parts) == 1 {
			// treat as default global message
			return parts[0]
		} else if len(parts) == 2 && containsFieldTag(parts[0], fieldError.Tag()) {
			return parts[1]
		}
	} else if len(validationTags) > 1 {
		defaultMsg := ""
		for _, tag := range validationTags {
			parts := strings.Split(tag, "=")
			if len(parts) == 2 && containsFieldTag(parts[0], fieldError.Tag()) {
				return parts[1]
			}
			if parts[0] == "__default" {
				defaultMsg = parts[1]

			}
		}
		return defaultMsg
	}
	return ""
}

// Checks if the errormsg tag whether or not contains the field error tag.
// This method helps with checking combined errormsg tags.
func containsFieldTag(errorMsgTag, fieldErrorTag string) bool {
	tags := strings.Split(errorMsgTag, "|")
	for _, tag := range tags {
		if tag == fieldErrorTag {
			return true
		}
	}
	return false
}
