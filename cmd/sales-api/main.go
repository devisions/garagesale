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

	"github.com/devisions/garagesale/cmd/sales-api/internal/handlers"
	"github.com/devisions/garagesale/internal/platform/conf"
	"github.com/devisions/garagesale/internal/platform/database"
)

func main() {

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost:54327"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}

	if err := conf.Parse(os.Args[1:], "sales", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error : generating config for output : %v", err)
	}
	log.Printf("main : Config :\n%v\n", out)

	// -----------------------------------------------------------------------
	// Setup Dependencies

	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		DBName:     cfg.DB.Name,
		Username:   cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatal("Failed to talk with the db:", err)
	}

	ps := handlers.Product{DB: db}

	srv := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(ps.List),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("main > error: Graceful shutdown did not complete in %v seconds: %v\n",
				cfg.Web.ShutdownTimeout, err)
			if err := srv.Close(); err != nil {
				log.Fatalf("main > error: While closing the server: %s", err)
			}
		}
		log.Println("main > Graceful shutdown complete.")
	}

}
