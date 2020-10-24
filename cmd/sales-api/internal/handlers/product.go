package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/devisions/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ProductHandlers has handler methods for dealing with Products.
type ProductHandlers struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ListProducts gives all products as a list
func (p *ProductHandlers) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve gives a single Product.
func (p *ProductHandlers) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(r.Context(), p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewWebError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewWebError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decodes a JSON from the POST request and create a new Product.
func (p *ProductHandlers) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}
