package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/devisions/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// ProductHandlers has handler methods for dealing with Products.
type ProductHandlers struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ListProducts gives all products as a list
func (p *ProductHandlers) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Error on querying products on db", err)
		return
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		p.Log.Println("Error responding", err)
		return
	}
}

// Retrieve gives a single Product.
func (p *ProductHandlers) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Error on querying a product on db", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusOK); err != nil {
		p.Log.Println("Error responding", err)
		return
	}
}

// Create decodes a JSON from the POST request and create a new Product.
func (p *ProductHandlers) Create(w http.ResponseWriter, r *http.Request) {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Printf("Error decoding NewProduct request body: %s\n", err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Error inserting into db", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("Error responding", err)
		return
	}
}
