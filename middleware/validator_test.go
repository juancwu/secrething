package middleware

import (
	"konbini/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestModel mirrors the model in types/errors_test.go
type TestModel struct {
	Name     string   `json:"name" validate:"required" errormsg:"Name is mandatory"`
	Email    string   `json:"email" validate:"required,email" errormsg:"required=Email is required;email=Please provide a valid email"`
	Age      int      `json:"age" validate:"min=18,max=100" errormsg:"min=Must be at least 18 years old;max=Must not be older than 100 years"`
	Hobbies  []string `json:"hobbies" validate:"required,min=1" errormsg:"__default=At least one hobby is required"`
	Password string   `json:"password" validate:"required,min=7" errormsg:"required|min=Password must be at least 8 characters long"`
}

func TestBindAndValidate(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedErrors map[string]string
	}{
		{
			name:           "Valid payload",
			payload:        `{"name":"John Doe","email":"john@example.com","age":25,"hobbies":["reading"],"password":"password123"}`,
			expectedStatus: http.StatusOK,
			expectedErrors: nil,
		},
		{
			name:           "Invalid payload",
			payload:        `{"email":"invalid","age":15,"password":"123"}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: map[string]string{
				"name":     "Name is mandatory",
				"email":    "Please provide a valid email",
				"age":      "Must be at least 18 years old",
				"hobbies":  "At least one hobby is required",
				"password": "Password must be at least 8 characters long",
			},
		},
		{
			name:           "Invalid JSON",
			payload:        `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}

			middleware := BindAndValidate(reflect.TypeOf(TestModel{}))
			err := middleware(handler)(c)

			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)

				// Verify the model was properly bound
				model, ok := c.Get(REQUEST_MODEL_CTX_KEY).(*TestModel)
				assert.True(t, ok)
				assert.NotNil(t, model)
				assert.Equal(t, "John Doe", model.Name)
				assert.Equal(t, "john@example.com", model.Email)
			} else {
				assert.Error(t, err)
				httpError, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, httpError.Code)

				if tt.expectedErrors != nil {
					errResp, ok := httpError.Message.(*types.ErrorResponse)
					assert.True(t, ok)

					validationErrors, ok := errResp.Details.([]types.ValidationError)
					assert.True(t, ok)

					// Verify each expected error message
					errorMap := make(map[string]string)
					for _, ve := range validationErrors {
						errorMap[ve.Field] = ve.Message
					}

					for field, expectedMsg := range tt.expectedErrors {
						assert.Equal(t, expectedMsg, errorMap[field])
					}
				}
			}
		})
	}
}

func TestBindAndValidateEdgeCases(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		payload        string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Empty payload",
			payload:        ``,
			contentType:    echo.MIMEApplicationJSON,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Wrong content type",
			payload:        `{"name": "John"}`,
			contentType:    echo.MIMETextPlain,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, tt.contentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}

			structType := reflect.TypeOf(TestModel{})

			middleware := BindAndValidate(structType)
			err := middleware(handler)(c)

			assert.Error(t, err)
			httpError, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedStatus, httpError.Code)
		})
	}
}
