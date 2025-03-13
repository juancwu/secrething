package validator

// setDefaultMessages sets commonly used validation error messages
func setDefaultMessages(translator *ErrorTranslator) {
	defaults := map[string]string{
		"required": "This field is required",
		"email":    "Must be a valid email address",
		"min":      "Value must be greater than or equal to the minimum",
		"max":      "Value must be less than or equal to the maximum",
		"len":      "Must have the exact required length",
		"eq":       "Value must be equal to the required value",
		"ne":       "Value cannot be equal to the specified value",
		"oneof":    "Must be one of the available options",
		"url":      "Must be a valid URL",
		"alpha":    "Must contain only letters",
		"alphanum": "Must contain only letters and numbers",
		"numeric":  "Must be a valid numeric value",
		"uuid":     "Must be a valid UUID",
		"datetime": "Must be a valid date/time",
		"password": "Password must be at least 8 characters long and contain uppercase, lowercase, digit, and at least one special character (!@#$%^&*()-_=+[]{}|;:'\",.<>/?)",
	}

	for tag, message := range defaults {
		translator.SetDefaultError(tag, message)
	}
}
