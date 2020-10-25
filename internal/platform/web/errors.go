package web

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
