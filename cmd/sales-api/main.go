package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "expvar"         // Register the /debug/vars handler.
	_ "net/http/pprof" // Register the /debug/pprof handlers.

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/devisions/garagesale/cmd/sales-api/internal/handlers"
	"github.com/devisions/garagesale/internal/platform/auth"
	"github.com/devisions/garagesale/internal/platform/conf"
	"github.com/devisions/garagesale/internal/platform/database"
	jwt "github.com/dgrijalva/jwt-go"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func main() {

	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	log := log.New(os.Stdout, "[sales] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		Authn struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost:54327"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			DebugAddress    string        `conf:"default:localhost:6060"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api/v2/spans"`
			Service     string  `conf:"default:sales-api"`
			Probability float64 `conf:"default:1"`
		}
	}

	if err := conf.Parse(os.Args[1:], "sales", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
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
		return errors.Wrap(err, "setting up the db conn")
	}
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "talking with db")
	}
	defer db.Close()

	// -----------------------------------------------------------------------
	// Authentication Support

	authenticator, err := createAuth(
		cfg.Authn.PrivateKeyFile,
		cfg.Authn.KeyID,
		cfg.Authn.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

	// -----------------------------------------------------------------------
	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()

	// -----------------------------------------------------------------------
	// Start Debug Server

	go func() {
		log.Printf("Debug service listening on %s", cfg.Web.DebugAddress)
		if err := http.ListenAndServe(cfg.Web.DebugAddress, http.DefaultServeMux); err != nil {
			log.Printf("Debug service ended with %s", err)
		}
	}()

	// -----------------------------------------------------------------------
	// Start API Server

	srv := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(db, authenticator, log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	srvErrs := make(chan error, 1)

	// Starting the server in the background.
	go func() {
		log.Printf("Server is listening on %s\n", srv.Addr)
		srvErrs <- srv.ListenAndServe()
	}()

	// Shutdown channel receives an interrupt or terminate signal from the OS.
	// It is a buffered channel since the signal package requires it.
	shutd := make(chan os.Signal, 1)
	signal.Notify(shutd, os.Interrupt, syscall.SIGTERM)

	// Everything has started. Just waiting for a shutdown signal.
	select {

	case err := <-srvErrs:
		return errors.Wrap(err, "on ListenAndServe")

	case <-shutd:
		log.Println("Shutting down ...")

		// Give existing requests a deadline to complete.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v seconds: %v\n",
				cfg.Web.ShutdownTimeout, err)
			if err := srv.Close(); err != nil {
				return errors.Wrap(err, "while closing the server")
			}
		}
		log.Println("Graceful shutdown complete.")
	}

	return nil
}

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	pubKeyLookupFunc := auth.NewSimpleKeyLookupFunc(keyID, &key.PublicKey)

	return auth.NewAuthenticator(key, keyID, algorithm, pubKeyLookupFunc)
}

func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {

	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
