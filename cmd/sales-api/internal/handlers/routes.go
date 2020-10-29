package handlers

import (
	"log"
	"net/http"

	"github.com/devisions/garagesale/internal/middleware"
	"github.com/devisions/garagesale/internal/platform/auth"
	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs a handler that knows about all API routes.
func API(db *sqlx.DB, authenticator *auth.Authenticator, logger *log.Logger) http.Handler {

	app := web.NewApp(logger,
		middleware.Logger(logger), middleware.Errors(logger), middleware.Metrics(),
	)

	hc := HealthCheck{DB: db}

	app.Handle(http.MethodGet, "/v1/health", hc.Health)

	phs := ProductHandlers{db: db, log: logger}

	uhs := UserHandlers{db: db, authenticator: authenticator}

	app.Handle(http.MethodGet, "/v1/users/token", uhs.Token)

	app.Handle(http.MethodGet, "/v1/products", phs.List)
	app.Handle(http.MethodPost, "/v1/products", phs.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", phs.Retrieve)
	app.Handle(http.MethodPut, "/v1/products/{id}", phs.Update)
	app.Handle(http.MethodDelete, "/v1/products/{id}", phs.Delete)

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", phs.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", phs.ListSales)

	return app
}
