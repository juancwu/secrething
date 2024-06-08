package router

// apiResponse is a generic json response for all APIs
type apiResponse struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}
