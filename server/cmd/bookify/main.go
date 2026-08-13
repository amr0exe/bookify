package main

import (
	"log"

	"github.com/amr0exe/bookify/internal/config"
	store "github.com/amr0exe/bookify/internal/db"
	"github.com/amr0exe/bookify/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application startup error: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := store.Connect(cfg.DB_URL)
	if err != nil {
		return nil
	}

	engine := router.SetUp(db, cfg)

	log.Printf("Server starting on port %s ...", cfg.PORT)
	return engine.Run(":" + cfg.PORT)
}
