package handlers

import (
	"log"
	"net/http"

	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger)

	phs := ProductHandlers{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/v1/products", phs.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", phs.Retrieve)
	return app
}
