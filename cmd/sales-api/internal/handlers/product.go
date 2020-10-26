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
	db  *sqlx.DB
	log *log.Logger
}

// ListProducts gives all products as a list
func (p *ProductHandlers) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(r.Context(), p.db)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve gives a single Product.
func (p *ProductHandlers) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(r.Context(), p.db, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}

	return web.Respond(r.Context(), w, prod, http.StatusOK)
}

// Create decodes a JSON from the POST request and create a new Product.
func (p *ProductHandlers) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(r.Context(), p.db, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, prod, http.StatusCreated)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (p *ProductHandlers) Update(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(r.Context(), p.db, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (p *ProductHandlers) Delete(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := product.Delete(r.Context(), p.db, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *ProductHandlers) AddSale(w http.ResponseWriter, r *http.Request) error {

	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(r.Context(), p.db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(r.Context(), w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *ProductHandlers) ListSales(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	list, err := product.ListSales(r.Context(), p.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}
