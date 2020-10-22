package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Decode looks for a JSON document in the request body
// and unmarshals it into the provided dest.
func Decode(r *http.Request, dest interface{}) error {

	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return errors.Wrap(err, "json decoding request body")
	}
	return nil
}
