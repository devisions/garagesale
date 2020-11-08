package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/devisions/garagesale/internal/middleware"
	"github.com/devisions/garagesale/internal/platform/auth"
	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs a handler that knows about all API routes.
func API(db *sqlx.DB, authenticator *auth.Authenticator, logger *log.Logger, shutdown chan os.Signal) http.Handler {

	app := web.NewApp(logger, shutdown,
		middleware.RequestLogger(logger),
		middleware.ErrorHandler(logger),
		middleware.Metrics(),
		middleware.Panics(),
	)

	hc := HealthCheck{DB: db}

	app.Handle(http.MethodGet, "/v1/health", hc.Health)

	phs := ProductHandlers{db: db, log: logger}

	uhs := UserHandlers{db: db, authenticator: authenticator}

	app.Handle(http.MethodGet, "/v1/users/token", uhs.Token)

	app.Handle(http.MethodGet, "/v1/products", phs.List, middleware.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/products", phs.Create, middleware.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/products/{id}", phs.Retrieve, middleware.Authenticate(authenticator))
	app.Handle(http.MethodPut, "/v1/products/{id}", phs.Update, middleware.Authenticate(authenticator))
	app.Handle(http.MethodDelete, "/v1/products/{id}", phs.Delete, middleware.Authenticate(authenticator), middleware.HasRole(auth.RoleAdmin))

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", phs.AddSale, middleware.Authenticate(authenticator), middleware.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", phs.ListSales, middleware.Authenticate(authenticator))

	return app
}
