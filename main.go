package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devisions/garagesale/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	log.Println("main > Starting up ...")
	defer log.Println("main > Exit.")

	// -----------------------------------------------------------------------
	// Setup Dependencies

	db, err := openDB()
	if err != nil {
		log.Fatal("Failed to talk with the db:", err)
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Failed to apply db migrations", err)
		}
		log.Println("Db migration complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Failed to seed data into db", err)
		}
		log.Println("Seed data into db complete")
		return
	}

	ps := ProductService{db: db}

	srv := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ps.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	srvErrs := make(chan error, 1)

	// Starting the server in the background.
	go func() {
		log.Printf("main > Server is listening on %s\n", srv.Addr)
		srvErrs <- srv.ListenAndServe()
	}()

	// Shutdown channel receives an interrupt or terminate signal from the OS.
	// It is a buffered channel since the signal package requires it.
	shutd := make(chan os.Signal, 1)
	signal.Notify(shutd, os.Interrupt, syscall.SIGTERM)

	// Everything has started. Just waiting for a shutdown signal.
	select {

	case err := <-srvErrs:
		log.Fatalf("main > error on ListenAndServe: %s", err)

	case <-shutd:
		log.Println("main > Shutting down ...")

		// Give existing requests a deadline to complete.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("main > error: Graceful shutdown did not complete in %v seconds: %v\n", timeout, err)
			if err := srv.Close(); err != nil {
				log.Fatalf("main > error: While closing the server: %s", err)
			}
		}
		log.Println("main > Graceful shutdown complete.")
	}
}

// Product is something we sell.
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Quantity    int       `json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// ProductService has handler methods for dealing with Products.
type ProductService struct {
	db *sqlx.DB
}

// ListProducts gives all products as a list
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {

	list := []Product{}

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created FROM products`
	if err := p.db.Select(&list, q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error on querying products on db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error on marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte(data)); err != nil {
		log.Println("error on responding", err)
	}
}

func openDB() (*sqlx.DB, error) {

	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost:54327",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
}
