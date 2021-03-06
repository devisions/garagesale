package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond marshals a value to JSON and send it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return errors.New("web values missing from context")
	}
	v.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	val, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "json encoding data: %v", data)
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(val); err != nil {
		return errors.Wrap(err, "writing the response to client")
	}
	return nil
}

// RespondError knows how to handle errors going out to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	// If the error is of type *RequestError, let's respond
	// with the included status code and error.
	if webErr, ok := errors.Cause(err).(*RequestError); ok {
		resp := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		return Respond(ctx, w, resp, webErr.Status)
	}

	resp := ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)}
	return Respond(ctx, w, resp, http.StatusInternalServerError)
}
