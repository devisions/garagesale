package web

import "github.com/pkg/errors"

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse represents how we respond to clients when something goes wrong.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request
// through the application with web specific context.
type RequestError struct {
	Err    error
	Status int
	Fields []FieldError
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

// shutdown is a support type for the graceful shutdown of the service.
type shutdown struct {
	Message string
}

// Error is what shutdown implements for being used as an error.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error implementation.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// IsShutdown checks and tells if a shutdown error error
// is contained in the specified error value.
func IsShutdown(err error) bool {

	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
