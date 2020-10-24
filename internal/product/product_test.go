package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/devisions/garagesale/internal/product"
	"github.com/devisions/garagesale/internal/schema"
	"github.com/devisions/garagesale/internal/tests"
	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {

	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()

	np := product.NewProduct{Name: "Comic Books", Cost: 10, Quantity: 20}
	now := time.Now().UTC()

	saved, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	fetched, err := product.Retrieve(ctx, db, saved.ID)
	if err != nil {
		t.Fatalf("could not retrieve product: %v", err)
	}

	if diff := cmp.Diff(saved, fetched); diff != "" {
		t.Fatalf("fetched product did not match saved. diff: %v", diff)
	}
}

func TestProductList(t *testing.T) {

	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("failed on listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
