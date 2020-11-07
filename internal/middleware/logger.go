package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/devisions/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// RequestLogger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.AppHandler) web.AppHandler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.middleware.RequestLogger")
			defer span.End()

			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web values missing from context")
			}

			// Run the handler chain and catch any propagated error.
			err := before(ctx, w, r)

			log.Printf("%s | %d | %s %s -> %s (%s)",
				v.TraceID, v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the (possible) error to be handled further up the chain.
			return err
		}
		return h
	}

	return f
}
