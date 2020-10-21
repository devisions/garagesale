package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/devisions/garagesale/internal/platform/conf"
	"github.com/devisions/garagesale/internal/platform/database"
	"github.com/devisions/garagesale/internal/schema"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// -----------------------------------------------------------------------
	// Configuration

	var cfg struct {
		DB struct {
			Username   string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost:54327"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// -----------------------------------------------------------------------
	// Dependencies

	db, err := database.Open(database.Config{
		Username:   cfg.DB.Username,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		DBName:     cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "tryting to talk with the db")
	}

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying db migrations")
		}
		log.Println("Db migration complete")
	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "seeding data into db")
		}
		log.Println("Seed data into db complete")
	}

	return nil
}
