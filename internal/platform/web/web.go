package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// AppHandler is the signature that all application handlers will implement.
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// App is the entrypoint for all web apps.
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp knows how to construct the internal state for an App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

// Handle connects a method and URL pattern to a particular HTTP handler.
func (a *App) Handle(method, pattern string, ah AppHandler) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := ah(w, r); err != nil {
			a.log.Printf("ERROR: %v", err)
			if err := RespondError(w, err); err != nil {
				a.log.Printf("ERROR: %v", err)
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
