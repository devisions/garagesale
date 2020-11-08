package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/devisions/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// ErrorHandler handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func ErrorHandler(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.AppHandler) web.AppHandler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.middleware.ErrorHandler")
			defer span.End()

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {
				// Log the error.
				log.Printf("ERROR: %+v", err)
				// Respond to the error.
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}
			}

			// Return nil to indicate the error has been handled.
			return nil
		}
		return h
	}

	return f
}
