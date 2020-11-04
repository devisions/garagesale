package middleware

import (
	"context"
	"errors"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"

	"github.com/devisions/garagesale/internal/platform/web"
)

// Logger will log a line for every request.
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.AppHandler) web.AppHandler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(r.Context(), "internal.middleware.Logger")
			defer span.End()

			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web values missing from context")
			}

			// Run the handler chain and catch any propagated error.
			err := before(ctx, w, r)

			log.Printf("%v %s %s (%v)", v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start))

			// Return the (possible) error to be handled further up the chain.
			return err
		}
		return h
	}

	return f
}
