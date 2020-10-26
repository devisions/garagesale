package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// -------------------------------------------------------
// Request Tracking

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// -------------------------------------------------------
// App specific Web Handler

// AppHandler is the signature that all application handlers will implement.
type AppHandler func(context.Context, http.ResponseWriter, *http.Request) error

// App is the entrypoint for all web apps.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mws []Middleware
}

// NewApp knows how to construct the internal state for an App.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mws: mw,
	}
}

// Handle connects a method and URL pattern to a particular HTTP handler.
func (a *App) Handle(method, pattern string, ah AppHandler) {

	// Wrapping the app middlewares around this handler.
	ah = wrapMiddlewares(a.mws, ah)

	// Create a function that conforms to the std lib definition of a handler.
	// This will be executed when the pattern's route is called.
	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			Start: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)
		r = r.WithContext(ctx)

		if err := ah(ctx, w, r); err != nil {
			a.log.Printf("ERROR: Unhandled error: %v", err)
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
