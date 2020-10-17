package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Println("main > Starting up ...")
	defer log.Println("main > Exit.")

	srv := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ListProducts),
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
		log.Fatalf("main > error: On ListenAndServe: %s", err)

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
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ListProducts gives all products as a list
func ListProducts(w http.ResponseWriter, r *http.Request) {

	list := []Product{}
	list = append(list,
		Product{Name: "Comic Books", Cost: 75, Quantity: 50},
		Product{Name: "McDonald's Toys", Cost: 25, Quantity: 120})

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error: marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte(data)); err != nil {
		log.Println("error: responding", err)
	}

}
