package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	healthyReport := HealthReport{
		DB: "healthy",
	}
	unhealthyReport := HealthReport{
		DB: "unhealthy",
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// connect to test db
	err := store.Connect(os.Getenv("DB_URL"))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v\n", err)
	}
	if assert.NoError(t, handleGetHealth(c), "handleGetHealth returned an error.") {
		// get response body
		data, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v\n", err)
		}
		var resBody HealthReport
		err = json.Unmarshal(data, &resBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v\n", err)
		}
		if !assert.Equal(t, healthyReport, resBody) {
			t.Errorf("Health report does not match the exptected report. Expected: %v\nReceived: %v\n", healthyReport, resBody)
		}
	}

	// test when db is not healhty
	err = store.Close()
	if err != nil {
		t.Fatalf("Failed to close database conntection: %v\n", err)
	}
	if assert.NoError(t, handleGetHealth(c), "handleGetHealth returned an error.") {
		// get response body
		data, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v\n", err)
		}
		var resBody HealthReport
		err = json.Unmarshal(data, &resBody)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %v\n", err)
		}
		if !assert.Equal(t, unhealthyReport, resBody) {
			t.Errorf("Health report does not match the exptected report. Expected: %v\nReceived: %v\n", healthyReport, resBody)
		}
	}
}
