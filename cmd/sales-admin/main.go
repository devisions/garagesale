package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/devisions/garagesale/internal/platform/conf"
	"github.com/devisions/garagesale/internal/platform/database"
	"github.com/devisions/garagesale/internal/schema"
)

func main() {

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
				log.Fatalf("main : generating usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
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
		log.Fatal("Failed to talk with the db:", err)
	}

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Failed to apply db migrations", err)
			os.Exit(1)
		}
		log.Println("Db migration complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Failed to seed data into db", err)
			os.Exit(1)
		}
		log.Println("Seed data into db complete")
		return
	}
}
