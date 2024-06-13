package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// writeApiErrorJSON is a helper function to write a generic api error response.
//
// Sample JSON:
//
//	{
//	  status_code: 500,
//	  message: "internal server error (requestId)"
//	}
func writeApiErrorJSON(c echo.Context, requestId string) error {
	return c.JSON(
		http.StatusInternalServerError,
		apiResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("internal server error (%s)", requestId),
		},
	)
}

func verifySignedChallenge(challenge, signature, pubKey string) error {
	return nil
}
