package main

import (
	"context"
	"fmt"
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
		Handler:      http.HandlerFunc(Echo),
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

// Echo just tells you about the request you made.
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You asked", r.Method, r.URL.Path)
	time.Sleep(3 * time.Second)
}
