package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devisions/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

// Product has handler methods for dealing with Products.
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ListProducts gives all products as a list
func (p *Product) List(w http.ResponseWriter, r *http.Request) {

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
