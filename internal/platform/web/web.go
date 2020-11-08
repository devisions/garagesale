package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
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
	TraceID    string
}

// -------------------------------------------------------
// App specific Web Handler

// AppHandler is the signature that all application handlers will implement.
type AppHandler func(context.Context, http.ResponseWriter, *http.Request) error

// App is the entrypoint for all web apps.
type App struct {
	mux      *chi.Mux
	log      *log.Logger
	mws      []Middleware
	och      *ochttp.Handler
	shutdown chan os.Signal
}

// NewApp knows how to construct the internal state for an App.
func NewApp(logger *log.Logger, shutdown chan os.Signal, mw ...Middleware) *App {

	app := App{
		mux:      chi.NewRouter(),
		log:      logger,
		mws:      mw,
		shutdown: shutdown,
	}
	// Create an OpenCensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	// Configured to use W3C TraceContext standard to set the remote parent
	// if a client request includes the appropriate headers.
	// (see https://w3c.github.com/trace-context)
	app.och = &ochttp.Handler{
		Handler:     app.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}
	return &app
}

// Handle connects a method and URL pattern to a particular HTTP handler.
func (a *App) Handle(method, pattern string, ah AppHandler, mw ...Middleware) {

	// First, wrap handler specific middleware around this app handler.
	ah = wrapMiddlewares(mw, ah)

	// Next, wrap the application's general middlewares around this app handler.
	ah = wrapMiddlewares(a.mws, ah)

	// Create a function that conforms to the std lib definition of a handler.
	// This will be executed when the pattern's route is called.
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		v := Values{Start: time.Now(), TraceID: span.SpanContext().TraceID.String()}
		ctx = context.WithValue(ctx, KeyValues, &v)

		r = r.WithContext(ctx)

		if err := ah(ctx, w, r); err != nil {
			a.log.Printf("ERROR: Unhandled: %v", err)
			if IsShutdown(err) {
				a.SignalShutdown()
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.och.ServeHTTP(w, r)
}

// SignalShutdown is used for gracefully shutdown the app
// when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.log.Println("A handler returned an integrity issue error. Shutting down now ...")
	a.shutdown <- syscall.SIGSTOP
}
