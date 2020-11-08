package middleware

import (
	"context"
	"net/http"

	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Panics recovers from panics and converts the panic to an error,
// so it is reported in Metrics and handled in Errors middlewares.
func Panics() web.Middleware {

	// The actual middleware function that will be executed.
	f := func(after web.AppHandler) web.AppHandler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			ctx, span := trace.StartSpan(ctx, "internal.middleware.Panics")
			defer span.End()

			// Defer a function to recover from a panic and set the err return variable after the fact.
			// Using the errors package it will generate a stack trace.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(ctx, w, r)
		}
		return h

	}
	return f
}
