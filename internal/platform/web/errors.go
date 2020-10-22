package web

// ErrorResponse represents how we respond to clients when something goes wrong.
type ErrorResponse struct {
	Error string `json:"error"`
}

// AppError is used to add web information to a request error
type WebError struct {
	Err    error
	Status int
}

func NewWebError(err error, status int) error {
	return &WebError{Err: err, Status: status}
}

func (e *WebError) Error() string {
	return e.Err.Error()
}
