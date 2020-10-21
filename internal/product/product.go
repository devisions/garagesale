package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List returns all known Products.
func List(db *sqlx.DB) ([]Product, error) {

	list := []Product{}

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created FROM products`
	if err := db.Select(&list, q); err != nil {
		return nil, err
	}
	return list, nil
}

// Retrieve returns a single Product.
func Retrieve(db *sqlx.DB, id string) (*Product, error) {

	var p Product
	const q = `SELECT * FROM products WHERE product_id = $1`
	if err := db.Get(&p, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting single product")
	}
	return &p, nil
}
