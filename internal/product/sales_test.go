package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/devisions/garagesale/internal/platform/auth"
	"github.com/devisions/garagesale/internal/product"
	"github.com/devisions/garagesale/internal/tests"
)

func TestSales(t *testing.T) {

	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()
	now := time.Now().UTC()

	newComics := product.NewProduct{Name: "Comic Books", Cost: 10, Quantity: 20}

	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // A random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	comics, err := product.Create(ctx, db, claims, newComics, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	newToys := product.NewProduct{Name: "Toys", Cost: 40, Quantity: 30}

	toys, err := product.Create(ctx, db, claims, newToys, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	{ // testing Sales: add and list

		ns := product.NewSale{
			Quantity: 3,
			Paid:     60,
		}

		s, err := product.AddSale(ctx, db, ns, comics.ID, now)
		if err != nil {
			t.Fatalf("adding sale: %s", err)
		}

		// Puzzles should show the 1 sale.
		sales, err := product.ListSales(ctx, db, comics.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 1, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}

		if exp, got := s.ID, sales[0].ID; exp != got {
			t.Fatalf("expected first sale ID %v, got %v", exp, got)
		}

		// Toys should have 0 sales.
		sales, err = product.ListSales(ctx, db, toys.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 0, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}
	}
}
