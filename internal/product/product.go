package product

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors for know failure scenarios.
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("provided id is not a valid UUID")
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

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	var p Product
	const q = `SELECT product_id, name, cost, quantity, date_created, date_updated
	           FROM products WHERE product_id = $1`
	if err := db.Get(&p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting single product")
	}
	return &p, nil
}

// Create makes a new Product.
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {

	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}
	const q = `INSERT INTO products 
		 (product_id, name, cost, quantity, date_created, date_updated)
		 VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}

	return &p, nil
}
