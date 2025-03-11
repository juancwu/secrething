package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// QueryParser provides type-safe methods to extract and convert query parameters
type QueryParser struct {
	ctx echo.Context
}

// NewQueryParser creates a new query parser from echo.Context
func NewQueryParser(ctx echo.Context) *QueryParser {
	return &QueryParser{ctx: ctx}
}

// String returns the value for the given key as a string
// If the key doesn't exist, it returns the defaultValue
func (qp *QueryParser) String(key string, defaultValue string) string {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// Int returns the value for the given key as an integer
// If the key doesn't exist or cannot be converted, it returns the defaultValue
func (qp *QueryParser) Int(key string, defaultValue int) int {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// Int64 returns the value for the given key as an int64
// If the key doesn't exist or cannot be converted, it returns the defaultValue
func (qp *QueryParser) Int64(key string, defaultValue int64) int64 {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}

	int64Val, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int64Val
}

// Float64 returns the value for the given key as a float64
// If the key doesn't exist or cannot be converted, it returns the defaultValue
func (qp *QueryParser) Float64(key string, defaultValue float64) float64 {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return floatVal
}

// Bool returns the value for the given key as a boolean
// If the key doesn't exist or cannot be converted, it returns the defaultValue
func (qp *QueryParser) Bool(key string, defaultValue bool) bool {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// Time returns the value for the given key as time.Time using the given layout
// If the key doesn't exist or cannot be parsed, it returns the defaultValue
func (qp *QueryParser) Time(key string, layout string, defaultValue time.Time) time.Time {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return defaultValue
	}

	timeVal, err := time.Parse(layout, val)
	if err != nil {
		return defaultValue
	}
	return timeVal
}

// StringSlice returns all values for the given key as a string slice
// If the key doesn't exist, it returns an empty slice
func (qp *QueryParser) StringSlice(key string) []string {
	return qp.ctx.QueryParams()[key]
}

// IntSlice returns all values for the given key as an int slice
// Invalid values are omitted from the result
func (qp *QueryParser) IntSlice(key string) []int {
	strValues := qp.ctx.QueryParams()[key]
	result := make([]int, 0, len(strValues))

	for _, v := range strValues {
		if val, err := strconv.Atoi(v); err == nil {
			result = append(result, val)
		}
	}

	return result
}

// Has checks if the query contains the specified key
func (qp *QueryParser) Has(key string) bool {
	return qp.ctx.QueryParam(key) != ""
}

// Required returns the value for the given key or an error if it doesn't exist
func (qp *QueryParser) Required(key string) (string, error) {
	val := qp.ctx.QueryParam(key)
	if val == "" {
		return "", fmt.Errorf("missing required query parameter: %s", key)
	}
	return val, nil
}

