package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
		p.Log.Println("error on querying products on db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error on marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte(data)); err != nil {
		p.Log.Println("error on responding", err)
	}
}

// Retrieve gives a single Product.
func (p *ProductHandlers) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error on querying a product on db", err)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error on marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte(data)); err != nil {
		p.Log.Println("error on responding", err)
	}
}

// Create decodes a JSON from the POST request and create a new Product.
func (p *ProductHandlers) Create(w http.ResponseWriter, r *http.Request) {

	var np product.NewProduct
	if err := json.NewDecoder(r.Body).Decode(&np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Printf("Error decoding NewProduct request body: %s\n", err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error on writing the response", err)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error on marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(data)); err != nil {
		p.Log.Println("error on writing the response", err)
	}
}
