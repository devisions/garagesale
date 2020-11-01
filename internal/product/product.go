package product

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/devisions/garagesale/internal/platform/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors for know failure scenarios.
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("provided id is not a valid UUID")
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List returns all known Products.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {

	list := []Product{}

	const q = `SELECT p.*,
			   COALESCE(SUM(s.quantity), 0) AS sold, 
			   COALESCE(SUM(s.paid), 0) AS revenue
			   FROM products AS p
			   LEFT JOIN sales AS s ON p.product_id = s.product_id
			   GROUP BY p.product_id`
	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, errors.Wrap(err, "selecting all products")
	}
	return list, nil
}

// Retrieve returns a single Product.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	var p Product
	const q = `SELECT p.*, 
			   COALESCE(SUM(s.quantity), 0) AS sold, 
			   COALESCE(SUM(s.paid), 0) AS revenue
			   FROM products AS p
			   LEFT JOIN sales AS s ON p.product_id = s.product_id
			   WHERE p.product_id = $1
			   GROUP BY p.product_id`
	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting single product")
	}
	return &p, nil
}

// Create makes a new Product.
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, np NewProduct, now time.Time) (*Product, error) {

	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      user.Subject,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}
	const q = `INSERT INTO products 
		 (product_id, user_id, name, cost, quantity, date_created, date_updated)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := db.ExecContext(ctx, q, p.ID, p.UserID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}

	return &p, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func Update(ctx context.Context, db *sqlx.DB, user auth.Claims, id string, update UpdateProduct, now time.Time) error {

	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	log.Printf("product Update > user.Subject='%+v' p.UserID='%v'", user.Subject, p.UserID)
	// If user is not admin and owner of that product, then action is forbidden.
	if !user.HasRole(auth.RoleAdmin) && user.Subject != p.UserID {
		return ErrForbidden
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE products SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE product_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// Delete removes the product identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM products WHERE product_id = $1`
	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting product %s", id)
	}

	return nil
}
