package utils

import (
	"os"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var logger *log.Logger

func init() {
	logger = log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006/01/02 15:04:05.0.000000000",
	})
}

func Logger() *log.Logger {
	return logger
}

func IsValidUUIDV4(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ValidateStringSlice(fl validator.FieldLevel) bool {
	value := fl.Field()
	if value.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < value.Len(); i++ {
		if value.Index(i).Kind() != reflect.String {
			return false
		}
	}

	return true
}
