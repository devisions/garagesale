package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond marshals a value to JSON and send it to the client.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	val, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshaling to json")
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(val); err != nil {
		return errors.Wrap(err, "writing the response to client")
	}
	return nil
}
