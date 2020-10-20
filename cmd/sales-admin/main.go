package main

import (
	"flag"
	"log"

	"github.com/devisions/garagesale/internal/platform/database"
	"github.com/devisions/garagesale/internal/schema"
)

func main() {

	// -----------------------------------------------------------------------
	// Setup Dependencies

	db, err := database.Open()
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
}
