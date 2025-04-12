package utils

import (
	"strings"
	"time"
)

// FormatRFC3339NanoFixed formats a time to RFC3339Nano with exactly 9 decimal places
// This ensures the output string is always 30 characters long. The time is always
// converted to UTC.
func FormatRFC3339NanoFixed(t time.Time) string {
	// Format the time to RFC3339Nano
	s := t.UTC().Format(time.RFC3339Nano)

	// Find the position of the decimal point
	dotIndex := strings.LastIndex(s, ".")
	if dotIndex == -1 {
		// If no decimal point is found (shouldn't happen with RFC3339Nano)
		return s[:len(s)-1] + ".000000000Z"
	}

	// Calculate how many digits we currently have
	currentDigits := len(s) - dotIndex - 2 // -2 for the dot and Z

	// If we have exactly 9 digits, return as is
	if currentDigits == 9 {
		return s
	}

	// If we have fewer than 9 digits, pad with zeros
	if currentDigits < 9 {
		return s[:len(s)-1] + strings.Repeat("0", 9-currentDigits) + "Z"
	}

	// If we have more than 9 digits (shouldn't happen), truncate
	return s[:dotIndex+10] + "Z"
}

// ParseRFC3339NanoStr parses the timestamp strings created with FormatRFC3339NanoFixed.
func ParseRFC3339NanoStr(t string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, t)
}
