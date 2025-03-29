package errors

import "fmt"

type AppError struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`

	PrivateMessage string `json:"-"`
	Err            error  `json:"-"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("AppError: %s - %s", e.Message, e.Code)
}

func NewBadRequest(message, privateMsg, code, requestID string, err error) AppError {
	return AppError{
		Success:   false,
		Code:      code,
		Message:   message,
		RequestID: requestID,

		PrivateMessage: privateMsg,
		Err:            err,
	}
}
