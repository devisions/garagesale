package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devisions/garagesale/cmd/sales-api/internal/handlers"
	"github.com/devisions/garagesale/internal/platform/database"
)

func main() {

	// -----------------------------------------------------------------------
	// Setup Dependencies

	db, err := database.Open()
	if err != nil {
		log.Fatal("Failed to talk with the db:", err)
	}

	ps := handlers.Product{DB: db}

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
